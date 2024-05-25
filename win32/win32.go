package win32

import (
	"os"
	"strings"
	"syscall"
	"unsafe"
)

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

func LOWORD(dw uintptr) int32 {
	return int32(dw & 0xffff)
}

func HIWORD(dw uintptr) int32 {
	return int32(dw >> 16 & 0xffff)
}

func StringToUTF16Ptr(s string) *uint16 {
	p, err := syscall.UTF16PtrFromString(s)
	if err != nil {
		panic(err)
	}
	return p
}

func Utf16PtrToString(p *uint16) string {
	if p == nil {
		return ""
	}
	end := unsafe.Pointer(p)
	n := 0
	for *(*uint16)(end) != 0 {
		end = unsafe.Pointer(uintptr(end) + unsafe.Sizeof(*p))
		n++
	}
	return syscall.UTF16ToString(unsafe.Slice(p, n))
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

func WaitCursor() HCURSOR {
	return loadCursorFromID(IDC_WAIT)
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

// func PtInRect(rect *RECT, pt POINT) (in bool) {
// 	// Pack X and Y into a single uintptr - each coordinate is 32 bits
// 	ptVal := uintptr(pt.X) | uintptr(pt.Y)<<32
// 	r0, _, _ := syscall.Syscall(procPtInRect.Addr(), 2, uintptr(unsafe.Pointer(rect)), ptVal, 0)
// 	in = r0 != 0
// 	return
// }

// var procPtInRect = moduser32.NewProc("PtInRect")

func PtInRect(lprc *RECT, pt POINT) bool {
	// Check if the rectangle is normalized
	if lprc.Right <= lprc.Left || lprc.Bottom <= lprc.Top {
		return false
	}

	// Check if the point lies within the rectangle
	if pt.X >= lprc.Left && pt.X < lprc.Right &&
		pt.Y >= lprc.Top && pt.Y < lprc.Bottom {
		return true
	}

	return false
}
