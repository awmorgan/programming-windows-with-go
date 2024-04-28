package main

import (
	"fmt"
	"x/win32"
)

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	return 0
}

func main() {

	wc := WNDCLASS{
		Style:         CS_HREDRAW | CS_VREDRAW,
		LpfnWndProc:   NewWndProc(wndproc),
		CbClsExtra:    0,
		CbWndExtra:    0,
		HInstance:     HInstance,
		HIcon:         0,
		HCursor:       0,
		HbrBackground: 0,
		LpszMenuName:  new(uint16),
		LpszClassName: new(uint16),
	}
	fmt.Printf("%#v\n", wc)
	// 	wc.HInstance = win.WinmainArgs.HInstance
	// 	wc.HIcon, _ = win.LoadIcon(0, win.IDI_APPLICATION)
	// 	wc.HCursor, _ = win.LoadCursor(0, win.IDC_ARROW)
	// 	wc.HbrBackground = win.HBRUSH(win.GetStockObject(win.WHITE_BRUSH))
	// 	wc.LpszMenuName = nil
	// 	wc.LpszClassName = win.StringToUTF16Ptr("HelloWin")
	// 	if _, err := win.RegisterClass(&wc); err != nil {
	// 		win.MessageBox(0, fmt.Sprintf("RegisterClass failed: %v", err), "Error", win.MB_ICONERROR)
	// 		os.Exit(1)
	// 	}
	// 	hwnd, err := win.CreateWindowEx(0, "HelloWin",
	// 		"Hello, Windows App",
	// 		win.WS_OVERLAPPEDWINDOW,
	// 		win.CW_USEDEFAULT,
	// 		win.CW_USEDEFAULT,
	// 		win.CW_USEDEFAULT,
	// 		win.CW_USEDEFAULT,
	// 		0, 0, win.WinmainArgs.HInstance, 0)
	// 	if err != nil {
	// 		win.MessageBox(0, fmt.Sprintf("CreateWindow failed: %v", err), "Error", win.MB_ICONERROR)
	// 		os.Exit(1)
	// 	}
	// 	win.ShowWindow(hwnd, win.SW_SHOWNORMAL)
	// 	win.UpdateWindow(hwnd)
	// 	msg := win.MSG{}
	// 	for {
	// 		if ret, _ := win.GetMessage(&msg, 0, 0, 0); ret == 0 {
	// 			break
	// 		}

	// 		win.TranslateMessage(&msg)
	// 		win.DispatchMessage(&msg)
	// 	}
	// 	os.Exit(int(msg.WParam))
}
