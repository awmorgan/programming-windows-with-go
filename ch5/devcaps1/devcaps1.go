package main

import (
	"fmt"
	"os"
	"runtime"
	"x/win32"
)

const NUMLINES = int32(len(devcaps))

var devcaps = [...]struct {
	iIndex  int32
	szLabel string
	szDesc  string
}{
	{win32.HORZSIZE, "HORZSIZE", "Width in millimeters:"},
	{win32.VERTSIZE, "VERTSIZE", "Height in millimeters:"},
	{win32.HORZRES, "HORZRES", "Width in pixels:"},
	{win32.VERTRES, "VERTRES", "Height in raster lines:"},
	{win32.BITSPIXEL, "BITSPIXEL", "Color bits per pixel:"},
	{win32.PLANES, "PLANES", "Number of color planes:"},
	{win32.NUMBRUSHES, "NUMBRUSHES", "Number of device brushes:"},
	{win32.NUMPENS, "NUMPENS", "Number of device pens:"},
	{win32.NUMMARKERS, "NUMMARKERS", "Number of device markers:"},
	{win32.NUMFONTS, "NUMFONTS", "Number of device fonts:"},
	{win32.NUMCOLORS, "NUMCOLORS", "Number of device colors:"},
	{win32.PDEVICESIZE, "PDEVICESIZE", "Size of device structure:"},
	{win32.ASPECTX, "ASPECTX", "Relative width of pixel:"},
	{win32.ASPECTY, "ASPECTY", "Relative height of pixel:"},
	{win32.ASPECTXY, "ASPECTXY", "Relative diagonal of pixel:"},
	{win32.LOGPIXELSX, "LOGPIXELSX", "Horizontal dots per inch:"},
	{win32.LOGPIXELSY, "LOGPIXELSY", "Vertical dots per inch:"},
	{win32.SIZEPALETTE, "SIZEPALETTE", "Number of palette entries:"},
	{win32.NUMRESERVED, "NUMRESERVED", "Reserved palette entries:"},
	{win32.COLORRES, "COLORRES", "Actual color resolution:"},
}

func main() {
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.
	appName := "Devcaps1"
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
		errMsg := "RegisterClass failed: " + err.Error()
		win32.MessageBox(0, errMsg, appName, win32.MB_ICONERROR)
		return
	}
	hwnd, _ := win32.CreateWindow(appName, "Device Capabilities",
		win32.WS_OVERLAPPEDWINDOW,
		win32.CW_USEDEFAULT, win32.CW_USEDEFAULT,
		win32.CW_USEDEFAULT, win32.CW_USEDEFAULT,
		0, 0, win32.HInstance(), 0)
	win32.ShowWindow(hwnd, win32.NCmdShow())
	win32.UpdateWindow(hwnd)
	msg := win32.MSG{}
	for {
		ret, err := win32.GetMessage(&msg, 0, 0, 0)
		if err != nil {
			errMsg := "GetMessage failed: " + err.Error()
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

var cxChar, cxCaps, cyChar int32

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	var hdc win32.HDC
	var i int32
	var ps win32.PAINTSTRUCT
	var tm win32.TEXTMETRIC

	switch msg {
	case win32.WM_CREATE:
		hdc = win32.GetDC(hwnd)
		win32.GetTextMetrics(hdc, &tm)
		cxChar = tm.TmAveCharWidth
		cxCaps = cxChar
		if tm.TmPitchAndFamily&1 == 1 {
			cxCaps = 3 * cxChar / 2
		}
		cyChar = tm.TmHeight + tm.TmExternalLeading
		win32.ReleaseDC(hwnd, hdc)
		return 0

	case win32.WM_PAINT:
		hdc = win32.BeginPaint(hwnd, &ps)

		for i = 0; i < NUMLINES; i++ {
			win32.TextOut(hdc, 0, cyChar*i,
				devcaps[i].szLabel,
				len(devcaps[i].szLabel))

			win32.TextOut(hdc, 14*cxCaps, cyChar*i,
				devcaps[i].szDesc,
				len(devcaps[i].szDesc))

			win32.SetTextAlign(hdc, win32.TA_RIGHT|win32.TA_TOP)

			s := fmt.Sprintf("%5d", win32.GetDeviceCaps(hdc, devcaps[i].iIndex))
			win32.TextOut(hdc, 14*cxCaps+35*cxChar, cyChar*i, s, len(s))

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
