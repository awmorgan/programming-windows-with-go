package win32

import (
	"fmt"
	"syscall"

	"github.com/lxn/win"
)

const (
	MB_OK = win.MB_OK
)

type HWND win.HWND

func Str(str string) *uint16 {
	s, err := syscall.UTF16PtrFromString(str)
	if err != nil {
		panic(err)
	}
	return s
}

func MessageBox(hwnd HWND, text, caption string, flags uint32) (int32, error) {
	t, err := syscall.UTF16PtrFromString(text)
	if err != nil {
		return 0, err
	}
	c, err := syscall.UTF16PtrFromString(caption)
	if err != nil {
		return 0, err
	}
	ret := win.MessageBox(win.HWND(hwnd), t, c, flags)
	if ret == 0 {
		return 0, syscall.GetLastError()
	}
	return ret, nil
}

const (
	SM_CXSCREEN = win.SM_CXSCREEN
	SM_CYSCREEN = win.SM_CYSCREEN
)

func GetSystemMetrics(nIndex int32) (int32, error) {
	ret := win.GetSystemMetrics(nIndex)
	if ret == 0 {
		return 0, fmt.Errorf("GetSystemMetrics failed")
	}
	return ret, nil
}

type HINSTANCE win.HINSTANCE
type HICON win.HICON
type HCURSOR win.HCURSOR
type HBRUSH win.HBRUSH

type WNDCLASS struct {
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     HINSTANCE
	HIcon         HICON
	HCursor       HCURSOR
	HbrBackground HBRUSH
	LpszMenuName  *uint16
	LpszClassName *uint16
}

const (
	CS_HREDRAW    = win.CS_HREDRAW
	CS_VREDRAW    = win.CS_VREDRAW
	CW_USEDEFAULT = win.CW_USEDEFAULT
	WS_OVERLAPPED = win.WS_OVERLAPPED
)

var HInstance HINSTANCE

func init() {
	HInstance := win.GetModuleHandle(nil)
	if HInstance == 0 {
		panic("GetModuleHandle")
	}
}

func NewWndProc(fn func(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr) uintptr {
	return syscall.NewCallback(fn)
}
