package coopstate

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type ClientState struct {
	s GameHostClient
}

func (c ClientState) Draw(image *ebiten.Image) {

}

func (c ClientState) Update() error {
	fmt.Println("fdfdfdf")
	return nil
}

func (c ClientState) End() bool {
	return false
}
