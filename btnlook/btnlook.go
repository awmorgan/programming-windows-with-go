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
	appName := "BtnLook"
	wc := win32.WNDCLASS{
		Style:         win32.CS_HREDRAW | win32.CS_VREDRAW,
		LpfnWndProc:   win32.NewWndProc(wndproc),
		HInstance:     win32.HInstance(),
		HIcon:         win32.LoadIcon(0, win32.IDI_APPLICATION),
		HCursor:       win32.LoadCursor(0, win32.IDC_ARROW),
		HbrBackground: win32.WhiteBrush(),
		LpszClassName: win32.StringToUTF16Ptr(appName),
	}
	if _, err := win32.RegisterClass(&wc); err != nil {
		errMsg := fmt.Sprintf("RegisterClass failed: %v", err)
		win32.MessageBox(0, errMsg, appName, win32.MB_ICONERROR)
		return
	}
	hwnd, _ := win32.CreateWindow(appName, "Button Look",
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

var button = [...]struct {
	style uint32
	text  string
}{
	{win32.BS_PUSHBUTTON, "PUSHBUTTON"},
	{win32.BS_DEFPUSHBUTTON, "DEFPUSHBUTTON"},
	{win32.BS_CHECKBOX, "CHECKBOX"},
	{win32.BS_AUTOCHECKBOX, "AUTOCHECKBOX"},
	{win32.BS_RADIOBUTTON, "RADIOBUTTON"},
	{win32.BS_3STATE, "3STATE"},
	{win32.BS_AUTO3STATE, "AUTO3STATE"},
	{win32.BS_GROUPBOX, "GROUPBOX"},
	{win32.BS_AUTORADIOBUTTON, "AUTORADIOBUTTON"},
	{win32.BS_OWNERDRAW, "OWNERDRAW"},
}

var (
	hwndButton     [len(button)]win32.HWND
	rect           win32.RECT
	top            string = "message            wParam       lParam\n"
	und            string = "_______            ______       ______\n"
	format         string = "%-16s%04X-%04X       %04x-%04x\n"
	cxChar, cyChar int32
)

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	switch msg {
	case win32.WM_CREATE:
		baseUnits := uintptr(win32.GetDialogBaseUnits())
		cxChar = win32.LOWORD(baseUnits)
		cyChar = win32.HIWORD(baseUnits)
		createStruct := (*win32.CREATESTRUCT)(unsafe.Pointer(lParam))
		for i, b := range button {
			hwndButton[i], _ = win32.CreateWindow("button", b.text,
				win32.WS_CHILD|win32.WS_VISIBLE|b.style,
				cxChar, cyChar*(1+2*int32(i)),
				20*cxChar, 7*cyChar/4,
				hwnd, win32.HMENU(i),
				createStruct.Instance, 0)
		}
		return 0

	case win32.WM_SIZE:
		rect.Left = 24 * cxChar
		rect.Top = 2 * cyChar
		rect.Right = win32.LOWORD(lParam)
		rect.Bottom = win32.HIWORD(lParam)
		return 0

	case win32.WM_PAINT:
		win32.InvalidateRect(hwnd, &rect, true)
		var ps win32.PAINTSTRUCT
		hdc := win32.BeginPaint(hwnd, &ps)
		win32.SelectObject(hdc, win32.GetStockObject(win32.SYSTEM_FIXED_FONT))
		win32.SetBkMode(hdc, win32.TRANSPARENT)
		win32.TextOut(hdc, 24*cxChar, cyChar, top, len(top))
		win32.TextOut(hdc, 24*cxChar, cyChar, und, len(und))
		win32.EndPaint(hwnd, &ps)
		return 0

	case win32.WM_DRAWITEM, win32.WM_COMMAND:
		win32.ScrollWindow(hwnd, 0, -cyChar, &rect, &rect)
		hdc := win32.GetDC(hwnd)
		win32.SelectObject(hdc, win32.GetStockObject(win32.SYSTEM_FIXED_FONT))
		var s string
		if msg == win32.WM_DRAWITEM {
			s = fmt.Sprintf(format, "WM_DRAWITEM",
				win32.HIWORD(wParam), win32.LOWORD(wParam),
				win32.HIWORD(lParam), win32.LOWORD(lParam))
		} else {
			s = fmt.Sprintf(format, "WM_COMMAND",
				win32.HIWORD(wParam), win32.LOWORD(wParam),
				win32.HIWORD(lParam), win32.LOWORD(lParam))
		}
		win32.TextOut(hdc, 24*cxChar, cyChar*(rect.Bottom/cyChar-1), s, len(s))
		win32.ReleaseDC(hwnd, hdc)
		win32.ValidateRect(hwnd, &rect)

	case win32.WM_DESTROY:
		win32.PostQuitMessage(0)
		return 0
	}
	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}
