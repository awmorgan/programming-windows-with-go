package main

import (
	"fmt"
	"os"
	"runtime"
	"x/win32"
)

func main() {
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.
	appName := "StokFont"
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
	hwnd, _ := win32.CreateWindow(appName, "Stock Fonts",
		win32.WS_OVERLAPPEDWINDOW|win32.WS_VSCROLL|win32.WS_HSCROLL,
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

var stockfont = []struct {
	idStockFont int32
	szStockFont string
}{
	{win32.OEM_FIXED_FONT, "OEM_FIXED_FONT"},
	{win32.ANSI_FIXED_FONT, "ANSI_FIXED_FONT"},
	{win32.ANSI_VAR_FONT, "ANSI_VAR_FONT"},
	{win32.SYSTEM_FONT, "SYSTEM_FONT"},
	{win32.DEVICE_DEFAULT_FONT, "DEVICE_DEFAULT_FONT"},
	{win32.SYSTEM_FIXED_FONT, "SYSTEM_FIXED_FONT"},
	{win32.DEFAULT_GUI_FONT, "DEFAULT_GUI_FONT"},
}
var (
	iFont  int32
	cFonts = int32(len(stockfont))
)

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	switch msg {
	case win32.WM_CREATE:
		win32.SetScrollRange(hwnd, win32.SB_VERT, 0, cFonts-1, true)
		return 0

	case win32.WM_DISPLAYCHANGE:
		win32.InvalidateRect(hwnd, nil, true)
		return 0

	case win32.WM_VSCROLL:
		switch win32.LOWORD(wParam) {
		case win32.SB_TOP:
			iFont = 0
		case win32.SB_BOTTOM:
			iFont = cFonts - 1
		case win32.SB_LINEUP, win32.SB_PAGEUP:
			iFont--
		case win32.SB_LINEDOWN, win32.SB_PAGEDOWN:
			iFont++
		case win32.SB_THUMBPOSITION:
			iFont = int32(win32.HIWORD(wParam))
		}
		iFont = max(0, min(iFont, cFonts-1))
		win32.SetScrollPos(hwnd, win32.SB_VERT, iFont, true)
		win32.InvalidateRect(hwnd, nil, true)
		return 0

	case win32.WM_KEYDOWN:
		switch wParam {
		case win32.VK_HOME:
			win32.SendMessage(hwnd, win32.WM_VSCROLL, win32.SB_TOP, 0)
		case win32.VK_END:
			win32.SendMessage(hwnd, win32.WM_VSCROLL, win32.SB_BOTTOM, 0)
		case win32.VK_PRIOR, win32.VK_LEFT, win32.VK_UP:
			win32.SendMessage(hwnd, win32.WM_VSCROLL, win32.SB_LINEUP, 0)
		case win32.VK_NEXT, win32.VK_RIGHT, win32.VK_DOWN:
			win32.SendMessage(hwnd, win32.WM_VSCROLL, win32.SB_PAGEDOWN, 0)
		}
		return 0

	case win32.WM_PAINT:
		ps := win32.PAINTSTRUCT{}
		hdc := win32.BeginPaint(hwnd, &ps)

		win32.SelectObject(hdc, win32.GetStockObject(stockfont[iFont].idStockFont))
		// win32.GetTextFace
		win32.EndPaint(hwnd, &ps)
		return 0

	case win32.WM_DESTROY:
		win32.PostQuitMessage(0)
		return 0
	}

	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}
