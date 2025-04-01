package server

import (
	"encoding/json"
	"fmt"

	"github.com/Cozzytree/nait/internal/model"
	"github.com/google/uuid"
)

type hub struct {
	blockStore     map[uuid.UUID]model.Block
	clients        map[uuid.UUID]*ws_client
	RegisterChan   chan *ws_client
	UnRegisterChan chan *ws_client
	BroadCastChan  chan []byte
}

func initWS() *hub {
	return &hub{
		blockStore:     map[uuid.UUID]model.Block{},
		clients:        map[uuid.UUID]*ws_client{},
		RegisterChan:   make(chan *ws_client),
		UnRegisterChan: make(chan *ws_client),
		BroadCastChan:  make(chan []byte),
	}
}

func (h *hub) run() {
	for {
		select {
		case client := <-h.RegisterChan:
			client.logger.Println("client registered", client.id)
			h.clients[client.Id()] = client
			fmt.Println(len(h.clients))
		case client := <-h.UnRegisterChan:
			if _, ok := h.clients[client.id]; ok {
				delete(h.clients, client.id)
			}
		case packet := <-h.BroadCastChan:
			var block []model.Block
			err := json.Unmarshal(packet, &block)
			if err != nil {
			}

			// for _, w := range h.clients {
			// 	if packet.User_Id != w.id {
			// 		select {
			// 		case w.sendChan <- packet:
			// 		default:
			// 			// If sendChan is blocked, log and skip the client
			// 			w.logger.Println("sendChan blocked, skipping client:", w.id)
			// 		}
			// 	}
			// }
		}
	}
}
