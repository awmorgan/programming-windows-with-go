package win32

import (
	"os"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

//sys	BeginPaint(hwnd HWND, ps *PAINTSTRUCT) (hdc HDC) = user32.BeginPaint
//sys	CreateWindowEx(exstyle uint32, className string, windowName string, style uint32, x int32, y int32, width int32, height int32, parent HWND, menu HMENU, instance HINSTANCE, param uintptr) (hwnd HWND, err error) [failretval==0] = user32.CreateWindowExW
//sys	DefWindowProc(hwnd HWND, msg uint32, wParam uintptr, lParam uintptr) (ret uintptr) = user32.DefWindowProcW
//sys	DispatchMessage(msg *MSG) = user32.DispatchMessageW
//sys	DrawText(hdc HDC, text string, n int32, rect *RECT, format uint32) (ret int32, err error) [failretval==0] = user32.DrawTextW
//sys	EndPaint(hwnd HWND, ps *PAINTSTRUCT) = user32.EndPaint
//sys	GetClientRect(hwnd HWND, rect *RECT) (err error) [failretval==0] = user32.GetClientRect
//sys	GetDC(hwnd HWND) (hdc HDC) = user32.GetDC
//sys	GetMessage(msg *MSG, hwnd HWND, msgFilterMin uint32, msgFilterMax uint32) (ret int32, err error) [failretval==-1] = user32.GetMessageW
//sys	getModuleHandle(moduleName *uint16) (hModule HMODULE, err error) [failretval==0] = kernel32.GetModuleHandleW
//sys	GetScrollPos(hwnd HWND, nBar int32) (ret int32, err error) [failretval==0] = user32.GetScrollPos
//sys	GetStockObject(fnObject int32) (ret HGDIOBJ) = gdi32.GetStockObject
//sys	GetSystemMetrics(nIndex int32) (ret int32) = user32.GetSystemMetrics
//sys	GetTextMetrics(hdc HDC, tm *TEXTMETRIC) (err error) [failretval==0] = gdi32.GetTextMetricsW
//sys	InvalidateRect(hwnd HWND, rect *RECT, erase bool) (err error) [failretval==0] = user32.InvalidateRect
//sys	LoadCursor(hInstance HINSTANCE, cursorName string) (hCursor HCURSOR, err error) [failretval==0] = user32.LoadCursorW
//sys	LoadIcon(hInstance HINSTANCE, iconName string) (hIcon HICON, err error) [failretval==0] = user32.LoadIconW
//sys	MessageBox(hwnd HWND, text string, caption string, boxtype uint32) (ret int32, err error) [failretval==0] = user32.MessageBoxW
//sys	PlaySound(sound string, hmod uintptr, flags uint32) (err error) [failretval==0] = winmm.PlaySoundW
//sys	PostQuitMessage(exitCode int32) = user32.PostQuitMessage
//sys	RegisterClass(wc *WNDCLASS) (atom ATOM, err error) [failretval==0] = user32.RegisterClassW
//sys	ReleaseDC(hwnd HWND, hdc HDC) (err error) [failretval==0] = user32.ReleaseDC
//sys	SetScrollPos(hwnd HWND, nBar int32, nPos int32, bRedraw bool) (ret int32, err error) [failretval==0] = user32.SetScrollPos
//sys	SetScrollRange(hwnd HWND, nBar int32, nMinPos int32, nMaxPos int32, bRedraw bool) (ret BOOL, err error) [failretval==0] = user32.SetScrollRange
//sys	SetTextAlign(hdc HDC, align uint32) (ret uint32) = gdi32.SetTextAlign
//sys	ShowWindow(hwnd HWND, nCmdShow int32) (wasVisible bool) = user32.ShowWindow
//sys	TextOut(hdc HDC, x int32, y int32, text string, n int) (err error) [failretval==0] = gdi32.TextOutW
//sys	TranslateMessage(msg *MSG) (translated bool) = user32.TranslateMessage
//sys	UpdateWindow(hwnd HWND) (ok bool) = user32.UpdateWindow

var WinmainArgs struct {
	HInstance HINSTANCE
	CmdLine   string
	NCmdShow  int32
}

func init() {
	h, _ := getModuleHandle(nil)
	WinmainArgs.HInstance = HINSTANCE(h)
	WinmainArgs.CmdLine = strings.Join(os.Args, " ")
	s := windows.StartupInfo{Cb: uint32(unsafe.Sizeof(windows.StartupInfo{}))}
	windows.GetStartupInfo(&s)
	if s.Flags&windows.STARTF_USESHOWWINDOW == windows.STARTF_USESHOWWINDOW {
		WinmainArgs.NCmdShow = int32(s.ShowWindow)
	} else {
		WinmainArgs.NCmdShow = SW_SHOWDEFAULT
	}
}

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

func NewWndProc(f func(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr) uintptr {
	return syscall.NewCallback(f)
}

func CreateWindow(className string, windowName string, style uint32, x int32, y int32, width int32, height int32, parent HWND, menu HMENU, instance HINSTANCE, param uintptr) (hwnd HWND, err error) {
	return CreateWindowEx(0, className, windowName, style, x, y, width, height, parent, menu, instance, param)
}

// Constants for flags parameter of PlaySound
const (
	SND_SYNC      uint32 = 0x0000     // play synchronously (default)
	SND_ASYNC     uint32 = 0x0001     // play asynchronously
	SND_NODEFAULT uint32 = 0x0002     // silence (!default) if sound not found
	SND_MEMORY    uint32 = 0x0004     // pszSound points to a memory file
	SND_LOOP      uint32 = 0x0008     // loop the sound until next sndPlaySound
	SND_NOSTOP    uint32 = 0x0010     // don't stop any currently playing sound
	SND_NOWAIT    uint32 = 0x00002000 // don't wait if the driver is busy
	SND_FILENAME  uint32 = 0x00020000 // name is a file name
	SND_RESOURCE  uint32 = 0x00040004 // name is a resource name or atom
)

// Constants for align parameter of SetTextAlign
const (
	TA_NOUPDATECP uint32 = 0
	TA_UPDATECP   uint32 = 1
	TA_LEFT       uint32 = 0
	TA_RIGHT      uint32 = 2
	TA_CENTER     uint32 = 6
	TA_TOP        uint32 = 0
	TA_BOTTOM     uint32 = 8
	TA_BASELINE   uint32 = 24
)

// const (
// 	S_OK           = 0x00000000
// 	S_FALSE        = 0x00000001
// 	E_UNEXPECTED   = 0x8000FFFF
// 	E_NOTIMPL      = 0x80004001
// 	E_OUTOFMEMORY  = 0x8007000E
// 	E_INVALIDARG   = 0x80070057
// 	E_NOINTERFACE  = 0x80004002
// 	E_POINTER      = 0x80004003
// 	E_HANDLE       = 0x80070006
// 	E_ABORT        = 0x80004004
// 	E_FAIL         = 0x80004005
// 	E_ACCESSDENIED = 0x80070005
// 	E_PENDING      = 0x8000000A
// )

// const (
// 	FALSE = 0
// 	TRUE  = 1
// )

type (
	BOOL int32

// HRESULT int32
)

// func SUCCEEDED(hr HRESULT) bool {
// 	return hr >= 0
// }

// func FAILED(hr HRESULT) bool {
// 	return hr < 0
// }

// func MAKEWORD(lo, hi byte) uint16 {
// 	return uint16(uint16(lo) | ((uint16(hi)) << 8))
// }

// func LOBYTE(w uint16) byte {
// 	return byte(w)
// }

// func HIBYTE(w uint16) byte {
// 	return byte(w >> 8 & 0xff)
// }

// func MAKELONG(lo, hi uint16) uint32 {
// 	return uint32(uint32(lo) | ((uint32(hi)) << 16))
// }

func LOWORD(dw uint32) uint16 {
	return uint16(dw)
}

func HIWORD(dw uint32) uint16 {
	return uint16(dw >> 16 & 0xffff)
}

func StringToUTF16Ptr(s string) *uint16 {
	p, err := syscall.UTF16PtrFromString(s)
	if err != nil {
		panic(err)
	}
	return p
}

func _MAKEINTRESOURCE(id uintptr) *uint16 {
	return (*uint16)(unsafe.Pointer(id))
}

func loadIconFromID(id uintptr) HICON {
	h, err := _LoadIcon(0, _MAKEINTRESOURCE(id))
	if err != nil {
		panic(err)
	}
	return h
}

func LoadApplicationIcon() HICON {
	return loadIconFromID(IDI_APPLICATION)
}

func loadCursorFromID(id uintptr) HCURSOR {
	h, err := _LoadCursor(0, _MAKEINTRESOURCE(id))
	if err != nil {
		panic(err)
	}
	return h
}

func LoadArrowCursor() HCURSOR {
	return loadCursorFromID(IDC_ARROW)
}

// func BoolToBOOL(value bool) BOOL {
// 	if value {
// 		return 1
// 	}

// 	return 0
// }
