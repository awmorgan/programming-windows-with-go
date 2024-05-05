package main

import (
	"x/win32"
)

func main() {
	win32.MessageBox(0, "Hello, Windows!", "HelloMsg", win32.MB_OK)
}
