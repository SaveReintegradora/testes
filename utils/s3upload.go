package utils

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Uploader define interface para upload S3 (real ou mock)
type S3Uploader interface {
	UploadToS3(ctx context.Context, fileName string, file io.Reader) (string, error)
}

// RealS3Uploader implementa S3Uploader usando AWS SDK
type RealS3Uploader struct{}

func (r *RealS3Uploader) UploadToS3(ctx context.Context, fileName string, file io.Reader) (string, error) {
	bucketName := strings.TrimSpace(os.Getenv("AWS_BUCKET_NAME"))
	bucketName = strings.Trim(bucketName, ".")
	endpoint := strings.TrimSpace(os.Getenv("AWS_ENDPOINT"))
	endpoint = strings.Trim(endpoint, ".")
	region := strings.TrimSpace(os.Getenv("AWS_REGION"))
	accessKey := strings.TrimSpace(os.Getenv("AWS_ACCESS_KEY_ID"))
	secretKey := strings.TrimSpace(os.Getenv("AWS_SECRET_ACCESS_KEY"))

	fmt.Println("[DEBUG] bucketName:", bucketName)
	fmt.Println("[DEBUG] endpoint:", endpoint)
	fmt.Println("[DEBUG] region:", region)
	fmt.Println("[DEBUG] accessKey:", accessKey)
	fmt.Println("[DEBUG] secretKey length:", len(secretKey))
	fmt.Println("[DEBUG] fileName:", fileName)

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return "", fmt.Errorf("erro ao carregar config AWS: %w", err)
	}

	customResolver := s3.EndpointResolverFunc(func(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           "https://" + endpoint,
			SigningRegion: region,
		}, nil
	})

	// Adiciona http.Client customizado para ignorar verificação de TLS
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.EndpointResolver = customResolver
		o.UsePathStyle = true // Necessário para Wasabi
		o.HTTPClient = httpClient
	})

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	fmt.Println("[DEBUG] Iniciando upload para:", bucketName, endpoint)

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &fileName,
		Body:   strings.NewReader(string(fileBytes)),
	})
	if err != nil {
		fmt.Println("[DEBUG] Erro ao enviar arquivo para S3:", err)
		return "", fmt.Errorf("erro ao enviar arquivo para S3: %w", err)
	}

	publicURL := fmt.Sprintf("https://%s.%s/%s", bucketName, endpoint, fileName)
	fmt.Println("[DEBUG] publicURL:", publicURL)
	return publicURL, nil
}

// UploadToS3 default aponta para RealS3Uploader (retrocompatibilidade)
func UploadToS3(ctx context.Context, fileName string, file io.Reader) (string, error) {
	uploader := &RealS3Uploader{}
	return uploader.UploadToS3(ctx, fileName, file)
}

// MockS3Uploader para testes automatizados (não faz upload real)
type MockS3Uploader struct {
	LastFileName string
	LastContent  string
	ShouldError  bool
}

func (m *MockS3Uploader) UploadToS3(ctx context.Context, fileName string, file io.Reader) (string, error) {
	m.LastFileName = fileName
	b, _ := io.ReadAll(file)
	m.LastContent = string(b)
	if m.ShouldError {
		return "", fmt.Errorf("erro simulado no mock S3")
	}
	return "https://mock-s3.local/" + fileName, nil
}
