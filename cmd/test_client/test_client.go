package main

import (
	"fmt"

	"google.golang.org/grpc"

	"github.com/gopher-co/td-game/models/coopstate"
)

func main() {
	conn, _ := grpc.Dial("localhost:8080", grpc.WithInsecure(), grpc.WithBlock())
	defer conn.Close()

	c := coopstate.NewGameHostClient(conn)

	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
