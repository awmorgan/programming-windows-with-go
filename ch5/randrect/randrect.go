package main

//go build -ldflags -H=windowsgui

import (
	"math/rand/v2"
	"os"
	"runtime"
	"x/win32"
)

func main() {
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.
	appName := "RandRect"
	wc := win32.WNDCLASS{
		Style:         win32.CS_HREDRAW | win32.CS_VREDRAW,
		LpfnWndProc:   win32.NewWndProc(wndproc),
		HInstance:     win32.HInstance(),
		HIcon:         win32.ApplicationIcon(),
		HCursor:       win32.ArrowCursor(),
		HbrBackground: win32.WhiteBrush(),
		LpszClassName: win32.StringToUTF16Ptr(appName),
	}
	if _, err := win32.RegisterClass(&wc); err != nil {
		errMsg := "RegisterClass failed: " + err.Error()
		win32.MessageBox(0, errMsg, appName, win32.MB_ICONERROR)
		return
	}
	hwnd, _ := win32.CreateWindow(appName, "Random Rectangles",
		win32.WS_OVERLAPPEDWINDOW,
		win32.CW_USEDEFAULT, win32.CW_USEDEFAULT,
		win32.CW_USEDEFAULT, win32.CW_USEDEFAULT,
		0, 0, win32.HInstance(), 0)
	win32.ShowWindow(hwnd, win32.NCmdShow())
	win32.UpdateWindow(hwnd)
	msg := win32.MSG{}
	for {
		ret := win32.PeekMessage(&msg, 0, 0, 0, win32.PM_REMOVE)
		if ret {
			if msg.Message == win32.WM_QUIT {
				os.Exit(int(msg.WParam))
			}
			win32.TranslateMessage(&msg)
			win32.DispatchMessage(&msg)
		} else {
			DrawRectangle(hwnd)
		}
	}
}

func DrawRectangle(hwnd win32.HWND) {
	if cxClient == 0 || cyClient == 0 {
		return
	}

	var rect win32.RECT
	win32.SetRect(&rect,
		rand.N(cxClient), rand.N(cyClient),
		rand.N(cxClient), rand.N(cyClient))

	r, g, b := byte(rand.N(256)), byte(rand.N(256)), byte(rand.N(256))
	hBrush := win32.CreateSolidBrush(win32.RGB(r, g, b))
	defer win32.DeleteObject(win32.HGDIOBJ(hBrush))
	hdc := win32.GetDC(hwnd)
	defer win32.ReleaseDC(hwnd, hdc)
	win32.FillRect(hdc, &rect, hBrush)
}

var cxClient, cyClient int32

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {

	switch msg {
	case win32.WM_SIZE:
		cxClient = int32(win32.LOWORD(lParam))
		cyClient = int32(win32.HIWORD(lParam))
		return 0

	case win32.WM_DESTROY:
		win32.PostQuitMessage(0)
		return 0
	}

	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}
