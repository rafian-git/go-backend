package centrifugo

import (
	"context"
	"github.com/centrifugal/gocent/v3"
	"github.com/rafian-git/go-backend/pkg/log"
)

type Socket interface {
	SendSocketData(ctx context.Context, channelId string, data interface{}) error
}

type Client struct {
	config *Config
	log    *log.Logger
	client *gocent.Client
}

func New(log *log.Logger, config *Config) (Socket, error) {
	cenCli := gocent.New(gocent.Config{
		Addr: config.Url,
		Key:  config.ApiKey,
	})

	client := &Client{
		log:    log.Named("centrifugo"),
		config: config,
		client: cenCli,
	}

	return client, nil
}
