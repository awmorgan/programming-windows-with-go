//go:build windows

package windows

import (
	"fmt"
	"os"
	"syscall"
)

func NewWndproc(fn func(hWnd HWND, msg uint32, wParam, lParam uintptr) uintptr) WNDPROC {
	return WNDPROC(syscall.NewCallback(fn))
}

type WinmainArgsType struct {
	HInstance HINSTANCE
	NCmdShow  int
}

var WinmainArgs WinmainArgsType

func init() {
	var hInstance Handle
	err := GetModuleHandleEx(GET_MODULE_HANDLE_EX_FLAG_UNCHANGED_REFCOUNT, nil, &hInstance)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WinmainArgs failed to get hInstance: %v\n", err)
		os.Exit(1)
	}
	WinmainArgs.HInstance = HINSTANCE(hInstance)
	var s StartupInfo
	err = GetStartupInfo(&s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WinmainArgs failed to get StartupInfo: %v\n", err)
		os.Exit(1)
	}
	WinmainArgs.NCmdShow = SW_SHOWDEFAULT
	if s.Flags&STARTF_USESHOWWINDOW != 0 {
		WinmainArgs.NCmdShow = int(s.ShowWindow)
	}
}
