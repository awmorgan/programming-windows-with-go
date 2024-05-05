//go:build windows

package win32

//sys	MessageBox(hwnd HWND, text string, caption string, boxtype uint32) (ret int32, err error) [failretval==0] = user32.MessageBoxW
//sys	GetSystemMetrics(nIndex int) (ret int, err error) [failretval==0] = user32.GetSystemMetrics
//sys	PlaySound(sound string, hmod uintptr, flags uint32) (ret bool, err error) [failretval==false] = winmm.PlaySoundW

