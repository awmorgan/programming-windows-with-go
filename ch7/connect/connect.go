package main

import (
	"fmt"
	"os"
	"runtime"
	"x/win32"
)

func main() {
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.
	appName := "Connect"
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
	hwnd, _ := win32.CreateWindow(appName, "Connect-the-Points Mouse Demo",
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

const maxpoints = 1000

var (
	pt     [maxpoints]win32.POINT
	nCount int
)

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	switch msg {
	case win32.WM_LBUTTONDOWN:
		nCount = 0
		win32.InvalidateRect(hwnd, nil, true)
		return 0

	case win32.WM_MOUSEMOVE:
		if wParam&win32.MK_LBUTTON == win32.MK_LBUTTON && nCount < maxpoints {
			lw := win32.LOWORD(lParam)
			hw := win32.HIWORD(lParam)
			pt[nCount].X = lw
			pt[nCount].Y = hw
			nCount++
			hdc := win32.GetDC(hwnd)
			win32.SetPixel(hdc, lw, hw, 0)
			win32.ReleaseDC(hwnd, hdc)
		}
		return 0

	case win32.WM_LBUTTONUP:
		win32.InvalidateRect(hwnd, nil, false)

	case win32.WM_PAINT:
		ps := win32.PAINTSTRUCT{}
		win32.BeginPaint(hwnd, &ps)

		win32.SetCursor(win32.WaitCursor())
		win32.ShowCursor(true)

		for i := 0; i < nCount-1; i++ {
			for j := i + 1; j < nCount; j++ {
				win32.MoveToEx(ps.Hdc, pt[i].X, pt[i].Y, nil)
				win32.LineTo(ps.Hdc, pt[j].X, pt[j].Y)
			}
		}

		win32.ShowCursor(false)
		win32.SetCursor(win32.ArrowCursor())

		win32.EndPaint(hwnd, &ps)
		return 0

	case win32.WM_DESTROY:
		win32.PostQuitMessage(0)
		return 0
	}

	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}
