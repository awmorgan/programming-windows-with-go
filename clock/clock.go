package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"x/win32"
)

const (
	ID_TIMER = 1
	TWOPI    = math.Pi * 2
)

func main() {
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.
	appName := "Clock"
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
	hwnd, _ := win32.CreateWindow(appName, "Analog Clock",
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

func SetIsotropic(hdc win32.HDC, cx, cy int32) {
	win32.SetMapMode(hdc, win32.MM_ISOTROPIC)
	win32.SetWindowExtEx(hdc, 1000, 1000, nil)
	win32.SetViewportExtEx(hdc, cx/2, -cy/2, nil)
	win32.SetViewportOrgEx(hdc, cx/2, cy/2, nil)
}

func RotatePoint(pt []win32.POINT, n, angle int32) {
	var temp win32.POINT
	for i := int32(0); i < n; i++ {
		temp.X = int32(float64(pt[i].X)*math.Cos(TWOPI*float64(angle)/360) +
			float64(pt[i].Y)*math.Sin(TWOPI*float64(angle)/360))

		temp.Y = int32(float64(pt[i].Y)*math.Cos(TWOPI*float64(angle)/360) -
			float64(pt[i].X)*math.Sin(TWOPI*float64(angle)/360))

		pt[i] = temp
	}
}

func DrawClock(hdc win32.HDC) {
	var pt [3]win32.POINT
	for angle := int32(0); angle < 360; angle += 6 {
		pt[0].X = 0
		pt[0].Y = 900

		RotatePoint(pt[:], 1, angle)
		if angle%5 != 0 {
			pt[2].X, pt[2].Y = 33, 33
		} else {
			pt[2].X, pt[2].Y = 100, 100
		}

		pt[0].X -= pt[2].X / 2
		pt[0].Y -= pt[2].Y / 2

		pt[1].X = pt[0].X + pt[2].X
		pt[1].Y = pt[0].Y + pt[2].Y

		win32.SelectObject(hdc, win32.GetStockObject(win32.BLACK_BRUSH))

		win32.Ellipse(hdc, pt[0].X, pt[0].Y, pt[1].X, pt[1].Y)
	}
}

func DrawHands(hdc win32.HDC, st *win32.SYSTEMTIME, change bool) {
	var pts = [3][5]win32.POINT{
		{{X: 0, Y: -150}, {X: 100, Y: 0}, {X: 0, Y: 600}, {X: -100, Y: 0}, {X: 0, Y: -150}},
		{{X: 0, Y: -200}, {X: 50, Y: 0}, {X: 0, Y: 800}, {X: -50, Y: 0}, {X: 0, Y: -200}},
		{{X: 0, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 800}},
	}

	var angle [3]int32
	var ptTemp [3][5]win32.POINT

	angle[0] = int32((st.WHour*30)%360 + st.WMinute/2)
	angle[1] = int32(st.WMinute * 6)
	angle[2] = int32(st.WSecond * 6)

	copy(ptTemp[:], pts[:])

	if change {
		for i := 0; i < 3; i++ {
			RotatePoint(ptTemp[i][:], 5, angle[i])
			win32.Polyline(hdc, ptTemp[i][:])
		}
	} else {
		for i := 2; i < 3; i++ {
			RotatePoint(ptTemp[i][:], 5, angle[i])
			win32.Polyline(hdc, ptTemp[i][:])
		}
	}
}

var cxClient, cyClient int32
var stPrev win32.SYSTEMTIME

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	switch msg {
	case win32.WM_CREATE:
		win32.SetTimer(hwnd, ID_TIMER, 1000, 0)
		win32.GetLocalTime(&stPrev)
		return 0

	case win32.WM_SIZE:
		cxClient = win32.LOWORD(lParam)
		cyClient = win32.HIWORD(lParam)
		return 0

	case win32.WM_TIMER:
		var st win32.SYSTEMTIME
		win32.GetLocalTime(&st)
		change := st.WHour != stPrev.WHour || st.WMinute != stPrev.WMinute
		hdc := win32.GetDC(hwnd)
		SetIsotropic(hdc, cxClient, cyClient)
		win32.SelectObject(hdc, win32.GetStockObject(win32.WHITE_PEN))
		DrawHands(hdc, &stPrev, change)
		win32.SelectObject(hdc, win32.GetStockObject(win32.BLACK_PEN))
		DrawHands(hdc, &st, true)
		win32.ReleaseDC(hwnd, hdc)
		stPrev = st
		return 0

	case win32.WM_PAINT:
		var ps win32.PAINTSTRUCT
		hdc := win32.BeginPaint(hwnd, &ps)
		SetIsotropic(hdc, cxClient, cyClient)
		DrawClock(hdc)
		DrawHands(hdc, &stPrev, true)
		win32.EndPaint(hwnd, &ps)
		return 0

	case win32.WM_DESTROY:
		win32.KillTimer(hwnd, ID_TIMER)
		win32.PostQuitMessage(0)
		return 0
	}
	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}
