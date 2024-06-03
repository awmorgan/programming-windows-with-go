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

// type WNDPROC func(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr
func NewWndProc(f func(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr) uintptr {
	return syscall.NewCallback(f)
}

func NewTimerProc(f func(hwnd HWND, msg uint32, timerID uintptr, time uintptr) uintptr) uintptr {
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

func MAKELONG(lo, hi int32) uint32 {
	return uint32(uint32(lo) | ((uint32(hi)) << 16))
}

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

func UTF16PtrToString(p *uint16) string {
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

func LoadIcon(hInstance HINSTANCE, id uintptr) HICON {
	h, err := loadIcon(hInstance, _MAKEINTRESOURCE(id))
	if err != nil {
		panic(err)
	}
	return h
}

func LoadCursor(hInstance HINSTANCE, id uintptr) HCURSOR {
	h, err := loadCursor(hInstance, _MAKEINTRESOURCE(id))
	if err != nil {
		panic(err)
	}
	return h
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
// directory, which is typically, though not always, `C:WindowsSystem32`.
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

type SecurityAttributes struct {
	Length             uint32
	SecurityDescriptor uintptr
	InheritHandle      uint32
}

type Overlapped struct {
	Internal     uintptr
	InternalHigh uintptr
	Offset       uint32
	OffsetHigh   uint32
	HEvent       HANDLE
}

const INVALID_HANDLE_VALUE = ^HANDLE(0)

// These are the generic rights.
const (
	GENERIC_READ    = 0x80000000
	GENERIC_WRITE   = 0x40000000
	GENERIC_EXECUTE = 0x20000000
	GENERIC_ALL     = 0x10000000
)

// Define access rights to files and directories
const (
	FILE_READ_DATA      = 0x0001 // file & pipe
	FILE_LIST_DIRECTORY = 0x0001 // directory

	FILE_WRITE_DATA = 0x0002 // file & pipe
	FILE_ADD_FILE   = 0x0002 // directory

	FILE_APPEND_DATA          = 0x0004 // file
	FILE_ADD_SUBDIRECTORY     = 0x0004 // directory
	FILE_CREATE_PIPE_INSTANCE = 0x0004 // named pipe

	FILE_READ_EA = 0x0008 // file & directory

	FILE_WRITE_EA = 0x0010 // file & directory

	FILE_EXECUTE  = 0x0020 // file
	FILE_TRAVERSE = 0x0020 // directory

	FILE_DELETE_CHILD = 0x0040 // directory

	FILE_READ_ATTRIBUTES = 0x0080 // all

	FILE_WRITE_ATTRIBUTES = 0x0100 // all

	FILE_ALL_ACCESS = STANDARD_RIGHTS_REQUIRED | SYNCHRONIZE | 0x1FF

	FILE_GENERIC_READ = STANDARD_RIGHTS_READ |
		FILE_READ_DATA |
		FILE_READ_ATTRIBUTES |
		FILE_READ_EA |
		SYNCHRONIZE

	FILE_GENERIC_WRITE = STANDARD_RIGHTS_WRITE |
		FILE_WRITE_DATA |
		FILE_WRITE_ATTRIBUTES |
		FILE_WRITE_EA |
		FILE_APPEND_DATA |
		SYNCHRONIZE

	FILE_GENERIC_EXECUTE = STANDARD_RIGHTS_EXECUTE |
		FILE_READ_ATTRIBUTES |
		FILE_EXECUTE |
		SYNCHRONIZE

	FILE_SHARE_READ                      = 0x00000001
	FILE_SHARE_WRITE                     = 0x00000002
	FILE_SHARE_DELETE                    = 0x00000004
	FILE_ATTRIBUTE_READONLY              = 0x00000001
	FILE_ATTRIBUTE_HIDDEN                = 0x00000002
	FILE_ATTRIBUTE_SYSTEM                = 0x00000004
	FILE_ATTRIBUTE_DIRECTORY             = 0x00000010
	FILE_ATTRIBUTE_ARCHIVE               = 0x00000020
	FILE_ATTRIBUTE_DEVICE                = 0x00000040
	FILE_ATTRIBUTE_NORMAL                = 0x00000080
	FILE_ATTRIBUTE_TEMPORARY             = 0x00000100
	FILE_ATTRIBUTE_SPARSE_FILE           = 0x00000200
	FILE_ATTRIBUTE_REPARSE_POINT         = 0x00000400
	FILE_ATTRIBUTE_COMPRESSED            = 0x00000800
	FILE_ATTRIBUTE_OFFLINE               = 0x00001000
	FILE_ATTRIBUTE_NOT_CONTENT_INDEXED   = 0x00002000
	FILE_ATTRIBUTE_ENCRYPTED             = 0x00004000
	FILE_ATTRIBUTE_INTEGRITY_STREAM      = 0x00008000
	FILE_ATTRIBUTE_VIRTUAL               = 0x00010000
	FILE_ATTRIBUTE_NO_SCRUB_DATA         = 0x00020000
	FILE_ATTRIBUTE_EA                    = 0x00040000
	FILE_ATTRIBUTE_PINNED                = 0x00080000
	FILE_ATTRIBUTE_UNPINNED              = 0x00100000
	FILE_ATTRIBUTE_RECALL_ON_OPEN        = 0x00040000
	FILE_ATTRIBUTE_RECALL_ON_DATA_ACCESS = 0x00400000
	TREE_CONNECT_ATTRIBUTE_PRIVACY       = 0x00004000
	TREE_CONNECT_ATTRIBUTE_INTEGRITY     = 0x00008000
	TREE_CONNECT_ATTRIBUTE_GLOBAL        = 0x00000004
	TREE_CONNECT_ATTRIBUTE_PINNED        = 0x00000002
	FILE_ATTRIBUTE_STRICTLY_SEQUENTIAL   = 0x20000000
	FILE_NOTIFY_CHANGE_FILE_NAME         = 0x00000001
	FILE_NOTIFY_CHANGE_DIR_NAME          = 0x00000002
	FILE_NOTIFY_CHANGE_ATTRIBUTES        = 0x00000004
	FILE_NOTIFY_CHANGE_SIZE              = 0x00000008
	FILE_NOTIFY_CHANGE_LAST_WRITE        = 0x00000010
	FILE_NOTIFY_CHANGE_LAST_ACCESS       = 0x00000020
	FILE_NOTIFY_CHANGE_CREATION          = 0x00000040
	FILE_NOTIFY_CHANGE_SECURITY          = 0x00000100
	FILE_ACTION_ADDED                    = 0x00000001
	FILE_ACTION_REMOVED                  = 0x00000002
	FILE_ACTION_MODIFIED                 = 0x00000003
	FILE_ACTION_RENAMED_OLD_NAME         = 0x00000004
	FILE_ACTION_RENAMED_NEW_NAME         = 0x00000005
	MAILSLOT_NO_MESSAGE                  = ^uint32(0)
	MAILSLOT_WAIT_FOREVER                = ^uint32(0)
	FILE_CASE_SENSITIVE_SEARCH           = 0x00000001
	FILE_CASE_PRESERVED_NAMES            = 0x00000002
	FILE_UNICODE_ON_DISK                 = 0x00000004
	FILE_PERSISTENT_ACLS                 = 0x00000008
	FILE_FILE_COMPRESSION                = 0x00000010
	FILE_VOLUME_QUOTAS                   = 0x00000020
	FILE_SUPPORTS_SPARSE_FILES           = 0x00000040
	FILE_SUPPORTS_REPARSE_POINTS         = 0x00000080
	FILE_SUPPORTS_REMOTE_STORAGE         = 0x00000100
	FILE_RETURNS_CLEANUP_RESULT_INFO     = 0x00000200
	FILE_SUPPORTS_POSIX_UNLINK_RENAME    = 0x00000400
	FILE_SUPPORTS_BYPASS_IO              = 0x00000800
	FILE_SUPPORTS_STREAM_SNAPSHOTS       = 0x00001000
	FILE_SUPPORTS_CASE_SENSITIVE_DIRS    = 0x00002000

	FILE_VOLUME_IS_COMPRESSED         = 0x00008000
	FILE_SUPPORTS_OBJECT_IDS          = 0x00010000
	FILE_SUPPORTS_ENCRYPTION          = 0x00020000
	FILE_NAMED_STREAMS                = 0x00040000
	FILE_READ_ONLY_VOLUME             = 0x00080000
	FILE_SEQUENTIAL_WRITE_ONCE        = 0x00100000
	FILE_SUPPORTS_TRANSACTIONS        = 0x00200000
	FILE_SUPPORTS_HARD_LINKS          = 0x00400000
	FILE_SUPPORTS_EXTENDED_ATTRIBUTES = 0x00800000
	FILE_SUPPORTS_OPEN_BY_FILE_ID     = 0x01000000
	FILE_SUPPORTS_USN_JOURNAL         = 0x02000000
	FILE_SUPPORTS_INTEGRITY_STREAMS   = 0x04000000
	FILE_SUPPORTS_BLOCK_REFCOUNTING   = 0x08000000
	FILE_SUPPORTS_SPARSE_VDL          = 0x10000000
	FILE_DAX_VOLUME                   = 0x20000000
	FILE_SUPPORTS_GHOSTING            = 0x40000000
)

// The following are masks for the predefined standard access types
const (
	DELETE       = 0x00010000
	READ_CONTROL = 0x00020000
	WRITE_DAC    = 0x00040000
	WRITE_OWNER  = 0x00080000
	SYNCHRONIZE  = 0x00100000

	STANDARD_RIGHTS_REQUIRED = 0x000F0000

	STANDARD_RIGHTS_READ    = READ_CONTROL
	STANDARD_RIGHTS_WRITE   = READ_CONTROL
	STANDARD_RIGHTS_EXECUTE = READ_CONTROL

	STANDARD_RIGHTS_ALL = 0x001F0000

	SPECIFIC_RIGHTS_ALL = 0x0000FFFF
)

// File api constants
const (
	CREATE_NEW        = 1
	CREATE_ALWAYS     = 2
	OPEN_EXISTING     = 3
	OPEN_ALWAYS       = 4
	TRUNCATE_EXISTING = 5

	INVALID_FILE_SIZE        = 0xFFFFFFFF
	INVALID_SET_FILE_POINTER = 0xFFFFFFFF
	INVALIDE_FILE_ATTRIBUTES = 0xFFFFFFFF
)
