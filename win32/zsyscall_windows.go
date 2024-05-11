// Code generated by 'go generate'; DO NOT EDIT.

package win32

import (
	"syscall"
	"unsafe"
)

var _ unsafe.Pointer

// Do the interface allocations only once for common
// Errno values.
const (
	errnoERROR_IO_PENDING = 997
)

var (
	errERROR_IO_PENDING error = syscall.Errno(errnoERROR_IO_PENDING)
	errERROR_EINVAL     error = syscall.EINVAL
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return errERROR_EINVAL
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	// TODO: add more here, after collecting data on the common
	// error values see on Windows. (perhaps when running
	// all.bat?)
	return e
}

var (
	modgdi32    = NewLazySystemDLL("gdi32.dll")
	modkernel32 = NewLazySystemDLL("kernel32.dll")
	moduser32   = NewLazySystemDLL("user32.dll")
	modwinmm    = NewLazySystemDLL("winmm.dll")

	procGetStockObject      = modgdi32.NewProc("GetStockObject")
	procGetTextMetricsW     = modgdi32.NewProc("GetTextMetricsW")
	procSetTextAlign        = modgdi32.NewProc("SetTextAlign")
	procTextOutW            = modgdi32.NewProc("TextOutW")
	procFreeLibrary         = modkernel32.NewProc("FreeLibrary")
	procGetModuleHandleW    = modkernel32.NewProc("GetModuleHandleW")
	procGetProcAddress      = modkernel32.NewProc("GetProcAddress")
	procGetStartupInfoW     = modkernel32.NewProc("GetStartupInfoW")
	procGetSystemDirectoryW = modkernel32.NewProc("GetSystemDirectoryW")
	procLoadLibraryExW      = modkernel32.NewProc("LoadLibraryExW")
	procBeginPaint          = moduser32.NewProc("BeginPaint")
	procCreateWindowExW     = moduser32.NewProc("CreateWindowExW")
	procDefWindowProcW      = moduser32.NewProc("DefWindowProcW")
	procDispatchMessageW    = moduser32.NewProc("DispatchMessageW")
	procDrawTextW           = moduser32.NewProc("DrawTextW")
	procEndPaint            = moduser32.NewProc("EndPaint")
	procGetClientRect       = moduser32.NewProc("GetClientRect")
	procGetDC               = moduser32.NewProc("GetDC")
	procGetMessageW         = moduser32.NewProc("GetMessageW")
	procGetScrollPos        = moduser32.NewProc("GetScrollPos")
	procGetSystemMetrics    = moduser32.NewProc("GetSystemMetrics")
	procInvalidateRect      = moduser32.NewProc("InvalidateRect")
	procLoadCursorW         = moduser32.NewProc("LoadCursorW")
	procLoadIconW           = moduser32.NewProc("LoadIconW")
	procMessageBoxW         = moduser32.NewProc("MessageBoxW")
	procPostQuitMessage     = moduser32.NewProc("PostQuitMessage")
	procRegisterClassW      = moduser32.NewProc("RegisterClassW")
	procReleaseDC           = moduser32.NewProc("ReleaseDC")
	procSetScrollPos        = moduser32.NewProc("SetScrollPos")
	procSetScrollRange      = moduser32.NewProc("SetScrollRange")
	procShowWindow          = moduser32.NewProc("ShowWindow")
	procTranslateMessage    = moduser32.NewProc("TranslateMessage")
	procUpdateWindow        = moduser32.NewProc("UpdateWindow")
	procPlaySoundW          = modwinmm.NewProc("PlaySoundW")
)

func GetStockObject(fnObject int32) (ret HGDIOBJ) {
	r0, _, _ := syscall.Syscall(procGetStockObject.Addr(), 1, uintptr(fnObject), 0, 0)
	ret = HGDIOBJ(r0)
	return
}

