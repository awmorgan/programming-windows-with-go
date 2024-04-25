package main

import (
	"os"
	"strings"
	"syscall"

	"github.com/awmorgan/programming-windows-with-go/internal/sys/windows"
)

func wndproc(hWnd windows.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	return 0
}

func winmain(hInstance windows.Handle, hPrevInstance windows.Handle, lpCmdLine *uint16, nCmdShow int) int {
	var wc windows.WNDCLASSW
	wc.Style = windows.CS_HREDRAW | windows.CS_VREDRAW
	wc.LpfnWndProc = syscall.NewCallback(wndproc)
	wc.CbClsExtra = 0
	wc.CbWndExtra = 0
	wc.HInstance = hInstance
	wc.hIcon = windows.LoadIcon(0, windows.MAKEINTRESOURCE(32512))

	return 0
}

func main() {
	var hInstance windows.Handle
	err := windows.GetModuleHandleEx(windows.GET_MODULE_HANDLE_EX_FLAG_UNCHANGED_REFCOUNT, nil, &hInstance)
	if err != nil {
		panic(err)
	}
	var s windows.StartupInfo
	err = windows.GetStartupInfo(&s)
	if err != nil {
		panic(err)
	}
	nCmdShow := windows.SW_SHOWDEFAULT
	if s.Flags&windows.STARTF_USESHOWWINDOW != 0 {
		nCmdShow = int(s.ShowWindow)
	}
	args := strings.Join(os.Args[1:], " ")
	lpCmdLine, err := windows.UTF16PtrFromString(args)
	if err != nil {
		panic(err)
	}

	winmain(hInstance, 0, lpCmdLine, nCmdShow)
}
