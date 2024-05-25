package main

import (
	"fmt"
	"os"
	"runtime"
	"x/win32"
)

const divisions = 5

var childClassName = "Checker3_Child"

func main() {
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.
	appName := "Checker3"
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
	wc.LpfnWndProc = win32.NewWndProc(childwndproc)
	wc.CbWndExtra = 8 // sizeof(uintptr)
	wc.HIcon = 0
	wc.LpszClassName = win32.StringToUTF16Ptr(childClassName)
	win32.RegisterClass(&wc)

	hwnd, _ := win32.CreateWindow(appName, "Cheker3 Mouse Hit-Test Demo",
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

var hwndChild [divisions][divisions]win32.HWND

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	var cxBlock, cyBlock int32
	switch msg {
	case win32.WM_CREATE:
		for x := int32(0); x < divisions; x++ {
			for y := int32(0); y < divisions; y++ {
				p, _ := win32.GetWindowLongPtr(hwnd, win32.GWLP_HINSTANCE)
				hwndChild[x][y], _ = win32.CreateWindow(childClassName, "",
					win32.WS_CHILD|win32.WS_VISIBLE,
					0, 0, 0, 0,
					hwnd, win32.HMENU(y<<8|x),
					win32.HINSTANCE(p),
					0)
			}
		}
		return 0

	case win32.WM_SIZE:
		cxBlock = win32.LOWORD(lParam) / divisions
		cyBlock = win32.HIWORD(lParam) / divisions
		for x := int32(0); x < divisions; x++ {
			for y := int32(0); y < divisions; y++ {
				win32.MoveWindow(hwndChild[x][y],
					x*cxBlock, y*cyBlock,
					cxBlock, cyBlock, true)
			}
		}
		return 0

	case win32.WM_LBUTTONDOWN:
		win32.MessageBeep(0)
		return 0

	case win32.WM_DESTROY:
		win32.PostQuitMessage(0)
		return 0
	}

	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}

func childwndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	switch msg {
	case win32.WM_CREATE:
		win32.SetWindowLongPtr(hwnd, 0, 0)
		return 0

	case win32.WM_LBUTTONDOWN:
		p, _ := win32.GetWindowLongPtr(hwnd, 0)
		p ^= 1
		win32.SetWindowLongPtr(hwnd, 0, p)
		win32.InvalidateRect(hwnd, nil, false)
		return 0

	case win32.WM_PAINT:
		ps := win32.PAINTSTRUCT{}
		win32.BeginPaint(hwnd, &ps)

		var r win32.RECT
		win32.GetClientRect(hwnd, &r)
		win32.Rectangle(ps.Hdc, 0, 0, r.Right, r.Bottom)

		p, _ := win32.GetWindowLongPtr(hwnd, 0)
		if p != 0 {
			win32.MoveToEx(ps.Hdc, 0, 0, nil)
			win32.LineTo(ps.Hdc, r.Right, r.Bottom)
			win32.MoveToEx(ps.Hdc, 0, r.Bottom, nil)
			win32.LineTo(ps.Hdc, r.Right, 0)
		}
		win32.EndPaint(hwnd, &ps)
		return 0
	}
	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}
