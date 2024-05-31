package main

import (
	"fmt"
	"os"
	"runtime"
	"x/win32"
)

func main() {
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.
	appName := "WhatSize"
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
	hwnd, _ := win32.CreateWindow(appName, "What Size is the Window?",
		win32.WS_OVERLAPPEDWINDOW,
		win32.CW_USEDEFAULT, win32.CW_USEDEFAULT,
		win32.CW_USEDEFAULT, win32.CW_USEDEFAULT,
		0, 0, win32.HInstance(), 0)
	win32.ShowWindow(hwnd, win32.NCmdShow())
	win32.UpdateWindow(hwnd)
	var msg win32.MSG
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

func Show(hwnd win32.HWND, hdc win32.HDC, xText int32, yText int32, iMapMode int32, mapMode string) {
	win32.SaveDC(hdc)
	win32.SetMapMode(hdc, iMapMode)
	var rect win32.RECT
	win32.GetClientRect(hwnd, &rect)
	var pts = []win32.POINT{{X: rect.Left, Y: rect.Top}, {X: rect.Right, Y: rect.Bottom}}
	win32.DPtoLP(hdc, pts)
	win32.RestoreDC(hdc, -1)
	s := fmt.Sprintf("%-20s %7d %7d %7d %7d", mapMode, pts[0].X, pts[1].X, pts[0].Y, pts[1].Y)
	//todo change textout to not pass len
	win32.TextOut(hdc, xText, yText, s, len(s))
}

var heading = "Map Mode                Left   Right     Top  Bottom"
var undline = "--------                ----   -----     ---  ------"
var cxChar, cyChar int32

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {

	switch msg {
	case win32.WM_CREATE:
		hdc := win32.GetDC(hwnd)
		defer win32.ReleaseDC(hwnd, hdc)
		win32.SelectObject(hdc, win32.GetStockObject(win32.SYSTEM_FIXED_FONT))
		var tm win32.TEXTMETRIC
		win32.GetTextMetrics(hdc, &tm)
		cxChar = tm.TmAveCharWidth
		cyChar = tm.TmHeight + tm.TmExternalLeading
		return 0

	case win32.WM_PAINT:
		var ps win32.PAINTSTRUCT
		hdc := win32.BeginPaint(hwnd, &ps)
		defer win32.EndPaint(hwnd, &ps)
		win32.SelectObject(hdc, win32.GetStockObject(win32.SYSTEM_FIXED_FONT))
		win32.SetMapMode(hdc, win32.MM_ANISOTROPIC)
		win32.SetWindowExtEx(hdc, 1, 1, nil)
		win32.SetViewportExtEx(hdc, cxChar, cyChar, nil)
		win32.TextOut(hdc, 1, 1, heading, len(heading))
		win32.TextOut(hdc, 1, 2, undline, len(undline))
		Show(hwnd, hdc, 1, 3, win32.MM_TEXT, "TEXT (pixels)")
		Show(hwnd, hdc, 1, 4, win32.MM_LOMETRIC, "LOMETRIC (.1 mm)")
		Show(hwnd, hdc, 1, 5, win32.MM_HIMETRIC, "HIMETRIC (.01 mm)")
		Show(hwnd, hdc, 1, 6, win32.MM_LOENGLISH, "LOENGLISH (.01 in)")
		Show(hwnd, hdc, 1, 7, win32.MM_HIENGLISH, "HIENGLISH (.001 in)")
		Show(hwnd, hdc, 1, 8, win32.MM_TWIPS, "TWIPS (1/1440 in)")
		return 0

	case win32.WM_DESTROY:
		win32.PostQuitMessage(0)
		return 0
	}

	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}
