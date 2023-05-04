// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	gp "github.com/god-jay/gools/otel/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log"
	"time"
	//_ "net/http/pprof"
)

var totalMsg int64
var pub int64

var (
	MHTTPServerRequestLatencyMs = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "ws",
		Name:      "http_server_request_latency_ms",
		Buckets:   []float64{10, 20, 50, 100, 150, 200, 300, 500, 750, 1000, 2000, 5000, 10000},
	}, []string{})

	MsgSent = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "ws",
		Name:      "msg_sent",
	})
	MsgCount = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "ws",
		Name:      "msg_count",
	})
)

func main() {
	config, err := gp.ResolveConf("etc/config_prometheus.yaml")
	if err != nil {
		panic(err)
	}
	gp.ServeMetrics(config)

	flag.Parse()

	hub := newHub()

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		elapsedTime := time.Since(startTime).Milliseconds()
		MHTTPServerRequestLatencyMs.WithLabelValues().Observe(float64(elapsedTime))
	})

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
	err = r.Run(":8081")
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
