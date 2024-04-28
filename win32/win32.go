package win32

import (
	"os"
	"strings"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

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
	NCmdShow  int
}

func init() {
	WinmainArgs.HInstance = win.GetModuleHandle(nil)
	WinmainArgs.HPrevInst = 0
	args := strings.Join(os.Args, " ")
	WinmainArgs.LpCmdLine = Str(args)
}
