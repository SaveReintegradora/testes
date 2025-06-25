package utils

import (
	"context"
	"fmt"
	"io"
	"os"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/s3blob" // Importa o driver S3
)

// UploadToS3 faz upload de um arquivo para o S3/Wasabi e retorna a URL pública
func UploadToS3(ctx context.Context, fileName string, file io.Reader) (string, error) {
	fmt.Println(">>> Entrou na função UploadToS3")
	// Pegue as variáveis de ambiente
	bucketName := os.Getenv("AWS_BUCKET_NAME")
	endpoint := os.Getenv("AWS_ENDPOINT")
	region := os.Getenv("AWS_REGION")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	// Logs para depuração
	fmt.Println("Bucket:", bucketName)
	fmt.Println("Endpoint:", endpoint)
	fmt.Println("Region:", region)
	fmt.Println("AccessKey:", accessKey)
	fmt.Println("SecretKey length:", len(secretKey))

	// Monta a URL do bucket para o gocloud.dev
	bucketURL := fmt.Sprintf("s3://%s?region=%s&endpoint=%s", bucketName, region, endpoint)
	fmt.Println("Abrindo bucket:", bucketURL)

	// Abre o bucket
	bucket, err := blob.OpenBucket(ctx, bucketURL)
	if err != nil {
		fmt.Println("Erro ao abrir bucket:", err)
		return "", fmt.Errorf("erro ao abrir bucket: %w", err)
	}
	defer bucket.Close()

	fmt.Println("Enviando arquivo:", fileName)
	// Lê o conteúdo do arquivo para um buffer []byte
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Erro ao ler arquivo:", err)
		return "", fmt.Errorf("erro ao ler arquivo: %w", err)
	}
	// Faz upload do arquivo
	err = bucket.WriteAll(ctx, fileName, fileBytes, nil)
	if err != nil {
		fmt.Println("Erro ao enviar arquivo:", err)
		return "", fmt.Errorf("erro ao enviar arquivo: %w", err)
	}
	fmt.Println("Upload concluído!")

	// Monta a URL pública (ajuste conforme sua configuração)
	publicURL := fmt.Sprintf("https://%s/%s/%s", endpoint, bucketName, fileName)
	return publicURL, nil
}
