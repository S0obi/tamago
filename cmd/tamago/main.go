package main

import (
	_ "image/png"
	"log"
	"tamago/pkg/game"
	"tamago/pkg/tamagotchi"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

func printTamagoStatus(tamago *tamagotchi.Tamagotchi) {
	for {
		tamago.PrintStatus()
		time.Sleep(2 * time.Second)
	}
}

func main() {
	g := game.Game{}
	g.Init()

	go g.Tamago.Live()
	go g.PlayMusic()
	go printTamagoStatus(g.Tamago)

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Tamago")
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
