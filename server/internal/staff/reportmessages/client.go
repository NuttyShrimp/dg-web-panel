package reportmessages

import (
	"bytes"
	"degrens/panel/internal/auth/authinfo"
	"degrens/panel/lib/errors"
	"degrens/panel/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	room *Room

	// Buffered channels of outbound network
	// Room sends data to this channel
	send     chan []byte
	userinfo *authinfo.AuthInfo
}

type ClientMessage struct {
	Message []byte
	Client  *Client
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

func JoinReportRoom(ctx *gin.Context, room *Room) {
	clientInfoPtr, exists := ctx.Get("userInfo")
	clientInfo := clientInfoPtr.(*authinfo.AuthInfo)
	if exists == false {
		room.logger.Error("Failed to retrieve userinfo when joining report room")
		ctx.JSON(http.StatusForbidden, errors.Unauthorized)
		return
	}
	conn, err := wsupgrades.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		room.logger.Error("Failed to upgrade connection", "error", err)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Websocket error",
			Description: "We failed to upgrade your connection to open a websocket",
		})
		return
	}
	client := &Client{
		conn:     conn,
		room:     room,
		send:     make(chan []byte, 256),
		userinfo: clientInfo,
	}
	client.room.register <- client
	go client.readRoutine()
	go client.writeRoutine()
}

func (c *Client) readRoutine() {
	defer func() {
		c.room.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.room.logger.Error("Report websocket was closed unexpectedly", "error", err)
			}
			break
		}
		message = bytes.TrimSpace(message)
		cMsg := ClientMessage{
			Message: message,
			Client:  c,
		}
		c.room.broadcast <- &cMsg
	}
}

func (c *Client) writeRoutine() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
