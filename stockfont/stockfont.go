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
			iFont = win32.HIWORD(wParam)
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
		var faceName [win32.LF_FACESIZE]uint16
		win32.GetTextFace(hdc, win32.LF_FACESIZE, &faceName[0])
		var tm win32.TEXTMETRIC
		win32.GetTextMetrics(hdc, &tm)
		cxGrid := max(3*tm.TmAveCharWidth, 2*tm.TmMaxCharWidth)
		cyGrid := tm.TmHeight + 3
		s := fmt.Sprintf("%s: Face Name: %s, CharSet = %d",
			stockfont[iFont].szStockFont, win32.Utf16PtrToString(&faceName[0]), tm.TmCharSet)
		win32.TextOut(hdc, 0, 0, s, len(s))
		win32.SetTextAlign(hdc, win32.TA_TOP|win32.TA_CENTER)

		// vertical and horizontal lines
		for i := int32(0); i < 17; i++ {
			win32.MoveToEx(hdc, (i+2)*cxGrid, 2*cyGrid, nil)
			win32.LineTo(hdc, (i+2)*cxGrid, 19*cyGrid)

			win32.MoveToEx(hdc, cxGrid, (i+3)*cyGrid, nil)
			win32.LineTo(hdc, 18*cxGrid, (i+3)*cyGrid)
		}

		// vertical and horizontal headings
		for i := int32(0); i < 16; i++ {
			s := fmt.Sprintf("%X-", i)
			win32.TextOut(hdc, (2*i+5)*cxGrid/2, 2*cyGrid+2, s, len(s))
			s = fmt.Sprintf("-%X", i)
			win32.TextOut(hdc, 3*cxGrid/2, (i+3)*cyGrid+2, s, len(s))
		}

		// characters
		for y := int32(0); y < 16; y++ {
			for x := int32(0); x < 16; x++ {
				ch := byte(x*16 + y)
				s := string(ch)
				win32.TextOut(hdc, (2*x+5)*cxGrid/2,
					(y+3)*cyGrid+2, s, len(s))
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
