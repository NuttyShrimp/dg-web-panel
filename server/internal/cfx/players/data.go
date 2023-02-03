package players

import (
	"degrens/panel/internal/api"
	"errors"
	"fmt"
)

func GetActiveCharacter(steamId string) (uint, error) {
	output := struct {
		Cid uint `json:"cid"`
	}{}
	ai, err := api.CfxApi.DoRequest("GET", fmt.Sprintf("/info/%s/active", steamId), nil, &output)
	if err != nil {
		return 0, err
	}
	if ai.Message != "" {
		return 0, errors.New(ai.Message)
	}
	return output.Cid, nil
}
