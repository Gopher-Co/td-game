package main

import (
	"github.com/gopher-co/td-game/models/coopstate"
	"net"
)

func main() {
	s := coopstate.NewServer("aboba", 3)
	l, _ := net.Listen("tcp", ":8080")

	if err := s.Serve(l); err != nil {
		panic(err)
	}
}
