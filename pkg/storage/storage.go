package storage

type Storage interface {
	UploadFile(bucket, key, filepath string) error
}
