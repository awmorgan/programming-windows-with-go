package main

import (
	"os"
	"runtime"
	"x/win32"
)

func main() {
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.
	appName := "winmain"
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
	hwnd, _ := win32.CreateWindow(appName, "winmain",
		win32.WS_OVERLAPPEDWINDOW,
		win32.CW_USEDEFAULT, win32.CW_USEDEFAULT,
		win32.CW_USEDEFAULT, win32.CW_USEDEFAULT,
		0, 0, win32.HInstance(), 0)
	win32.ShowWindow(hwnd, win32.NCmdShow())
	win32.UpdateWindow(hwnd)
	msg := win32.MSG{}
	for {
		ret, err := win32.GetMessage(&msg, 0, 0, 0)
		if err != nil {
			errMsg := "GetMessage failed: " + err.Error()
			win32.MessageBox(0, errMsg, appName, win32.MB_ICONERROR)
			os.Exit(1)
		}
		if ret == 0 {
			os.Exit(int(msg.WParam))
		}
		win32.TranslateMessage(&msg)
		win32.DispatchMessage(&msg)
	}
}

var cxChar, cxCaps, cyChar int32

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	var hdc win32.HDC
	var ps win32.PAINTSTRUCT
	var tm win32.TEXTMETRIC

	switch msg {
	case win32.WM_CREATE:
		hdc = win32.GetDC(hwnd)
		win32.GetTextMetrics(hdc, &tm)
		cxChar = tm.TmAveCharWidth
		cxCaps = cxChar
		if tm.TmPitchAndFamily&1 == 1 {
			cxCaps = 3 * cxChar / 2
		}
		cyChar = tm.TmHeight + tm.TmExternalLeading
		win32.ReleaseDC(hwnd, hdc)
		return 0

	case win32.WM_PAINT:
		hdc = win32.BeginPaint(hwnd, &ps)
		win32.EndPaint(hwnd, &ps)
		return 0

	case win32.WM_DESTROY:
		win32.PostQuitMessage(0)
		return 0
	}

	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}
