//go:build generate

package win32

//go:generate go run golang.org/x/sys/windows/mkwinsyscall -output zsyscall_windows.go win32.go
