package main

import (
	"fmt"
	"os"
	"runtime"
	"unsafe"
	"x/win32"
)

var appName = "PopPad1"

func main() {
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.
	wc := win32.WNDCLASS{
		Style:         win32.CS_HREDRAW | win32.CS_VREDRAW,
		LpfnWndProc:   win32.NewWndProc(wndproc),
		HInstance:     win32.HInstance(),
		HIcon:         win32.ApplicationIcon(),
		HCursor:       win32.ArrowCursor(),
		HbrBackground: win32.CreateSolidBrush(0),
		LpszClassName: win32.StringToUTF16Ptr(appName),
	}
	if _, err := win32.RegisterClass(&wc); err != nil {
		errMsg := fmt.Sprintf("RegisterClass failed: %v", err)
		win32.MessageBox(0, errMsg, appName, win32.MB_ICONERROR)
		return
	}
	hwnd, _ := win32.CreateWindow(appName, appName,
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

var hwndEdit win32.HWND

const ID_EDIT = 1

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	switch msg {
	case win32.WM_CREATE:
		createStruct := (*win32.CREATESTRUCT)(unsafe.Pointer(lParam))
		hwndEdit, _ = win32.CreateWindow("edit", "",
			win32.WS_CHILD|win32.WS_VISIBLE|win32.WS_HSCROLL|
				win32.WS_VSCROLL|win32.WS_BORDER|win32.ES_LEFT|
				win32.ES_MULTILINE|win32.ES_AUTOHSCROLL|win32.ES_AUTOVSCROLL,
			0, 0, 0, 0,
			hwnd, ID_EDIT, createStruct.Instance, 0)
		return 0

	case win32.WM_SETFOCUS:
		win32.SetFocus(hwndEdit)
		return 0

	case win32.WM_SIZE:
		win32.MoveWindow(hwndEdit, 0, 0, win32.LOWORD(lParam), win32.HIWORD(lParam), true)
		return 0

	case win32.WM_COMMAND:
		if win32.LOWORD(wParam) == ID_EDIT {
			if win32.HIWORD(wParam) == win32.EN_ERRSPACE || win32.HIWORD(wParam) == win32.EN_MAXTEXT {
				win32.MessageBox(hwnd, "Edit control out of space.", appName, win32.MB_OK|win32.MB_ICONSTOP)
			}
		}
		return 0

	case win32.WM_DESTROY:
		win32.PostQuitMessage(0)
		return 0
	}
	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}