func GetTextMetrics(hdc HDC, tm *TEXTMETRIC) (err error) {
	r1, _, e1 := syscall.Syscall(procGetTextMetricsW.Addr(), 2, uintptr(hdc), uintptr(unsafe.Pointer(tm)), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func SetTextAlign(hdc HDC, align uint32) (ret uint32) {
	r0, _, _ := syscall.Syscall(procSetTextAlign.Addr(), 2, uintptr(hdc), uintptr(align), 0)
	ret = uint32(r0)
	return
}

func TextOut(hdc HDC, x int32, y int32, text string, n int) (err error) {
	var _p0 *uint16
	_p0, err = syscall.UTF16PtrFromString(text)
	if err != nil {
		return
	}
	return _TextOut(hdc, x, y, _p0, n)
}

func _TextOut(hdc HDC, x int32, y int32, text *uint16, n int) (err error) {
	r1, _, e1 := syscall.Syscall6(procTextOutW.Addr(), 5, uintptr(hdc), uintptr(x), uintptr(y), uintptr(unsafe.Pointer(text)), uintptr(n), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func FreeLibrary(handle HANDLE) (err error) {
	r1, _, e1 := syscall.Syscall(procFreeLibrary.Addr(), 1, uintptr(handle), 0, 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func getModuleHandle(moduleName *uint16) (hModule HMODULE, err error) {
	r0, _, e1 := syscall.Syscall(procGetModuleHandleW.Addr(), 1, uintptr(unsafe.Pointer(moduleName)), 0, 0)
	hModule = HMODULE(r0)
	if hModule == 0 {
		err = errnoErr(e1)
	}
	return
}

func GetProcAddress(module HANDLE, procname string) (proc uintptr, err error) {
	var _p0 *byte
	_p0, err = syscall.BytePtrFromString(procname)
	if err != nil {
		return
	}
	return _GetProcAddress(module, _p0)
}

func _GetProcAddress(module HANDLE, procname *byte) (proc uintptr, err error) {
	r0, _, e1 := syscall.Syscall(procGetProcAddress.Addr(), 2, uintptr(module), uintptr(unsafe.Pointer(procname)), 0)
	proc = uintptr(r0)
	if proc == 0 {
		err = errnoErr(e1)
	}
	return
}

func GetStartupInfo(startupInfo *StartupInfo) {
	syscall.Syscall(procGetStartupInfoW.Addr(), 1, uintptr(unsafe.Pointer(startupInfo)), 0, 0)
	return
}

func getSystemDirectory(dir *uint16, dirLen uint32) (len uint32, err error) {
	r0, _, e1 := syscall.Syscall(procGetSystemDirectoryW.Addr(), 2, uintptr(unsafe.Pointer(dir)), uintptr(dirLen), 0)
	len = uint32(r0)
	if len == 0 {
		err = errnoErr(e1)
	}
	return
}

func LoadLibraryEx(libname string, zero HANDLE, flags uintptr) (handle HANDLE, err error) {
	var _p0 *uint16
	_p0, err = syscall.UTF16PtrFromString(libname)
	if err != nil {
		return
	}
	return _LoadLibraryEx(_p0, zero, flags)
}

func _LoadLibraryEx(libname *uint16, zero HANDLE, flags uintptr) (handle HANDLE, err error) {
	r0, _, e1 := syscall.Syscall(procLoadLibraryExW.Addr(), 3, uintptr(unsafe.Pointer(libname)), uintptr(zero), uintptr(flags))
	handle = HANDLE(r0)
	if handle == 0 {
		err = errnoErr(e1)
	}
	return
}

func BeginPaint(hwnd HWND, ps *PAINTSTRUCT) (hdc HDC) {
	r0, _, _ := syscall.Syscall(procBeginPaint.Addr(), 2, uintptr(hwnd), uintptr(unsafe.Pointer(ps)), 0)
	hdc = HDC(r0)
	return
}

func CreateWindowEx(exstyle uint32, className string, windowName string, style uint32, x int32, y int32, width int32, height int32, parent HWND, menu HMENU, instance HINSTANCE, param uintptr) (hwnd HWND, err error) {
	var _p0 *uint16
	_p0, err = syscall.UTF16PtrFromString(className)
	if err != nil {
		return
	}
	var _p1 *uint16
	_p1, err = syscall.UTF16PtrFromString(windowName)
	if err != nil {
		return
	}
	return _CreateWindowEx(exstyle, _p0, _p1, style, x, y, width, height, parent, menu, instance, param)
}

func _CreateWindowEx(exstyle uint32, className *uint16, windowName *uint16, style uint32, x int32, y int32, width int32, height int32, parent HWND, menu HMENU, instance HINSTANCE, param uintptr) (hwnd HWND, err error) {
	r0, _, e1 := syscall.Syscall12(procCreateWindowExW.Addr(), 12, uintptr(exstyle), uintptr(unsafe.Pointer(className)), uintptr(unsafe.Pointer(windowName)), uintptr(style), uintptr(x), uintptr(y), uintptr(width), uintptr(height), uintptr(parent), uintptr(menu), uintptr(instance), uintptr(param))
	hwnd = HWND(r0)
	if hwnd == 0 {
		err = errnoErr(e1)
	}
	return
}

func DefWindowProc(hwnd HWND, msg uint32, wParam uintptr, lParam uintptr) (ret uintptr) {
	r0, _, _ := syscall.Syscall6(procDefWindowProcW.Addr(), 4, uintptr(hwnd), uintptr(msg), uintptr(wParam), uintptr(lParam), 0, 0)
	ret = uintptr(r0)
	return
}

func DispatchMessage(msg *MSG) {
	syscall.Syscall(procDispatchMessageW.Addr(), 1, uintptr(unsafe.Pointer(msg)), 0, 0)
	return
}

func DrawText(hdc HDC, text string, n int32, rect *RECT, format uint32) (ret int32, err error) {
	var _p0 *uint16
	_p0, err = syscall.UTF16PtrFromString(text)
	if err != nil {
		return
	}
	return _DrawText(hdc, _p0, n, rect, format)
}

func _DrawText(hdc HDC, text *uint16, n int32, rect *RECT, format uint32) (ret int32, err error) {
	r0, _, e1 := syscall.Syscall6(procDrawTextW.Addr(), 5, uintptr(hdc), uintptr(unsafe.Pointer(text)), uintptr(n), uintptr(unsafe.Pointer(rect)), uintptr(format), 0)
	ret = int32(r0)
	if ret == 0 {
		err = errnoErr(e1)
	}
	return
}

func EndPaint(hwnd HWND, ps *PAINTSTRUCT) {
	syscall.Syscall(procEndPaint.Addr(), 2, uintptr(hwnd), uintptr(unsafe.Pointer(ps)), 0)
	return
}

func GetClientRect(hwnd HWND, rect *RECT) (err error) {
	r1, _, e1 := syscall.Syscall(procGetClientRect.Addr(), 2, uintptr(hwnd), uintptr(unsafe.Pointer(rect)), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func GetDC(hwnd HWND) (hdc HDC) {
	r0, _, _ := syscall.Syscall(procGetDC.Addr(), 1, uintptr(hwnd), 0, 0)
	hdc = HDC(r0)
	return
}

func GetMessage(msg *MSG, hwnd HWND, msgFilterMin uint32, msgFilterMax uint32) (ret int32, err error) {
	r0, _, e1 := syscall.Syscall6(procGetMessageW.Addr(), 4, uintptr(unsafe.Pointer(msg)), uintptr(hwnd), uintptr(msgFilterMin), uintptr(msgFilterMax), 0, 0)
	ret = int32(r0)
	if ret == -1 {
		err = errnoErr(e1)
	}
	return
}

func GetScrollPos(hwnd HWND, nBar int32) (ret int32, err error) {
	r0, _, e1 := syscall.Syscall(procGetScrollPos.Addr(), 2, uintptr(hwnd), uintptr(nBar), 0)
	ret = int32(r0)
	if ret == 0 {
		err = errnoErr(e1)
	}
	return
}

func GetSystemMetrics(nIndex int32) (ret int32) {
	r0, _, _ := syscall.Syscall(procGetSystemMetrics.Addr(), 1, uintptr(nIndex), 0, 0)
	ret = int32(r0)
	return
}

func InvalidateRect(hwnd HWND, rect *RECT, erase bool) (err error) {
	var _p0 uint32
	if erase {
		_p0 = 1
	}
	r1, _, e1 := syscall.Syscall(procInvalidateRect.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(rect)), uintptr(_p0))
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func LoadCursor(hInstance HINSTANCE, cursorName string) (hCursor HCURSOR, err error) {
	var _p0 *uint16
	_p0, err = syscall.UTF16PtrFromString(cursorName)
	if err != nil {
		return
	}
	return _LoadCursor(hInstance, _p0)
}

func _LoadCursor(hInstance HINSTANCE, cursorName *uint16) (hCursor HCURSOR, err error) {
	r0, _, e1 := syscall.Syscall(procLoadCursorW.Addr(), 2, uintptr(hInstance), uintptr(unsafe.Pointer(cursorName)), 0)
	hCursor = HCURSOR(r0)
	if hCursor == 0 {
		err = errnoErr(e1)
	}
	return
}

func LoadIcon(hInstance HINSTANCE, iconName string) (hIcon HICON, err error) {
	var _p0 *uint16
	_p0, err = syscall.UTF16PtrFromString(iconName)
	if err != nil {
		return
	}
	return _LoadIcon(hInstance, _p0)
}

func _LoadIcon(hInstance HINSTANCE, iconName *uint16) (hIcon HICON, err error) {
	r0, _, e1 := syscall.Syscall(procLoadIconW.Addr(), 2, uintptr(hInstance), uintptr(unsafe.Pointer(iconName)), 0)
	hIcon = HICON(r0)
	if hIcon == 0 {
		err = errnoErr(e1)
	}
	return
}

func MessageBox(hwnd HWND, text string, caption string, boxtype uint32) (ret int32, err error) {
	var _p0 *uint16
	_p0, err = syscall.UTF16PtrFromString(text)
	if err != nil {
		return
	}
	var _p1 *uint16
	_p1, err = syscall.UTF16PtrFromString(caption)
	if err != nil {
		return
	}
	return _MessageBox(hwnd, _p0, _p1, boxtype)
}

func _MessageBox(hwnd HWND, text *uint16, caption *uint16, boxtype uint32) (ret int32, err error) {
	r0, _, e1 := syscall.Syscall6(procMessageBoxW.Addr(), 4, uintptr(hwnd), uintptr(unsafe.Pointer(text)), uintptr(unsafe.Pointer(caption)), uintptr(boxtype), 0, 0)
	ret = int32(r0)
	if ret == 0 {
		err = errnoErr(e1)
	}
	return
}

func PostQuitMessage(exitCode int32) {
	syscall.Syscall(procPostQuitMessage.Addr(), 1, uintptr(exitCode), 0, 0)
	return
}

func RegisterClass(wc *WNDCLASS) (atom ATOM, err error) {
	r0, _, e1 := syscall.Syscall(procRegisterClassW.Addr(), 1, uintptr(unsafe.Pointer(wc)), 0, 0)
	atom = ATOM(r0)
	if atom == 0 {
		err = errnoErr(e1)
	}
	return
}

func ReleaseDC(hwnd HWND, hdc HDC) (err error) {
	r1, _, e1 := syscall.Syscall(procReleaseDC.Addr(), 2, uintptr(hwnd), uintptr(hdc), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func SetScrollPos(hwnd HWND, nBar int32, nPos int32, bRedraw bool) (ret int32, err error) {
	var _p0 uint32
	if bRedraw {
		_p0 = 1
	}
	r0, _, e1 := syscall.Syscall6(procSetScrollPos.Addr(), 4, uintptr(hwnd), uintptr(nBar), uintptr(nPos), uintptr(_p0), 0, 0)
	ret = int32(r0)
	if ret == 0 {
		err = errnoErr(e1)
	}
	return
}

func SetScrollRange(hwnd HWND, nBar int32, nMinPos int32, nMaxPos int32, bRedraw bool) (ret BOOL, err error) {
	var _p0 uint32
	if bRedraw {
		_p0 = 1
	}
	r0, _, e1 := syscall.Syscall6(procSetScrollRange.Addr(), 5, uintptr(hwnd), uintptr(nBar), uintptr(nMinPos), uintptr(nMaxPos), uintptr(_p0), 0)
	ret = BOOL(r0)
	if ret == 0 {
		err = errnoErr(e1)
	}
	return
}

func ShowWindow(hwnd HWND, nCmdShow int32) (wasVisible bool) {
	r0, _, _ := syscall.Syscall(procShowWindow.Addr(), 2, uintptr(hwnd), uintptr(nCmdShow), 0)
	wasVisible = r0 != 0
	return
}

func TranslateMessage(msg *MSG) (translated bool) {
	r0, _, _ := syscall.Syscall(procTranslateMessage.Addr(), 1, uintptr(unsafe.Pointer(msg)), 0, 0)
	translated = r0 != 0
	return
}

func UpdateWindow(hwnd HWND) (ok bool) {
	r0, _, _ := syscall.Syscall(procUpdateWindow.Addr(), 1, uintptr(hwnd), 0, 0)
	ok = r0 != 0
	return
}

func PlaySound(sound string, hmod uintptr, flags uint32) (err error) {
	var _p0 *uint16
	_p0, err = syscall.UTF16PtrFromString(sound)
	if err != nil {
		return
	}
	return _PlaySound(_p0, hmod, flags)
}

func _PlaySound(sound *uint16, hmod uintptr, flags uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procPlaySoundW.Addr(), 3, uintptr(unsafe.Pointer(sound)), uintptr(hmod), uintptr(flags))
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}
