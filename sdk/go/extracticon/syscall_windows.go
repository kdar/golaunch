package extracticon

const (
	LOAD_LIBRARY_AS_DATAFILE uint32 = 0x0002

	MAX_PATH = 260

	SHGFI_SYSICONINDEX      = 0x4000
	SHGFI_ICON              = 0x000000100
	SHGFI_LARGEICON         = 0x000000000
	SHGFI_USEFILEATTRIBUTES = 0x10
	SHIL_JUMBO              = 0x4
	SHIL_EXTRALARGE         = 0x2

	FILE_ATTRIBUTE_NORMAL = 0x80

	ILD_TRANSPARENT = 1
)

var (
	RT_ICON       uintptr = 3
	RT_GROUP_ICON         = RT_ICON + 11
)

type SHFILEINFO struct {
	HIcon         Handle
	IIcon         int32
	DwAttributes  uint32
	SzDisplayName [MAX_PATH]uint16
	SzTypeName    [80]uint16
}

type GUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

// Handle to a the OS specific event log.
type Handle uintptr

const InvalidHandle = ^Handle(0)

//go:generate go run $GOROOT/src/syscall/mksyscall_windows.go -output zsyscall_windows.go syscall_windows.go

//sys   shGetFileInfo(pszPath *uint16, dwFileAttributes uint32, psfi uintptr, cbFileInfo uint32, uFlags uint32) (i uintptr, err error) = shell32.SHGetFileInfoW
//sys   getDC(hWnd Handle) (h Handle, err error) = user32.GetDC
//sys   releaseDC(hWnd Handle, hDC Handle) (ret int, err error) = user32.ReleaseDC
//sys   deleteDC(hdc Handle) (b int, err error) = gdi32.DeleteDC
//sys   deleteObject(hObject Handle) (b int, err error) = gdi32.DeleteObject
//sys   getIconInfo(icon Handle, pinfo uintptr) (b int, err error) = user32.GetIconInfo
//sys   createIconFromResourceEx(pbIconBits *byte, cbIconBits int, fIcon bool, dwVersion int, cxDesired int32, cyDesired int32, flags uint32) (h Handle, err error) = user32.CreateIconFromResourceEx
//sys   loadLibraryEx(filename *uint16, file Handle, flags uint32) (handle Handle, err error) = kernel32.LoadLibraryExW
////sys   enumResourceNames(hModule Handle, lpszType uintptr, lpEnumFunc uintptr,  lParam uintptr) (b int) = kernel32.EnumResourceNamesW
//sys   loadResource(hModule Handle, hResInfo Handle) (h Handle, err error) = kernel32.LoadResource
//sys   lockResource(hResData Handle) (u uintptr, err error) = kernel32.LockResource
//sys   findResourceEx(hModule Handle, lpType uintptr, lpName uintptr, wLanguage int) (h Handle) = kernel32.FindResourceExW
//sys   findResource(hModule Handle, lpName uintptr, lpType uintptr) (h Handle, err error) = kernel32.FindResourceW
////sys   lookupIconIdFromDirectoryEx(presbits uintptr, fIcon bool, cxDesired int, cyDesired int, flags uint) (i int) = user32.LookupIconIdFromDirectoryEx
//sys   sizeofResource(hModule Handle, hResInfo Handle) (i int) = kernel32.SizeofResource
////sys   loadBitmap(hInstance Handle, lpBitmapName *uint16) (h Handle) = user32.LoadBitmapW
////sys   createCompatibleBitmap(hdc Handle, nWidth int, nHeight int) (h Handle) = gdi32.CreateCompatibleBitmap
//sys   createCompatibleDC(hdc Handle) (h Handle, err error) = gdi32.CreateCompatibleDC
//sys   createDIBSection(hdc Handle, pbmih uintptr, iUsage uint32, ppvBits *unsafe.Pointer, hSection Handle, dwOffset uint32) (h Handle, err error) = gdi32.CreateDIBSection
//sys   selectObject(hdc Handle, hgdiobj Handle) (h Handle, err error) = gdi32.SelectObject
//sys   drawIconEx(hdc Handle, xLeft int32, yTop int32, hIcon Handle, cxWidth int32, cyWidth int32, istepIfAniCur uint32, hbrFlickerFreeDraw Handle, diFlags uint32) (b int, err error) = user32.DrawIconEx
//sys   extractIconEx(lpszFile *uint16, nIconIndex int, phiconLarge *Handle, phiconSmall *Handle, nIcons int) (i int, err error) = shell32.ExtractIconExW
////sys   shGetImageList(iImageList int, riid *GUID, ppv **uintptr) (r uintptr, err error) [failretval!=0] = shell32.SHGetImageList
////sys   imageList_GetIcon(himl uintptr, i int, flags uint) (h Handle, err error) = comctl32.ImageList_GetIcon
////sys   extractAssociatedIcon(hInst Handle, lpIconPath *uint16, lpiIcon *uint16) (h Handle, err error) = shell32.ExtractAssociatedIconW
//sys   destroyIcon(hIcon Handle) (b bool, err error) [failretval==false] = user32.DestroyIcon
