package minio

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"go-storage/internal/domain"
	"go-storage/pkg/errors"
)

type StorageRepository struct {
	client     *minio.Client
	bucketName string
}

func NewStorageRepository(client *minio.Client, bucketName string) *StorageRepository {
	return &StorageRepository{
		client:     client,
		bucketName: bucketName,
	}
}

func (r *StorageRepository) StoreFile(ctx context.Context, key string, reader io.Reader, size int64, mimeType string) (string, error) {
	opts := minio.PutObjectOptions{
		ContentType: mimeType,
	}

	info, err := r.client.PutObject(ctx, r.bucketName, key, reader, size, opts)
	if err != nil {
		return "", errors.InternalServer("failed to store file in storage")
	}

	return info.ETag, nil
}

func (r *StorageRepository) GetFile(ctx context.Context, key string) (io.ReadCloser, error) {
	object, err := r.client.GetObject(ctx, r.bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.NotFound("file not found in storage")
	}

	return object, nil
}

func (r *StorageRepository) DeleteFile(ctx context.Context, key string) error {
	err := r.client.RemoveObject(ctx, r.bucketName, key, minio.RemoveObjectOptions{})
	if err != nil {
		return errors.InternalServer("failed to delete file from storage")
	}

	return nil
}

func (r *StorageRepository) GetFileInfo(ctx context.Context, key string) (*domain.StorageFileInfo, error) {
	info, err := r.client.StatObject(ctx, r.bucketName, key, minio.StatObjectOptions{})
	if err != nil {
		return nil, errors.NotFound("file not found in storage")
	}

	return &domain.StorageFileInfo{
		Key:          key,
		Size:         info.Size,
		MimeType:     info.ContentType,
		ETag:         info.ETag,
		LastModified: info.LastModified,
	}, nil
}

func (r *StorageRepository) InitChunkedUpload(ctx context.Context, key string, mimeType string) (string, error) {
	return uuid.NewString(), nil
}

func (r *StorageRepository) UploadChunk(ctx context.Context, uploadID, key string, chunkIndex int, reader io.Reader, size int64) (string, error) {
	chunkKey := fmt.Sprintf("%s.chunk.%d", key, chunkIndex)
	_, err := r.client.PutObject(ctx, r.bucketName, chunkKey, reader, size, minio.PutObjectOptions{})
	if err != nil {
		return "", errors.InternalServer("failed to upload chunk")
	}

	return fmt.Sprintf("chunk-%d", chunkIndex), nil
}

func (r *StorageRepository) CompleteChunkedUpload(ctx context.Context, uploadID, key string, parts []string) error {
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()
		for i := range parts {
			chunkKey := fmt.Sprintf("%s.chunk.%d", key, i)
			obj, err := r.client.GetObject(ctx, r.bucketName, chunkKey, minio.GetObjectOptions{})
			if err != nil {
				pw.CloseWithError(err)
				return
			}
			_, err = io.Copy(pw, obj)
			obj.Close()
			if err != nil {
				pw.CloseWithError(err)
				return
			}
		}
	}()

	_, err := r.client.PutObject(ctx, r.bucketName, key, pr, -1, minio.PutObjectOptions{})
	if err != nil {
		return errors.InternalServer("failed to combine chunks")
	}

	for i := range parts {
		chunkKey := fmt.Sprintf("%s.chunk.%d", key, i)
		r.client.RemoveObject(ctx, r.bucketName, chunkKey, minio.RemoveObjectOptions{})
	}

	return nil
}

func (r *StorageRepository) AbortChunkedUpload(ctx context.Context, uploadID, key string) error {
	obj := r.client.ListObjects(ctx, r.bucketName, minio.ListObjectsOptions{
		Prefix: key + ".chunk.",
	})

	for object := range obj {
		if object.Err != nil {
			continue
		}
		r.client.RemoveObject(ctx, r.bucketName, object.Key, minio.RemoveObjectOptions{})
	}

	return nil
}

func (r *StorageRepository) EnsureBucket(ctx context.Context) error {
	exists, err := r.client.BucketExists(ctx, r.bucketName)
	if err != nil {
		return errors.InternalServer("failed to check bucket existence")
	}

	if !exists {
		err = r.client.MakeBucket(ctx, r.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return errors.InternalServer("failed to create bucket")
		}
	}

	return nil
}

func (r *StorageRepository) ListFiles(ctx context.Context, prefix string) (<-chan minio.ObjectInfo, error) {
	objectCh := r.client.ListObjects(ctx, r.bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	return objectCh, nil
}

func (r *StorageRepository) GetPresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	url, err := r.client.PresignedGetObject(ctx, r.bucketName, key, expiry, nil)
	if err != nil {
		return "", errors.InternalServer("failed to generate presigned URL")
	}

	return url.String(), nil
}
