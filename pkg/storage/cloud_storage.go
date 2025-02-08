package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// CloudStorage представляет интерфейс для работы с облачным хранилищем.
type CloudStorage interface {
	UploadFile(bucket, key, filePath string) error
}

// S3Storage реализует CloudStorage для AWS S3.
type S3Storage struct{}

// NewS3Storage создает новый экземпляр S3Storage.
func NewS3Storage() *S3Storage {
	return &S3Storage{}
}

// UploadFile загружает файл в S3.
func (s *S3Storage) UploadFile(bucket, key, filePath string) error {
	// Загрузка конфигурации AWS
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("your-region"))
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Создание клиента S3
	client := s3.NewFromConfig(cfg)

	// Открываем файл для загрузки
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Получаем размер файла
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}
	size := fileInfo.Size()

	// Создаем контекст для операции
	ctx := context.Background()

	// Формируем параметры для загрузки
	input := &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		Body:          file,
		ContentLength: aws.Int64(int64(size)),
	}

	// Выполняем загрузку
	_, err = client.PutObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return nil
}
