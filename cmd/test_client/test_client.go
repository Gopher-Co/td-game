package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/gopher-co/td-game/models/coopstate"
	"google.golang.org/grpc"
)

func main() {
	conn, _ := grpc.Dial("localhost:8080", grpc.WithInsecure(), grpc.WithBlock())
	defer conn.Close()

	c := coopstate.NewServerClient(conn)
	resp, err := c.JoinLobby(context.Background(), &coopstate.JoinLobbyRequest{
		Player: &coopstate.Player{
			Id: &coopstate.PlayerId{
				Uuid:     uuid.NewString(),
				Nickname: "abobas",
			},
		},
		Lobby: &coopstate.LobbyId{
			Name: "8d67d227-b7eb-451f-a8aa-6a359ea65607",
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
