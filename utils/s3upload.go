package utils

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

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

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &fileName,
		Body:   strings.NewReader(string(fileBytes)),
	})
	if err != nil {
		return "", fmt.Errorf("erro ao enviar arquivo para S3: %w", err)
	}

	publicURL := fmt.Sprintf("https://%s.%s/%s", bucketName, endpoint, fileName)
	return publicURL, nil
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

// S3Presigner define interface para geração de link pré-assinado
// Pode ser implementada por um mock nos testes

type S3Presigner interface {
	PresignGetObject(ctx context.Context, bucket, key string, expires time.Duration) (string, error)
}

// RealS3Presigner implementa S3Presigner usando AWS SDK
// (código real pode ser movido do controller)
type RealS3Presigner struct{}

func (r *RealS3Presigner) PresignGetObject(ctx context.Context, bucket, key string, expires time.Duration) (string, error) {
	endpoint := os.Getenv("AWS_ENDPOINT")
	region := os.Getenv("AWS_REGION")
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return "", err
	}
	customResolver := s3.EndpointResolverFunc(func(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           "https://" + endpoint,
			SigningRegion: region,
		}, nil
	})
	httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.EndpointResolver = customResolver
		o.UsePathStyle = true
		o.HTTPClient = httpClient
	})
	presignClient := s3.NewPresignClient(s3Client)
	presignInput := &s3.GetObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)}
	presignResult, err := presignClient.PresignGetObject(ctx, presignInput, func(opts *s3.PresignOptions) { opts.Expires = expires })
	if err != nil {
		return "", err
	}
	return presignResult.URL, nil
}

// MockS3Presigner para testes
// Retorna sempre uma URL fake

type MockS3Presigner struct{}

func (m *MockS3Presigner) PresignGetObject(ctx context.Context, bucket, key string, expires time.Duration) (string, error) {
	return "https://mock-s3.local/" + key + "?mock-presigned", nil
}
