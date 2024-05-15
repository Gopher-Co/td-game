package main

import (
	"github.com/gopher-co/td-game/models/coopstate"
	"log"
	"net"
)

func main() {
	s, id := coopstate.NewServer("aboba", 3)
	l, _ := net.Listen("tcp", ":8080")
	log.Println(id)
	if err := s.Serve(l); err != nil {
		panic(err)
	}
}
