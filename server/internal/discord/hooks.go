package discord

import (
	"degrens/panel/internal/config"
	"fmt"

	disgo "github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
)

var clients = make(map[string]webhook.Client)

func getClient(key string) (webhook.Client, error) {
	client, ok := clients[key]
	if !ok {
		url, ok := config.GetConfig().Discord.Hooks[key]
		if !ok {
			return nil, fmt.Errorf("no webhook found for %s", key)
		}
		client, err := webhook.NewWithURL(url)
		if err != nil {
			return nil, err
		}
		clients[key] = client
		return client, nil
	}
	return client, nil
}

func SendToReportWebHook(msg *disgo.WebhookMessageCreate) error {
	client, err := getClient("reportlog")
	if err != nil {
		return err
	}
	_, err = client.CreateMessage(*msg)
	return err
}
