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
	appName := "Beeper2"
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
	hwnd, _ := win32.CreateWindow(appName, "Beeper2 Timer Demo",
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

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	switch msg {
	case win32.WM_CREATE:
		win32.SetTimer(hwnd, ID_TIMER, 1000, win32.NewTimerProc(timerProc))
		return 0

	case win32.WM_DESTROY:
		win32.KillTimer(hwnd, ID_TIMER)
		win32.PostQuitMessage(0)
		return 0
	}
	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}

var flipFlop bool

func timerProc(hwnd win32.HWND, msg uint32, timerID uintptr, time uintptr) uintptr {

	win32.MessageBeep(^uint32(0))
	flipFlop = !flipFlop
	var rect win32.RECT
	win32.GetClientRect(hwnd, &rect)
	hdc := win32.GetDC(hwnd)
	var hbrush win32.HBRUSH
	if flipFlop {
		hbrush = win32.CreateSolidBrush(win32.RGB(255, 0, 0))
	} else {
		hbrush = win32.CreateSolidBrush(win32.RGB(0, 0, 255))
	}
	win32.FillRect(hdc, &rect, hbrush)
	win32.ReleaseDC(hwnd, hdc)
	win32.DeleteObject(win32.HGDIOBJ(hbrush))
	return 0
}
