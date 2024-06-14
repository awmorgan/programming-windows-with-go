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

	var appNameArray [16]uint16
	var captionArray [64]uint16
	var errMsgArray [64]uint16

	win32.LoadString(win32.HInstance(), IDS_APPNAME, appNameArray[:])
	win32.LoadString(win32.HInstance(), IDS_CAPTION, captionArray[:])
	win32.LoadString(win32.HInstance(), IDS_ERRMSG, errMsgArray[:])

	appName := win32.UTF16PtrToString(&appNameArray[0])
	caption := win32.UTF16PtrToString(&captionArray[0])

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

var (
	pText string
	hResource win32.HGLOBAL
	hScroll win32.HWND
	iPosition, cxChar, cyChar, cyClient, iNumLines, xScroll int32
)

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	switch msg {
	case win32.WM_CREATE:
		hdc := win32.GetDC(hwnd)
		var tm win32.TEXTMETRIC
		win32.GetTextMetrics(hdc, &tm)
		cxChar = tm.TmAveCharWidth
		cyChar = tm.TmHeight + tm.TmExternalLeading
		win32.ReleaseDC(hwnd, hdc)

		xScroll = win32.GetSystemMetrics(win32.SM_CXVSCROLL)

		hScroll, _ = win32.CreateWindow("scrollbar", "",
			win32.WS_CHILD | win32.WS_VISIBLE | win32.SBS_VERT,
			0, 0, 0, 0,
			hwnd, win32.HMENU(1), win32.HInstance(), 0)

		res, _ := win32.FindResource(win32.HInstance(), "Annabellee", "TEXT")
		hResource, _ = win32.LoadResource( win32.HInstance(), res)
		pText = 

	}
	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}
