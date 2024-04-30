package main

import (
	"fmt"
	"os"
	"runtime"
	"x/win32"

	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

func wndproc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_CREATE:
		win32.PlaySound(win32.Str("hellowin.wav"), 0, win32.SND_FILENAME|win32.SND_ASYNC)
		return 0
	case win.WM_PAINT:
		rect := win.RECT{}
		ps := win.PAINTSTRUCT{}
		hdc := win.BeginPaint(hwnd, &ps)
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
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.
	appName := win32.Str("HelloWin")
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
	hwnd := win32.CreateWindow(appName, win32.Str("The Hello Program"), win.WS_OVERLAPPEDWINDOW,
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
