package centrifugo

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.techetronventures.com/core/backend/pkg/apierror"
	"go.uber.org/zap"
)

func (socket *Client) SendSocketData(ctx context.Context, channelId string, data interface{}) error {
	socket.log.Info(ctx, "", zap.Any("config", socket.config))
	chs, err := socket.client.Channels(ctx)

	if err != nil {
		msg := "failed to fetch channels"
		socket.log.Error(ctx, fmt.Sprintf("%s : %v", msg, err.Error()))
		return apierror.New(apierror.DataLoss, msg)
	}

	_, found := chs.Channels[channelId]
	if !found {
		socket.log.Warn(ctx, fmt.Sprintf("channel Id(%s) not found", channelId))
		return nil
	}

	json, err := json.Marshal(data)

	result, err := socket.client.Publish(ctx, channelId, []byte(json))

	if err != nil {
		socket.log.Error(ctx, fmt.Sprintf("Error calling publish: %v", err))
		return err
	}
	socket.log.Info(ctx, fmt.Sprintf("Publish into channel %s successful, stream position {offset: %d, epoch: %s}", channelId, result.Offset, result.Epoch))
	return nil
}
