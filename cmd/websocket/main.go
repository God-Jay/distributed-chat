// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	//_ "net/http/pprof"
)

var totalMsg int64
var pub int64

func main() {
	flag.Parse()

	hub := newHub()

	r := gin.Default()
	r.GET("/ws/:roomId", func(c *gin.Context) {
		hub.clientEnter(c.Writer, c.Request, c.Param("roomId"))
	})
	r.GET("/roomInfo", func(c *gin.Context) {
		roomInfo := hub.getRoomInfo()
		c.JSON(200, struct {
			DestroyedRooms []*RoomInfo      `json:"destroyed_rooms"`
			Rooms          []ReportRoomInfo `json:"rooms"`
			TotalMsg       int64            `json:"total_msg"`
			Pub            int64            `json:"pub"`
		}{
			DestroyedRooms: hub.destroyedRooms,
			Rooms:          roomInfo,
			TotalMsg:       totalMsg,
			Pub:            pub,
		})
	})
	err := r.Run(":8081")
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
