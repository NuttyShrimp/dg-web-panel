package staff

import (
	"degrens/panel/internal/api"
	"degrens/panel/lib/graylogger"
	"degrens/panel/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aidenwallis/go-utils/utils"
)

func FetchCfxLogs(page int, query string) ([]*graylogger.GraylogMessage[string], int, error) {
	options := models.QueryRequestInput{
		Fields: "message,logtype,full_message,timestamp",
		Filter: fmt.Sprint("streams:", api.GraylogApi.Config.StreamId),
		Sort:   "timestamp:desc",
		Query:  query,
		Limit:  100,
		Range:  0,
		Offset: page * 100,
	}
	messages := &models.Message{}
	// TODO: refactor this to use the same fetching method as the graylog webui
	ei, err := api.GraylogApi.DoRequest(http.MethodGet, "/search/universal/relative", &options, messages)
	if err != nil {
		api.GraylogApi.Logger.Error("A query to graylog failed", "error", err, "query", query, "ErrorMsg", ei.Message, "ErrorType", ei.Type)
		return nil, 0, err
	}
	// serialize, deserialize to also get the json strings in the msgs to be deserialized for our struct
	return utils.SliceMap(messages.Messages, func(v models.ResultMessage) *graylogger.GraylogMessage[string] {
		v.Message["short_message"] = v.Message["message"]
		v.Message["_logtype"] = v.Message["logtype"]
		timestamp := strings.Clone(v.Message["timestamp"])
		delete(v.Message, "timestamp")
		msgStr, _ := json.Marshal(v.Message)
		msg := graylogger.GraylogMessage[string]{}

		err := json.Unmarshal(msgStr, &msg)
		if err != nil {
			logger.Error("Failed to decode panel message", "error", err)
			return nil
		}
		msgTime, _ := time.Parse(time.RFC3339, timestamp)
		msg.Timestamp = msgTime.Unix()

		return &msg
	}), messages.TotalResults, nil
}
