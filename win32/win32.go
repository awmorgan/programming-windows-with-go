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
