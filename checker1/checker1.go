package main

import (
	"fmt"
	"os"
	"runtime"
	"x/win32"
)

const divisions = 5

func main() {
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.
	appName := "Checker1"
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
	hwnd, _ := win32.CreateWindow(appName, "Cheker1 Mouse Hit-Test Demo",
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

var fstate [divisions][divisions]bool
var cxBlock, cyBlock int32

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	switch msg {
	case win32.WM_SIZE:
		cxBlock = win32.LOWORD(lParam) / divisions
		cyBlock = win32.HIWORD(lParam) / divisions
		return 0

	case win32.WM_LBUTTONDOWN:
		x := win32.LOWORD(lParam) / cxBlock
		y := win32.HIWORD(lParam) / cyBlock

		if x < divisions && y < divisions {
			fstate[x][y] = !fstate[x][y]
			r := win32.RECT{
				Left:   x * cxBlock,
				Top:    y * cyBlock,
				Right:  (x + 1) * cxBlock,
				Bottom: (y + 1) * cyBlock,
			}
			win32.InvalidateRect(hwnd, &r, false)
		} else {
			win32.MessageBeep(0)
		}
		return 0

	case win32.WM_PAINT:
		ps := win32.PAINTSTRUCT{}
		win32.BeginPaint(hwnd, &ps)
		for x := int32(0); x < divisions; x++ {
			for y := int32(0); y < divisions; y++ {
				win32.Rectangle(ps.Hdc, x*cxBlock, y*cyBlock,
					(x+1)*cxBlock, (y+1)*cyBlock)
				if fstate[x][y] {
					win32.MoveToEx(ps.Hdc, x*cxBlock, y*cyBlock, nil)
					win32.LineTo(ps.Hdc, (x+1)*cxBlock, (y+1)*cyBlock)
					win32.MoveToEx(ps.Hdc, x*cxBlock, (y+1)*cyBlock, nil)
					win32.LineTo(ps.Hdc, (x+1)*cxBlock, y*cyBlock)
				}
			}
		}
		win32.EndPaint(hwnd, &ps)
		return 0

	case win32.WM_DESTROY:
		win32.PostQuitMessage(0)
		return 0
	}

	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}
