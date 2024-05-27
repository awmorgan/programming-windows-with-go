package main

import (
	"fmt"
	"os"
	"runtime"
	"x/win32"
)

const divisions = 5

var childClassName = "Checker4_Child"

func main() {
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.
	appName := "Checker4"
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
	wc.LpfnWndProc = win32.NewWndProc(childwndproc)
	wc.CbWndExtra = 8 // sizeof(uintptr)
	wc.HIcon = 0
	wc.LpszClassName = win32.StringToUTF16Ptr(childClassName)
	win32.RegisterClass(&wc)

	hwnd, _ := win32.CreateWindow(appName, "Checker4 Mouse Hit-Test Demo",
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

var hwndChild [divisions][divisions]win32.HWND
var idFocus uintptr

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	var cxBlock, cyBlock int32
	switch msg {
	case win32.WM_CREATE:
		for x := int32(0); x < divisions; x++ {
			for y := int32(0); y < divisions; y++ {
				p, _ := win32.GetWindowLongPtr(hwnd, win32.GWLP_HINSTANCE)
				hwndChild[x][y], _ = win32.CreateWindow(childClassName, "",
					win32.WS_CHILD|win32.WS_VISIBLE,
					0, 0, 0, 0,
					hwnd, win32.HMENU(y<<8|x),
					win32.HINSTANCE(p),
					0)
			}
		}
		return 0

	case win32.WM_SIZE:
		cxBlock = win32.LOWORD(lParam) / divisions
		cyBlock = win32.HIWORD(lParam) / divisions
		for x := int32(0); x < divisions; x++ {
			for y := int32(0); y < divisions; y++ {
				win32.MoveWindow(hwndChild[x][y],
					x*cxBlock, y*cyBlock,
					cxBlock, cyBlock, true)
			}
		}
		return 0

	case win32.WM_LBUTTONDOWN:
		win32.MessageBeep(0)
		return 0

	case win32.WM_SETFOCUS:
		id, _ := win32.GetDlgItem(hwnd, int32(idFocus))
		win32.SetFocus(id)
		return 0

	case win32.WM_KEYDOWN:
		x := idFocus & 0xff
		y := idFocus >> 8
		switch wParam {
		case win32.VK_UP:
			y--
		case win32.VK_DOWN:
			y++
		case win32.VK_LEFT:
			x--
		case win32.VK_RIGHT:
			x++
		case win32.VK_HOME:
			x, y = 0, 0
		case win32.VK_END:
			x, y = divisions-1, divisions-1
		default:
			return 0
		}
		x = (x + divisions) % divisions
		y = (y + divisions) % divisions

		idFocus = y<<8 | x
		id, _ := win32.GetDlgItem(hwnd, int32(idFocus))
		win32.SetFocus(id)
		return 0

	case win32.WM_DESTROY:
		win32.PostQuitMessage(0)
		return 0
	}

	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}

func childwndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	switch msg {
	case win32.WM_CREATE:
		win32.SetWindowLongPtr(hwnd, 0, 0)
		return 0

	case win32.WM_KEYDOWN:
		if wParam != win32.VK_RETURN && wParam != win32.VK_SPACE {
			parent, _ := win32.GetParent(hwnd)
			win32.SendMessage(parent, win32.WM_KEYDOWN, wParam, lParam)
			return 0
		}
		fallthrough

	case win32.WM_LBUTTONDOWN:
		p, _ := win32.GetWindowLongPtr(hwnd, 0)
		p ^= 1
		win32.SetWindowLongPtr(hwnd, 0, p)
		win32.InvalidateRect(hwnd, nil, false)
		return 0

	case win32.WM_SETFOCUS:
		idFocus, _ = win32.GetWindowLongPtr(hwnd, win32.GWLP_ID)
		fallthrough

	case win32.WM_KILLFOCUS:
		win32.InvalidateRect(hwnd, nil, false)
		return 0

	case win32.WM_PAINT:
		ps := win32.PAINTSTRUCT{}
		hdc := win32.BeginPaint(hwnd, &ps)

		var r win32.RECT
		win32.GetClientRect(hwnd, &r)
		win32.Rectangle(hdc, 0, 0, r.Right, r.Bottom)

		p, _ := win32.GetWindowLongPtr(hwnd, 0)
		if p != 0 {
			win32.MoveToEx(hdc, 0, 0, nil)
			win32.LineTo(hdc, r.Right, r.Bottom)
			win32.MoveToEx(hdc, 0, r.Bottom, nil)
			win32.LineTo(hdc, r.Right, 0)
		}

		if hwnd == win32.GetFocus() {
			r.Left += r.Right / 10
			r.Right -= r.Left
			r.Top += r.Bottom / 10
			r.Bottom -= r.Top

			win32.SelectObject(hdc, win32.GetStockObject(win32.NULL_BRUSH))
			win32.SelectObject(hdc, win32.HGDIOBJ(win32.CreatePen(win32.PS_DASH, 0, 0)))
			win32.Rectangle(hdc, r.Left, r.Top, r.Right, r.Bottom)
			win32.DeleteObject(win32.SelectObject(hdc, win32.GetStockObject(win32.BLACK_PEN)))
		}

		win32.EndPaint(hwnd, &ps)
		return 0
	}
	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}
