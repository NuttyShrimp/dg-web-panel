package admin

import (
	"degrens/panel/internal/api"
	"degrens/panel/lib/graylogger"
	"degrens/panel/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aidenwallis/go-utils/utils"
)

func FetchPanelLogs(offset int, query string) ([]*graylogger.GraylogMessage[map[string]interface{}], error) {
	options := models.QueryRequestInput{
		Fields: "message,logtype,full_message,timestamp",
		Filter: fmt.Sprint("streams:", api.GraylogApi.Config.PanelStreamId),
		Sort:   "timestamp:desc",
		Query:  query,
		Limit:  150,
		Range:  0,
		Offset: offset,
	}
	messages := &models.Message{}
	ei, err := api.GraylogApi.DoRequest(http.MethodGet, "/search/universal/relative", &options, messages)
	if err != nil {
		api.GraylogApi.Logger.Error("A query to graylog failed", "error", err, "query", query, "ErrorMsg", ei.Message, "ErrorType", ei.Type)
		return nil, err
	}
	// serialize, deserialize to also get the json strings in the msgs to be deserialized for our struct
	return utils.SliceMap(messages.Messages, func(v models.ResultMessage) *graylogger.GraylogMessage[map[string]interface{}] {
		v.Message["short_message"] = v.Message["message"]
		v.Message["_logtype"] = v.Message["logtype"]
		msgStr, _ := json.Marshal(v.Message)
		msg := graylogger.GraylogMessage[map[string]interface{}]{}
		err := json.Unmarshal(msgStr, &msg)
		if err == nil {
			logger.Error("Failed to decode panel message", "error", err)
			return nil
		}
		msgTime, _ := time.Parse(time.RFC3339, v.Message["timestamp"])
		msg.Timestamp = msgTime.Unix()

		// For some reason is this not getting deserialized to the map
		msg.FullMessage = make(map[string]interface{})
		err = json.Unmarshal([]byte(v.Message["full_message"]), &msg.FullMessage)
		if err == nil {
			logger.Error("Failed to decode panel message", "error", err)
			return nil
		}
		return &msg
	}), nil
}
