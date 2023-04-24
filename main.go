package main

import (
	"flag"
	"log"
	"os"

	"cellony/game"
	"cellony/game/assets"
	"cellony/game/config"

	"runtime/pprof"

	"github.com/hajimehoshi/ebiten/v2"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, _ := os.Create(*cpuprofile)
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

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
