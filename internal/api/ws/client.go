package ws

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = pongWait * 9 / 10
)

// Client represents a websocket connection.
type Client struct {
	hub            *Hub
	conn           *websocket.Conn
	send           chan []byte
	accountID      string
	conversationID string
	closeOnce      sync.Once
	onMessage      func(context.Context, *Client, []byte)
}

// NewClient constructs a websocket client for the hub.
func NewClient(hub *Hub, conn *websocket.Conn, accountID, conversationID string, onMessage func(context.Context, *Client, []byte)) *Client {
	return &Client{
		hub:            hub,
		conn:           conn,
		send:           make(chan []byte, 16),
		accountID:      accountID,
		conversationID: conversationID,
		onMessage:      onMessage,
	}
}

// Run starts read and write loops for the client.
func (c *Client) Run() {
	c.hub.Register(c)
	go c.writeLoop()
	c.readLoop()
}

// Close shuts down the client connection.
func (c *Client) Close() {
	c.closeOnce.Do(func() {
		c.hub.Unregister(c)
		close(c.send)
		_ = c.conn.Close()
	})
}

func (c *Client) readLoop() {
	defer c.Close()

	c.conn.SetReadLimit(512)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		messageType, payload, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		if messageType != websocket.TextMessage {
			continue
		}
		if c.onMessage != nil {
			go c.onMessage(context.Background(), c, payload)
		}
	}
}

func (c *Client) writeLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SendJSON enqueues a JSON message to the client.
func (c *Client) SendJSON(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	select {
	case c.send <- data:
		return nil
	default:
		go c.Close()
		return errors.New("client send buffer full")
	}
}
