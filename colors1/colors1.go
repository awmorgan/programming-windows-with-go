package main

import (
	"fmt"
	"os"
	"runtime"
	"x/win32"
)

var (
	idFocus   int32
	oldScroll [3]uintptr
)

func main() {
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.
	appName := "Colors1"
	wc := win32.WNDCLASS{
		Style:         win32.CS_HREDRAW | win32.CS_VREDRAW,
		LpfnWndProc:   win32.NewWndProc(wndproc),
		HInstance:     win32.HInstance(),
		HIcon:         win32.LoadIcon(0, win32.IDI_APPLICATION),
		HCursor:       win32.LoadCursor(0, win32.IDC_ARROW),
		HbrBackground: win32.CreateSolidBrush(0),
		LpszClassName: win32.StringToUTF16Ptr(appName),
	}
	if _, err := win32.RegisterClass(&wc); err != nil {
		errMsg := fmt.Sprintf("RegisterClass failed: %v", err)
		win32.MessageBox(0, errMsg, appName, win32.MB_ICONERROR)
		return
	}
	hwnd, _ := win32.CreateWindow(appName, "Color Scroll",
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
	crPrimaries = [3]win32.COLORREF{
		win32.RGB(255, 0, 0), // Red
		win32.RGB(0, 255, 0), // Green
		win32.RGB(0, 0, 255), // Blue
	}
	hBrushes                         [3]win32.HBRUSH
	hBrushStatic                     win32.HBRUSH
	hwndScroll, hwndLabel, hwndValue [3]win32.HWND
	hwndRect                         win32.HWND
	color                            [3]int32
	cyChar                           int32
	rcColor                          win32.RECT
	colorLabel                       = [3]string{"Red", "Green", "Blue"}
)

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	var hInstance win32.HINSTANCE
	var i, cxClient, cyClient int32
	switch msg {
	case win32.WM_CREATE:
		p, _ := win32.GetWindowLongPtr(hwnd, win32.GWLP_HINSTANCE)
		hInstance = win32.HINSTANCE(p)

		// create the white rectangle window against which the
		// scroll bars will be positioned. the child window id is 9.
		hwndRect, _ = win32.CreateWindow("static", "",
			win32.WS_CHILD|win32.WS_VISIBLE|win32.SS_WHITERECT,
			0, 0, 0, 0,
			hwnd, win32.HMENU(9), hInstance, 0)
		for i = range 3 {
			// the three scroll bars have ids 0, 1, and 2 with
			// scroll bar ranges from 0 to 255.
			hwndScroll[i], _ = win32.CreateWindow("scrollbar", "",
				win32.WS_CHILD|win32.WS_VISIBLE|win32.WS_TABSTOP|win32.SBS_VERT,
				0, 0, 0, 0,
				hwnd, win32.HMENU(i), hInstance, 0)
			win32.SetScrollRange(hwndScroll[i], win32.SB_CTL, 0, 255, false)
			win32.SetScrollPos(hwndScroll[i], win32.SB_CTL, 0, false)

			// the three color-name labels have ids 3, 4, and 5
			// and text "Red", "Green", and "Blue".
			hwndLabel[i], _ = win32.CreateWindow("static", colorLabel[i],
				win32.WS_CHILD|win32.WS_VISIBLE|win32.SS_CENTER,
				0, 0, 0, 0,
				hwnd, win32.HMENU(i+3), hInstance, 0)

			// the three color-value text fields have ids 6, 7, and 8
			// and initial text strings of "0".
			hwndValue[i], _ = win32.CreateWindow("static", "0",
				win32.WS_CHILD|win32.WS_VISIBLE|win32.SS_CENTER,
				0, 0, 0, 0,
				hwnd, win32.HMENU(i+6), hInstance, 0)

			oldScroll[i], _ = win32.SetWindowLongPtr(hwndScroll[i], win32.GWLP_WNDPROC, win32.NewWndProc(scrollproc))

			hBrushes[i] = win32.CreateSolidBrush(crPrimaries[i])
		}
		hBrushStatic = win32.CreateSolidBrush(win32.GetSysColor(win32.COLOR_BTNHIGHLIGHT))
		cyChar = win32.HIWORD(uintptr(win32.GetDialogBaseUnits()))
		return 0

	case win32.WM_SIZE:
		cxClient = win32.LOWORD(lParam)
		cyClient = win32.HIWORD(lParam)

		win32.SetRect(&rcColor, cxClient/2, 0, cxClient, cyClient)
		win32.MoveWindow(hwndRect, 0, 0, cxClient/2, cyClient, true)

		for i = range 3 {
			win32.MoveWindow(hwndScroll[i],
				(2*i+1)*cxClient/14, 2*cyChar,
				cxClient/14, cyClient-4*cyChar, true)

			win32.MoveWindow(hwndLabel[i],
				(4*i+1)*cxClient/28, cyChar/2,
				cxClient/7, cyChar, true)

			win32.MoveWindow(hwndValue[i],
				(4*i+1)*cxClient/28, cyClient-3*cyChar/2,
				cxClient/7, cyChar, true)
		}
		win32.SetFocus(hwnd)
		return 0

	case win32.WM_SETFOCUS:
		win32.SetFocus(hwndScroll[idFocus])
		return 0

	case win32.WM_VSCROLL:
		u, _ := win32.GetWindowLongPtr(win32.HWND(lParam), win32.GWLP_ID)
		i = int32(u)
		switch win32.LOWORD(wParam) {
		case win32.SB_PAGEDOWN:
			color[i] += 15
			fallthrough

		case win32.SB_LINEDOWN:
			color[i] = min(255, color[i]+1)

		case win32.SB_PAGEUP:
			color[i] -= 15
			fallthrough

		case win32.SB_LINEUP:
			color[i] = max(0, color[i]-1)

		case win32.SB_TOP:
			color[i] = 0

		case win32.SB_BOTTOM:
			color[i] = 255

		case win32.SB_THUMBPOSITION, win32.SB_THUMBTRACK:
			color[i] = int32(win32.HIWORD(wParam))
		}

		win32.SetScrollPos(hwndScroll[i], win32.SB_CTL, color[i], true)
		win32.SetWindowText(hwndValue[i], fmt.Sprintf("%d", color[i]))
		b := win32.CreateSolidBrush(win32.RGB(byte(color[0]), byte(color[1]), byte(color[2])))
		p, _ := win32.SetClassLongPtr(hwnd, win32.GCLP_HBRBACKGROUND, uintptr(b))
		win32.DeleteObject(win32.HGDIOBJ(p))

		win32.InvalidateRect(hwnd, &rcColor, true)
		return 0

	case win32.WM_CTLCOLORSCROLLBAR:
		u, _ := win32.GetWindowLongPtr(win32.HWND(lParam), win32.GWLP_ID)
		i = int32(u)
		return uintptr(hBrushes[i])

	case win32.WM_CTLCOLORSTATIC:
		u, _ := win32.GetWindowLongPtr(win32.HWND(lParam), win32.GWLP_ID)
		i = int32(u)
		if 3 <= i && i <= 8 {
			win32.SetTextColor(win32.HDC(wParam), crPrimaries[i%3])
			win32.SetBkColor(win32.HDC(wParam), win32.GetSysColor(win32.COLOR_BTNHIGHLIGHT))
			return uintptr(hBrushStatic)
		}

	case win32.WM_SYSCOLORCHANGE:
		win32.DeleteObject(win32.HGDIOBJ(hBrushStatic))
		hBrushStatic = win32.CreateSolidBrush(win32.GetSysColor(win32.COLOR_BTNHIGHLIGHT))
		return 0

	case win32.WM_DESTROY:
		b := win32.GetStockObject(win32.WHITE_BRUSH)
		p, _ := win32.SetClassLongPtr(hwnd, win32.GCLP_HBRBACKGROUND, uintptr(b))
		win32.DeleteObject(win32.HGDIOBJ(p))
		for i = range 3 {
			win32.DeleteObject(win32.HGDIOBJ(hBrushes[i]))
		}
		win32.PostQuitMessage(0)
		return 0
	}
	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}

func scrollproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	u, _ := win32.GetWindowLongPtr(hwnd, win32.GWLP_ID)
	id := int32(u)
	switch msg {
	case win32.WM_KEYDOWN:
		if wParam == win32.VK_TAB {
			ks := win32.GetKeyState(win32.VK_SHIFT)
			if ks < 0 {
				id += 2
			} else {
				id++
			}
			id %= 3
			p, _ := win32.GetParent(hwnd)
			dlgItem, _ := win32.GetDlgItem(p, id)
			win32.SetFocus(dlgItem)
		}
	case win32.WM_SETFOCUS:
		idFocus = id
	}
	return win32.CallWindowProc(oldScroll[id], hwnd, msg, wParam, lParam)
}
