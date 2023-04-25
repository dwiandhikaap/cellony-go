package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"cellony/game"
	"cellony/game/assets"
	"cellony/game/config"
	"cellony/game/util"

	"runtime/pprof"

	"github.com/hajimehoshi/ebiten/v2"

	"time"

	_ "github.com/silbinarywolf/preferdiscretegpu"
)

var prof = flag.Int("prof", -1, "enable profiling")

func main() {
	flag.Parse()
	if *prof > 0 {
		// name is current date
		name := time.Now().Format("2006-01-02_15-04-05")
		f, _ := os.Create(fmt.Sprintf("profiling/cpu_%s.out", name))
		pprof.StartCPUProfile(f)

		go _shutdownTimer(*prof)

		println("Running on profiling mode for", *prof, "seconds")
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

	gpu := util.GpuInfo()
	log.Printf("GPU: %s", gpu)

	g := game.CreateGame()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func _shutdownTimer(second int) {
	time.Sleep(time.Duration(second) * time.Second)
	pprof.StopCPUProfile()
	os.Exit(0)
}
