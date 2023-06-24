package staff

import (
	"degrens/panel/internal/api"
	"degrens/panel/models"
	"errors"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type TimelineEvent struct {
	Title string `json:"title"`
	Time  int64  `json:"time"`
	Type  string `json:"type"`
}

type QueuedPlayer struct {
	Name        string            `json:"name"`
	ServerId    uint              `json:"source"`
	Identifiers map[string]string `json:"identifiers"`
}

type DashboardInfo struct {
	ActivePlayers  int             `json:"activePlayers"`
	PlayersQueue   []QueuedPlayer  `json:"queue"`
	PlayersInQueue int             `json:"queuedPlayers"`
	JoinEvents     []TimelineEvent `json:"joinEvents"`
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
		logrus.WithError(err).Error("An error occurred while fetching from the Cfx Api")
		return &models.RouteError{
			Message: models.RouteErrorMessage{
				Title:       "Unexpected error from the fivem server",
				Description: "Nothing you can do about it. Error is reported to devs",
			},
			Code: 500,
		}
	}
	if ei.Message != "" {
		logrus.WithError(errors.New(ei.Message)).Error("Failed to fetch info from CfxServer")
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
			logrus.WithError(err).Error("Failed to parse " + message.Message["timestamp"] + " to a time struct")
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
