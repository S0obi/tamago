package main

import (
	_ "image/png"
	"log"
	"tamago/pkg/game"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

func main() {
	g := game.Game{}
	g.Init()

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Tamago")
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
