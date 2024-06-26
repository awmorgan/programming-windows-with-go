package win32

//go:generate go run x/win32/mkwinsyscall -output zsyscall_win32.go syscall.go

//sys	BeginPaint(hwnd HWND, ps *PAINTSTRUCT) (hdc HDC) = user32.BeginPaint
//sys	CallWindowProc(lpPrevWndFunc uintptr, hwnd HWND, msg uint32, wParam uintptr, lParam uintptr) (ret uintptr) = user32.CallWindowProcW
//sys	ClientToScreen(hwnd HWND, pt *POINT) (ok bool) = user32.ClientToScreen
//sys	CloseHandle(handle HANDLE) (err error) [failretval==0] = CloseHandle
//sys	CombineRgn(dest HRGN, src1 HRGN, src2 HRGN, mode int32) (ret int32) = gdi32.CombineRgn
//sys	CopyRect(dst *RECT, src *RECT) (ok bool) = user32.CopyRect
//sys	CreateCaret(hwnd HWND, hBitmap HBITMAP, width int32, height int32) (err error) [failretval==0] = user32.CreateCaret
//sys	CreateDC(driver *uint16, device *uint16, port *uint16, pdm *DEVMODE) (hdc HDC) = gdi32.CreateDCW
//sys	CreateEllipticRgn(x1 int32, y1 int32, x2 int32, y2 int32) (hrgn HRGN) = gdi32.CreateEllipticRgn
//sys	CreateEllipticRgnIndirect(rect *RECT) (hrgn HRGN) = gdi32.CreateEllipticRgnIndirect
//sys	CreateFile(filename string, access uint32, sharemode uint32, sa *SecurityAttributes, creation uint32, flags uint32, template HANDLE) (handle HANDLE, err error) [failretval==INVALID_HANDLE_VALUE] = CreateFileW
//sys	CreateFont(height int32, width int32, escapement int32, orientation int32, weight int32, italic int32, underline int32, strikeOut int32, charset int32, outputPrecision int32, clipPrecision int32, quality int32, pitchAndFamily int32, face *uint16) (hfont HFONT) = gdi32.CreateFontW
//sys	CreateIC( driver *uint16, device *uint16, port *uint16, pdm *DEVMODE) (hdc HDC) = gdi32.CreateICW
//sys	CreatePen(style int32, width int32, color COLORREF) (hpen HPEN) = gdi32.CreatePen
//sys	CreatePolygonRgn(pt []POINT, cPoints int32, fnPolyFillMode int32) (hrgn HRGN) = gdi32.CreatePolygonRgn
//sys	CreatePolyPolygonRgn(pt []POINT, lpPolyCounts *int32, nCount int32, fnPolyFillMode int32) (hrgn HRGN) = gdi32.CreatePolyPolygonRgn
//sys	CreateRectRgn(x1 int32, y1 int32, x2 int32, y2 int32) (hrgn HRGN) = gdi32.CreateRectRgn
//sys	CreateRectRgnIndirect(rect *RECT) (hrgn HRGN) = gdi32.CreateRectRgnIndirect
//sys	CreateRoundRectRgn(x1 int32, y1 int32, x2 int32, y2 int32, width int32, height int32) (hrgn HRGN) = gdi32.CreateRoundRectRgn
//sys	CreateSolidBrush(color COLORREF) (hbr HBRUSH) = gdi32.CreateSolidBrush
//sys	CreateWindowEx(exstyle uint32, className string, windowName string, style uint32, x int32, y int32, width int32, height int32, parent HWND, menu HMENU, instance HINSTANCE, param uintptr) (hwnd HWND, err error) [failretval==0] = user32.CreateWindowExW
//sys	DefWindowProc(hwnd HWND, msg uint32, wParam uintptr, lParam uintptr) (ret uintptr) = user32.DefWindowProcW
//sys	DeleteDC(hdc HDC) (ok bool) = gdi32.DeleteDC
//sys	DeleteObject(hObject HGDIOBJ) (ok bool) = gdi32.DeleteObject
//sys	DestroyCaret() (destroyed bool) = user32.DestroyCaret
//sys	DispatchMessage(msg *MSG) = user32.DispatchMessageW
//sys	DPtoLP(hdc HDC, pt []POINT ) (ok bool) = gdi32.DPtoLP
//sys	DrawFocusRect(hdc HDC, lprc *RECT) (ok bool) = user32.DrawFocusRect
//sys	DrawIcon(hdc HDC, x int32, y int32, hIcon HICON) (err error) [failretval==0] = user32.DrawIcon
//sys	DrawText(hdc HDC, text string, n int32, rect *RECT, format uint32) (ret int32, err error) [failretval==0] = user32.DrawTextW
//sys	DrawTextA(hdc HDC, text string, n int32, rect *RECT, format uint32) (ret int32, err error) [failretval==0] = user32.DrawTextA
//sys	Ellipse(hdc HDC, left int32, top int32, right int32, bottom int32) (ok bool) = gdi32.Ellipse
//sys	EndPaint(hwnd HWND, ps *PAINTSTRUCT) = user32.EndPaint
//sys	ExcludeClipRect(hdc HDC, left int32, top int32, right int32, bottom int32) (ret int32) = gdi32.ExcludeClipRect
//sys	FillRect( hdc HDC, lprc *RECT, hbr HBRUSH ) (ok bool) = user32.FillRect
//sys	FillRgn(hdc HDC, hrgn HRGN, hbr HBRUSH) (ok bool) = gdi32.FillRgn
//sys	FindResource(hModule HMODULE, name string, typ string) (hResInfo HRSRC, err error) [failretval==0] = kernel32.FindResourceW
//sys	FrameRect( hdc HDC, lprc *RECT, hbr HBRUSH ) (ok bool) = user32.FrameRect
//sys	FrameRgn(hdc HDC, hrgn HRGN, hbr HBRUSH, width int32, height int32) (ok bool) = gdi32.FrameRgn
//sys	FreeLibrary(handle HANDLE) (err error)
//sys	GetClientRect(hwnd HWND, rect *RECT) (err error) [failretval==0] = user32.GetClientRect
//sys	GetCursorPos(pt *POINT) (err error) [failretval==0] = user32.GetCursorPos
//sys	GetDC(hwnd HWND) (hdc HDC) = user32.GetDC
//sys	GetDeviceCaps(hdc HDC, index int32) (ret int32) = gdi32.GetDeviceCaps
//sys	GetDialogBaseUnits() (base int32) = user32.GetDialogBaseUnits
//sys	GetDlgItem(hwnd HWND, id int32) (hwndChild HWND, err error) [failretval==0] = user32.GetDlgItem
//sys	GetFocus() (hwnd HWND) = user32.GetFocus
//sys	GetKeyNameText(lparam uintptr, buffer []uint16) (ret int32) = user32.GetKeyNameTextW
//sys	GetKeyState(nVirtKey int32) (state int16) = user32.GetKeyState
//sys	GetLocaleInfo(locale LCID, lctype LCTYPE, lcdata *uint16, n int) (cch int32, err error) [failretval==0] = kernel32.GetLocaleInfoW
//sys	GetLocalTime(lpSystemTime *SYSTEMTIME)  = kernel32.GetLocalTime
//sys	GetMessage(msg *MSG, hwnd HWND, msgFilterMin uint32, msgFilterMax uint32) (ret int32, err error) [failretval==-1] = user32.GetMessageW
//sys	getModuleHandle(moduleName *uint16) (hModule HMODULE, err error) [failretval==0] = kernel32.GetModuleHandleW
//sys	GetParent(hwnd HWND) (parent HWND, err error) [failretval==0] = user32.GetParent
//sys	GetPixel(hdc HDC, x int32, y int32) (color COLORREF) = gdi32.GetPixel
//sys	GetProcAddress(module HANDLE, procname string) (proc uintptr, err error)
//sys	GetScrollInfo(hwnd HWND, nBar int32, si *SCROLLINFO) (err error) [failretval==0] = user32.GetScrollInfo
//sys	GetScrollPos(hwnd HWND, nBar int32) (ret int32, err error) [failretval==0] = user32.GetScrollPos
//sys	GetStartupInfo(startupInfo *StartupInfo) = GetStartupInfoW
//sys	GetStockObject(fnObject int32) (ret HGDIOBJ) = gdi32.GetStockObject
//sys	GetSysColor(nIndex int32) (color COLORREF) = user32.GetSysColor
//sys	getSystemDirectory(dir *uint16, dirLen uint32) (len uint32, err error) = kernel32.GetSystemDirectoryW
//sys	GetSystemMetrics(nIndex int32) (ret int32) = user32.GetSystemMetrics
//sys	GetTextFace(hdc HDC, n int32, faceName *uint16) (nOut int32) = gdi32.GetTextFaceW
//sys	GetTextMetrics(hdc HDC, tm *TEXTMETRIC) (err error) [failretval==0] = gdi32.GetTextMetricsW
//sys	GetUpdateRect(hwnd HWND, rect *RECT, erase bool) (notEmpty bool) = user32.GetUpdateRect
//sys	GetWindowLongPtr(hwnd HWND, index int32) (ret uintptr, err error) [failretval==0] = user32.GetWindowLongPtrW
//sys	GetWindowRect(hwnd HWND, rect *RECT) (ok bool, err error) [failretval==false] = user32.GetWindowRect
//sys	HideCaret(hwnd HWND) (err error) [failretval==0] = user32.HideCaret
//sys	InflateRect(rect *RECT, x int32, y int32) (ok bool) = user32.InflateRect
//sys	IntersectClipRect(hdc HDC, left int32, top int32, right int32, bottom int32) (ret int32) = gdi32.IntersectClipRect
//sys	IntersectRect(dst *RECT, src1 *RECT, src2 *RECT) (intersect bool) = user32.IntersectRect
//sys	InvalidateRect(hwnd HWND, rect *RECT, erase bool) (err error) [failretval==0] = user32.InvalidateRect
//sys	InvalidateRgn(hwnd HWND, hrgn HRGN, erase bool) = user32.InvalidateRgn
//sys	InvertRect( hdc HDC, lprc *RECT ) (ok bool) = user32.InvertRect
//sys	InvertRgn(hdc HDC, hrgn HRGN) (ok bool) = gdi32.InvertRgn
//sys	IsRectEmpty(rect *RECT) (empty bool) = user32.IsRectEmpty
//sys	KillTimer(hwnd HWND, id uintptr) (err error) [failretval==0] = user32.KillTimer
//sys	LineTo(hdc HDC, x int32, y int32) (ok bool) = gdi32.LineTo
//sys	loadCursor(hInstance HINSTANCE, cursorName *uint16) (hCursor HCURSOR, err error) [failretval==0] = user32.LoadCursorW
//sys	loadIcon(hInstance HINSTANCE, iconName *uint16) (hIcon HICON, err error) [failretval==0] = user32.LoadIconW
//sys	LoadLibraryEx(libname string, zero HANDLE, flags uintptr) (handle HANDLE, err error) = LoadLibraryExW
//sys	LoadResource(hModule HMODULE, hResInfo HRSRC) (hResData HGLOBAL, err error) [failretval==0] = kernel32.LoadResource
//sys	LoadString(hInstance HINSTANCE, id uint32, buffer []uint16) (n int32, err error) [failretval==0] = user32.LoadStringW
//sys	MessageBeep(uType uint32) (err error) [failretval==0] = user32.MessageBeep
//sys	MessageBox(hwnd HWND, text string, caption string, boxtype uint32) (ret int32, err error) [failretval==0] = user32.MessageBoxW
//sys	MoveToEx(hdc HDC, x int32, y int32, lpPoint *POINT) (ok bool) = gdi32.MoveToEx
//sys	MoveWindow(hwnd HWND, x int32, y int32, width int32, height int32, repaint bool) (err error) [failretval==0] = user32.MoveWindow
//sys	OffsetClipRgn(hdc HDC, x int32, y int32) (ret int32) = gdi32.OffsetClipRgn
//sys	OffsetRect(rect *RECT, x int32, y int32) (ok bool) = user32.OffsetRect
//sys	OffsetWindowOrgEx(hdc HDC, x int32, y int32, pt *POINT) (ok bool) = gdi32.OffsetWindowOrgEx
//sys	PaintRgn(hdc HDC, hrgn HRGN) (ok bool) = gdi32.PaintRgn
//sys	PeekMessage(msg *MSG, hwnd HWND, msgFilterMin uint32, msgFilterMax uint32, removeMsg uint32) (msgAvail bool) = user32.PeekMessageW
//sys	PlaySound(sound string, hmod uintptr, flags uint32) (err error) [failretval==0] = winmm.PlaySoundW
//sys	PolyBezier(hdc HDC, pt []POINT) (ok bool) = gdi32.PolyBezier
//sys	Polygon(hdc HDC, pt []POINT) (ok bool) = gdi32.Polygon
//sys	Polyline(hdc HDC, pt []POINT) (ok bool) = gdi32.Polyline
//sys	PostQuitMessage(exitCode int32) = user32.PostQuitMessage
//sys	ReadFile(handle HANDLE, buffer []byte, bytesRead *uint32, overlapped *Overlapped) (err error) [failretval==0] = ReadFile
//sys	Rectangle(hdc HDC, left int32, top int32, right int32, bottom int32) (ok bool) = gdi32.Rectangle
//sys	RegisterClass(wc *WNDCLASS) (atom ATOM, err error) [failretval==0] = user32.RegisterClassW
//sys	ReleaseCapture() (err error) [failretval==0] = user32.ReleaseCapture
//sys	ReleaseDC(hwnd HWND, hdc HDC) (err error) [failretval==0] = user32.ReleaseDC
//sys	RestoreDC(hdc HDC, saved int32) (ok bool) = gdi32.RestoreDC
//sys	RoundRect(hdc HDC, left int32, top int32, right int32, bottom int32, width int32, height int32) (ok bool) = gdi32.RoundRect
//sys	SaveDC(hdc HDC) (ret int32) = gdi32.SaveDC
//sys	ScreenToClient(hwnd HWND, pt *POINT) (ok bool) = user32.ScreenToClient
//sys	ScrollWindow(hwnd HWND, dx int32, dy int32, rect *RECT, clipRect *RECT) (err error) [failretval==0] = user32.ScrollWindow
//sys	SelectClipRgn(hdc HDC, hrgn HRGN) (mode int32) = gdi32.SelectClipRgn
//sys	SelectObject(hdc HDC, h HGDIOBJ) (ret HGDIOBJ) = gdi32.SelectObject
//sys	SendMessage(hwnd HWND, msg uint32, wParam uintptr, lParam uintptr) (lResult uintptr) = user32.SendMessageW
//sys	SetBkColor( hdc HDC, color COLORREF ) (prevColor COLORREF) = gdi32.SetBkColor
//sys	SetBkMode(hdc HDC, mode int32) (prevMode int32) = gdi32.SetBkMode
//sys	SetCapture(hwnd HWND) (prev HWND) = user32.SetCapture
//sys	SetCaretPos(x int32, y int32) (err error) [failretval==0] = user32.SetCaretPos
//sys	SetClassLongPtr(hwnd HWND, index int32, value uintptr) (prev uintptr, err error) [failretval==0] = user32.SetClassLongPtrW
//sys	SetCursor(hCursor HCURSOR) (hCursorOld HCURSOR) = user32.SetCursor
//sys	SetCursorPos(x int32, y int32) (err error) [failretval==0] = user32.SetCursorPos
//sys	SetFocus(hwnd HWND) (hwndPrev HWND, err error) [failretval==0] = user32.SetFocus
//sys	SetMapMode(hdc HDC, iMapMode int32) (ret int32) = gdi32.SetMapMode
//sys	SetPixel(hdc HDC, x int32, y int32, color COLORREF) (prevColor COLORREF) = gdi32.SetPixel
//sys	SetPolyFillMode(hdc HDC, mode int32) (ret int32) = gdi32.SetPolyFillMode
//sys	SetRect(rect *RECT, left int32, top int32, right int32, bottom int32) (ok bool) = user32.SetRect
//sys	SetRectEmpty(rect *RECT) (ok bool) = user32.SetRectEmpty
//sys	SetROP2(hdc HDC, mode int32) (prevMode int32) = gdi32.SetROP2
//sys	SetScrollInfo(hwnd HWND, nBar int32, si *SCROLLINFO, redraw bool) (pos int32) = user32.SetScrollInfo
//sys	SetScrollPos(hwnd HWND, nBar int32, nPos int32, bRedraw bool) (ret int32, err error) [failretval==0] = user32.SetScrollPos
//sys	SetScrollRange(hwnd HWND, nBar int32, nMinPos int32, nMaxPos int32, bRedraw bool) (err error) [failretval==0] = user32.SetScrollRange
//sys	SetTextAlign(hdc HDC, align uint32) (ret uint32) = gdi32.SetTextAlign
//sys	SetTextColor(hdc HDC, color COLORREF) (prevColor COLORREF) = gdi32.SetTextColor
//sys	SetTimer(hwnd HWND, id uintptr, elapse uint32, proc uintptr) (timerID uintptr, err error) [failretval==0] = user32.SetTimer
//sys	SetViewportExtEx(hdc HDC, x int32, y int32, size *SIZE) (ok bool) = gdi32.SetViewportExtEx
//sys	SetViewportOrgEx(hdc HDC, x int32, y int32, pt *POINT) (ok bool) = gdi32.SetViewportOrgEx
//sys	SetWindowExtEx(hdc HDC, x int32, y int32, size *SIZE) (ok bool) = gdi32.SetWindowExtEx
//sys	SetWindowLongPtr(hwnd HWND, index int32, value uintptr) (prev uintptr, err error) [failretval==0] = user32.SetWindowLongPtrW
//sys	SetWindowOrgEx(hdc HDC, x int32, y int32, pt *POINT) (ok bool) = gdi32.SetWindowOrgEx
//sys	SetWindowText( hwnd HWND, text string ) (ok bool, err error) [failretval==false] = user32.SetWindowTextW
//sys	ShowCaret(hwnd HWND) (err error) [failretval==0] = user32.ShowCaret
//sys	ShowCursor(show bool) (count int32) = user32.ShowCursor
//sys	ShowWindow(hwnd HWND, nCmdShow int32) (wasVisible bool) = user32.ShowWindow
//sys	SystemParametersInfo(uiAction uint32, uiParam uint32, pvParam uintptr, fWinIni uint32) (err error) [failretval==0] = user32.SystemParametersInfoW
//sys	TextOut(hdc HDC, x int32, y int32, text string, n int) (err error) [failretval==0] = gdi32.TextOutW
//sys	TranslateMessage(msg *MSG) (translated bool) = user32.TranslateMessage
//sys	UnionRect(dst *RECT, src1 *RECT, src2 *RECT) (nonempty bool) = user32.UnionRect
//sys	UpdateWindow(hwnd HWND) (ok bool) = user32.UpdateWindow
//sys	ValidateRect(hwnd HWND, rect *RECT) (ok bool) = user32.ValidateRect
//sys	ValidateRgn(hwnd HWND, hrgn HRGN) (ok bool) = user32.ValidateRgn