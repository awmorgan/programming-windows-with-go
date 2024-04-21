package main

import (
	"os"
	"strings"

	"golang.org/x/sys/windows"
)

// make a wrapper for winmain
// make a wrapper for winproc
func winmain(hInstance windows.Handle, hPrevInstance windows.Handle, lpCmdLine *uint16, nCmdShow int) int {
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

	// convert os.Args to lpCmdLine
	args := strings.Join(os.Args[1:], " ")
	lpCmdLine := windows.StringToUTF16Ptr(args)

	winmain(hInstance, 0, lpCmdLine, nCmdShow)
	// MessageBox("Hello, Windows!", "HelloMsg")
	// println("hInstance:", hInstance)
}

// func MessageBox(caption, text string) int {
// 	t, _ := syscall.BytePtrFromString(text)
// 	c, _ := syscall.BytePtrFromString(caption)

// 	ret, _, _ := MessageBoxA.Call(0, uintptr(unsafe.Pointer(t)), uintptr(unsafe.Pointer(c)), 0)
// 	return int(ret)
// }
