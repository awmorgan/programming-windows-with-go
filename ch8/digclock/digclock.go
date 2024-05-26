package main

import (
	"fmt"
	"os"
	"runtime"
	"x/win32"
)

const ID_TIMER = 1

func main() {
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.
	appName := "DigClock"
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
	hwnd, _ := win32.CreateWindow(appName, "Digital Clock",
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

func displayDigit(hdc win32.HDC, number int) {
	sevenSegment := [10][7]int{
		{1, 1, 1, 1, 1, 1, 0}, // 0
		{0, 1, 1, 0, 0, 0, 0}, // 1
		{1, 1, 0, 1, 1, 0, 1}, // 2
		{1, 1, 1, 1, 0, 0, 1}, // 3
		{0, 1, 1, 0, 0, 1, 1}, // 4
		{1, 0, 1, 1, 0, 1, 1}, // 5
		{1, 0, 1, 1, 1, 1, 1}, // 6
		{1, 1, 1, 0, 0, 0, 0}, // 7
		{1, 1, 1, 1, 1, 1, 1}, // 8
		{1, 1, 1, 1, 0, 1, 1}, // 9
	}
	segments := [7][6]win32.POINT{
		{{X: 7, Y: 6}, {X: 11, Y: 2}, {X: 31, Y: 2}, {X: 35, Y: 6}, {X: 31, Y: 10}, {X: 11, Y: 10}},
		{{X: 6, Y: 7}, {X: 10, Y: 11}, {X: 10, Y: 31}, {X: 6, Y: 35}, {X: 2, Y: 31}, {X: 2, Y: 11}},
		{{X: 36, Y: 7}, {X: 40, Y: 11}, {X: 40, Y: 31}, {X: 36, Y: 35}, {X: 32, Y: 31}, {X: 32, Y: 11}},
		{{X: 7, Y: 36}, {X: 11, Y: 32}, {X: 31, Y: 32}, {X: 35, Y: 36}, {X: 31, Y: 40}, {X: 11, Y: 40}},
		{{X: 6, Y: 37}, {X: 10, Y: 41}, {X: 10, Y: 61}, {X: 6, Y: 65}, {X: 2, Y: 61}, {X: 2, Y: 41}},
		{{X: 36, Y: 37}, {X: 40, Y: 41}, {X: 40, Y: 61}, {X: 36, Y: 65}, {X: 32, Y: 61}, {X: 32, Y: 41}},
		{{X: 7, Y: 66}, {X: 11, Y: 62}, {X: 31, Y: 62}, {X: 35, Y: 66}, {X: 31, Y: 70}, {X: 11, Y: 70}},
	}
	for seg := 0; seg < 7; seg++ {
		if sevenSegment[number][seg] == 1 {
			win32.Polygon(hdc, segments[seg][:])
		}
	}
}

func displayTwoDigits(hdc win32.HDC, number int, suppress bool) {
	if !suppress || number/10 > 0 {
		displayDigit(hdc, number/10)
	}
	win32.OffsetWindowOrgEx(hdc, -42, 0, nil)
	displayDigit(hdc, number%10)
	win32.OffsetWindowOrgEx(hdc, -42, 0, nil)
}

func displayColon(hdc win32.HDC) {
	colon := [2][4]win32.POINT{
		{{X: 2, Y: 21}, {X: 6, Y: 17}, {X: 10, Y: 21}, {X: 6, Y: 25}},
		{{X: 2, Y: 51}, {X: 6, Y: 47}, {X: 10, Y: 51}, {X: 6, Y: 55}},
	}
	win32.Polygon(hdc, colon[0][:])
	win32.Polygon(hdc, colon[1][:])
	win32.OffsetWindowOrgEx(hdc, -12, 0, nil)
}

func displayTime(hdc win32.HDC, _24h bool, suppress bool) {
	var t win32.SYSTEMTIME
	win32.GetLocalTime(&t)
	if _24h {
		displayTwoDigits(hdc, int(t.WHour), suppress)
	} else {
		t.WHour %= 12
		if t.WHour != 0 {
			displayTwoDigits(hdc, int(t.WHour), suppress)
		} else {
			displayTwoDigits(hdc, 12, suppress)
		}
	}
	displayColon(hdc)
	displayTwoDigits(hdc, int(t.WMinute), false)
	displayColon(hdc)
	displayTwoDigits(hdc, int(t.WSecond), false)
}

var _24hour, suppress bool
var hbrushRed win32.HBRUSH
var cxClient, cyClient int32

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	switch msg {
	case win32.WM_CREATE:
		hbrushRed = win32.CreateSolidBrush(win32.RGB(255, 0, 0))
		win32.SetTimer(hwnd, ID_TIMER, 1000, 0)
		fallthrough
	case win32.WM_SETTINGCHANGE:
		var buf [2]uint16
		win32.GetLocaleInfo(win32.LOCALE_USER_DEFAULT, win32.LOCALE_ITIME, &buf[0], 2)
		if buf[0] == '1' {
			_24hour = true
		} else {
			_24hour = false
		}
		win32.GetLocaleInfo(win32.LOCALE_USER_DEFAULT, win32.LOCALE_ITLZERO, &buf[0], 2)
		if buf[0] == '0' {
			suppress = true
		} else {
			suppress = false
		}
		win32.InvalidateRect(hwnd, nil, true)
		return 0

	case win32.WM_SIZE:
		cxClient = win32.LOWORD(lParam)
		cyClient = win32.HIWORD(lParam)
		return 0

	case win32.WM_TIMER:
		win32.InvalidateRect(hwnd, nil, true)
		return 0

	case win32.WM_PAINT:
		var ps win32.PAINTSTRUCT
		hdc := win32.BeginPaint(hwnd, &ps)
		win32.SetMapMode(hdc, win32.MM_ISOTROPIC)
		win32.SetWindowExtEx(hdc, 276, 72, nil)
		win32.SetViewportExtEx(hdc, cxClient, cyClient, nil)

		win32.SetWindowOrgEx(hdc, 138, 36, nil)
		win32.SetViewportOrgEx(hdc, cxClient/2, cyClient/2, nil)
		win32.SelectObject(hdc, win32.GetStockObject(win32.NULL_PEN))
		win32.SelectObject(hdc, win32.HGDIOBJ(hbrushRed))

		displayTime(hdc, _24hour, suppress)
		win32.EndPaint(hwnd, &ps)
		return 0

	case win32.WM_DESTROY:
		win32.KillTimer(hwnd, ID_TIMER)
		win32.PostQuitMessage(0)
		return 0
	}
	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}
