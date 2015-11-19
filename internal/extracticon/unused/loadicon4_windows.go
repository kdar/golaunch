package extracticon

import (
	"image"
	"image/color"
	"syscall"
	"unsafe"

	"github.com/ansel1/merry"
	"github.com/lxn/win"
)

type GUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

var (
	// System IImageList object:
	iidImageList *GUID = &GUID{0x46EB5926, 0x582E, 0x4017, [8]byte{0x9F, 0xDF, 0xE8, 0x99, 0x8D, 0xAA, 0x09, 0x50}}
)

func From4(path string) (image.Image, error) {
	var shfi SHFILEINFO

	_, err := shGetFileInfo(
		syscall.StringToUTF16Ptr(path),
		FILE_ATTRIBUTE_NORMAL,
		uintptr(unsafe.Pointer(&shfi)),
		uint32(unsafe.Sizeof(shfi)),
		SHGFI_SYSICONINDEX)
	if err != nil {
		return nil, merry.Wrap(err)
	}

	var himl *uintptr
	_, err = shGetImageList(SHIL_JUMBO, iidImageList, &himl)
	if err != nil {
		return nil, merry.Wrap(err)
	}

	hIcon, err := imageList_GetIcon(uintptr(unsafe.Pointer(himl)), int(shfi.IIcon), ILD_TRANSPARENT)
	if err != nil {
		return nil, merry.Wrap(err)
	}

	var iconInfo win.ICONINFO
	_, err = getIconInfo(Handle(hIcon), uintptr(unsafe.Pointer(&iconInfo)))
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

	if _, err := drawIconEx(hdc, 0, 0, Handle(hIcon), w, h, 0, 0, win.DI_NORMAL); err != nil {
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

	return rgba, nil
}
