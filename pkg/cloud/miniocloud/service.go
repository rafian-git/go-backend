package miniocloud

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/rafian-git/go-backend/pkg/apierror"
	"io"
	"strings"
	"time"
)

func (m *MinioBucket) GeneratePresignedUploadURL(ctx context.Context, path string) (string, *time.Time, error) {
	if len(path) == 0 {
		msg := "object can't be empty"
		m.log.Error(ctx, msg)
		return "", nil, apierror.New(apierror.InvalidArgument, msg)
	}

	exists, err := m.client.BucketExists(ctx, m.BucketConfig.BucketName)
	if err != nil {
		m.log.Error(ctx, err.Error())
		return "", nil, apierror.New(apierror.InvalidArgument, err.Error())
	}

	if !exists {
		msg := fmt.Sprintf("bucket %s does not exist", m.BucketConfig.BucketName)
		m.log.Error(ctx, msg)
		return "", nil, apierror.New(apierror.InvalidArgument, msg)
	}
	expire := time.Duration(m.BucketConfig.SignedUrlExpDuration) * time.Second
	presignedURL, err := m.client.PresignedPutObject(ctx, m.BucketConfig.BucketName, m.BucketConfig.ObjectPathPrefix+path, expire)
	if err != nil {
		m.log.Error(ctx, err.Error())
		return "", nil, apierror.New(apierror.Internal, err.Error())
	}
	expireat := time.Now().Add(expire)
	return presignedURL.String(), &expireat, nil

}

func (m *MinioBucket) ParseXmlFile(ctx context.Context, path string, data interface{}) error {
	obj, err := m.client.GetObject(ctx, m.BucketConfig.BucketName, path, minio.GetObjectOptions{})
	if err != nil {
		m.log.Error(ctx, err.Error())
		return err
	}

	decoder := xml.NewDecoder(obj)
	err = decoder.Decode(data)
	if err != nil {
		m.log.Error(ctx, err.Error())
		return err
	}
	return nil
}

func (m *MinioBucket) GetRawData(ctx context.Context, path string) (*[]byte, error) {
	t1 := time.Now()
	defer func() {
		t2 := time.Now()
		m.log.Info(ctx, fmt.Sprintf("total time took for reading data: %v seconds", t2.Sub(t1).Seconds()))
	}()

	exists, err := m.client.BucketExists(ctx, m.BucketConfig.BucketName)
	if err != nil {
		m.log.Error(ctx, err.Error())
		return nil, apierror.New(apierror.InvalidArgument, err.Error())
	}

	if !exists {
		msg := fmt.Sprintf("bucket %s does not exist", m.BucketConfig.BucketName)
		m.log.Error(ctx, msg)
		return nil, apierror.New(apierror.InvalidArgument, msg)
	}

	objectReader, err := m.client.GetObject(ctx, m.BucketConfig.BucketName, m.BucketConfig.ObjectPathPrefix+path, minio.GetObjectOptions{})
	if err != nil {
		m.log.Error(ctx, err.Error())
		return nil, err
	}

	defer objectReader.Close()
	data, err := io.ReadAll(objectReader)
	if err != nil {
		m.log.Error(ctx, err.Error())
		return nil, err
	}

	return &data, nil
}

func (m *MinioBucket) StreamFileUpload(ctx context.Context, path, fileBase64 string) error {
	t1 := time.Now()
	defer func() {
		t2 := time.Now()
		msg := fmt.Sprintf("total time took for File Upload: %v seconds\n", t2.Sub(t1).Seconds())
		m.log.Error(ctx, msg)
	}()

	// removing data URI part: ex: data:image/jpeg;base64,/
	b64data := fileBase64[strings.IndexByte(fileBase64, ',')+1:]
	decodedFile, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		m.log.Error(ctx, err.Error())
		return fmt.Errorf("DecodeString: %v", err)
	}

	if len(decodedFile) == 0 {
		msg := "empty file detected"
		m.log.Error(ctx, msg)
		return fmt.Errorf(msg)
	}
	if len(decodedFile) > m.BucketConfig.MaxFileSizeObject {
		return errors.New("file size must be equal or less than 1 mb")
	}

	decodedReader := bytes.NewReader(decodedFile)

	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	// Upload the object
	_, err = m.client.PutObject(ctx, m.BucketConfig.BucketName, m.BucketConfig.ObjectPathPrefix+path, decodedReader, int64(decodedReader.Len()), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		m.log.Error(ctx, err.Error())
		return fmt.Errorf("PutObject: %v", err)
	}

	return nil
}
