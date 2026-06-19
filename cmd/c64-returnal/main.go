package main

import (
	"log"

	"c64returnal/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	handleRunError(run())
}

var (
	runEbitenGame  = ebiten.RunGame
	handleRunError = logRunError
	fatal          = log.Fatal
)

func logRunError(err error) {
	if err != nil {
		fatal(err)
	}
}

func run() error {
	ebiten.SetWindowTitle("C64-Returnal")
	ebiten.SetWindowIcon(game.WindowIcons())
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(game.TargetTPS)

	return runEbitenGame(game.New())
}
