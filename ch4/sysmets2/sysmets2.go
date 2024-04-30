package main

import (
	"fmt"
	"os"
	"runtime"
	"unicode/utf8"
	"x/win32"

	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

var cxChar, cxCaps, cyChar int32

func wndproc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_CREATE:
		hdc := win.GetDC(hwnd)
		tm := win.TEXTMETRIC{}
		win.GetTextMetrics(hdc, &tm)
		cxChar = int32(tm.TmAveCharWidth)
		cxCaps = int32((int32(tm.TmPitchAndFamily) & 1) * 3 * cxChar / 2)
		cyChar = int32(tm.TmHeight + tm.TmExternalLeading)
		win.ReleaseDC(hwnd, hdc)
		return 0
	case win.WM_PAINT:
		var ps win.PAINTSTRUCT
		hdc := win.BeginPaint(hwnd, &ps)
		for i, sm := range win32.Sysmetrics {
			win.TextOut(hdc, 0, cyChar*int32(i), sm.Label, int32(sm.LabelLen))
			win.TextOut(hdc, 22*cxCaps, cyChar*int32(i), sm.Desc, int32(sm.DescLen))
			win32.SetTextAlign(hdc, win32.TA_RIGHT|win32.TA_TOP)
			s := fmt.Sprintf("%5d", win.GetSystemMetrics(sm.Index))
			l := utf8.RuneCountInString(s)
			win.TextOut(hdc, 22*cxCaps+40*cxChar, cyChar*int32(i), win32.Str(s), int32(l))
			win32.SetTextAlign(hdc, win32.TA_LEFT|win32.TA_TOP)
		}
		win.EndPaint(hwnd, &ps)
		return 0
	case win.WM_DESTROY:
		win.PostQuitMessage(0)
		return 0
	}
	return win.DefWindowProc(hwnd, msg, wParam, lParam)
}

func main() {
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.
	appName := win32.Str("Sysmets2")
	wc := win32.WNDCLASS{
		Style:         win.CS_HREDRAW | win.CS_VREDRAW,
		LpfnWndProc:   win32.NewWndProc(wndproc),
		HInstance:     win32.WinmainArgs.HInstance,
		HIcon:         win.LoadIcon(0, win.MAKEINTRESOURCE(win.IDI_APPLICATION)),
		HCursor:       win.LoadCursor(0, win.MAKEINTRESOURCE(win.IDC_ARROW)),
		HbrBackground: win.HBRUSH(win.GetStockObject(win.WHITE_BRUSH)),
		LpszClassName: appName,
	}
	if win32.RegisterClass(&wc) == 0 {
		win.MessageBox(0, win32.Str("RegisterClass failed"), appName, win.MB_ICONERROR)
		return
	}
	hwnd := win32.CreateWindow(appName, win32.Str("Get System Metrics No. 1"), win.WS_OVERLAPPEDWINDOW|win.WS_VSCROLL,
		win.CW_USEDEFAULT, win.CW_USEDEFAULT, win.CW_USEDEFAULT, win.CW_USEDEFAULT,
		0, 0, win32.WinmainArgs.HInstance, nil)

	win.ShowWindow(hwnd, win32.WinmainArgs.NCmdShow)
	win.UpdateWindow(hwnd)
	msg := win.MSG{}
	for {
		ret := win.GetMessage(&msg, 0, 0, 0)
		switch ret {
		case 0: // WM_QUIT
			os.Exit(int(msg.WParam))
		case -1: // error
			errMsg := fmt.Sprintf("GetMessage failed: %v", windows.GetLastError())
			win.MessageBox(0, win32.Str(errMsg), appName, win.MB_ICONERROR)
			os.Exit(1)
		}
		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)
	}

}
