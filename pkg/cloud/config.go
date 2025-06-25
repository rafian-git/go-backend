package cloud

import (
	"time"
)

type Config struct {
	GCPBucketConfig *GCPBucketConfig `json:"gcp_bucket_config" yaml:"gcp_bucket_config" toml:"gcp_bucket_config" mapstructure:"gcp_bucket_config"`
}

type GCPBucketConfig struct {
	BucketName           string        `json:"bucket_name" yaml:"bucket_name" toml:"bucket_name" mapstructure:"bucket_name"`
	ObjectPathPrefix     string        `json:"object_path_prefix" yaml:"object_path_prefix" toml:"object_path_prefix" mapstructure:"object_path_prefix"`
	MaxFileSizeObject    int           `json:"max_file_size_object" yaml:"max_file_size_object" toml:"max_file_size_object" mapstructure:"max_file_size_object"`
	SignedUrlExpDuration time.Duration `json:"signed_url_exp_duration" yaml:"signed_url_exp_duration" toml:"signed_url_exp_duration" mapstructure:"signed_url_exp_duration"`
}

// NewConfig returns new default configurations.
func NewConfig() (conf *Config) {
	conf = new(Config)
	return
}
