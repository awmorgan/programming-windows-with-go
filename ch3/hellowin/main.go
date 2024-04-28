package main

import (
	"x/win32"

	"github.com/lxn/win"
)

func wndproc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	var hdc win.HDC
	var ps win.PAINTSTRUCT
	var rect win.RECT
	switch msg {
	case win.WM_CREATE:
		win32.PlaySound(win32.Str("hellowin.wav"), 0, win32.SND_FILENAME|win32.SND_ASYNC)
		return 0
	case win.WM_PAINT:
		hdc = win.BeginPaint(hwnd, &ps)
		win.GetClientRect(hwnd, &rect)
		win32.DrawText(hdc, win32.Str("Hello, Windows 98!"), -1, &rect, win.DT_SINGLELINE|win.DT_CENTER|win.DT_VCENTER)
		win.EndPaint(hwnd, &ps)
		return 0
	case win.WM_DESTROY:
		win.PostQuitMessage(0)
		return 0
	}
	return win.DefWindowProc(hwnd, msg, wParam, lParam)
}

func main() {
	var appName = win32.Str("HelloWin")
	var wc win32.WNDCLASS
	wc.Style = win.CS_HREDRAW | win.CS_VREDRAW
	wc.LpfnWndProc = win32.NewWndProc(wndproc)
	wc.CbClsExtra = 0
	wc.CbWndExtra = 0
	wc.HInstance = win32.WinmainArgs.HInstance
	wc.HIcon = win.LoadIcon(0, win.MAKEINTRESOURCE(win.IDI_APPLICATION))
	wc.HCursor = win.LoadCursor(0, win.MAKEINTRESOURCE(win.IDC_ARROW))
	wc.HbrBackground = win.HBRUSH(win.GetStockObject(win.WHITE_BRUSH))
	wc.LpszMenuName = nil
	wc.LpszClassName = appName
	if win32.RegisterClass(&wc) == 0 {
		win.MessageBox(0, win32.Str("RegisterClass failed"), appName, win.MB_ICONERROR)
		return
	}
	hwnd := win32.CreateWindow(appName, win32.Str("The Hello Program"), win.WS_OVERLAPPEDWINDOW,
		win.CW_USEDEFAULT, win.CW_USEDEFAULT, win.CW_USEDEFAULT, win.CW_USEDEFAULT,
		0, 0, win32.WinmainArgs.HInstance, nil)

	win.ShowWindow(hwnd, win.SW_SHOWNORMAL)
	win.UpdateWindow(hwnd)
	msg := win.MSG{}
	for win.GetMessage(&msg, 0, 0, 0) != 0 {
		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)
	}
}
