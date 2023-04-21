package main

import (
	"log"

	"cellony/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("Hello, World!")

	g := game.CreateGame()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
