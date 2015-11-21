package extracticon

// possible use PrivateExtractIcons to extract higher quality icons

import (
	"image"
	"image/color"
	"syscall"
	"unsafe"

	"github.com/ansel1/merry"
	"github.com/lxn/win"
)

type Extract struct {
	cache map[string]image.Image
}

func New() *Extract {
	return &Extract{
		cache: make(map[string]image.Image),
	}
}

func (e *Extract) FromExt(ext string) (image.Image, error) {
	if v, ok := e.cache[ext]; ok {
		return v, nil
	}

	var shfi SHFILEINFO

	_, err := shGetFileInfo(
		syscall.StringToUTF16Ptr(ext),
		0,
		uintptr(unsafe.Pointer(&shfi)),
		uint32(unsafe.Sizeof(shfi)),
		SHGFI_ICON|SHGFI_USEFILEATTRIBUTES)
	if err != nil {
		return nil, merry.Wrap(err)
	}

	img, err := e.handleToImage(shfi.HIcon)
	if err != nil {
		return nil, err
	}

	e.cache[ext] = img
	return img, nil
}

func (e *Extract) From(path string) (image.Image, error) {
	large := []Handle{0}
	_, err := extractIconEx(syscall.StringToUTF16Ptr(path), 0, &large[0], nil, 1)
	if err != nil {
		return nil, err
	}

	img, err := e.handleToImage(Handle(large[0]))
	if err != nil {
		return nil, err
	}

	destroyIcon(large[0])
	return img, nil

	// var lpiIcon uint16 = 0
	// icon, err := extractAssociatedIcon(0, syscall.StringToUTF16Ptr(path), &lpiIcon)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// img, err := e.handleToImage(Handle(icon))
	// if err != nil {
	// 	return nil, err
	// }
	//
	// destroyIcon(icon)
	// return img, nil
}

// func (e *Extract) From2(path string) (image.Image, error) {
// 	sptr, err := syscall.UTF16PtrFromString(path)
// 	if err != nil {
// 		return nil, merry.Wrap(err)
// 	}
//
// 	hExe, err := loadLibraryEx(sptr, 0, LOAD_LIBRARY_AS_DATAFILE)
// 	if err != nil {
// 		return nil, merry.Wrap(err)
// 	}
//
// 	var hRes Handle
// 	var sizeIcon int
// 	for l := 1; l <= 500; l++ {
// 		hRes, err = findResource(hExe, uintptr(l), RT_ICON)
// 		if err == nil {
// 			sizeIcon = sizeofResource(hExe, hRes)
// 			if sizeIcon > 0 {
// 				break
// 			}
// 		}
// 	}
//
// 	if hRes == 0 {
// 		return nil, merry.Wrap(errors.New("could not find icon"))
// 	}
//
// 	hResLoad, err := loadResource(hExe, hRes)
// 	if err != nil {
// 		return nil, merry.Wrap(err)
// 	}
//
// 	lpResLock, err := lockResource(hResLoad)
// 	if err != nil {
// 		return nil, merry.Wrap(err)
// 	}
//
// 	data := (*[1 << 30]byte)(unsafe.Pointer(lpResLock))[0:sizeIcon]
// 	hIconRet, _ := createIconFromResourceEx(&data[0], sizeIcon, true, 0x00030000, 0, 0, 0x00000000)
//
// 	img, err := e.handleToImage(Handle(hIconRet))
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	destroyIcon(hIconRet)
// 	return img, nil
// }

func (e *Extract) handleToImage(handle Handle) (image.Image, error) {
	var iconInfo win.ICONINFO
	_, err := getIconInfo(handle, uintptr(unsafe.Pointer(&iconInfo)))
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

	if _, err := selectObject(hdc, winBitmap); err != nil {
		return nil, merry.Wrap(err)
	}

	if _, err := drawIconEx(hdc, 0, 0, Handle(handle), w, h, 0, 0, win.DI_NORMAL); err != nil {
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

	if _, err := drawIconEx(hdc, 0, 0, Handle(handle), w, h, 0, 0, win.DI_MASK); err != nil {
		return nil, merry.Wrap(err)
	}

	for y := int32(0); y < h; y++ {
		for x := int32(0); x < w; x++ {
			if pixels[((y*w+x)*4)+3] == 0 && pixels[((y*w+x)*4)+2] == 0 && pixels[((y*w+x)*4)+1] == 0 && pixels[((y*w+x)*4)+0] == 0 {
				tmp := rgba.RGBAAt(int(x), int(y))
				tmp.A = 0xFF
				rgba.SetRGBA(int(x), int(y), tmp)
			}
		}
	}

	return rgba, nil
}
