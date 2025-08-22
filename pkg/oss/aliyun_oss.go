package oss

import (
	"mime/multipart"

	"hinoob.net/learn-go/internal/config"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var (
	ossClient *oss.Client
)

// InitOSS initializes the OSS client.
func InitOSS() {
	var err error
	ossClient, err = oss.New(
		config.AppConfig.OSS.Endpoint,
		config.AppConfig.OSS.AccessKeyID,
		config.AppConfig.OSS.AccessKeySecret,
	)
	if err != nil {
		panic("Failed to initialize OSS client: " + err.Error())
	}
}

// UploadFile uploads a file to the configured OSS bucket and returns its URL.
func UploadFile(fileHeader *multipart.FileHeader) (string, error) {
	// 1. Get the bucket
	bucket, err := ossClient.Bucket(config.AppConfig.OSS.BucketName)
	if err != nil {
		return "", err
	}

	// 2. Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 3. Upload the file to the bucket
	// We use the original filename. In a real-world scenario, you'd want to
	// generate a unique name (e.g., using a UUID) to prevent overwrites.
	objectKey := "uploads/" + fileHeader.Filename
	err = bucket.PutObject(objectKey, file)
	if err != nil {
		return "", err
	}

	// 4. The public URL of the object
	// Note: This assumes your bucket has public read access.
	// If not, you would generate a signed URL instead.
	fileURL := "https://" + config.AppConfig.OSS.BucketName + "." + config.AppConfig.OSS.Endpoint + "/" + objectKey
	return fileURL, nil
}
