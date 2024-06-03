package main

import (
	"fmt"
	"os"
	"runtime"
	"x/win32"
)

func main() {
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.
	appName := "BlokOut1"
	wc := win32.WNDCLASS{
		Style:         win32.CS_HREDRAW | win32.CS_VREDRAW,
		LpfnWndProc:   win32.NewWndProc(wndproc),
		HInstance:     win32.HInstance(),
		HIcon:         win32.LoadIcon(0, win32.IDI_APPLICATION),
		HCursor:       win32.LoadCursor(0, win32.IDC_ARROW),
		HbrBackground: win32.WhiteBrush(),
		LpszClassName: win32.StringToUTF16Ptr(appName),
	}
	if _, err := win32.RegisterClass(&wc); err != nil {
		errMsg := fmt.Sprintf("RegisterClass failed: %v", err)
		win32.MessageBox(0, errMsg, appName, win32.MB_ICONERROR)
		return
	}

	hwnd, _ := win32.CreateWindow(appName, "Mouse Button Demo",
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

func DrawBoxOutline(hwnd win32.HWND, ptBeg, ptEnd win32.POINT) {
	hdc := win32.GetDC(hwnd)
	win32.SetROP2(hdc, win32.R2_NOT)
	win32.SelectObject(hdc, win32.GetStockObject(win32.NULL_BRUSH))
	win32.Rectangle(hdc, ptBeg.X, ptBeg.Y, ptEnd.X, ptEnd.Y)
	win32.ReleaseDC(hwnd, hdc)
}

var blocking, validbox bool
var ptBeg, ptEnd, ptBoxBeg, ptBoxEnd win32.POINT

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	switch msg {
	case win32.WM_LBUTTONDOWN:
		ptBeg.X, ptEnd.X = win32.LOWORD(lParam), win32.LOWORD(lParam)
		ptBeg.Y, ptEnd.Y = win32.HIWORD(lParam), win32.HIWORD(lParam)
		DrawBoxOutline(hwnd, ptBeg, ptEnd)
		win32.SetCapture(hwnd)
		win32.SetCursor(win32.LoadCursor(0, win32.IDC_CROSS))
		blocking = true
		return 0

	case win32.WM_MOUSEMOVE:
		if blocking {
			win32.SetCursor(win32.LoadCursor(0, win32.IDC_CROSS))
			DrawBoxOutline(hwnd, ptBeg, ptEnd)
			ptEnd.X, ptEnd.Y = win32.LOWORD(lParam), win32.HIWORD(lParam)
			DrawBoxOutline(hwnd, ptBeg, ptEnd)
		}
		return 0

	case win32.WM_LBUTTONUP:
		if blocking {
			DrawBoxOutline(hwnd, ptBeg, ptEnd)
			ptBoxBeg = ptBeg
			ptBoxEnd.X, ptBoxEnd.Y = win32.LOWORD(lParam), win32.HIWORD(lParam)
			win32.ReleaseCapture()
			win32.SetCursor(win32.LoadCursor(0, win32.IDC_ARROW))
			blocking = false
			validbox = true
			win32.InvalidateRect(hwnd, nil, true)
		}
		return 0

	case win32.WM_CHAR:
		if blocking && wParam == '\x1B' {
			DrawBoxOutline(hwnd, ptBeg, ptEnd)
			win32.ReleaseCapture()
			win32.SetCursor(win32.LoadCursor(0, win32.IDC_ARROW))
			blocking = false
		}
		return 0

	case win32.WM_PAINT:
		ps := win32.PAINTSTRUCT{}
		hdc := win32.BeginPaint(hwnd, &ps)
		if validbox {
			win32.SelectObject(hdc, win32.GetStockObject(win32.BLACK_BRUSH))
			win32.Rectangle(hdc, ptBoxBeg.X, ptBoxBeg.Y, ptBoxEnd.X, ptBoxEnd.Y)
		}
		if blocking {
			win32.SetROP2(hdc, win32.R2_NOT)
			win32.SelectObject(hdc, win32.GetStockObject(win32.NULL_BRUSH))
			win32.Rectangle(hdc, ptBeg.X, ptBeg.Y, ptEnd.X, ptEnd.Y)
		}
		win32.EndPaint(hwnd, &ps)
		return 0

	case win32.WM_DESTROY:
		win32.PostQuitMessage(0)
		return 0
	}

	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}
