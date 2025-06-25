package cloud

import (
	"archive/zip"
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"io"
	"os"
	"strings"
	"time"
)

type fileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi fileInfo) Name() string       { return fi.name }
func (fi fileInfo) Size() int64        { return fi.size }
func (fi fileInfo) Mode() os.FileMode  { return fi.mode }
func (fi fileInfo) ModTime() time.Time { return fi.modTime }
func (fi fileInfo) IsDir() bool        { return false }
func (fi fileInfo) Sys() interface{}   { return nil }

func (m *GCPBucket) GetFilesFromGCS(ctx context.Context, bucketName string, fileNames []string) ([]*storage.ObjectHandle, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)

	var objectHandles []*storage.ObjectHandle
	for _, fileName := range fileNames {
		object := bucket.Object(fileName)
		objectHandles = append(objectHandles, object)
	}

	return objectHandles, nil
}

func (g *GCPBucket) CreateZipFromGCS(ctx context.Context, bucketName string, fileNames []string, fileNameReplaceWith map[string]string) ([]byte, error) {
	objectHandles, err := g.GetFilesFromGCS(ctx, bucketName, fileNames)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for _, objectHandle := range objectHandles {
		// Open the file
		reader, err := objectHandle.NewReader(ctx)
		if err != nil {
			return nil, err
		}
		defer reader.Close()

		// Get file info
		attrs, err := objectHandle.Attrs(ctx)
		if err != nil {
			return nil, err
		}

		// Remove directory information from the file name
		parts := strings.Split(attrs.Name, "/")
		fileName := parts[len(parts)-1]

		newFileName, ok := fileNameReplaceWith[attrs.Name]
		if ok {
			fileName = newFileName
		}

		// Create a new file header for the file
		header, err := zip.FileInfoHeader(fileInfo{
			name:    fileName,
			size:    attrs.Size,
			mode:    0600,
			modTime: attrs.Updated,
		})

		if err != nil {
			return nil, err
		}

		header.Method = zip.Deflate

		// Write the file to the zip archive
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(writer, reader)
		if err != nil {
			return nil, err
		}
	}

	zipWriter.Close()
	return buf.Bytes(), nil
}

func (g *GCPBucket) UploadZipToGCS(ctx context.Context, bucketName string, zipBytes []byte, objectName string) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	wc := bucket.Object(objectName).NewWriter(ctx)
	wc.ContentType = "application/zip"

	if _, err := wc.Write(zipBytes); err != nil {
		return err
	}

	return wc.Close()
}
