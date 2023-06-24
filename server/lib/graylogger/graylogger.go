package graylogger

import (
	"bytes"
	"degrens/panel/internal/config"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type GraylogMessage[M interface{}] struct {
	Id           string `json:"_id,omitempty"`
	Version      string `json:"version"`
	Host         string `json:"host"`
	ShortMessage string `json:"short_message"`
	Type         string `json:"_logtype"`
	FullMessage  M      `json:"full_message"`
	Timestamp    int64  `json:"timestamp"`
}

var gelfURL string

func InitGrayLogger(cfg *config.ConfigGraylog) {
	gelfURL = cfg.Gelf
}

func pushLog(msg *GraylogMessage[interface{}]) {
	msgStr, err := json.Marshal(msg)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"message": &msg,
		}).Error("Failed to serialize graylogger message")
		return
	}
	resp, err := http.Post(gelfURL, "application/json", bytes.NewBuffer(msgStr))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"message": &msg,
		}).WithError(err).Error("Failed to send graylogger message")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"message": &msg,
		}).WithError(err).Error("Failed to close graylogger request")
	}
}

func createMessage() *GraylogMessage[interface{}] {
	return &GraylogMessage[interface{}]{
		Version:   "1.0",
		Host:      "panel.degrensrp.be",
		Timestamp: time.Now().Unix(),
	}
}

func Log(logtype, message string, kvpPair ...interface{}) {
	if len(kvpPair)%2 == 1 {
		logrus.Errorf("log of %s has uneven key-value pairs", logtype)
		return
	}
	fields := make(map[string]interface{})
	if len(kvpPair) != 0 {
		for i := 0; i < len(kvpPair)-1; i += 2 {
			fields[kvpPair[i].(string)] = kvpPair[i+1]
		}
	}

	fieldsStr, err := json.Marshal(&fields)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"logtype": logtype,
			"fields":  fields,
		}).WithError(err).Error("failed to serialize log fields", "error", err, "logtype", logtype, "fields", fields)
		return
	}

	msg := createMessage()
	msg.ShortMessage = message
	msg.Type = logtype
	msg.FullMessage = string(fieldsStr)
	pushLog(msg)
}
