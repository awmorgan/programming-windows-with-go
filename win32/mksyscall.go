//go:build generate

package win32

//go:generate go run x/win32/mkwinsyscall -output zsyscall_windows.go win32.go
