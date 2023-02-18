package api

import (
	"degrens/panel/lib/log"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AddInput func(req *http.Request, input interface{}) io.Reader
type DoAuthentication func(req *http.Request)

type Api interface {
	DoRequest(method, endpoint string, input, output interface{}) (*ErrorInfo, error)
	addInput(req *http.Request, input interface{})
	doAuthentication(req *http.Request)
}

type api struct {
	baseURL string
	Logger  log.Logger
}

func (a *api) doInternalRequest(method, endpoint string, input, output interface{}, inputHandler AddInput, authHandler DoAuthentication) (*ErrorInfo, error) {
	// prepare request
	var (
		req *http.Request
		err error
	)
	endpoint = fmt.Sprint(a.baseURL, endpoint)

	var body io.Reader = http.NoBody
	if input != nil {
		body = inputHandler(req, input)
	}

	req, err = http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to call http.NewRequest: %s %s: %w", method, endpoint, err)
	}

	req.Header.Set("User-Agent", "go-degrens-panel/1.0")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Requested-By", "go-graylog")

	ei := &ErrorInfo{Request: req}

	authHandler(req)

	hc := http.DefaultClient

	// request
	resp, err := hc.Do(req)
	if err != nil {
		return ei, fmt.Errorf(
			"failed to call API: %s %s: %w", method, endpoint, err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			a.Logger.Error("Failed to close API request body", "error", err)
		}
	}()
	ei.Response = resp

	if resp.StatusCode >= 400 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return ei, fmt.Errorf(
				"API error: failed to read the response body: %s %s %d",
				method, endpoint, resp.StatusCode)
		}
		if err := json.Unmarshal(b, ei); err != nil {
			return ei, fmt.Errorf(
				"failed to parse response body as ErrorInfo: %s %s %d %s: %w",
				method, endpoint, resp.StatusCode, string(b), err)
		}
		return ei, fmt.Errorf(
			"API error: %s %s %d: "+string(b),
			method, endpoint, resp.StatusCode)
	}
	if output != nil {
		if err := json.NewDecoder(ei.Response.Body).Decode(output); err != nil {
			return ei, fmt.Errorf(
				"failed to decode API response body: %s %s: %w",
				method, endpoint, err)
		}
	}
	return ei, nil
}
