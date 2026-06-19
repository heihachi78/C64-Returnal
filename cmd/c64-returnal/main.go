package main

import (
	"log"

	"c64returnal/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowTitle("C64-Returnal")
	ebiten.SetWindowIcon(game.WindowIcons())
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(game.TargetTPS)

	if err := ebiten.RunGame(game.New()); err != nil {
		log.Fatal(err)
	}
}
