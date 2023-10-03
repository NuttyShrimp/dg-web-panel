package discord

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	dgerrors "degrens/panel/lib/errors"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func fetchFromDiscordAPI(token *oauth2.Token, path string, target interface{}) error {
	res, err := info.Conf.Client(context.Background(), token).Get("https://discord.com/api/" + path)

	if err != nil || res.StatusCode != 200 {
		if err == nil {
			err = fmt.Errorf("Discord API responded with status code %d: %s", res.StatusCode, res.Status)
		}
		logrus.WithField("statusCode", res.StatusCode).WithError(err).Error("Error while getting user info")
		return errors.New("error while getting user info")
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	dec := json.NewDecoder(res.Body)

	if err := dec.Decode(&target); err != nil {
		dgerrors.HandleJsonError(err, logrus.StandardLogger())
		return errors.New("error while getting member info")
	}
	return nil
}
