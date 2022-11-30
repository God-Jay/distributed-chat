// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	flag.Parse()

	hub := newHub()

	r := mux.NewRouter()

	r.HandleFunc("/ws/{roomId}", func(w http.ResponseWriter, r *http.Request) {
		log.Println("get...", r.URL)
		hub.clientEnter(w, r)
	})

	r.HandleFunc("/roomInfo", func(w http.ResponseWriter, r *http.Request) {
		roomInfo := hub.getRoomInfo()
		marshal, _ := json.Marshal(roomInfo)
		w.Write(marshal)
	})

	http.Handle("/", r)

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
