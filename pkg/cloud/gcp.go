package cloud

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"
)

// StreamFileUpload uploads an object to GCP via a stream.
func (m *GCPBucket) StreamFileUpload(ctx context.Context, bucket, object string, fileBase64 string) error {
	t1 := time.Now()
	defer func() {
		t2 := time.Now()
		log.Printf("total time took for File Upload: %v\n", t2.Sub(t1).Seconds())
	}()

	client, err := storage.NewClient(ctx)
	if err != nil {
		m.log.Error(ctx, err.Error())
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	//removing data uri: ex: data:image/jpeg;base64,/
	b64data := fileBase64[strings.IndexByte(fileBase64, ',')+1:]
	decodedFile, err := base64.StdEncoding.DecodeString(b64data)

	if err != nil {
		m.log.Error(ctx, err.Error())
		return fmt.Errorf("DecodeString: %v", err)
	}

	if len(decodedFile) == 0 {
		return errors.New("empty file detected")
	}
	if len(decodedFile) > m.GCPBucketConfig.MaxFileSizeObject {
		return errors.New("file size must be equal or less than 1 mb")
	}

	decodedReader := bytes.NewReader(decodedFile)

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	wc.ChunkSize = 0 // note retries are not supported for chunk size 0.
	wc.CacheControl = "no-cache, max-age=0"

	if _, err = io.Copy(wc, decodedReader); err != nil {
		m.log.Error(ctx, err.Error())
		return fmt.Errorf("io.Copy: %v", err)
	}
	// Data can continue to be added to the file until the writer is closed.
	if err := wc.Close(); err != nil {
		m.log.Error(ctx, err.Error())
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}

// GenerateV4GetObjectSignedURL generates object signed URL with GET method.
func (m *GCPBucket) GenerateV4GetObjectSignedURL(bucket, object string) (string, error) {
	if len(object) == 0 {
		return "", errors.New("object can't be empty")
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	SignedURLExpDuration := m.GCPBucketConfig.SignedUrlExpDuration
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(SignedURLExpDuration * time.Minute),
	}

	// Get object metadata to retrieve file size
	//attrs, err := client.Bucket(bucket).Object(object).Attrs(ctx)
	//if err != nil {
	//	return "", fmt.Errorf("Object(%q).Attrs: %v", object, err)
	//}
	//log.Println("getting image url for :", object)
	u, err := client.Bucket(bucket).SignedURL(object, opts)

	if err != nil {
		return "", fmt.Errorf("Bucket(%q).SignedURL: %v", bucket, err)
	}

	//log.Println("url is: ", u)

	return u, nil
}

// GenerateV4GetObjectSignedURLs generates signed URLs from object paths with GET method.
func (m *GCPBucket) GenerateV4GetObjectSignedURLs(bucket string, objects []string) ([]string, error) {
	if len(objects) == 0 {
		return nil, errors.New("no object paths provided")
	}
	var urls []string
	for _, obj := range objects {
		url, err := m.GenerateV4GetObjectSignedURL(bucket, obj)
		if err != nil {
			return nil, fmt.Errorf("GenerateV4GetObjectSignedURL error: %v", err)
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func (m *GCPBucket) GetRawData(bucket, object string) (*[]byte, error) {

	t1 := time.Now()
	defer func() {
		t2 := time.Now()
		log.Printf("total time took for writting data: %v\n", t2.Sub(t1).Seconds())
	}()

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	defer client.Close()

	reader, err := client.Bucket(bucket).Object(object).NewReader(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create reader: %v", err)
	}

	reader.Close()

	data, err := io.ReadAll(reader)

	if err != nil {
		return nil, fmt.Errorf("failed to read data from reader: %v", err)
	}

	return &data, nil
}

// WriteFile writes a file into gcp bucket
func (m *GCPBucket) WriteFile(bucket, object, fileContent string) error {

	t1 := time.Now()
	defer func() {
		t2 := time.Now()
		log.Printf("total time took for writting data: %v\n", t2.Sub(t1).Seconds())
	}()

	// creating client
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}

	defer client.Close()

	fileContentByte := []byte(fileContent)
	fileReader := bytes.NewReader(fileContentByte)

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	wc.ChunkSize = 0 // note retries are not supported for chunk size 0.
	//wc.CacheControl = "no-cache, max-age=0"

	if _, err = io.Copy(wc, fileReader); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	// Data can continue to be added to the file until the writer is closed.
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}

// GenerateSignedUrlForFileUpload generates signed URL for file upload
func (m *GCPBucket) GenerateSignedUrlForFileUpload(object string, bucket, contentType string) (string, error) {
	if len(object) == 0 {
		return "", errors.New("object can't be empty")
	}
	if len(bucket) == 0 {
		return "", errors.New("bucket can't be empty")
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	contentLength := "1"

	SignedURLExpDuration := m.GCPBucketConfig.SignedUrlExpDuration
	r, err := client.Bucket(bucket).SignedURL(object, &storage.SignedURLOptions{
		Method:  "PUT",
		Expires: time.Now().Add(SignedURLExpDuration * time.Minute),
		Headers: []string{
			"Content-Length: " + contentLength,
		},
		ContentType: contentType,
	})
	return r, err
}

// FindFilesByPrefix finds files by prefix
func (m *GCPBucket) FindFilesByPrefix(bucket, prefix string) ([]string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	var files []string

	it := client.Bucket(bucket).Objects(ctx, &storage.Query{Prefix: prefix})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Bucket(%q).Objects(): %v", bucket, err)
		}
		files = append(files, attrs.Name)
	}
	return files, nil
}

// checkIfFileExist checks if file exist
func (m *GCPBucket) CheckIfFileExist(bucket, objectName string) (bool, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return false, fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	_, err = client.Bucket(bucket).Object(objectName).Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return false, nil
		} else if gErr, ok := err.(*googleapi.Error); ok && gErr.Code == 404 {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

// GetPublicURL gets a public URL for a file
func (m *GCPBucket) GetPublicURL(bucketName, fileName string) string {
	baseURL := "https://storage.googleapis.com"
	escapedBucketName := url.PathEscape(bucketName)
	escapedFileName := url.PathEscape(fileName)
	publicURL := fmt.Sprintf("%s/%s/%s", baseURL, escapedBucketName, escapedFileName)
	return publicURL
}

// MakeFilePublic makes a file public
func (m *GCPBucket) MakeFilePublic(ctx context.Context, bucketName, fileName string) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	object := bucket.Object(fileName)

	acl := object.ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return err
	}

	return nil
}

// UploadLocalFile uploads a local file to gcp bucket
func (m *GCPBucket) UploadLocalFile(ctx context.Context, bucketName string, objectPath string, localFilepath string) error {
	// creating client
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()
	// open local file
	f, err := os.Open(localFilepath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()
	// Create a new object handle
	bucket := client.Bucket(bucketName)
	obj := bucket.Object(objectPath)
	writer := obj.NewWriter(ctx)
	// Copy the contents of the local file to the writer
	if _, err := io.Copy(writer, f); err != nil {
		fmt.Println(err)
		return err
	}

	// Close the writer
	if err := writer.Close(); err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Local File uploaded to GCP successfully.")
	return nil
}

// GenerateV4GetObjectDownloadURLs generates signed URLs for objects in a folder for download with GET method and content-disposition parameter.
func (m *GCPBucket) GenerateV4GetObjectDownloadURLs(bucket, folder string) ([]string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	var downloadURLs []string

	// List objects in the specified folder.
	it := client.Bucket(bucket).Objects(ctx, &storage.Query{Prefix: folder})
	for {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("iterator.Next: %v", err)
		}

		filename := objAttrs.Name[len(folder):] // Extract the filename without the folder path

		SignedURLExpDuration := m.GCPBucketConfig.SignedUrlExpDuration
		opts := &storage.SignedURLOptions{
			Scheme:  storage.SigningSchemeV4,
			Method:  "GET",
			Expires: time.Now().Add(SignedURLExpDuration * time.Minute),
			// Add the content-disposition parameter to the query.
			// It enables downloadable signed url
			QueryParameters: map[string][]string{
				"response-content-disposition": {"attachment;filename=" + filename},
			},
		}

		u, err := client.Bucket(bucket).SignedURL(objAttrs.Name, opts)
		if err != nil {
			return nil, fmt.Errorf("Object(%q).SignedURL: %v", objAttrs.Name, err)
		}

		downloadURLs = append(downloadURLs, u)
	}

	return downloadURLs, nil
}

func (m *GCPBucket) DownloadToFile(bucket, object, localPath string) error {
	t1 := time.Now()
	defer func() {
		t2 := time.Now()
		log.Printf("total time took for downloading file: %.2f seconds\n", t2.Sub(t1).Seconds())
	}()

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create GCS client: %v", err)
	}
	defer client.Close()

	reader, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("failed to create GCS reader: %v", err)
	}
	defer reader.Close()

	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %v", err)
	}
	defer file.Close()

	// Stream from GCS to file
	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("failed to copy GCS object to file: %v", err)
	}

	log.Println("✅ File downloaded to", localPath)
	return nil
}
