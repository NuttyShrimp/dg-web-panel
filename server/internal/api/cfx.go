package api

import (
	"bytes"
	"degrens/panel/internal/config"
	"degrens/panel/models"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

type cfxApi struct {
	api
	apiKey string
}

var CfxApi cfxApi

func CreateCfxApi(c *config.ConfigCfx) {
	regex := regexp.MustCompile(`/$`)
	baseApi := api{
		baseURL: regex.ReplaceAllString(c.Server, "") + "/dg-api",
		Logger:  logrus.WithField("api", "Cfx"),
	}
	CfxApi = cfxApi{
		api:    baseApi,
		apiKey: c.ApiKey,
	}
}

func (ca *cfxApi) DoRequest(method, endpoint string, input, output interface{}) (*ErrorInfo, error) {
	ei, err := ca.doInternalRequest(method, endpoint, input, output, ca.addInput, ca.doAuthentication)
	if err != nil && strings.HasSuffix(err.Error(), "connection refused") {
		return ei, &models.RouteError{
			Message: models.RouteErrorMessage{
				Title:       "The fivem server is currently down!",
				Description: "We try make as much services available for you but this is not always possible",
			},
			Code: http.StatusUnauthorized,
		}
	}
	return ei, err
}

func (ca *cfxApi) addInput(req *http.Request, input interface{}) {
	body, err := json.Marshal(input)
	if err != nil {
		ca.Logger.WithError(err).Error("Failed to encode input to JSON for Cfx request")
		return
	}
	buf := bytes.NewBuffer(body)
	req.ContentLength = int64(buf.Len())
	req.Body = io.NopCloser(buf)
}

func (ca *cfxApi) doAuthentication(req *http.Request) {
	req.Header.Add("Authorization", "Bearer "+ca.apiKey)
}

func (ca *cfxApi) Post(endpoint string, input, output interface{}) error {
	ai, err := ca.DoRequest("POST", endpoint, input, output)
	if err != nil {
		return err
	}
	if ai.Message != "" {
		return errors.New(ai.Message)
	}
	return nil
}
