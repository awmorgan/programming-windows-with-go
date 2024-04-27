package main

import (
	"fmt"
	"os"

	"github.com/awmorgan/programming-windows-with-go/internal/sys/windows"
)

func wndproc(hWnd windows.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	return 0
}

func main() {
	var wc windows.WNDCLASS
	wc.Style = windows.CS_HREDRAW | windows.CS_VREDRAW
	wc.LpfnWndProc = windows.NewWndproc(wndproc)
	wc.CbClsExtra = 0
	wc.CbWndExtra = 0
	wc.HInstance = windows.WinmainArgs.HInstance
	wc.HIcon, _ = windows.LoadIcon(0, windows.IDI_APPLICATION)
	wc.HCursor, _ = windows.LoadCursor(0, windows.IDC_ARROW)
	wc.HbrBackground = windows.HBRUSH(windows.GetStockObject(windows.WHITE_BRUSH))
	wc.LpszMenuName = nil
	wc.LpszClassName = windows.StringToUTF16Ptr("HelloWin")
	if _, err := windows.RegisterClass(&wc); err != nil {
		windows.MessageBox(0, fmt.Sprintf("RegisterClass failed: %v", err), "Error", windows.MB_ICONERROR)
		os.Exit(1)
	}

}
