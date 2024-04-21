package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/windows"
)

// var (
// 	user32      = syscall.NewLazyDLL("user32.dll")
// 	MessageBoxA = user32.NewProc("MessageBoxA")

// 	kernel32        = syscall.NewLazyDLL("kernel32.dll")
// )

// make a wrapper for win32 api
// make a wrapper for winmain
// make a wrapper for winproc
func winmain(hInstance windows.Handle, hPrevInstance windows.Handle, lpCmdLine *uint16, nCmdShow int) int {
	return 0
}
func main() {
	fmt.Println("Hello, Windows!")
	var hInstance windows.Handle
	err := windows.GetModuleHandleEx(windows.GET_MODULE_HANDLE_EX_FLAG_UNCHANGED_REFCOUNT, nil, &hInstance)
	if err != nil {
		panic(err)
	}
	fmt.Printf("hInstance: %#v\n", hInstance)
	var s windows.StartupInfo
	err = windows.GetStartupInfo(&s)
	if err != nil {
		panic(err)
	}
	fmt.Printf("s: %#v\n", s)
	var nCmdShow int
	if s.Flags&windows.STARTF_USESHOWWINDOW != 0 {
		fmt.Printf("use show window: %#v\n", s.ShowWindow)
		nCmdShow = int(s.ShowWindow)
	} else {
		fmt.Printf("use show default\n")
		nCmdShow = windows.SW_SHOWDEFAULT
	}
	desktop := windows.UTF16PtrToString(s.Desktop)
	fmt.Printf("desktop: %s\n", desktop)
	title := windows.UTF16PtrToString(s.Title)
	fmt.Printf("title: %s\n", title)

	// convert os.Args to lpCmdLine
	args := strings.Join(os.Args[1:], " ")
	lpCmdLine := windows.StringToUTF16Ptr(args)

	winmain(hInstance, 0, lpCmdLine, nCmdShow)

	// 4th param is show style
	// extern "C" WORD __cdecl __scrt_get_show_window_mode()
	// {
	//     STARTUPINFOW startup_info{};
	//     GetStartupInfoW(&startup_info);
	//     return startup_info.dwFlags & STARTF_USESHOWWINDOW
	// 	? startup_info.wShowWindow
	// 	: SW_SHOWDEFAULT;
	// }

	// int APIENTRY wWinMain(_In_ HINSTANCE hInstance,
	// 	_In_opt_ HINSTANCE hPrevInstance,
	// 	_In_ LPWSTR    lpCmdLine,
	// 	_In_ int       nCmdShow)
	// {
	// UNREFERENCED_PARAMETER(hPrevInstance);
	// UNREFERENCED_PARAMETER(lpCmdLine);
	//func((char*)NULL); // Calls func(int), even if you meant to pass a null pointer
	//func(nullptr); // Unambiguously calls func(char *ptr)

	// TODO: Place code here.
	// HMODULE h = GetModuleHandleW(NULL);
	// if (h != hInstance)

	// MessageBox("Hello, Windows!", "HelloMsg")
	// println("hInstance:", hInstance)
}

// func MessageBox(caption, text string) int {
// 	t, _ := syscall.BytePtrFromString(text)
// 	c, _ := syscall.BytePtrFromString(caption)

// 	ret, _, _ := MessageBoxA.Call(0, uintptr(unsafe.Pointer(t)), uintptr(unsafe.Pointer(c)), 0)
// 	return int(ret)
// }
