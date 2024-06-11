package main

import (
	"fmt"
	"os"
	"runtime"
	"x/win32"
)

// make the resources with visual studio 2022 community edition
// compile the rc file with visual studio
// then run llvm windres --output res.syso res.res
// where res.res is the compiled rc file

var hIcon win32.HICON
const (
	IDS_APPNAME = 1
	IDS_CAPTION = 2
	IDS_ERRMSG  = 3
)
func main() {
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.

	var appName [16]uint16
	var caption [64]uint16
	var errMsg [64]uint16

	win32.LoadString(win32.HInstance(), IDS_APPNAME, appName[:])
	win32.LoadString(win32.HInstance(), IDS_CAPTION, caption[:])


	hIcon = win32.LoadIconFromString(win32.HInstance(), appName)
	wc := win32.WNDCLASS{
		Style:         win32.CS_HREDRAW | win32.CS_VREDRAW,
		LpfnWndProc:   win32.NewWndProc(wndproc),
		HInstance:     win32.HInstance(),
		HIcon:         hIcon,
		HCursor:       win32.LoadCursor(0, win32.IDC_ARROW),
		HbrBackground: win32.HBRUSH(win32.GetStockObject(win32.WHITE_BRUSH)),
		LpszClassName: win32.StringToUTF16Ptr(appName),
	}
	if _, err := win32.RegisterClass(&wc); err != nil {
		errMsg := fmt.Sprintf("RegisterClass failed: %v", err)
		win32.MessageBox(0, errMsg, appName, win32.MB_ICONERROR)
		return
	}
	hwnd, _ := win32.CreateWindow(appName, "Icon Demo",
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

var cxIcon, cyIcon, cxClient, cyClient int32

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	switch msg {
	case win32.WM_CREATE:
		cxIcon = win32.GetSystemMetrics(win32.SM_CXICON)
		cyIcon = win32.GetSystemMetrics(win32.SM_CYICON)
		return 0

	case win32.WM_SIZE:
		cxClient = win32.LOWORD(lParam)
		cyClient = win32.HIWORD(lParam)
		return 0

	case win32.WM_PAINT:
		var ps win32.PAINTSTRUCT
		hdc := win32.BeginPaint(hwnd, &ps)
		for y := int32(0); y < cyClient; y += cyIcon {
			for x := int32(0); x < cxClient; x += cxIcon {
				err := win32.DrawIcon(hdc, x, y, hIcon)
				if err != nil {
					errMsg := fmt.Sprintf("DrawIcon failed: %v", err)
					win32.MessageBox(0, errMsg, "Icon Demo", win32.MB_ICONERROR)
					return 0
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
