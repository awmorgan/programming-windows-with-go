package main

import (
	"github.com/awmorgan/programming-windows-with-go/internal/windows"
)

func main() {
	windows.MessageBox(0, "Hello, Windows!", "HelloMsg", windows.MB_OK)
}
