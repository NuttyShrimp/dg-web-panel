package api

import (
	"degrens/panel/internal/config"
	"degrens/panel/models"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"

	"github.com/sirupsen/logrus"
)

type graylogApi struct {
	api
	Config config.ConfigGraylog
}

var GraylogApi graylogApi

func CreateGraylogApi(c *config.ConfigGraylog) {
	regex := regexp.MustCompile(`/$`)
	baseApi := api{
		Logger:  logrus.WithField("api", "Graylog"),
		baseURL: regex.ReplaceAllString(c.URL, "") + "/api",
	}
	GraylogApi = graylogApi{
		api:    baseApi,
		Config: *c,
	}
}

func ValidateGraylogApi() bool {
	version := &models.SystemInfo{}
	_, err := GraylogApi.DoRequest(http.MethodGet, "/system", nil, &version)
	if err != nil {
		GraylogApi.Logger.WithError(err).Error("Failed to connect to graylog")
		return false
	}
	stream := &models.Stream{}
	_, err = GraylogApi.DoRequest(http.MethodGet, fmt.Sprint("/streams/", GraylogApi.Config.StreamId), nil, stream)
	if err != nil {
		GraylogApi.Logger.WithError(err).Error("Failed to fetch the stream")
		return false
	}
	GraylogApi.Logger.Info("Did a successful request to graylog")
	return true
}

func (ga *graylogApi) DoRequest(method, endpoint string, input, output interface{}) (*ErrorInfo, error) {
	return ga.doInternalRequest(method, endpoint, input, output, ga.addInput, ga.doAuthentication)
}

func (ga *graylogApi) addInput(req *http.Request, input interface{}) {
	// This generates Query parameters based on the josn tags from input
	// TODO: Check if works for POST requests
	inputVal := reflect.ValueOf(input)
	if inputVal.Kind() == reflect.Ptr {
		inputVal = inputVal.Elem()
	}
	typeOfV := inputVal.Type()

	q := req.URL.Query()

	for i := 0; i < inputVal.NumField(); i++ {
		// Get json  tag
		tag := typeOfV.Field(i).Tag.Get("json")
		if tag != "" && tag != "-" {
			if typeOfV.Field(i).Type.Kind() == reflect.Int {
				field := strconv.Itoa(int(inputVal.Field(i).Int()))
				q.Add(tag, field)
			} else {
				q.Add(tag, inputVal.Field(i).String())
			}
		}
	}
	req.URL.RawQuery = q.Encode()
}

func (ga *graylogApi) doAuthentication(req *http.Request) {
	req.SetBasicAuth(GraylogApi.Config.Token, "token")
}

// should only fetch: source, message, logtype, full_message and timestamp fields
func FetchQuery(query string, limit, timeRange int) (*[]models.ResultMessage, error) {
	options := models.QueryRequestInput{
		Fields: "source,message,logtype,full_message,timestamp",
		Filter: fmt.Sprint("streams:", GraylogApi.Config.StreamId),
		Sort:   "timestamp:desc",
		Query:  query,
		Limit:  limit,
		Range:  timeRange,
	}
	messages := &models.Message{}
	ei, err := GraylogApi.DoRequest(http.MethodGet, "/search/universal/relative", &options, messages)
	if err != nil {
		GraylogApi.Logger.WithFields(logrus.Fields{
			"query":     query,
			"errorMsg":  ei.Message,
			"errorType": ei.Type}).WithError(err).Error("A query to graylog failed")
		return nil, err
	}
	return &messages.Messages, nil
}
