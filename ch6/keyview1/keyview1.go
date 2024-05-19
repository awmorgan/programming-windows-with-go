package main

import (
	"fmt"
	"os"
	"runtime"
	"x/win32"
)

func main() {
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.
	appName := "KeyView1"
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
	hwnd, _ := win32.CreateWindow(appName, "Keyboard Message Viewer #1",
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

var (
	cyClientMax        int32
	cxClient, cyClient int32
	cyChar             int32
	cLinesMax, cLines  int32
	pmsg               []win32.MSG
	rectScroll         win32.RECT
)

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	switch msg {
	case win32.WM_CREATE, win32.WM_DISPLAYCHANGE:
		// Get maximum size of client area.
		// cxClientMax = win32.GetSystemMetrics(win32.SM_CXMAXIMIZED)
		cyClientMax = win32.GetSystemMetrics(win32.SM_CYMAXIMIZED)

		// Get character size for the fixed-pitch font.
		hdc := win32.GetDC(hwnd)
		win32.SelectObject(hdc, win32.GetStockObject(win32.SYSTEM_FIXED_FONT))
		textMetric := win32.TEXTMETRIC{}
		win32.GetTextMetrics(hdc, &textMetric)
		// cxChar = textMetric.TmAveCharWidth
		cyChar = textMetric.TmHeight

		win32.ReleaseDC(hwnd, hdc)

		// Allocate memory for display lines.
		cLinesMax = cyClientMax / cyChar
		pmsg = make([]win32.MSG, cLinesMax)
		cLines = 0
		fallthrough

	case win32.WM_SIZE:
		if msg == win32.WM_SIZE {
			cxClient = int32(win32.LOWORD(lParam))
			cyClient = int32(win32.HIWORD(lParam))
		}

		// Calculate scrolling rectangle.
		rectScroll.Left = 0
		rectScroll.Right = cxClient
		rectScroll.Top = cyChar
		rectScroll.Bottom = cyChar * (cyClient / cyChar)

		win32.InvalidateRect(hwnd, nil, true)
		return 0

	case win32.WM_KEYDOWN, win32.WM_KEYUP, win32.WM_CHAR, win32.WM_DEADCHAR,
		win32.WM_SYSKEYDOWN, win32.WM_SYSKEYUP, win32.WM_SYSCHAR, win32.WM_SYSDEADCHAR:
		for i := cLinesMax - 1; i > 0; i-- {
			pmsg[i] = pmsg[i-1]
		}
		pmsg[0].HWnd = hwnd
		pmsg[0].Message = msg
		pmsg[0].WParam = wParam
		pmsg[0].LParam = lParam

		cLines = min(cLines+1, cLinesMax)

		// Scroll up the display.
		win32.ScrollWindow(hwnd, 0, -cyChar, nil, nil)

	case win32.WM_PAINT:
		ps := win32.PAINTSTRUCT{}
		hdc := win32.BeginPaint(hwnd, &ps)

		win32.SelectObject(hdc, win32.GetStockObject(win32.SYSTEM_FIXED_FONT))
		win32.SetBkMode(hdc, win32.TRANSPARENT)
		top := "Message        Key       Char     Repeat Scan Ext ALT Prev Tran"
		und := "_______        ___       ____     ______ ____ ___ ___ ____ ____"
		win32.TextOut(hdc, 0, 0, top, len(top))
		win32.TextOut(hdc, 0, 0, und, len(und))

		var fmtStr = [...]string{
			"%-13s %3d %-15s%c%6d %4d %3d %3d %4d %4d",
			"%-13s            0x%04X%1s%c %6d %4d %3d %3d %4s %4s",
		}

		var messages = [...]string{
			"WM_KEYDOWN", "WM_KEYUP", "WM_CHAR", "WM_DEADCHAR",
			"WM_SYSKEYDOWN", "WM_SYSKEYUP", "WM_SYSCHAR", "WM_SYSDEADCHAR",
		}

		yesNo := func(flag uint32) string {
			if flag != 0 {
				return "Yes"
			}
			return "No"
		}

		upDown := func(flag uint32) string {
			if flag != 0 {
				return "Up"
			}
			return "Down"
		}

		for i := int32(0); i < min(cLines, cyClient/cyChar-1); i++ {
			var szKeyName [32]uint16
			iType := 0
			if pmsg[i].Message == win32.WM_CHAR || pmsg[i].Message == win32.WM_SYSCHAR ||
				pmsg[i].Message == win32.WM_DEADCHAR || pmsg[i].Message == win32.WM_SYSDEADCHAR {
				iType = 1
			}

			win32.GetKeyNameText(pmsg[i].LParam, &szKeyName[0], int32(len(szKeyName)))
			s := fmt.Sprintf(fmtStr[iType], messages[pmsg[i].Message-win32.WM_KEYDOWN],
				pmsg[i].WParam, func() string {
					if iType == 1 {
						return " "
					} else {
						return win32.Utf16PtrToString(&szKeyName[0])
					}
				}(),
				func() uint16 {
					if iType == 1 {
						return uint16(pmsg[i].WParam)
					} else {
						return ' '
					}
				}(),
				win32.LOWORD(pmsg[i].LParam), win32.HIWORD(pmsg[i].LParam),
				yesNo(uint32(pmsg[i].LParam)&0x0100_0000),
				yesNo(uint32(pmsg[i].LParam)&0x2000_0000),
				upDown(uint32(pmsg[i].LParam)&0x4000_0000),
				upDown(uint32(pmsg[i].LParam)&0x8000_0000))

			win32.TextOut(hdc, 0, (cyClient/cyChar-1-i)*cyChar, s, len(s))
			fmt.Println(s)
		}

		win32.EndPaint(hwnd, &ps)
		return 0

	case win32.WM_DESTROY:
		win32.PostQuitMessage(0)
		return 0
	}

	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}
