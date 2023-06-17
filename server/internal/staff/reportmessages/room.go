package reportmessages

import (
	"degrens/panel/internal/db"
	panel_models "degrens/panel/internal/db/models/panel"
	"degrens/panel/internal/staff/reports"
	"degrens/panel/internal/users"
	dgerrors "degrens/panel/lib/errors"
	"degrens/panel/lib/log"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Room struct {
	// Registered clients
	clients map[*Client]bool

	// Buffered channel with messages from clients
	broadcast chan *ClientMessage

	register   chan *Client
	unregister chan *Client

	logger log.Logger

	report *reports.Report
}

// The messages we send as server to clients
type WebsocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

var (
	wsupgrades = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Code is mostly stolen from gin-contrib/cors
			origin := r.Header.Get("Origin")
			return origin != ""
		},
	}
	rooms = make(map[uint]*Room)
)

func GetRoom(report *reports.Report, logger log.Logger) *Room {
	room, exists := rooms[report.Data.ID]
	if !exists {
		// Create room
		room = &Room{
			broadcast:  make(chan *ClientMessage),
			register:   make(chan *Client),
			unregister: make(chan *Client),
			clients:    make(map[*Client]bool),
			logger:     logger.With("roomId", report.Data.ID),
			report:     report,
		}
		rooms[report.Data.ID] = room
		go room.run()
	}
	return room
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
			r.sendMessages(client, 0)
			r.logger.Debug("registered new client")
		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
			}
		// Client messages == a new message to the report
		case clientMessage := <-r.broadcast:
			message, err := r.parseIncomingMessage(clientMessage.Message)
			if err != nil {
				clientMessage.Client.send <- r.generateError(err.Error())
				return
			}
			err = r.handleIncomingMessage(*message, clientMessage.Client)
			if err != nil {
				clientMessage.Client.send <- r.generateError(err.Error())
				return
			}
		}
	}
}

func (r *Room) sendToClients(message []byte) {
	for client := range r.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(r.clients, client)
		}
	}
}

func (r *Room) generateError(message string) []byte {
	msg := WebsocketMessage{
		Type: "error",
		Data: gin.H{
			"title":       "Websocket error",
			"description": message,
		},
	}
	byteMsg, err := json.Marshal(msg)
	if err != nil {
		r.logger.Error("Failed to generate error message for websocket", "error", err)
		return []byte{}
	}
	return byteMsg
}

func (r *Room) parseIncomingMessage(msgArr []byte) (*WebsocketMessage, error) {
	msg := WebsocketMessage{}
	marshalErr := json.Unmarshal(msgArr, &msg)
	if marshalErr != nil {
		dgerrors.HandleJsonError(marshalErr, r.logger)
		return nil, errors.New("Failed to parse message")
	}
	return &msg, nil
}

func (r *Room) sendMessages(c *Client, offset int) {
	var msgs []panel_models.ReportMessage
	err := db.MariaDB.Client.Order("id DESC").Offset(offset*50).Limit(50).Where("report_id = ?", r.report.Data.ID).Find(&msgs).Error
	if err != nil {
		r.logger.Error("Failed to fetch messages", "error", err)
		c.send <- r.generateError("Failed to fetch message batch")
	}

	for i := range msgs {
		// seed messages to separate array
		// If error happens the mesasge is replaced with a placeholder
		// indicating that there is an issue with
		// that message
		err = SeedReportMessageMember(&msgs[i])
		if err != nil {
			r.logger.Error("Failed to seed reportmessage", "messageId", msgs[i].ID, "error", err)
		}
	}

	response := WebsocketMessage{
		Type: "addMessages",
		Data: msgs,
	}
	responseStr, err := json.Marshal(response)
	if err != nil {
		r.logger.Error("Failed to encode websocket message while trying to load messages", "error", err)
		c.send <- r.generateError("Failed to fetch message batch")
	}
	c.send <- responseStr
}

// TODO: accept uploading images
// images will be stored in minio buckets
// Each report will have its unique bucket
// Id for bucket will be a salt + reportId hashed with sha256 or equiv.

func (r *Room) handleIncomingMessage(msg WebsocketMessage, origin *Client) error {
	switch msg.Type {
	case "addMessage":
		// Prevent ghost messages from crashing the server
		if msg.Data == nil || msg.Data.(string) == "" {
			return nil
		}
		reportMsg, err := r.report.AddMessage(r.report.Data.ID, msg.Data, origin.authInfo)
		if err != nil {
			r.logger.Error("Failed to save new report message", "error", err, "message", msg.Data)
			return errors.New("Failed to save message")
		}
		err = SeedReportMessageMember(reportMsg)
		if err != nil {
			r.logger.Error("Failed to seed new report message", "error", err, "message", msg.Data)
			return errors.New("Failed to seed message")
		}
		// announce new message for all clients
		response := WebsocketMessage{
			Type: "addMessage",
			Data: reportMsg,
		}
		responseStr, err := json.Marshal(response)
		if err != nil {
			r.logger.Error("Failed to encode websocket message while trying to announce new report message", "error", err)
			return errors.New("Failed to spread new message")
		}
		r.sendToClients(responseStr)
		return nil
	case "removeMember":
		if !users.DoesUserHaveRole(origin.authInfo.Roles, "staff") {
			return errors.New("missing permissions to do this")
		}
		steamId, ok := msg.Data.(string)
		if !ok {
			return errors.New("Failed to convert data to a valid steamId")
		}
		err := r.report.RemoveMember(steamId)
		if err != nil {
			r.logger.Error("Failed to remove member from report", "error", err)
			return errors.New("Failed to remove member")
		}

		response := WebsocketMessage{
			Type: "setMembers",
			Data: nil,
		}
		responseStr, err := json.Marshal(response)
		if err != nil {
			r.logger.Error("Failed to encode websocket message while removing a member", "error", err)
			return errors.New("Failed to set new members")
		}
		r.sendToClients(responseStr)
		return nil
	case "addMember":
		// Currently allow players to add others to reports, if this is getting to much fucked with
		// we will add restrictions to it
		// if !users.DoesUserHaveRole(origin.authInfo.Roles, "staff") {
		// 	return errors.New("missing permissions to do this")
		// }
		steamId, ok := msg.Data.(string)
		if !ok {
			return errors.New("Failed to convert data to a valid steamId")
		}
		err := r.report.AddMember(steamId)
		if err != nil {
			r.logger.Error("Failed to add member", "error", err)
			return errors.New("Failed to add member")
		}

		response := WebsocketMessage{
			Type: "setMembers",
			Data: nil,
		}
		responseStr, err := json.Marshal(response)
		if err != nil {
			r.logger.Error("Failed to encode websocket message while adding a member", "error", err)
			return errors.New("Failed to set new members")
		}
		r.sendToClients(responseStr)
		return nil
	case "toggleReportState":
		toggle, ok := msg.Data.(bool)
		if !ok {
			return errors.New("Failed to convert data to a boolean")
		}
		err := r.report.ToggleState(toggle)
		if err != nil {
			r.logger.Error("Failed to toggle report state", "state", toggle)
			return errors.New("Failed to change report state")
		}

		response := WebsocketMessage{
			Type: "toggleState",
			Data: toggle,
		}
		responseStr, err := json.Marshal(response)
		if err != nil {
			r.logger.Error("Failed to encode websocket message while chaning state", "error", err)
			return errors.New("Failed to announce report state change")
		}
		r.sendToClients(responseStr)
		return nil
	default:
		return errors.New("Invalid action")
	}
}
