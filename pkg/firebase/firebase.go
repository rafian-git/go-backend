package firebase

import (
	"context"
	"encoding/json"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"go.uber.org/zap"
	"time"
)

func (service *FirebaseService) Send(ctx context.Context, title, body, deviceSignature string, data map[string]string) (string, error) {
	var startTime int64 = time.Now().Unix()

	data["title"] = title
	data["body"] = body

	message := &messaging.Message{
		Data:    data,
		Token:   deviceSignature,
		Android: &messaging.AndroidConfig{},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					MutableContent: true,
					Alert: &messaging.ApsAlert{
						Title: title,
						Body:  body,
					},
					Sound: "default",
				},
			},
		},
	}

	response, err := service.client.Send(ctx, message)
	if err != nil {
		service.log.Error(ctx, err.Error())
		return "", err
	}

	service.log.Info(ctx, "", zap.Any("response", response))

	defer func() {
		var endTime int64 = time.Now().Unix()
		service.log.Info(ctx, fmt.Sprintf("took time: %d", (endTime-startTime)))
	}()

	return response, nil
}

func (service *FirebaseService) SendMultiMessage(ctx context.Context, title, body string, tokens []string, data map[string]string) (string, error) {
	var startTime int64 = time.Now().Unix()
	data["title"] = title
	data["body"] = body
	message := &messaging.MulticastMessage{
		Data:    data,
		Tokens:  tokens,
		Android: &messaging.AndroidConfig{},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					MutableContent: true,
					Alert: &messaging.ApsAlert{
						Title: title,
						Body:  body,
					},
					Sound: "default",
				},
			},
		},
	}
	response, err := service.client.SendEachForMulticast(ctx, message)
	if err != nil {
		service.log.Error(ctx, err.Error())
		return "", err
	}
	service.log.Info(ctx, fmt.Sprintf("response SuccessCount: %d  FailureCount: %d", response.SuccessCount, response.FailureCount))
	resp := &BatchResponse{
		SuccessCount: int32(response.SuccessCount),
		FailureCount: int32(response.FailureCount),
	}
	for _, item := range response.Responses {
		obj := &SendResponse{
			Success:   item.Success,
			MessageID: item.MessageID,
		}
		if item.Error != nil {
			obj.Error = item.Error.Error()
		}
		resp.Responses = append(resp.Responses, obj)
	}
	bytes, err := json.Marshal(resp)
	if err != nil {
		service.log.Error(ctx, err.Error())
		return "", err
	}

	defer func() {
		var endTime int64 = time.Now().Unix()
		service.log.Info(ctx, fmt.Sprintf("took time: %d", endTime-startTime))
	}()

	return string(bytes), nil
}

type BatchResponse struct {
	SuccessCount int32           `json:"success_count"`
	FailureCount int32           `json:"failure_count"`
	Responses    []*SendResponse `json:"responses"`
}

type SendResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"message_id,omitempty"`
	Error     string `json:"error,omitempty"`
}
