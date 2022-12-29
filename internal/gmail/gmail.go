package gmail

import (
	"context"

	"github.com/robsztal/FanMyMail/cmd/config"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	v1 "google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type Client struct {
	svc *v1.Service
}

var tokenCh = make(chan *oauth2.Token, 1)

func NewClient(ctx context.Context, cfg config.Config) (Client, error) {
	token := getToken(ctx, cfg.OAuthCfg)
	client := cfg.OAuthCfg.Client(ctx, token)
	svc, err := v1.NewService(ctx, option.WithHTTPClient(client))
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
