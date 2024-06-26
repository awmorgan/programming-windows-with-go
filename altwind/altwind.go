package main

import (
	"os"
	"runtime"
	"x/win32"
)

func main() {
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.
	appName := "AltWind"
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
		errMsg := "RegisterClass failed: " + err.Error()
		win32.MessageBox(0, errMsg, appName, win32.MB_ICONERROR)
		return
	}
	hwnd, _ := win32.CreateWindow(appName, "Alternate and Winding Fill Modes",
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

var cxClient, cyClient int32
var aptFigure = [...]win32.POINT{
	{X: 10, Y: 70},
	{X: 50, Y: 70},
	{X: 50, Y: 10},
	{X: 90, Y: 10},
	{X: 90, Y: 50},
	{X: 30, Y: 50},
	{X: 30, Y: 90},
	{X: 70, Y: 90},
	{X: 70, Y: 30},
	{X: 10, Y: 30},
}

const nPts = len(aptFigure)

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {

	switch msg {
	case win32.WM_SIZE:
		cxClient = win32.LOWORD(lParam)
		cyClient = win32.HIWORD(lParam)
		return 0

	case win32.WM_PAINT:
		ps := win32.PAINTSTRUCT{}
		hdc := win32.BeginPaint(hwnd, &ps)
		defer win32.EndPaint(hwnd, &ps)
		apt := [nPts]win32.POINT{}
		win32.SelectObject(hdc, win32.GetStockObject(win32.GRAY_BRUSH))
		for i := range nPts {
			apt[i].X = cxClient * aptFigure[i].X / 200
			apt[i].Y = cyClient * aptFigure[i].Y / 100
		}
		win32.SetPolyFillMode(hdc, win32.ALTERNATE)
		win32.Polygon(hdc, apt[:])

		for i := range nPts {
			apt[i].X += cxClient / 2
		}
		win32.SetPolyFillMode(hdc, win32.WINDING)
		win32.Polygon(hdc, apt[:])
		return 0

	case win32.WM_DESTROY:
		win32.PostQuitMessage(0)
		return 0
	}

	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}
