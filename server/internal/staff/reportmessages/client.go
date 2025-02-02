package reportmessages

import (
	"bytes"
	"degrens/panel/internal/auth/authinfo"
	dgerrors "degrens/panel/lib/errors"
	"degrens/panel/models"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Client struct {
	conn *websocket.Conn
	room *Room

	// Buffered channels of outbound network
	// Room sends data to this channel
	send     chan []byte
	authInfo *authinfo.AuthInfo
	logger   *logrus.Entry
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
	if !exists {
		room.logger.Error("Failed to retrieve userinfo when joining report room")
		ctx.JSON(http.StatusForbidden, dgerrors.Unauthorized)
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
		authInfo: clientInfo,
		logger:   room.logger,
	}
	client.room.register <- client
	go client.readRoutine()
	go client.writeRoutine()
}

func (c *Client) readRoutine() {
	defer func() {
		c.room.unregister <- c
		err := c.conn.Close()
		if err != nil {
			c.logger.WithError(err).Error("Failed to properly close report WS")
		}
	}()

	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		c.logger.WithError(err).Error("Failed to set read deadline in report WS")
	}
	c.conn.SetPongHandler(func(appData string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
				c.room.logger.WithError(err).Error("Report websocket was closed unexpectedly")
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
		err := c.conn.Close()
		if err != nil && !errors.Is(err, net.ErrClosed) {
			c.logger.WithError(err).Error("Failed to properly close report WS")
		}
	}()
	doRoutine := false
	for !doRoutine {
		select {
		case message, ok := <-c.send:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				c.logger.WithError(err).Error("Failed to set Write Deadline for report WS")
			}
			if !ok {
				// The hub closed the channel.
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil && !errors.Is(err, websocket.ErrCloseSent) {
					c.logger.WithError(err).Error("Failed to close a report WS")
				}
				doRoutine = true
			}

			c.sendQueuedMsg(message)
		case <-ticker.C:
			c.sendPingMsg()
		}
	}
}

func (c *Client) sendQueuedMsg(message []byte) {
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}
	_, err = w.Write(message)
	if err != nil {
		c.logger.WithError(err).Error("Failed to send a message to the report WS")
		return
	}
	// Add queued chat messages to the current websocket message.
	n := len(c.send)
	for i := 0; i < n; i++ {
		_, err := w.Write([]byte{'\n'})
		if err != nil {
			c.logger.WithError(err).Error("Failed to send a newline through a report WS")
			return
		}
		_, err = w.Write(<-c.send)
		if err != nil {
			c.logger.WithError(err).Error("Failed to send a chat message through a report WS")
			return
		}
	}

	if err := w.Close(); err != nil {
		return
	}
}

func (c *Client) sendPingMsg() {
	if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return
	}
	if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
		return
	}
}
