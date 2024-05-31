package main

import (
	"os"
	"runtime"
	"x/win32"
)

func main() {
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.
	appName := "Bezier"
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
	hwnd, _ := win32.CreateWindow(appName, "Bezier Splines",
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

func DrawBezier(hdc win32.HDC, pt []win32.POINT) {
	win32.PolyBezier(hdc, pt)
	win32.MoveToEx(hdc, pt[0].X, pt[0].Y, nil)
	win32.LineTo(hdc, pt[1].X, pt[1].Y)
	win32.MoveToEx(hdc, pt[2].X, pt[2].Y, nil)
	win32.LineTo(hdc, pt[3].X, pt[3].Y)
}

var cxClient, cyClient int32
var apt = [4]win32.POINT{}

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {

	switch msg {
	case win32.WM_SIZE:
		cxClient = win32.LOWORD(lParam)
		cyClient = win32.HIWORD(lParam)

		apt[0].X = cxClient / 4
		apt[0].Y = cyClient / 2

		apt[1].X = cxClient / 2
		apt[1].Y = cyClient / 4

		apt[2].X = cxClient / 2
		apt[2].Y = 3 * cyClient / 4

		apt[3].X = 3 * cxClient / 4
		apt[3].Y = cyClient / 2
		return 0

	case win32.WM_LBUTTONDOWN, win32.WM_RBUTTONDOWN, win32.WM_MOUSEMOVE:
		if wParam&win32.MK_LBUTTON != 0 || wParam&win32.MK_RBUTTON != 0 {
			hdc := win32.GetDC(hwnd)

			win32.SelectObject(hdc, win32.GetStockObject(win32.WHITE_PEN))
			DrawBezier(hdc, apt[:])

			if wParam&win32.MK_LBUTTON != 0 {
				apt[1].X = win32.LOWORD(lParam)
				apt[1].Y = win32.HIWORD(lParam)
			}
			if wParam&win32.MK_RBUTTON != 0 {
				apt[2].X = win32.LOWORD(lParam)
				apt[2].Y = win32.HIWORD(lParam)
			}

			win32.SelectObject(hdc, win32.GetStockObject(win32.BLACK_PEN))
			DrawBezier(hdc, apt[:])

			win32.ReleaseDC(hwnd, hdc)
		}

	case win32.WM_PAINT:
		win32.InvalidateRect(hwnd, nil, true)

		ps := win32.PAINTSTRUCT{}
		hdc := win32.BeginPaint(hwnd, &ps)
		defer win32.EndPaint(hwnd, &ps)

		DrawBezier(hdc, apt[:])
		return 0

	case win32.WM_DESTROY:
		win32.PostQuitMessage(0)
		return 0
	}

	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}
