package win32

//go:generate go run x/win32/mkwinsyscall -output zsyscall_windows.go win32.go

import (
	"os"
	"strings"
	"syscall"
	"unsafe"
)

//sys	BeginPaint(hwnd HWND, ps *PAINTSTRUCT) (hdc HDC) = user32.BeginPaint
//sys	CreateWindowEx(exstyle uint32, className string, windowName string, style uint32, x int32, y int32, width int32, height int32, parent HWND, menu HMENU, instance HINSTANCE, param uintptr) (hwnd HWND, err error) [failretval==0] = user32.CreateWindowExW
//sys	DefWindowProc(hwnd HWND, msg uint32, wParam uintptr, lParam uintptr) (ret uintptr) = user32.DefWindowProcW
//sys	DispatchMessage(msg *MSG) = user32.DispatchMessageW
//sys	DrawText(hdc HDC, text string, n int32, rect *RECT, format uint32) (ret int32, err error) [failretval==0] = user32.DrawTextW
//sys	Ellipse(hdc HDC, left int32, top int32, right int32, bottom int32) (ok bool) = gdi32.Ellipse
//sys	EndPaint(hwnd HWND, ps *PAINTSTRUCT) = user32.EndPaint
//sys	FreeLibrary(handle HANDLE) (err error)
//sys	GetClientRect(hwnd HWND, rect *RECT) (err error) [failretval==0] = user32.GetClientRect
//sys	GetDC(hwnd HWND) (hdc HDC) = user32.GetDC
//sys	GetDeviceCaps(hdc HDC, index int32) (ret int32) = gdi32.GetDeviceCaps
//sys	GetMessage(msg *MSG, hwnd HWND, msgFilterMin uint32, msgFilterMax uint32) (ret int32, err error) [failretval==-1] = user32.GetMessageW
//sys	getModuleHandle(moduleName *uint16) (hModule HMODULE, err error) [failretval==0] = kernel32.GetModuleHandleW
//sys	GetProcAddress(module HANDLE, procname string) (proc uintptr, err error)
//sys	GetScrollInfo(hwnd HWND, nBar int32, si *SCROLLINFO) (ok bool, err error) [failretval==false] = user32.GetScrollInfo
//sys	GetScrollPos(hwnd HWND, nBar int32) (ret int32, err error) [failretval==0] = user32.GetScrollPos
//sys	GetStartupInfo(startupInfo *StartupInfo) = GetStartupInfoW
//sys	GetStockObject(fnObject int32) (ret HGDIOBJ) = gdi32.GetStockObject
//sys	getSystemDirectory(dir *uint16, dirLen uint32) (len uint32, err error) = kernel32.GetSystemDirectoryW
//sys	GetSystemMetrics(nIndex int32) (ret int32) = user32.GetSystemMetrics
//sys	GetTextMetrics(hdc HDC, tm *TEXTMETRIC) (err error) [failretval==0] = gdi32.GetTextMetricsW
//sys	InvalidateRect(hwnd HWND, rect *RECT, erase bool) (err error) [failretval==0] = user32.InvalidateRect
//sys	LineTo(hdc HDC, x int32, y int32) (ok bool) = gdi32.LineTo
//sys	LoadCursor(hInstance HINSTANCE, cursorName string) (hCursor HCURSOR, err error) [failretval==0] = user32.LoadCursorW
//sys	LoadIcon(hInstance HINSTANCE, iconName string) (hIcon HICON, err error) [failretval==0] = user32.LoadIconW
//sys	LoadLibraryEx(libname string, zero HANDLE, flags uintptr) (handle HANDLE, err error) = LoadLibraryExW
//sys	MessageBox(hwnd HWND, text string, caption string, boxtype uint32) (ret int32, err error) [failretval==0] = user32.MessageBoxW
//sys	MoveToEx(hdc HDC, x int32, y int32, lpPoint *POINT) (ok bool) = gdi32.MoveToEx
//sys	PlaySound(sound string, hmod uintptr, flags uint32) (err error) [failretval==0] = winmm.PlaySoundW
//sys	PolyBezier(hdc HDC, pt []POINT) (ok bool) = gdi32.PolyBezier
//sys	Polyline(hdc HDC, pt []POINT) (ok bool) = gdi32.Polyline
//sys	PostQuitMessage(exitCode int32) = user32.PostQuitMessage
//sys	Rectangle(hdc HDC, left int32, top int32, right int32, bottom int32) (ok bool) = gdi32.Rectangle
//sys	RegisterClass(wc *WNDCLASS) (atom ATOM, err error) [failretval==0] = user32.RegisterClassW
//sys	ReleaseDC(hwnd HWND, hdc HDC) (err error) [failretval==0] = user32.ReleaseDC
//sys	RoundRect(hdc HDC, left int32, top int32, right int32, bottom int32, width int32, height int32) (ok bool) = gdi32.RoundRect
//sys	ScrollWindow(hwnd HWND, dx int32, dy int32, rect *RECT, clipRect *RECT) (ok bool, err error) [failretval==false] = user32.ScrollWindow
//sys	SelectObject(hdc HDC, h HGDIOBJ) (ret HGDIOBJ) = gdi32.SelectObject
//sys	SetScrollInfo(hwnd HWND, nBar int32, si *SCROLLINFO, redraw bool) (pos int32) = user32.SetScrollInfo
//sys	SetScrollPos(hwnd HWND, nBar int32, nPos int32, bRedraw bool) (ret int32, err error) [failretval==0] = user32.SetScrollPos
//sys	SetScrollRange(hwnd HWND, nBar int32, nMinPos int32, nMaxPos int32, bRedraw bool) (ret BOOL, err error) [failretval==0] = user32.SetScrollRange
//sys	SetTextAlign(hdc HDC, align uint32) (ret uint32) = gdi32.SetTextAlign
//sys	ShowWindow(hwnd HWND, nCmdShow int32) (wasVisible bool) = user32.ShowWindow
//sys	TextOut(hdc HDC, x int32, y int32, text string, n int) (err error) [failretval==0] = gdi32.TextOutW
//sys	TranslateMessage(msg *MSG) (translated bool) = user32.TranslateMessage
//sys	UpdateWindow(hwnd HWND) (ok bool) = user32.UpdateWindow

