package main

import (
	"log"

	"cellony/game"
	"cellony/game/assets"
	"cellony/game/config"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	err := assets.InitializeAssets()
	if err != nil {
		panic(err)
	}

	err = config.LoadConfig()
	if err != nil {
		panic(err)
	}

	ebiten.SetWindowSize(int(config.Video.Width), int(config.Video.Height))
	ebiten.SetWindowTitle("Hello, World!")

	g := game.CreateGame()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
