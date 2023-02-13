package api

import (
	"bytes"
	"degrens/panel/internal/config"
	"degrens/panel/lib/log"
	"degrens/panel/models"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type cfxApi struct {
	api
	apiKey string
}

var CfxApi cfxApi

func CreateCfxApi(c *config.ConfigCfx, logger log.Logger) {
	regex := regexp.MustCompile(`/$`)
	baseApi := api{
		baseURL: regex.ReplaceAllString(c.Server, "") + "/dg-api",
		Logger:  logger.With("api", "Cfx"),
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
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	if err := encoder.Encode(input); err != nil {
		ca.Logger.Error("Failed to encode input to JSON for Cfx request", "error", err.Error())
		return
	}
	req.Body = io.NopCloser(buf)
}

func (ca *cfxApi) doAuthentication(req *http.Request) {
	req.Header.Add("Authorization", "Bearer "+ca.apiKey)
}
