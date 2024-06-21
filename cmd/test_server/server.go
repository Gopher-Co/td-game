package main

import (
	"net"

	"github.com/gopher-co/td-game/models/coopstate"
)

func main() {
	s := coopstate.NewServer("1. Tutorial", 1)
	l, _ := net.Listen("tcp", ":8080")

	if err := s.Serve(l); err != nil {
		panic(err)
	}
}
