package main

import (
	"fmt"
	"os"
	"runtime"
	"unsafe"
	"x/win32"
)

func main() {
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.
	appName := "Sysmets3"
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
	hwnd, _ := win32.CreateWindow(appName, "Get System Metrics No. 3",
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

var cxChar, cxCaps, cyChar, cxClient, cyClient, iMaxWidth int32

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	var hdc win32.HDC
	var i, x, y, iVertPos, iHorzPos, iPaintBeg, iPaintEnd int32
	var ps win32.PAINTSTRUCT
	var si win32.SCROLLINFO
	var tm win32.TEXTMETRIC
	switch msg {
	case win32.WM_CREATE:
		hdc = win32.GetDC(hwnd)

		win32.GetTextMetrics(hdc, &tm)
		cxChar = tm.TmAveCharWidth
		cxCaps = cxChar
		if tm.TmPitchAndFamily&1 == 1 {
			cxCaps += cxChar / 2
		}
		cyChar = tm.TmHeight + tm.TmExternalLeading

		win32.ReleaseDC(hwnd, hdc)
		// save the width of the three columns
		iMaxWidth = 40*cxChar + 22*cxCaps
		return 0
	case win32.WM_SIZE:
		cxClient = win32.LOWORD(lParam)
		cyClient = win32.HIWORD(lParam)

		// set vertical scroll bar range and page size
		si.CbSize = uint32(unsafe.Sizeof(si))
		si.FMask = win32.SIF_RANGE | win32.SIF_PAGE
		si.NMin = 0
		si.NMax = int32(win32.NUMLINES - 1)
		si.NPage = uint32(cyClient / cyChar)
		win32.SetScrollInfo(hwnd, win32.SB_VERT, &si, true)

		// set horizontal scroll bar range and page size
		si.CbSize = uint32(unsafe.Sizeof(si))
		si.FMask = win32.SIF_RANGE | win32.SIF_PAGE
		si.NMin = 0
		si.NMax = 2 + iMaxWidth/cxChar
		si.NPage = uint32(cxClient / cxChar)
		win32.SetScrollInfo(hwnd, win32.SB_HORZ, &si, true)
		return 0

	case win32.WM_VSCROLL:
		// Get all the vertical scroll bar information
		si.CbSize = uint32(unsafe.Sizeof(si))
		si.FMask = win32.SIF_ALL
		win32.GetScrollInfo(hwnd, win32.SB_VERT, &si)

		// save the position for comparison later on
		iVertPos = si.NPos

		switch win32.LOWORD(wParam) {
		case win32.SB_TOP:
			si.NPos = si.NMin
		case win32.SB_BOTTOM:
			si.NPos = si.NMax
		case win32.SB_LINEUP:
			si.NPos -= 1
		case win32.SB_LINEDOWN:
			si.NPos += 1
		case win32.SB_PAGEUP:
			si.NPos -= int32(si.NPage)
		case win32.SB_PAGEDOWN:
			si.NPos += int32(si.NPage)
		case win32.SB_THUMBTRACK:
			si.NPos = si.NTrackPos
		default:
		}

		// Set the position and then retrieve it.  Due to adjustments
		// by Windows it may not be the same as the value set.
		si.FMask = win32.SIF_POS
		win32.SetScrollInfo(hwnd, win32.SB_VERT, &si, true)
		win32.GetScrollInfo(hwnd, win32.SB_VERT, &si)

		// If the position has changed, scroll the window and update it
		if si.NPos != iVertPos {
			win32.ScrollWindow(hwnd, 0, cyChar*(iVertPos-si.NPos), nil, nil)
			win32.UpdateWindow(hwnd)
		}
		return 0
	case win32.WM_HSCROLL:
		// Get all the horizontal scroll bar information
		si.CbSize = uint32(unsafe.Sizeof(si))
		si.FMask = win32.SIF_ALL

		// save the position for comparison later on
		win32.GetScrollInfo(hwnd, win32.SB_HORZ, &si)
		iHorzPos = si.NPos

		switch win32.LOWORD(wParam) {
		case win32.SB_LINELEFT:
			si.NPos -= 1
		case win32.SB_LINERIGHT:
			si.NPos += 1
		case win32.SB_PAGELEFT:
			si.NPos -= int32(si.NPage)
		case win32.SB_PAGERIGHT:
			si.NPos += int32(si.NPage)
		case win32.SB_THUMBPOSITION:
			si.NPos = si.NTrackPos
		default:
		}
		// Set the position and then retrieve it.  Due to adjustments
		// by Windows it may not be the same as the value set.
		si.FMask = win32.SIF_POS
		win32.SetScrollInfo(hwnd, win32.SB_HORZ, &si, true)
		win32.GetScrollInfo(hwnd, win32.SB_HORZ, &si)

		// If the position has changed, scroll the window
		if si.NPos != iHorzPos {
			win32.ScrollWindow(hwnd, cxChar*(iHorzPos-si.NPos), 0, nil, nil)
		}
		return 0
	case win32.WM_PAINT:
		hdc = win32.BeginPaint(hwnd, &ps)

		// get vertical scroll bar position
		si.CbSize = uint32(unsafe.Sizeof(si))
		si.FMask = win32.SIF_POS
		win32.GetScrollInfo(hwnd, win32.SB_VERT, &si)
		iVertPos = si.NPos

		// get horizontal scroll bar position
		win32.GetScrollInfo(hwnd, win32.SB_HORZ, &si)
		iHorzPos = si.NPos

		// find painting limits
		iPaintBeg = max(0, iVertPos+int32(ps.RcPaint.Top/cyChar))
		iPaintEnd = min(int32(win32.NUMLINES)-1,
			iVertPos+int32(ps.RcPaint.Bottom/cyChar))

		for i = iPaintBeg; i <= iPaintEnd; i++ {
			x = cxChar * (1 - iHorzPos)
			y = cyChar * (i - iVertPos)

			win32.TextOut(hdc, x, y,
				win32.Sysmetrics[i].Label,
				len(win32.Sysmetrics[i].Label))

			win32.TextOut(hdc, x+22*cxChar, y,
				win32.Sysmetrics[i].Desc,
				len(win32.Sysmetrics[i].Desc))

			win32.SetTextAlign(hdc, win32.TA_RIGHT|win32.TA_TOP)

			s := fmt.Sprintf("%5d", win32.GetSystemMetrics(win32.Sysmetrics[i].Index))

			win32.TextOut(hdc, x+22*cxChar+40*cxCaps, y, s, len(s))

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