var winmainArgs struct {
	hinstance HINSTANCE
	cmdLine   string
	nCmdShow  int32
}

func init() {
	h, _ := getModuleHandle(nil)
	winmainArgs.hinstance = HINSTANCE(h)
	winmainArgs.cmdLine = strings.Join(os.Args, " ")
	s := StartupInfo{Cb: uint32(unsafe.Sizeof(StartupInfo{}))}
	GetStartupInfo(&s)
	if s.Flags&STARTF_USESHOWWINDOW == STARTF_USESHOWWINDOW {
		winmainArgs.nCmdShow = int32(s.ShowWindow)
	} else {
		winmainArgs.nCmdShow = SW_SHOWDEFAULT
	}
}

func HInstance() HINSTANCE {
	return winmainArgs.hinstance
}

func NCmdShow() int32 {
	return winmainArgs.nCmdShow
}

const (
	STARTF_USESTDHANDLES = 0x00000100
	STARTF_USESHOWWINDOW = 0x00000001
)

type StartupInfo struct {
	Cb            uint32
	_             *uint16
	Desktop       *uint16
	Title         *uint16
	X             uint32
	Y             uint32
	XSize         uint32
	YSize         uint32
	XCountChars   uint32
	YCountChars   uint32
	FillAttribute uint32
	Flags         uint32
	ShowWindow    uint16
	_             uint16
	_             *byte
	StdInput      HANDLE
	StdOutput     HANDLE
	StdErr        HANDLE
}

type StartupInfoEx struct {
	StartupInfo
	ProcThreadAttributeList ProcThreadAttributeList
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

func LOWORD(dw uintptr) uint16 {
	return uint16(dw)
}

func HIWORD(dw uintptr) uint16 {
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

func ApplicationIcon() HICON {
	return loadIconFromID(IDI_APPLICATION)
}

func loadCursorFromID(id uintptr) HCURSOR {
	h, err := _LoadCursor(0, _MAKEINTRESOURCE(id))
	if err != nil {
		panic(err)
	}
	return h
}

func ArrowCursor() HCURSOR {
	return loadCursorFromID(IDC_ARROW)
}

// func BoolToBOOL(value bool) BOOL {
// 	if value {
// 		return 1
// 	}

// 	return 0
// }

// ProcThreadAttributeList is a placeholder type to represent a PROC_THREAD_ATTRIBUTE_LIST.
//
// To create a *ProcThreadAttributeList, use NewProcThreadAttributeList, update
// it with ProcThreadAttributeListContainer.Update, free its memory using
// ProcThreadAttributeListContainer.Delete, and access the list itself using
// ProcThreadAttributeListContainer.List.
type ProcThreadAttributeList struct{}

// LoadLibrary flags for determining from where to search for a DLL
const (
	DONT_RESOLVE_DLL_REFERENCES               = 0x1
	LOAD_LIBRARY_AS_DATAFILE                  = 0x2
	LOAD_WITH_ALTERED_SEARCH_PATH             = 0x8
	LOAD_IGNORE_CODE_AUTHZ_LEVEL              = 0x10
	LOAD_LIBRARY_AS_IMAGE_RESOURCE            = 0x20
	LOAD_LIBRARY_AS_DATAFILE_EXCLUSIVE        = 0x40
	LOAD_LIBRARY_REQUIRE_SIGNED_TARGET        = 0x80
	LOAD_LIBRARY_SEARCH_DLL_LOAD_DIR          = 0x100
	LOAD_LIBRARY_SEARCH_APPLICATION_DIR       = 0x200
	LOAD_LIBRARY_SEARCH_USER_DIRS             = 0x400
	LOAD_LIBRARY_SEARCH_SYSTEM32              = 0x800
	LOAD_LIBRARY_SEARCH_DEFAULT_DIRS          = 0x1000
	LOAD_LIBRARY_SAFE_CURRENT_DIRS            = 0x00002000
	LOAD_LIBRARY_SEARCH_SYSTEM32_NO_FORWARDER = 0x00004000
	LOAD_LIBRARY_OS_INTEGRITY_CONTINUITY      = 0x00008000
)

// GetSystemDirectory retrieves the path to current location of the system
// directory, which is typically, though not always, `C:\Windows\System32`.
func GetSystemDirectory() (string, error) {
	n := uint32(MAX_PATH)
	for {
		b := make([]uint16, n)
		l, e := getSystemDirectory(&b[0], n)
		if e != nil {
			return "", e
		}
		if l <= n {
			return syscall.UTF16ToString(b[:l]), nil
		}
		n = l
	}
}
