package main

import (
	"fmt"
	"os"
	"runtime"
	"unsafe"
	"x/win32"
)

const (
	ID_LIST = 1
	ID_TEXT = 2
)
const (
	MAXREAD = 8192
	DIRATTR = win32.DDL_READWRITE | win32.DDL_READONLY | win32.DDL_HIDDEN | win32.DDL_SYSTEM | win32.DDL_DIRECTORY | win32.DDL_ARCHIVE | win32.DDL_DRIVES
	DTFLAGS = win32.DT_WORDBREAK | win32.DT_EXPANDTABS | win32.DT_NOCLIP | win32.DT_NOPREFIX
)

func main() {
	runtime.LockOSThread() // Windows messages are delivered to the thread that created the window.
	appName := "head"
	wc := win32.WNDCLASS{
		Style:         win32.CS_HREDRAW | win32.CS_VREDRAW,
		LpfnWndProc:   win32.NewWndProc(wndproc),
		HInstance:     win32.HInstance(),
		HIcon:         win32.ApplicationIcon(),
		HCursor:       win32.ArrowCursor(),
		HbrBackground: win32.COLOR_BTNFACE + 1,
		LpszClassName: win32.StringToUTF16Ptr(appName),
	}
	if _, err := win32.RegisterClass(&wc); err != nil {
		errMsg := fmt.Sprintf("RegisterClass failed: %v", err)
		win32.MessageBox(0, errMsg, appName, win32.MB_ICONERROR)
		return
	}
	hwnd, _ := win32.CreateWindow(appName, "head",
		win32.WS_OVERLAPPEDWINDOW,
		win32.CW_USEDEFAULT, win32.CW_USEDEFAULT,
		win32.CW_USEDEFAULT, win32.CW_USEDEFAULT,
		0, 0, win32.HInstance(), 0)

	win32.ShowWindow(hwnd, win32.NCmdShow())
	win32.UpdateWindow(hwnd)
	var msg win32.MSG
	for {
		ret, err := win32.GetMessage(&msg, 0, 0, 0)
		if err != nil {
			errMsg := fmt.Sprintf("GetMessage failed: %v", err)
			win32.MessageBox(0, errMsg, appName, win32.MB_ICONERROR)
			os.Exit(1)
		}
		if ret == 0 {
			os.Exit(int(msg.WParam))
		}

		win32.TranslateMessage(&msg)
		win32.DispatchMessage(&msg)
	}

}

var (
	bValidFile         bool
	buffer             [MAXREAD]byte
	hwndList, hwndText win32.HWND
	rect               win32.RECT
	szFile             string
	oldList            uintptr
)

func wndproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	switch msg {
	case win32.WM_CREATE:
		cxChar := win32.LOWORD(uintptr(win32.GetDialogBaseUnits()))
		cyChar := win32.HIWORD(uintptr(win32.GetDialogBaseUnits()))
		rect.Left = 20 * cxChar
		rect.Top = 3 * cyChar

		hwndList, _ = win32.CreateWindow("listbox", "",
			win32.WS_CHILDWINDOW|win32.WS_VISIBLE|win32.LBS_STANDARD,
			cxChar, cyChar*3,
			cxChar*13+win32.GetSystemMetrics(win32.SM_CXVSCROLL),
			cyChar*10,
			hwnd, ID_LIST, win32.HInstance(), 0)
		cwd, _ := os.Getwd()
		hwndText, _ = win32.CreateWindow("static", cwd,
			win32.WS_CHILDWINDOW|win32.WS_VISIBLE|win32.SS_LEFT,
			cxChar, cyChar, cxChar*win32.MAX_PATH, cyChar,
			hwnd, ID_TEXT, win32.HInstance(), 0)
		oldList, _ = win32.SetWindowLongPtr(hwndList, win32.GWLP_WNDPROC, win32.NewWndProc(listproc))
		win32.SendMessage(hwndList, win32.LB_DIR, DIRATTR, uintptr(unsafe.Pointer(win32.StringToUTF16Ptr("*.*"))))
		return 0

	case win32.WM_SIZE:
		rect.Right = win32.LOWORD(lParam)
		rect.Bottom = win32.HIWORD(lParam)
		return 0

	case win32.WM_SETFOCUS:
		win32.SetFocus(hwndList)
		return 0

	case win32.WM_COMMAND:
		if win32.LOWORD(wParam) == ID_LIST && win32.HIWORD(wParam) == win32.LBN_DBLCLK {
			i := int32(win32.SendMessage(hwndList, win32.LB_GETCURSEL, 0, 0))
			if i == win32.LB_ERR {
				break
			}
			var buf [win32.MAX_PATH + 1]uint16
			win32.SendMessage(hwndList, win32.LB_GETTEXT, uintptr(i), uintptr(unsafe.Pointer(&buf[0])))
			hfile, err := win32.CreateFile(win32.UTF16PtrToString(&buf[0]), win32.GENERIC_READ, win32.FILE_SHARE_READ, nil, win32.OPEN_EXISTING, 0, 0)
			if err == nil {
				win32.CloseHandle(hfile)
				bValidFile = true
				szFile = win32.UTF16PtrToString(&buf[0])
				wd, _ := os.Getwd()
				if wd[len(wd)-1] != '\\' {
					wd += "\\"
				}
				win32.SetWindowText(hwndText, wd+szFile)
			} else {
				bValidFile = false
				// try setting the directory
				err := os.Chdir(win32.UTF16PtrToString(&buf[1]))
				if err != nil {
					buf[3] = ':'
					buf[4] = 0
					os.Chdir(win32.UTF16PtrToString(&buf[2]))
				}
				// get the new directory name and fill the list box
				wd, _ := os.Getwd()
				win32.SetWindowText(hwndText, wd)
				win32.SendMessage(hwndList, win32.LB_RESETCONTENT, 0, 0)
				win32.SendMessage(hwndList, win32.LB_DIR, DIRATTR, uintptr(unsafe.Pointer(win32.StringToUTF16Ptr("*.*"))))
			}
			win32.InvalidateRect(hwnd, nil, true)
		}
		return 0

	case win32.WM_PAINT:
		if !bValidFile {
			break
		}
		hfile, err := win32.CreateFile(szFile, win32.GENERIC_READ, win32.FILE_SHARE_READ, nil, win32.OPEN_EXISTING, 0, 0)
		if err != nil {
			bValidFile = false
			break
		}
		var i uint32
		win32.ReadFile(hfile, buffer[:], &i, nil)
		win32.CloseHandle(hfile)
		// i is the number of bytes read
		var ps win32.PAINTSTRUCT
		hdc := win32.BeginPaint(hwnd, &ps)
		win32.SelectObject(hdc, win32.GetStockObject(win32.SYSTEM_FIXED_FONT))
		win32.SetTextColor(hdc, win32.GetSysColor(win32.COLOR_BTNTEXT))
		win32.SetBkColor(hdc, win32.GetSysColor(win32.COLOR_BTNFACE))
		// assume the file is ASCII text
		win32.DrawTextA(hdc, string(buffer[:i]), int32(i), &rect, DTFLAGS)
		win32.EndPaint(hwnd, &ps)
		return 0

	case win32.WM_DESTROY:
		win32.PostQuitMessage(0)
		return 0
	}
	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}

func listproc(hwnd win32.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	if msg == win32.WM_KEYDOWN && wParam == win32.VK_RETURN {
		parent, _ := win32.GetParent(hwnd)
		win32.SendMessage(parent, win32.WM_COMMAND, uintptr(win32.MAKELONG(1, win32.LBN_DBLCLK)), uintptr(hwnd))
	}
	return win32.CallWindowProc(oldList, hwnd, msg, wParam, lParam)
}
