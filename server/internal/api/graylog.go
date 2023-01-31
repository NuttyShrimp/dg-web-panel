package api

import (
	"degrens/panel/internal/config"
	"degrens/panel/lib/log"
	"degrens/panel/models"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
)

type graylogApi struct {
	api
	Config config.ConfigGraylog
}

var GraylogApi graylogApi

func CreateGraylogApi(c *config.ConfigGraylog, logger *log.Logger) {
	regexp, _ := regexp.Compile(`/$`)
	baseApi := api{
		Logger:  (*logger).With("api", "Graylog"),
		baseURL: regexp.ReplaceAllString(c.URL, "") + "/api",
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
		GraylogApi.Logger.Error("Failed to connect to graylog", "error", err.Error())
		return false
	}
	stream := &models.Stream{}
	_, err = GraylogApi.DoRequest(http.MethodGet, fmt.Sprint("/streams/", GraylogApi.Config.StreamId), nil, stream)
	if err != nil {
		GraylogApi.Logger.Error("Failed to fetch the stream", "error", err.Error())
		return false
	}
	GraylogApi.Logger.Info("Did a successfull request to graylog")
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
func FetchQuery(query string, limit int, timeRange int) (*[]models.ResultMessage, error) {
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
		GraylogApi.Logger.Error("A query to graylog failed", "error", err, "query", query, "ErrorMsg", ei.Message, "ErrorType", ei.Type)
		return nil, err
	}
	return &messages.Messages, nil
}
