package gmail

import (
	"context"

	"github.com/rs/zerolog/log"
	v1 "google.golang.org/api/gmail/v1"
)

type Client struct {
	svc *v1.Service
}

func NewClient(ctx context.Context) (Client, error) {
	svc, err := v1.NewService(ctx)
	if err != nil {
		return Client{}, err
	}
	return Client{svc: svc}, nil
}

func (c Client) Fetch(ctx context.Context, userID string) ([]*v1.Message, error) {
	msgList, err := c.svc.Users.Messages.List(userID).Do()
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("user_id", userID).Msg("Cannot fetch list of messages")
		return nil, err
	}

	msgContent := make([]*v1.Message, 0, len(msgList.Messages))
	for _, msg := range msgList.Messages {
		content, err := c.GetMessageContent(userID, msg.Id)
		if err != nil {
			log.Ctx(ctx).Warn().Str("message_id", msg.Id).Msg("Failed to get content for message")
			continue
		}
		msgContent = append(msgContent, content)
	}
	log.Ctx(ctx).Info().Int("messages", len(msgContent)).Msg("Fetched messages.")
	return msgContent, nil
}

func (c Client) GetMessageContent(userID, messageID string) (*v1.Message, error) {
	return c.svc.Users.Messages.Get(userID, messageID).Do()
}
