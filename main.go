package main

import (
	"log"

	"cellony/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")

	g := game.Game{}

	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
