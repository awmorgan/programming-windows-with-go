package main

import (
	"x/win32"

	"github.com/lxn/win"
)

func main() {
	win32.MessageBox(0, win32.Str("Hello, Windows!"), win32.Str("HelloMsg"), win.MB_OK)
}
