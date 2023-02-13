package discord

import (
	"context"
	"encoding/json"
	"errors"

	"golang.org/x/oauth2"
)

func fetchFromDiscordAPI(token *oauth2.Token, path string) (*json.Decoder, error) {
	res, err := info.Conf.Client(context.Background(), token).Get("https://discord.com/api/" + path)

	if err != nil || res.StatusCode != 200 {
		logger.Error("Error while getting user info", "error", err, "statusCode", res.StatusCode)
		return nil, errors.New("error while getting user info")
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	dec := json.NewDecoder(res.Body)
	return dec, nil
}
