package main

//go build -ldflags -H=windowsgui

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"x/win32"
)

func main() {
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.
	appName := "Clover"
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
	hwnd, _ := win32.CreateWindow(appName, "Draw a Clover",
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

var cxClient, cyClient int32
var hRgnClip win32.HRGN

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {

	switch msg {
	case win32.WM_SIZE:
		cxClient = int32(win32.LOWORD(lParam))
		cyClient = int32(win32.HIWORD(lParam))

		hCursor := win32.SetCursor(win32.WaitCursor())
		win32.ShowCursor(true)

		if hRgnClip != 0 {
			win32.DeleteObject(win32.HGDIOBJ(hRgnClip))
		}

		var hRgnTemp [6]win32.HRGN

		hRgnTemp[0] = win32.CreateEllipticRgn(0, cyClient/3,
			cxClient/2, 2*cyClient/3)
		hRgnTemp[1] = win32.CreateEllipticRgn(cxClient/2, cyClient/3,
			cxClient, 2*cyClient/3)
		hRgnTemp[2] = win32.CreateEllipticRgn(cxClient/3, 0,
			2*cxClient/3, cyClient/2)
		hRgnTemp[3] = win32.CreateEllipticRgn(cxClient/3, cyClient/2,
			2*cxClient/3, cyClient)
		hRgnTemp[4] = win32.CreateRectRgn(0, 0, 1, 1)
		hRgnTemp[5] = win32.CreateRectRgn(0, 0, 1, 1)
		hRgnClip = win32.CreateRectRgn(0, 0, 1, 1)

		win32.CombineRgn(hRgnTemp[4], hRgnTemp[0], hRgnTemp[1], win32.RGN_OR)
		win32.CombineRgn(hRgnTemp[5], hRgnTemp[2], hRgnTemp[3], win32.RGN_OR)
		win32.CombineRgn(hRgnClip, hRgnTemp[4], hRgnTemp[5], win32.RGN_XOR)

		for i := 0; i < 6; i++ {
			win32.DeleteObject(win32.HGDIOBJ(hRgnTemp[i]))
		}
		win32.SetCursor(hCursor)
		win32.ShowCursor(false)
		return 0

	case win32.WM_PAINT:
		ps := win32.PAINTSTRUCT{}
		hdc := win32.BeginPaint(hwnd, &ps)
		win32.SetViewportExtEx(hdc, cxClient/2, cyClient/2, nil)
		win32.SelectClipRgn(hdc, hRgnClip)

		fRadius := math.Hypot(float64(cxClient)/2.0, float64(cyClient)/2.0)
		const TWO_PI = 2.0 * math.Pi
		for fAngle := 0.0; fAngle < TWO_PI; fAngle += TWO_PI / 360.0 {
			win32.MoveToEx(hdc, 0, 0, nil)
			win32.LineTo(hdc, int32(fRadius*math.Cos(fAngle)+0.5),
				int32(-fRadius*math.Sin(fAngle)+0.5))
		}
		win32.EndPaint(hwnd, &ps)
		return 0

	case win32.WM_DESTROY:
		win32.PostQuitMessage(0)
		return 0
	}

	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}
