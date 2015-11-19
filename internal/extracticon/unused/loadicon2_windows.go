package extracticon

import (
	"image"
	"image/color"
	"syscall"
	"unsafe"

	"github.com/ansel1/merry"
	"github.com/lxn/win"
)

var (
	doinit = true
)

func From2(path string) (image.Image, error) {
	var shfi SHFILEINFO

	_, err := shGetFileInfo(
		syscall.StringToUTF16Ptr(path),
		0,
		uintptr(unsafe.Pointer(&shfi)),
		uint32(unsafe.Sizeof(shfi)),
		SHGFI_ICON|SHGFI_LARGEICON)
	if err != nil || shfi.HIcon == 0 {
		return nil, merry.Wrap(err)
	}

	// var large []win.HICON
	// fmt.Println(ExtractIconEx(syscall.StringToUTF16Ptr(exe), 7, large, nil, 1))
	// fmt.Println(large)
	//fmt.Println(ExtractIcon(exe, 0))

	// hIcon := ImageList_GetIcon(hIml, int(shfi.IIcon), 0)
	// fmt.Println(shfi.IIcon)

	var iconInfo win.ICONINFO
	_, err = getIconInfo(Handle(shfi.HIcon), uintptr(unsafe.Pointer(&iconInfo)))
	if err != nil {
		return nil, merry.Wrap(err)
	}

	w := int32(iconInfo.XHotspot * 2)
	h := int32(iconInfo.YHotspot * 2)

	var bitmapInfo win.BITMAPINFOHEADER
	bitmapInfo.BiSize = uint32(unsafe.Sizeof(bitmapInfo))
	bitmapInfo.BiWidth = w
	bitmapInfo.BiHeight = -h
	bitmapInfo.BiPlanes = 1
	bitmapInfo.BiBitCount = 32
	bitmapInfo.BiCompression = win.BI_RGB
	bitmapInfo.BiSizeImage = uint32(w * h * 4)
	bitmapInfo.BiXPelsPerMeter = 0
	bitmapInfo.BiYPelsPerMeter = 0
	bitmapInfo.BiClrUsed = 0
	bitmapInfo.BiClrImportant = 0

	screenDevice, err := getDC(0)
	if err != nil {
		return nil, merry.Wrap(err)
	}

	hdc, err := createCompatibleDC(screenDevice)
	if err != nil {
		return nil, merry.Wrap(err)
	}
	if _, err := releaseDC(0, screenDevice); err != nil {
		return nil, merry.Wrap(err)
	}
	defer deleteDC(hdc)

	var bits unsafe.Pointer
	winBitmap, err := createDIBSection(hdc, uintptr(unsafe.Pointer(&bitmapInfo)), win.DIB_RGB_COLORS, &bits, 0, 0)
	if err != nil {
		return nil, merry.Wrap(err)
	}
	defer deleteObject(winBitmap)

	var pixels = (*[1 << 30]byte)(bits)[0:bitmapInfo.BiSizeImage]
	// for i := uint32(0); i < bitmapInfo.BiSizeImage; i++ {
	// 	pixels[i] = 0xFF
	// }

	if _, err := selectObject(hdc, winBitmap); err != nil {
		return nil, merry.Wrap(err)
	}

	if _, err := drawIconEx(hdc, 0, 0, Handle(shfi.HIcon), w, h, 0, 0, win.DI_NORMAL); err != nil {
		return nil, merry.Wrap(err)
	}

	rgba := image.NewRGBA(image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: 0,
		},
		Max: image.Point{
			X: int(w),
			Y: int(h),
		},
	})
	for y := int32(0); y < h; y++ {
		for x := int32(0); x < w; x++ {
			rgba.SetRGBA(int(x), int(y), color.RGBA{
				A: uint8(pixels[((y*w+x)*4)+3]),
				R: uint8(pixels[((y*w+x)*4)+2]),
				G: uint8(pixels[((y*w+x)*4)+1]),
				B: uint8(pixels[((y*w+x)*4)+0]),
			})
		}
	}

	// var bitmap win.BITMAP
	// win.GetObject(
	// 	win.HGDIOBJ(iconInfo.HbmColor),
	// 	uintptr(unsafe.Sizeof(bitmap)),
	// 	unsafe.Pointer(&bitmap),
	// )

	// fmt.Printf("%#+v\n", bitmap)
	//
	// buf := make([]byte, 91)
	// fmt.Println(GetBitmapBits(iconInfo.HbmColor, 90, buf))
	// fmt.Println(buf)

	//win.GlobalLock(win.HGLOBAL(iconInfo.HbmColor))
	//var pixels = (*[1 << 30]byte)(unsafe.Pointer(&iconInfo.HbmColor))[0:6]
	// pixels := (*[20]byte)(unsafe.Pointer(&shfi.HIcon))
	// fmt.Printf("%#+v\n", pixels)
	//
	// fmt.Printf("%#+v\n", bitmap)

	//win.GlobalUnlock(win.HGLOBAL(iconInfo.HbmColor))

	// var bitmap win.GpBitmap
	// win.GdipCreateBitmapFromHBITMAP(iconInfo.HbmColor, hpal win.HPALETTE, bitmap **win.GpBitmap)

	//
	//
	//
	// p := *(**int)(unsafe.Pointer(&iconInfo.HbmColor))
	// fmt.Printf("%#+v\n", p[0])
	//
	//
	//
	// fmt.Printf("%#+v\n", iconInfo)

	return rgba, nil
}
