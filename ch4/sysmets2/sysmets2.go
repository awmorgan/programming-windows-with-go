package main

import (
	"fmt"
	"os"
	"runtime"
	"x/sysmets"
	"x/win32"
)

var cxChar, cxCaps, cyChar, cyClient, iVscrollPos int32

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	var hdc win32.HDC
	var ps win32.PAINTSTRUCT
	var tm win32.TEXTMETRIC
	switch msg {
	case win32.WM_CREATE:
		hdc = win32.GetDC(hwnd)
		win32.GetTextMetrics(hdc, &tm)
		cxChar = int32(tm.TmAveCharWidth)
		cxCaps = cxChar
		if tm.TmPitchAndFamily&1 == 1 {
			cxCaps += cxChar / 2
		}
		cyChar = int32(tm.TmHeight + tm.TmExternalLeading)
		win32.ReleaseDC(hwnd, hdc)
		win32.SetScrollRange(hwnd, win32.SB_VERT, 0, int32(sysmets.NUMLINES-1), false)
		win32.SetScrollPos(hwnd, win32.SB_VERT, iVscrollPos, true)
		return 0
	case win32.WM_SIZE:
		cyClient = win32.HIWORD(lParam)
		return 0
	case win32.WM_VSCROLL:
		switch win32.LOWORD(wParam) {
		case win32.SB_LINEUP:
			iVscrollPos--
		case win32.SB_LINEDOWN:
			iVscrollPos++
		case win32.SB_PAGEUP:
			iVscrollPos -= cyClient / cyChar
		case win32.SB_PAGEDOWN:
			iVscrollPos += cyClient / cyChar
		case win32.SB_THUMBPOSITION:
			iVscrollPos = win32.HIWORD(wParam)
		}
		iVscrollPos = max(0, min(iVscrollPos, int32(sysmets.NUMLINES-1)))
		currentPos, _ := win32.GetScrollPos(hwnd, win32.SB_VERT)
		if iVscrollPos != currentPos {
			win32.SetScrollPos(hwnd, win32.SB_VERT, iVscrollPos, true)
			win32.InvalidateRect(hwnd, nil, true)
		}
		return 0
	case win32.WM_PAINT:
		hdc = win32.BeginPaint(hwnd, &ps)
		for i, sm := range sysmets.Sysmetrics {
			y := cyChar * (int32(i) - iVscrollPos)
			win32.TextOut(hdc, 0, y, sm.Label, len(sm.Label))
			win32.TextOut(hdc, 22*cxCaps, y, sm.Desc, len(sm.Desc))
			win32.SetTextAlign(hdc, win32.TA_RIGHT|win32.TA_TOP)
			s := fmt.Sprintf("%5d", win32.GetSystemMetrics(sm.Index))
			win32.TextOut(hdc, 22*cxCaps+40*cxChar, y, s, len(s))
			win32.SetTextAlign(hdc, win32.TA_LEFT|win32.TA_TOP)
		}
		win32.EndPaint(hwnd, &ps)
		return 0
	case win32.WM_DESTROY:
		win32.PostQuitMessage(0)
		return 0
	}
	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}

func main() {
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.
	appName := "Sysmets2"
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
	hwnd, _ := win32.CreateWindow(appName, "Get System Metrics No. 1", win32.WS_OVERLAPPEDWINDOW|win32.WS_VSCROLL,
		win32.CW_USEDEFAULT, win32.CW_USEDEFAULT, win32.CW_USEDEFAULT, win32.CW_USEDEFAULT,
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
