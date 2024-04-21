package main

import (
	"os"
	"strings"

	"golang.org/x/sys/windows"
)

func winmain(hInstance windows.Handle, hPrevInstance windows.Handle, lpCmdLine *uint16, nCmdShow int) int {
	text := windows.StringToUTF16Ptr("Hello, Windows!")
	caption := windows.StringToUTF16Ptr("HelloMsg")
	windows.MessageBox(0, text, caption, 0)
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
	lpCmdLine := windows.StringToUTF16Ptr(args)
	
	winmain(hInstance, 0, lpCmdLine, nCmdShow)
}
