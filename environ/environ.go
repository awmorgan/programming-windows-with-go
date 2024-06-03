package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"unsafe"
	"x/win32"
)

const (
	ID_LIST = 1
	ID_TEXT = 2
)

func main() {
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.
	appName := "Envrion"
	wc := win32.WNDCLASS{
		Style:         win32.CS_HREDRAW | win32.CS_VREDRAW,
		LpfnWndProc:   win32.NewWndProc(wndproc),
		HInstance:     win32.HInstance(),
		HIcon:         win32.LoadIcon(0, win32.IDI_APPLICATION),
		HCursor:       win32.LoadCursor(0, win32.IDC_ARROW),
		HbrBackground: win32.COLOR_WINDOW + 1,
		LpszClassName: win32.StringToUTF16Ptr(appName),
	}
	if _, err := win32.RegisterClass(&wc); err != nil {
		errMsg := fmt.Sprintf("RegisterClass failed: %v", err)
		win32.MessageBox(0, errMsg, appName, win32.MB_ICONERROR)
		return
	}
	hwnd, _ := win32.CreateWindow(appName, "Environment List Box",
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

func fillListBox(hwndList win32.HWND) {
	env := os.Environ()

	for _, v := range env {
		// split the name=value pair
		nameStr := v[:strings.Index(v, "=")]
		name := win32.StringToUTF16Ptr(nameStr)
		win32.SendMessage(hwndList, win32.LB_ADDSTRING, 0, uintptr(unsafe.Pointer(name)))
	}
}

var hwndList, hwndText win32.HWND

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	switch msg {
	case win32.WM_CREATE:
		cxChar := win32.LOWORD(uintptr(win32.GetDialogBaseUnits()))
		cyChar := win32.HIWORD(uintptr(win32.GetDialogBaseUnits()))

		// create list box and static text windows
		hwndList, _ = win32.CreateWindow("listbox", "",
			win32.WS_CHILD|win32.WS_VISIBLE|win32.LBS_STANDARD,
			cxChar, cyChar*3,
			cxChar*16+win32.GetSystemMetrics(win32.SM_CXVSCROLL),
			cyChar*5,
			hwnd, ID_LIST,
			win32.HInstance(), 0)
		hwndText, _ = win32.CreateWindow("static", "",
			win32.WS_CHILD|win32.WS_VISIBLE|win32.SS_LEFT,
			cxChar, cyChar,
			win32.GetSystemMetrics(win32.SM_CXSCREEN), cyChar,
			hwnd, ID_TEXT,
			win32.HInstance(), 0)
		fillListBox(hwndList)
		return 0

	case win32.WM_SETFOCUS:
		win32.SetFocus(hwndList)
		return 0

	case win32.WM_COMMAND:
		if win32.LOWORD(wParam) == ID_LIST && win32.HIWORD(wParam) == win32.LBN_SELCHANGE {
			index := win32.SendMessage(hwndList, win32.LB_GETCURSEL, 0, 0)
			len := win32.SendMessage(hwndList, win32.LB_GETTEXTLEN, index, 0)
			text := win32.StringToUTF16Ptr(strings.Repeat(" ", int(len)))
			win32.SendMessage(hwndList, win32.LB_GETTEXT, index, uintptr(unsafe.Pointer(text)))

			// get the environment string
			envStr := os.Getenv(win32.UTF16PtrToString(text))
			win32.SetWindowText(hwndText, envStr)
		}
		return 0

	case win32.WM_DESTROY:
		win32.PostQuitMessage(0)
		return 0
	}
	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}
