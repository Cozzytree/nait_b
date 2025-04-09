package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Cozzytree/nait/internal/database"
	"github.com/Cozzytree/nait/internal/model"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second
)

type packet struct {
}

type ws_client struct {
	hub      *hub
	id       uuid.UUID
	conn     *websocket.Conn
	sendChan chan *model.BlockPacket
	logger   log.Logger
}

var upgrader = websocket.Upgrader{}

func (c *ws_client) Id() uuid.UUID {
	return c.id
}

func (c *ws_client) readPump() {
	var block []model.BlockPacket
	defer func() {
		c.close()
	}()

	for {
		err := c.conn.ReadJSON(&block)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Println("unexpected close error:", err)
			} else {
				c.logger.Println("connection closed:", err)
			}
			break // Exit the loop to prevent further reads
		}

		// Echo the message back
		data, err := json.Marshal(&block)
		if err != nil {
			c.logger.Println(err.Error())
		} else {
			c.hub.BroadCastChan <- data

		}
	}
}

func (c *ws_client) writePump() {
	defer func() {
		c.close()
	}()
	for {
		select {
		case block, ok := <-c.sendChan:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteJSON(block)
			if err != nil {
				c.logger.Println("error while writing to client :", c.id)
				continue
			}
		}
	}
}

func (c *ws_client) close() {
	c.logger.Printf("connect closed for client : %v", c.id)
	err := c.conn.Close()
	if err != nil {
		c.logger.Println("error while closing")
	}
	c.hub.UnRegisterChan <- c
}

func new_wsClient(h *hub, w http.ResponseWriter, r *http.Request, user database.User) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := ws_client{
		id:   user.ID,
		conn: conn,
		hub:  h,
		logger: *log.New(log.Writer(),
			fmt.Sprintf(""), log.Flags()),
		sendChan: make(chan *model.BlockPacket),
	}

	client.hub.RegisterChan <- &client

	go client.readPump()
	// go client.writePump()
}
