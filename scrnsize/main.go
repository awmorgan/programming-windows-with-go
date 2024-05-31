package main

import (
	"fmt"
	"x/win32"
)

func main() {
	cxScreen := win32.GetSystemMetrics(win32.SM_CXSCREEN)
	cyScreen := win32.GetSystemMetrics(win32.SM_CYSCREEN)
	text := fmt.Sprintf("The screen is %d pixels wide by %d pixels high.", cxScreen, cyScreen)
	win32.MessageBox(0, text, "ScreenSize", win32.MB_OK)
}
