package reportmessages

import (
	"degrens/panel/internal/db"
	panel_models "degrens/panel/internal/db/models/panel"
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

	reportId uint
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
			if len(origin) == 0 {
				// request is not a CORS request
				return false
			}
			return true
		},
	}
	rooms = make(map[uint]*Room)
)

func GetRoom(reportId uint, logger log.Logger) *Room {
	room, exists := rooms[reportId]
	if !exists {
		// Create room
		room = &Room{
			broadcast:  make(chan *ClientMessage),
			register:   make(chan *Client),
			unregister: make(chan *Client),
			clients:    make(map[*Client]bool),
			logger:     logger.With("roomId", reportId),
			reportId:   reportId,
		}
		go room.run()
	}
	return room
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.register:
			{
				r.clients[client] = true
				r.sendMessages(client, 0)
				r.logger.Debug("registered new client")
			}
		case client := <-r.unregister:
			{
				if _, ok := r.clients[client]; ok {
					delete(r.clients, client)
					close(client.send)
				}
			}
		// Client messages == a new message to the report
		case clientMessage := <-r.broadcast:
			{
				message, err := r.parseIncomingMessage(clientMessage.Message)
				if err != nil {
					clientMessage.Client.send <- r.generateError(err.Error())
					return
				}
				r.logger.Debug("received a new message", "message", message)
				err = r.handleIncomingMessage(*message, clientMessage.Client)
				if err != nil {
					clientMessage.Client.send <- r.generateError(err.Error())
					return
				}
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
	err := db.MariaDB.Client.Order("id DESC").Offset(offset*50).Limit(50).Where("report_id = ?", r.reportId).Find(&msgs).Error
	if err != nil {
		r.logger.Error("Failed to fetch messages", "error", err)
		c.send <- r.generateError("Failed to fetch message batch")
	}

	for i := range msgs {
		// seed messages to seperate array
		// If error happens the mesasge is replaced with a placeholder
		// indicating that there is an issue with
		// that message
		err := SeedReportMessageMember(&msgs[i])
		if err != nil {
			r.logger.Error("Failed to seed reportmessage", "messageId", msgs[i].ID)
		}
	}

	// TODO replace imageId with link to minio
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
		// TODO data should be marshaled to str not
		reportMsg, err := saveMessage(r.reportId, msg.Data, origin.userinfo)
		if err != nil {
			r.logger.Error("Failed to save new report message", "error", err, "message", msg.Data)
			return errors.New("Failed to save message")
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
	default:
		return errors.New("Invalid action")
	}
}
