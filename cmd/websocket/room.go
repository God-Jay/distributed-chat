package main

import (
	"sync"
	"sync/atomic"
	"time"
)

type RoomInfo struct {
	Id        string
	Name      string
	ClientNum uint
	MsgNum    uint

	StartAt time.Time
}
type Room struct {
	sync.RWMutex

	hub *Hub

	info    *RoomInfo
	clients map[*Client]struct{}

	msgSystem MsgSystem

	broadcast chan []byte
	enter     chan *Client
	leave     chan *Client
}

func (r *Room) run() {
	r.msgSystem.subscribe()
	defer r.msgSystem.unsubscribe()

	for {
		select {
		case client := <-r.enter:
			r.Lock()
			r.clients[client] = struct{}{}
			r.info.ClientNum++
			r.Unlock()
		case client := <-r.leave:
			//TODO check client exists?
			r.Lock()
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
				r.info.ClientNum--
			}

			// TODO
			if r.info.ClientNum == 0 {
				//r.hub.deleteRoom(r.info.Id, r.info)
				r.Unlock()
				//return
			} else {
				r.Unlock()
			}

		// TODO improve
		case msg := <-r.broadcast:
			r.info.MsgNum++
			r.RLock()
			for client := range r.clients {
				atomic.AddInt64(&totalMsg, 1)
				select {
				case client.send <- msg:
				default:
					delete(r.clients, client)
					close(client.send)
					r.info.ClientNum--
				}
			}
			r.RUnlock()
		}
	}
}
