package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"autocell/game"
	"autocell/game/assets"
	"autocell/game/config"
	"autocell/game/util"
	"autocell/generator"

	"runtime/pprof"

	"github.com/hajimehoshi/ebiten/v2"

	"time"
)

var prof = flag.Int("prof", -1, "enable profiling")
var generateTexture = flag.String("generateTexture", "", "generate texture shades")

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

	if *generateTexture != "" {
		generator.GenerateTextureShades(*generateTexture, 10)
		os.Exit(0)
	}

	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	err = assets.InitializeAssets()
	if err != nil {
		panic(err)
	}

	ebiten.SetWindowSize(int(config.Video.Width), int(config.Video.Height))
	ebiten.SetWindowTitle("Autocell")

	gpu := util.GpuInfo()

	g := game.CreateGame()

	runOp := &ebiten.RunGameOptions{}
	runOp.GraphicsLibrary = util.GetGraphicsLibrary()

	log.Printf("GPU: %s", gpu)
	log.Printf("Graphics API: %s", runOp.GraphicsLibrary.String())

	err = ebiten.RunGameWithOptions(g, runOp)
	if err != nil {
		log.Println(err)
	}
}

func _shutdownTimer(second int) {
	time.Sleep(time.Duration(second) * time.Second)
	pprof.StopCPUProfile()
	os.Exit(0)
}
