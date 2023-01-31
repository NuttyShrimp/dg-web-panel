package staff

import (
	"degrens/panel/internal/api"
	"degrens/panel/models"
	"errors"
	"net/http"
	"time"
)

type TimelineEvent struct {
	Title string `json:"title"`
	Time  int64  `json:"time"`
	Type  string `json:"type"`
}

type DashboardInfo struct {
	ActivePlayers  int                   `json:"activePlayers"`
	PlayersQueue   []models.CfxCharacter `json:"queue"`
	PlayersInQueue int                   `json:"queuedPlayers"`
	JoinEvents     []TimelineEvent       `json:"joinEvents"`
}

var lastFetch time.Time
var savedInfo *DashboardInfo

func fetchDashboardInfo() error {
	info := DashboardInfo{}
	// Fetch real-time info from Cfx
	ei, err := api.CfxApi.DoRequest(http.MethodGet, "/info", nil, &info)
	if err != nil {
		if errors.Is(err, &models.RouteError{}) {
			return err
		}
		logger.Error("An error occured while fetching from the Cfx Api", "endpoint", "/info", "error", err.Error())
		return &models.RouteError{
			Message: models.RouteErrorMessage{
				Title:       "Unexpected error from the fivem server",
				Description: "Nothing you can do about it. Error is reported to devs",
			},
			Code: 500,
		}
	}
	if ei.Message != "" {
		logger.Error("Failed to fetch info from CfxServer", "error", ei.Message)
		return nil
	}
	// Get timeline Events from graylog
	results, err := api.FetchQuery("logtype:(\"chars:select\" OR \"core:left\" OR \"core:joined\" OR \"chars:created\")", 30, 0)
	if results == nil {
		return err
	}
	events := []TimelineEvent{}
	for _, message := range *results {
		t, err := time.Parse(time.RFC3339, message.Message["timestamp"])
		if err != nil {
			logger.Error("Failed to parse "+message.Message["timestamp"]+" to a time struct", "error", err)
			continue
		}
		event := TimelineEvent{
			Title: message.Message["message"],
			Time:  t.Unix(),
			Type:  message.Message["logtype"],
		}
		events = append(events, event)
	}
	info.JoinEvents = events
	savedInfo = &info
	lastFetch = time.Now()
	return nil
}

func GetDashboardInfo() (*DashboardInfo, *models.RouteError) {
	if time.Now().After(lastFetch.Add(20 * time.Second)) {
		if err := fetchDashboardInfo(); err != nil {
			return nil, err.(*models.RouteError)
		}
	}
	return savedInfo, nil
}
