package main

import (
	"fmt"
	"os"
	"runtime"
	"x/win32"
)

func main() {
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.
	appName := "Typer"
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
	hwnd, _ := win32.CreateWindow(appName, "Typing Program",
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

var (
	cxClient, cyClient int32
	cxChar, cyChar     int32
	pBuffer            []uint16
	xCaret, yCaret     int32
	charset            int32 = win32.DEFAULT_CHARSET
	cxBuffer, cyBuffer int32
)

func buffer(x, y int32) *uint16 {
	if x < cxBuffer && y < cyBuffer {
		return &pBuffer[y*cxBuffer+x]
	}
	panic("buffer: index out of range")
	// return ' '
}

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	switch msg {
	case win32.WM_INPUTLANGCHANGE:
		fallthrough

	case win32.WM_CREATE:
		hdc := win32.GetDC(hwnd)
		f := win32.CreateFont(0, 0, 0, 0, 0, 0, 0, 0, charset, 0, 0, 0, win32.FIXED_PITCH, nil)
		win32.SelectObject(hdc, win32.HGDIOBJ(f))
		textMetric := win32.TEXTMETRIC{}
		win32.GetTextMetrics(hdc, &textMetric)
		cxChar = textMetric.TmAveCharWidth
		cyChar = textMetric.TmHeight
		win32.DeleteObject(win32.SelectObject(hdc, win32.GetStockObject(win32.SYSTEM_FONT)))
		win32.ReleaseDC(hwnd, hdc)
		fallthrough

	case win32.WM_SIZE:
		if msg == win32.WM_SIZE {
			cxClient = win32.LOWORD(lParam)
			cyClient = win32.HIWORD(lParam)
		}
		cxBuffer = max(1, cxClient/cxChar)
		cyBuffer = max(1, cyClient/cyChar)
		// allocate memory for buffer and clear it
		pBuffer = make([]uint16, cxBuffer*cyBuffer)
		for i := range pBuffer {
			pBuffer[i] = ' '
		}
		// set the caret to the upper left corner
		xCaret = 0
		yCaret = 0
		if hwnd == win32.GetFocus() {
			win32.SetCaretPos(xCaret*cxChar, yCaret*cyChar)
		}

		win32.InvalidateRect(hwnd, nil, true)

		return 0

	case win32.WM_SETFOCUS:
		// create and show the caret
		win32.CreateCaret(hwnd, 0, cxChar, cyChar)
		win32.SetCaretPos(xCaret*cxChar, yCaret*cyChar)
		win32.ShowCaret(hwnd)
		return 0

	case win32.WM_KILLFOCUS:
		// hide and destroy the caret
		win32.HideCaret(hwnd)
		win32.DestroyCaret()
		return 0

	case win32.WM_KEYDOWN:
		switch wParam {
		case win32.VK_HOME:
			xCaret = 0
		case win32.VK_END:
			xCaret = cxBuffer - 1
		case win32.VK_PRIOR:
			yCaret = 0
		case win32.VK_NEXT:
			yCaret = cyBuffer - 1
		case win32.VK_LEFT:
			xCaret = max(xCaret-1, 0)
		case win32.VK_RIGHT:
			xCaret = min(xCaret+1, cxBuffer-1)
		case win32.VK_UP:
			yCaret = max(yCaret-1, 0)
		case win32.VK_DOWN:
			yCaret = min(yCaret+1, cyBuffer-1)
		case win32.VK_DELETE:
			for x := xCaret; x < cxBuffer-1; x++ {
				*buffer(x, yCaret) = *buffer(x+1, yCaret)
			}
			*buffer(cxBuffer-1, yCaret) = ' '
			win32.HideCaret(hwnd)
			hdc := win32.GetDC(hwnd)
			f := win32.CreateFont(0, 0, 0, 0, 0, 0, 0, 0, charset, 0, 0, 0, win32.FIXED_PITCH, nil)
			win32.SelectObject(hdc, win32.HGDIOBJ(f))
			s := win32.UTF16PtrToString(buffer(xCaret, yCaret))
			win32.TextOut(hdc, xCaret*cxChar, yCaret*cyChar, s, int(cxBuffer-xCaret))
			win32.ReleaseDC(hwnd, hdc)
			win32.ShowCaret(hwnd)
		}
		win32.SetCaretPos(xCaret*cxChar, yCaret*cyChar)
		return 0

	case win32.WM_CHAR:
		for range win32.LOWORD(lParam) {
			switch wParam {
			case '\b': // backspace
				if xCaret > 0 {
					xCaret--
					win32.SendMessage(hwnd, win32.WM_KEYDOWN, win32.VK_DELETE, 1)
				}
			case '\t': // tab
				for {
					win32.SendMessage(hwnd, win32.WM_CHAR, ' ', 1)
					if xCaret%8 == 0 {
						break
					}
				}
			case '\n': // line feed
				yCaret++
				if yCaret == cyBuffer {
					yCaret = 0
				}
			case '\r': // carriage return
				xCaret = 0
				yCaret++
				if yCaret == cyBuffer {
					yCaret = 0
				}
			case '\x1B': // escape
				for y := int32(0); y < cyBuffer; y++ {
					for x := int32(0); x < cxBuffer; x++ {
						*buffer(x, y) = ' '
					}
				}
				xCaret = 0
				yCaret = 0
				win32.InvalidateRect(hwnd, nil, false)
			default: // character codes
				*buffer(xCaret, yCaret) = uint16(wParam)
				win32.HideCaret(hwnd)
				hdc := win32.GetDC(hwnd)
				f := win32.CreateFont(0, 0, 0, 0, 0, 0, 0, 0, charset, 0, 0, 0, win32.FIXED_PITCH, nil)
				win32.SelectObject(hdc, win32.HGDIOBJ(f))
				s := win32.UTF16PtrToString(buffer(xCaret, yCaret))
				win32.TextOut(hdc, xCaret*cxChar, yCaret*cyChar, s, 1)
				win32.DeleteObject(win32.SelectObject(hdc, win32.GetStockObject(win32.SYSTEM_FONT)))
				win32.ReleaseDC(hwnd, hdc)
				win32.ShowCaret(hwnd)
				xCaret++
				if xCaret == cxBuffer {
					xCaret = 0
					yCaret++
					if yCaret == cyBuffer {
						yCaret = 0
					}
				}
			}
		}

	case win32.WM_PAINT:
		ps := win32.PAINTSTRUCT{}
		hdc := win32.BeginPaint(hwnd, &ps)
		f := win32.CreateFont(0, 0, 0, 0, 0, 0, 0, 0, charset, 0, 0, 0, win32.FIXED_PITCH, nil)
		win32.SelectObject(hdc, win32.HGDIOBJ(f))
		for y := int32(0); y < cyBuffer; y++ {
			s := win32.UTF16PtrToString(buffer(0, y))
			win32.TextOut(hdc, 0, y*cyChar, s, len(s))
		}

		win32.EndPaint(hwnd, &ps)
		return 0

	case win32.WM_DESTROY:
		win32.PostQuitMessage(0)
		return 0
	}

	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}
