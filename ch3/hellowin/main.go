package main

import (
	"syscall"
	"x/win32"

	"github.com/lxn/win"
)

func wndproc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	return 0
}

func main() {
	var wc win32.WNDCLASS
	wc.Style = win.CS_HREDRAW | win.CS_VREDRAW
	wc.LpfnWndProc = syscall.NewCallback(wndproc)
	wc.CbClsExtra = 0
	wc.CbWndExtra = 0
	wc.HInstance = win32.WinmainArgs.HInstance
	wc.HIcon = win.LoadIcon(0, win.MAKEINTRESOURCE(win.IDI_APPLICATION))
	wc.HCursor = win.LoadCursor(0, win.MAKEINTRESOURCE(win.IDC_ARROW))
	// wc.HbrBackground =

}
