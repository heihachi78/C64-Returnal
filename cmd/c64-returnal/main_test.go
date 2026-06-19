package main

import (
	"errors"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestRunConfiguresWindowAndReturnsRunGameError(t *testing.T) {
	wantErr := errors.New("run")
	oldRunGame := runEbitenGame
	defer func() { runEbitenGame = oldRunGame }()

	called := false
	runEbitenGame = func(g ebiten.Game) error {
		called = true
		if g == nil {
			t.Fatal("runEbitenGame received nil game")
		}
		return wantErr
	}

	if err := run(); !errors.Is(err, wantErr) {
		t.Fatalf("run err = %v, want %v", err, wantErr)
	}
	if !called {
		t.Fatal("runEbitenGame was not called")
	}
}

func TestMainDelegatesRunErrorHandling(t *testing.T) {
	oldRunGame := runEbitenGame
	oldHandleRunError := handleRunError
	defer func() {
		runEbitenGame = oldRunGame
		handleRunError = oldHandleRunError
	}()

	wantErr := errors.New("main")
	runEbitenGame = func(ebiten.Game) error { return wantErr }
	handled := false
	handleRunError = func(err error) {
		handled = true
		if !errors.Is(err, wantErr) {
			t.Fatalf("handled err = %v, want %v", err, wantErr)
		}
	}

	main()
	if !handled {
		t.Fatal("main did not delegate to handleRunError")
	}
}

func TestLogRunErrorCallsFatalOnlyForErrors(t *testing.T) {
	oldFatal := fatal
	defer func() { fatal = oldFatal }()

	calls := 0
	fatal = func(v ...any) {
		calls++
		if len(v) != 1 {
			t.Fatalf("fatal args = %v, want one error", v)
		}
	}

	logRunError(nil)
	if calls != 0 {
		t.Fatalf("fatal calls after nil error = %d, want 0", calls)
	}
	logRunError(errors.New("boom"))
	if calls != 1 {
		t.Fatalf("fatal calls after error = %d, want 1", calls)
	}
}
