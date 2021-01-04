package main

import (
	_ "image/png"
	"log"
	"tamago/pkg/game"
	"tamago/pkg/tamagotchi"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func printTamagoStatus(tamago *tamagotchi.Tamagotchi) {
	for tamago.IsAlive() {
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

	ebiten.SetWindowSize(game.ScreenWidth*2, game.ScreenHeight*2)
	ebiten.SetWindowTitle("Tamago")
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
