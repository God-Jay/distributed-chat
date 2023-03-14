package main

type MsgSystem interface {
	publish(msg []byte)
	//subscribe must send []byte msg to the Room.broadcast chan
	subscribe()
	unsubscribe()
}
