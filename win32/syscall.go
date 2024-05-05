//go:build windows

package win32

//sys	MessageBox(hwnd HWND, text string, caption string, boxtype uint32) (ret int32, err error) [failretval==0] = user32.MessageBoxW
