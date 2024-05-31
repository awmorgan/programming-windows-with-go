package main

import (
	"fmt"
	"os"
	"runtime"
	"x/win32"
)

const (
	ID_TIMER = 1
)

func main() {
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.
	appName := "WhatClr"
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
		errMsg := fmt.Sprintf("RegisterClass failed: %v", err)
		win32.MessageBox(0, errMsg, appName, win32.MB_ICONERROR)
		return
	}

	cx, cy := FindWindowSize()
	hwnd, _ := win32.CreateWindow(appName, "What Color",
		win32.WS_OVERLAPPEDWINDOW|win32.WS_CAPTION|win32.WS_SYSMENU|win32.WS_BORDER,
		win32.CW_USEDEFAULT, win32.CW_USEDEFAULT,
		cx, cy,
		0, 0, win32.HInstance(), 0)

	win32.ShowWindow(hwnd, win32.NCmdShow())
	win32.UpdateWindow(hwnd)
	var msg win32.MSG
	for {
		ret, err := win32.GetMessage(&msg, 0, 0, 0)
		if err != nil {
			errMsg := fmt.Sprintf("GetMessage failed: %v", err)
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

func FindWindowSize() (int32, int32) {
	hdcScreen := win32.CreateIC(win32.StringToUTF16Ptr("DISPLAY"), nil, nil, nil)
	var tm win32.TEXTMETRIC
	win32.GetTextMetrics(hdcScreen, &tm)
	win32.DeleteDC(hdcScreen)
	cx := 2*win32.GetSystemMetrics(win32.SM_CXBORDER) +
		12*tm.TmAveCharWidth
	cy := 2*win32.GetSystemMetrics(win32.SM_CYBORDER) +
		win32.GetSystemMetrics(win32.SM_CYCAPTION) +
		2*tm.TmHeight
	return cx, cy
}

var (
	cr, crLast win32.COLORREF
	hdcScreen  win32.HDC
)

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	switch msg {
	case win32.WM_CREATE:
		hdcScreen = win32.CreateDC(win32.StringToUTF16Ptr("DISPLAY"), nil, nil, nil)
		win32.SetTimer(hwnd, ID_TIMER, 1000, 0)
		return 0

	case win32.WM_TIMER:
		var pt win32.POINT
		win32.GetCursorPos(&pt)
		cr = win32.GetPixel(hdcScreen, pt.X, pt.Y)
		win32.SetPixel(hdcScreen, pt.X, pt.Y, 0)
		if cr != crLast {
			crLast = cr
			win32.InvalidateRect(hwnd, nil, false)
		}
		return 0

	case win32.WM_PAINT:
		var ps win32.PAINTSTRUCT
		hdc := win32.BeginPaint(hwnd, &ps)
		var rc win32.RECT
		win32.GetClientRect(hwnd, &rc)
		s := fmt.Sprintf("  %02X %02X %02X  ",
			win32.GetRValue(cr), win32.GetGValue(cr), win32.GetBValue(cr))
		win32.DrawText(hdc, s, -1, &rc, win32.DT_SINGLELINE|win32.DT_CENTER|win32.DT_VCENTER)
		win32.EndPaint(hwnd, &ps)
		return 0

	case win32.WM_DESTROY:
		win32.KillTimer(hwnd, ID_TIMER)
		win32.PostQuitMessage(0)
		return 0
	}
	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}
