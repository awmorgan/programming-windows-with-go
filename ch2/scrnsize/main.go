package main

import (
	"fmt"

	"github.com/awmorgan/programming-windows-with-go/internal/windows"
)

func main() {
	cxScreen, _ := windows.GetSystemMetrics(windows.SM_CXSCREEN)
	cyScreen, _ := windows.GetSystemMetrics(windows.SM_CYSCREEN)
	text := fmt.Sprintf("The screen is %d pixels wide by %d pixels high.", cxScreen, cyScreen)
	windows.MessageBox(0, text, "Screen Size", 0)
}
