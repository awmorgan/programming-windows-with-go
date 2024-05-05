package main

import (
	"x/win32"
)

func main() {
	win32.MessageBox(0, win32.Str("Hello, Windows!"), win32.Str("HelloMsg"), win32.MB_OK)
}
