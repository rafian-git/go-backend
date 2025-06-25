package pubsub

import (
	"github.com/rafian-git/go-backend/pkg/log"
	"github.com/segmentio/kafka-go"
)

type Topic struct {
	Name  string `json:"name" yaml:"name" toml:"name" mapstructure:"name"`
	Group string `json:"group" yaml:"group" toml:"group" mapstructure:"group"`
}

type Config struct {
	Address string `json:"address" yaml:"address" toml:"address" mapstructure:"address"`
	// for Consumer/Subscriber
	MaxWait      int `json:"max_wait" yaml:"max_wait" toml:"max_wait" mapstructure:"max_wait"`
	MaxBytes     int `json:"max_bytes" yaml:"max_bytes" toml:"max_bytes" mapstructure:"max_bytes"`
	MaxAttempts  int `json:"max_attempts" yaml:"max_attempts" toml:"max_attempts" mapstructure:"max_attempts"`
	BatchSize    int `json:"batch_size" yaml:"batch_size" toml:"batch_size" mapstructure:"batch_size"`
	BatchTimeout int `json:"batch_timeout" yaml:"batch_timeout" toml:"batch_timeout" mapstructure:"batch_timeout"`
	// for Producer/Publisher
	Balancer               string        `json:"balancer" yaml:"balancer" toml:"balancer" mapstructure:"balancer"`
	AllowAutoTopicCreation bool          `json:"allow_auto_topic_creation" yaml:"allow_auto_topic_creation" toml:"allow_auto_topic_creation" mapstructure:"allow_auto_topic_creation"`
	EmailPusherTopic       *Topic        `json:"email_pusher_topic" yaml:"email_pusher_topic" toml:"email_pusher_topic" mapstructure:"email_pusher_topic"`
	SMSPusherTopic         *Topic        `json:"sms_pusher_topic" yaml:"sms_pusher_topic" toml:"sms_pusher_topic" mapstructure:"sms_pusher_topic"`
	BankOtpVerifierTopic   *Topic        `json:"bank_otp_verifier_topic" yaml:"bank_otp_verifier_topic" toml:"bank_otp_verifier_topic" mapstructure:"bank_otp_verifier_topic"`
	Writer                 *kafka.Writer `json:"writer" yaml:"writer" toml:"writer" mapstructure:"writer"`
	Log                    *log.Logger
	// DseIndexTopic		  *Topic `json:"dse_index_topic" yaml:"dse_index_topic" toml:"dse_index_topic" mapstructure:"dse_index_topic"`
	// TickerStatusTopic		  *Topic `json:"ticker_status_topic" yaml:"ticker_status_topic" toml:"ticker_status_topic" mapstructure:"ticker_status_topic"`
	// MarketDepthTopic		  *Topic `json:"market_depth_topic" yaml:"market_depth_topic" toml:"market_depth_topic" mapstructure:"market_depth_topic"`
}

// NewConfig returns new default configurations.
func NewConfig() (conf *Config) {
	conf = new(Config)
	return
}
