package miniocloud

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gitlab.techetronventures.com/core/backend/pkg/log"
	"time"
)

type Cloud interface {
	GeneratePresignedUploadURL(ctx context.Context, objectName string) (string, *time.Time, error)
	ParseXmlFile(ctx context.Context, path string, data interface{}) error
	GetRawData(ctx context.Context, path string) (*[]byte, error)
	StreamFileUpload(ctx context.Context, path, fileBase64 string) error
}

type MinioBucket struct {
	BucketConfig *BucketConfig

	log    *log.Logger
	client *minio.Client
}

func New(cnf *Config) (Cloud, error) {
	logger := log.New().Named("backend_minio_cloud")

	client, err := minio.New(cnf.Credential.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cnf.Credential.MinioAccessKey, cnf.Credential.MinioSecretKey, ""),
		Secure: cnf.Credential.MinioUseSSL,
	})
	if err != nil {
		return nil, err
	}
	return &MinioBucket{
		BucketConfig: cnf.BucketConfig,
		client:       client,
		log:          logger,
	}, nil
}
