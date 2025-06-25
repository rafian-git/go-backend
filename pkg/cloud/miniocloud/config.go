package miniocloud

type Config struct {
	BucketConfig *BucketConfig `json:"bucket_config" yaml:"bucket_config" toml:"bucket_config" mapstructure:"bucket_config"`
	Credential   *Credential   `json:"credential" yaml:"credential" toml:"credential" mapstructure:"credential"`
}

type Credential struct {
	MinioEndpoint  string `json:"minio_endpoint" yaml:"minio_endpoint" toml:"minio_endpoint" mapstructure:"minio_endpoint"`
	MinioAccessKey string `json:"minio_access_key" yaml:"minio_access_key" toml:"minio_access_key" mapstructure:"minio_access_key"`
	MinioSecretKey string `json:"minio_secret_key" yaml:"minio_secret_key" toml:"minio_secret_key" mapstructure:"minio_secret_key"`
	MinioUseSSL    bool   `json:"minio_use_ssl" yaml:"minio_use_ssl" toml:"minio_use_ssl" mapstructure:"minio_use_ssl"`
}

type BucketConfig struct {
	BucketName           string `json:"bucket_name" yaml:"bucket_name" toml:"bucket_name" mapstructure:"bucket_name"`
	ObjectPathPrefix     string `json:"object_path_prefix" yaml:"object_path_prefix" toml:"object_path_prefix" mapstructure:"object_path_prefix"`
	MaxFileSizeObject    int    `json:"max_file_size_object" yaml:"max_file_size_object" toml:"max_file_size_object" mapstructure:"max_file_size_object"`
	SignedUrlExpDuration int    `json:"signed_url_exp_duration" yaml:"signed_url_exp_duration" toml:"signed_url_exp_duration" mapstructure:"signed_url_exp_duration"`
}

// NewConfig returns new default configurations.
func NewConfig() (conf *Config) {
	conf = new(Config)
	return
}

type MinioEvent struct {
	EventName string   `json:"EventName"`
	Key       string   `json:"Key"`
	Records   []Record `json:"Records"`
}

type Record struct {
	EventVersion string    `json:"eventVersion"`
	EventSource  string    `json:"eventSource"`
	EventTime    string    `json:"eventTime"`
	EventName    string    `json:"eventName"`
	S3           S3Details `json:"s3"`
}

type S3Details struct {
	Bucket BucketDetails `json:"bucket"`
	Object ObjectDetails `json:"object"`
}

type BucketDetails struct {
	Name string `json:"name"`
}

type ObjectDetails struct {
	Key          string            `json:"key"`
	ETag         string            `json:"eTag"`
	ContentType  string            `json:"contentType"`
	UserMetadata map[string]string `json:"userMetadata"`
}
