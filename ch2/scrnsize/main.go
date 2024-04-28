package main

import (
	"fmt"
	"x/win32"

	"github.com/lxn/win"
)

func main() {
	cxScreen := win.GetSystemMetrics(win.SM_CXSCREEN)
	cyScreen := win.GetSystemMetrics(win.SM_CYSCREEN)
	text := fmt.Sprintf("The screen is %d pixels wide by %d pixels high.", cxScreen, cyScreen)
	win.MessageBox(0, win32.Str(text), win32.Str("ScreenSize"), win.MB_OK)
}
