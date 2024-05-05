package win32

import (
	"os"
	"strings"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

type HWND = win.HWND

func MessageBox(hwnd win.HWND, text, caption *uint16, flags uint32) int32 {
	return win.MessageBox(hwnd, text, caption, flags)
}

func Str(s string) *uint16 {
	p, err := syscall.UTF16PtrFromString(s)
	if err != nil {
		panic(err)
	}
	return p
}

type WNDCLASS struct {
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     win.HINSTANCE
	HIcon         win.HICON
	HCursor       win.HCURSOR
	HbrBackground win.HBRUSH
	LpszMenuName  *uint16
	LpszClassName *uint16
}

func RegisterClass(wc *WNDCLASS) win.ATOM {
	wcex := win.WNDCLASSEX{
		CbSize:        uint32(unsafe.Sizeof(win.WNDCLASSEX{})),
		Style:         wc.Style,
		LpfnWndProc:   wc.LpfnWndProc,
		CbClsExtra:    wc.CbClsExtra,
		CbWndExtra:    wc.CbWndExtra,
		HInstance:     wc.HInstance,
		HIcon:         wc.HIcon,
		HCursor:       wc.HCursor,
		HbrBackground: wc.HbrBackground,
		LpszMenuName:  wc.LpszMenuName,
		LpszClassName: wc.LpszClassName,
		HIconSm:       0,
	}
	return win.RegisterClassEx(&wcex)
}

var WinmainArgs struct {
	HInstance win.HINSTANCE
	HPrevInst win.HINSTANCE
	LpCmdLine *uint16
	NCmdShow  int32
}

func init() {
	WinmainArgs.HInstance = win.GetModuleHandle(nil)
	WinmainArgs.HPrevInst = 0
	args := strings.Join(os.Args, " ")
	WinmainArgs.LpCmdLine = Str(args)
	s := windows.StartupInfo{}
	err := windows.GetStartupInfo(&s)
	if err != nil {
		panic(err)
	}
	if s.Flags&windows.STARTF_USESHOWWINDOW == windows.STARTF_USESHOWWINDOW {
		WinmainArgs.NCmdShow = int32(s.ShowWindow)
	} else {
		WinmainArgs.NCmdShow = win.SW_SHOWDEFAULT
	}
}

func NewWndProc(f func(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr) uintptr {
	return syscall.NewCallback(f)
}

func CreateWindow(lpClassName, lpWindowName *uint16, dwStyle uint32, x, y, nWidth, nHeight int32, hWndParent win.HWND, hMenu win.HMENU, hInstance win.HINSTANCE, lpParam unsafe.Pointer) win.HWND {
	return win.CreateWindowEx(0, lpClassName, lpWindowName, dwStyle, x, y, nWidth, nHeight, hWndParent, hMenu, hInstance, lpParam)
}

func DrawText(hdc win.HDC, lpchText *uint16, cchText int32, lprc *win.RECT, dwDTFormat uint32) int32 {
	return win.DrawTextEx(hdc, lpchText, cchText, lprc, dwDTFormat, nil)
}

var (
	libwinmm  = syscall.NewLazyDLL("winmm.dll")
	libgdi32  = syscall.NewLazyDLL("gdi32.dll")
	libuser32 = syscall.NewLazyDLL("user32.dll")

	procPlaySound  = libwinmm.NewProc("PlaySoundW")
	setTextAlign   = libgdi32.NewProc("SetTextAlign")
	setScrollRange = libuser32.NewProc("SetScrollRange")
	setScrollPos   = libuser32.NewProc("SetScrollPos")
	getScrollPos   = libuser32.NewProc("GetScrollPos")
)

// PlaySound plays a sound from a file, resource, or system event.
// Parameters:
// - soundName is the name of the sound to play, or the resource identifier.
// - hMod specifies the executable module (use 0 for a file or system event).
// - flags specify how to play the sound (use SND_ASYNC, SND_FILENAME, SND_RESOURCE, etc.).
func PlaySound(soundName *uint16, hMod win.HMODULE, flags uint32) bool {
	ret, _, _ := procPlaySound.Call(
		uintptr(unsafe.Pointer(soundName)),
		uintptr(hMod),
		uintptr(flags),
	)
	return ret != 0
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

func SetTextAlign(hdc win.HDC, align uint32) uint32 {
	ret, _, _ := setTextAlign.Call(
		uintptr(hdc),
		uintptr(align),
	)
	return uint32(ret)
}

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

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func SetScrollRange(hwnd win.HWND, fnBar int32, nMinPos, nMaxPos int32, redraw bool) bool {
	ret, _, _ := setScrollRange.Call(
		uintptr(hwnd),
		uintptr(fnBar),
		uintptr(nMinPos),
		uintptr(nMaxPos),
		uintptr(boolToInt(redraw)),
	)
	return ret != 0
}

func SetScrollPos(hwnd win.HWND, fnBar, nPos int32, redraw bool) int32 {
	r1, _, _ := setScrollPos.Call(
		uintptr(hwnd),
		uintptr(fnBar),
		uintptr(nPos),
		uintptr(boolToInt(redraw)),
	)
	return int32(r1)
}

func GetScrollPos(hwnd win.HWND, fnBar int32) int32 {
	r1, _, _ := getScrollPos.Call(
		uintptr(hwnd),
		uintptr(fnBar),
	)
	return int32(r1)
}
