package oss

import "context"

// Client defines minimal interface for interacting with Aliyun OSS.
type Client interface {
	GenerateUploadCredentials(ctx context.Context, objectKey string) (*UploadCredentials, error)
}

// UploadCredentials represent temporary upload data for clients.
type UploadCredentials struct {
	Endpoint        string `json:"endpoint"`
	Bucket          string `json:"bucket"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	SecurityToken   string `json:"security_token"`
	ExpireAt        int64  `json:"expire_at"`
}

// StaticClient returns preconfigured credentials (placeholder for integration).
type StaticClient struct {
	Endpoint string
	Bucket   string
}

// GenerateUploadCredentials returns static credentials. Replace with real STS integration.
func (c *StaticClient) GenerateUploadCredentials(ctx context.Context, objectKey string) (*UploadCredentials, error) {
	return &UploadCredentials{
		Endpoint: c.Endpoint,
		Bucket:   c.Bucket,
	}, nil
}
