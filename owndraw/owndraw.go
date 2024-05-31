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
	appName := "OwnDraw"
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
	hwnd, _ := win32.CreateWindow(appName, "Owner-Draw Button Demo",
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
func triangle(hdc win32.HDC, pt []win32.POINT) {
	win32.SelectObject(hdc, win32.GetStockObject(win32.BLACK_BRUSH))
	win32.Polygon(hdc, pt)
	win32.SelectObject(hdc, win32.GetStockObject(win32.WHITE_BRUSH))
}

var (
	hwndSmaller, hwndLarger win32.HWND
	cxClient, cyClient      int32
	cxChar, cyChar          int32
)

const (
	ID_SMALLER = 1
	ID_LARGER  = 2
)

func btnWidth() int32 {
	return 8 * cxChar
}

func btnHeight() int32 {
	return 4 * cyChar
}

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	switch msg {
	case win32.WM_CREATE:
		cxChar = win32.LOWORD(uintptr(win32.GetDialogBaseUnits()))
		cyChar = win32.HIWORD(uintptr(win32.GetDialogBaseUnits()))
		// create the owner draw pushbuttons
		hwndSmaller, _ = win32.CreateWindow("button", "",
			win32.WS_CHILD|win32.WS_VISIBLE|win32.BS_OWNERDRAW,
			0, 0, btnWidth(), btnHeight(),
			hwnd, win32.HMENU(ID_SMALLER), win32.HInstance(), 0)
		hwndLarger, _ = win32.CreateWindow("button", "",
			win32.WS_CHILD|win32.WS_VISIBLE|win32.BS_OWNERDRAW,
			0, 0, btnWidth(), btnHeight(),
			hwnd, win32.HMENU(ID_LARGER), win32.HInstance(), 0)
		return 0

	case win32.WM_SIZE:
		cxClient = win32.LOWORD(lParam)
		cyClient = win32.HIWORD(lParam)
		// move the buttons to the center
		win32.MoveWindow(hwndSmaller, cxClient/2-3*btnWidth()/2,
			cyClient/2-btnHeight()/2,
			btnWidth(), btnHeight(), true)
		win32.MoveWindow(hwndLarger, cxClient/2+btnWidth()/2,
			cyClient/2-btnHeight()/2,
			btnWidth(), btnHeight(), true)
		return 0
	case win32.WM_COMMAND:
		var rc win32.RECT
		win32.GetWindowRect(hwnd, &rc)
		// make the window 10% smaller or larger
		switch wParam {
		case ID_SMALLER:
			rc.Left += cxClient / 20
			rc.Right -= cxClient / 20
			rc.Top += cyClient / 20
			rc.Bottom -= cyClient / 20
		case ID_LARGER:
			rc.Left -= cxClient / 20
			rc.Right += cxClient / 20
			rc.Top -= cyClient / 20
			rc.Bottom += cyClient / 20
		}
		win32.MoveWindow(hwnd, rc.Left, rc.Top, rc.Right-rc.Left, rc.Bottom-rc.Top, true)
		return 0

	case win32.WM_DRAWITEM:
		dis := (*win32.DRAWITEMSTRUCT)(unsafe.Pointer(lParam))
		// fill area with white and frame it in black
		win32.FillRect(dis.HDC, &dis.RcItem, win32.HBRUSH(win32.GetStockObject(win32.WHITE_BRUSH)))
		win32.FrameRect(dis.HDC, &dis.RcItem, win32.HBRUSH(win32.GetStockObject(win32.BLACK_BRUSH)))

		// draw inward and outward black triangles
		cx := dis.RcItem.Right - dis.RcItem.Left
		cy := dis.RcItem.Bottom - dis.RcItem.Top
		var pt [3]win32.POINT
		switch dis.CtlID {
		case ID_SMALLER:
			pt[0].X = 3 * cx / 8
			pt[0].Y = 1 * cy / 8
			pt[1].X = 5 * cx / 8
			pt[1].Y = 1 * cy / 8
			pt[2].X = 4 * cx / 8
			pt[2].Y = 3 * cy / 8
			triangle(dis.HDC, pt[:])
			pt[0].X = 7 * cx / 8
			pt[0].Y = 3 * cy / 8
			pt[1].X = 7 * cx / 8
			pt[1].Y = 5 * cy / 8
			pt[2].X = 5 * cx / 8
			pt[2].Y = 4 * cy / 8
			triangle(dis.HDC, pt[:])
			pt[0].X = 5 * cx / 8
			pt[0].Y = 7 * cy / 8
			pt[1].X = 3 * cx / 8
			pt[1].Y = 7 * cy / 8
			pt[2].X = 4 * cx / 8
			pt[2].Y = 5 * cy / 8
			triangle(dis.HDC, pt[:])
			pt[0].X = 1 * cx / 8
			pt[0].Y = 5 * cy / 8
			pt[1].X = 1 * cx / 8
			pt[1].Y = 3 * cy / 8
			pt[2].X = 3 * cx / 8
			pt[2].Y = 4 * cy / 8
			triangle(dis.HDC, pt[:])
		case ID_LARGER:
			pt[0].X = 5 * cx / 8
			pt[0].Y = 3 * cy / 8
			pt[1].X = 3 * cx / 8
			pt[1].Y = 3 * cy / 8
			pt[2].X = 4 * cx / 8
			pt[2].Y = 1 * cy / 8
			triangle(dis.HDC, pt[:])
			pt[0].X = 5 * cx / 8
			pt[0].Y = 5 * cy / 8
			pt[1].X = 5 * cx / 8
			pt[1].Y = 3 * cy / 8
			pt[2].X = 7 * cx / 8
			pt[2].Y = 4 * cy / 8
			triangle(dis.HDC, pt[:])
			pt[0].X = 3 * cx / 8
			pt[0].Y = 5 * cy / 8
			pt[1].X = 5 * cx / 8
			pt[1].Y = 5 * cy / 8
			pt[2].X = 4 * cx / 8
			pt[2].Y = 7 * cy / 8
			triangle(dis.HDC, pt[:])
			pt[0].X = 3 * cx / 8
			pt[0].Y = 3 * cy / 8
			pt[1].X = 3 * cx / 8
			pt[1].Y = 5 * cy / 8
			pt[2].X = 1 * cx / 8
			pt[2].Y = 4 * cy / 8
			triangle(dis.HDC, pt[:])
		}
		// invert the rectangle if the button is selected
		if dis.ItemState&win32.ODS_SELECTED != 0 {
			win32.InvertRect(dis.HDC, &dis.RcItem)
		}
		// draw the focus rectangle if the button has the focus
		if dis.ItemState&win32.ODS_FOCUS != 0 {
			dis.RcItem.Left += cx / 16
			dis.RcItem.Top += cy / 16
			dis.RcItem.Right -= cx / 16
			dis.RcItem.Bottom -= cy / 16
			win32.DrawFocusRect(dis.HDC, &dis.RcItem)
		}
		return 0

	case win32.WM_DESTROY:
		win32.PostQuitMessage(0)
		return 0
	}
	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}
