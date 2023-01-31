package api

import "net/http"

type ErrorInfo struct {
	Type     string         `json:"type,omitempty"`
	Message  string         `json:"message"`
	Request  *http.Request  `json:"request"`
	Response *http.Response `json:"response"`
}
