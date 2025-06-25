package firebase

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"gitlab.techetronventures.com/core/backend/pkg/log"
	"google.golang.org/api/option"
)

type FCM interface {
	Send(ctx context.Context, title, body, deviceSignature string, data map[string]string) (string, error)
	SendMultiMessage(ctx context.Context, title, body string, tokens []string, data map[string]string) (string, error)
}

type FirebaseService struct {
	log    *log.Logger
	client *messaging.Client
}

func New(ctx context.Context, log *log.Logger, credentialPath string) (FCM, error) {
	service := &FirebaseService{
		log: log.Named("firebase"),
	}

	opts := []option.ClientOption{option.WithCredentialsFile(credentialPath)}
	app, err := firebase.NewApp(ctx, nil, opts...)

	if err != nil {
		service.log.Error(ctx, err.Error())
		return nil, err
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		service.log.Error(ctx, err.Error())
		return nil, err
	}

	service.client = client

	return service, nil
}
