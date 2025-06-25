package cloud

import (
	"context"
	"gitlab.techetronventures.com/core/backend/pkg/log"
)

type Cloud interface {
	StreamFileUpload(ctx context.Context, bucket, object string, imgBase64 string) error
	GenerateV4GetObjectSignedURL(bucket, object string) (string, error)
	GenerateV4GetObjectSignedURLs(bucket string, objects []string) ([]string, error)
	GetRawData(bucket, object string) (*[]byte, error)
	WriteFile(bucket, object, fileContent string) error
	CreateZipFromGCS(ctx context.Context, bucketName string, fileNames []string, fileNameReplaceWith map[string]string) ([]byte, error)
	UploadZipToGCS(ctx context.Context, bucketName string, zipBytes []byte, objectName string) error
	GenerateSignedUrlForFileUpload(object string, bucket, contentType string) (string, error)
	FindFilesByPrefix(bucket, prefix string) ([]string, error)
	CheckIfFileExist(bucket, objectName string) (bool, error)
	GetPublicURL(bucketName, fileName string) string
	MakeFilePublic(ctx context.Context, bucketName, fileName string) error
	UploadLocalFile(ctx context.Context, bucketName string, objectPath string, localFilepath string) error
	GenerateV4GetObjectDownloadURLs(bucket, folder string) ([]string, error)
	DownloadToFile(bucket, object, localPath string) error
}

type GCPBucket struct {
	GCPBucketConfig *GCPBucketConfig
	log             *log.Logger
}

func New(cnf *Config) (Cloud, error) {
	logger := log.New().Named("backend_cloud")
	return &GCPBucket{
		GCPBucketConfig: cnf.GCPBucketConfig,
		log:             logger,
	}, nil
}
