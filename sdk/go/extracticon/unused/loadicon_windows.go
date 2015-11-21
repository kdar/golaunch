package extracticon

import (
	"errors"
	"image"
	"image/color"
	"syscall"
	"unsafe"

	"github.com/ansel1/merry"
	"github.com/lxn/win"
)

var (
	enumResourceNamesCb = syscall.NewCallback(func(hModule Handle, lpszType *uint16, lpszName *uint16, lParam uintptr) uintptr {
		*((*uintptr)(unsafe.Pointer(lParam))) = uintptr(unsafe.Pointer(lpszName))
		return 0
	})
)

func From(path string, sizes []int) ([]image.Image, error) {
	var images []image.Image

	sptr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return nil, merry.Wrap(err)
	}

	hLib, err := loadLibraryEx(sptr, 0, LOAD_LIBRARY_AS_DATAFILE)
	if err != nil {
		return nil, merry.Wrap(err)
	}

	var groupName uintptr
	enumResourceNames(hLib, RT_GROUP_ICON, enumResourceNamesCb, uintptr(unsafe.Pointer(&groupName)))

	hRes := loadResource(hLib, findResourceEx(hLib, RT_GROUP_ICON, groupName, 0))
	memIconDir := lockResource(hRes)

	for _, iconSize := range sizes {
		iconName := lookupIconIdFromDirectoryEx(memIconDir, true, iconSize, iconSize, 0x00000000)
		hResInfo := findResourceEx(hLib, RT_ICON, uintptr(iconName), 0)
		size := sizeofResource(hLib, hResInfo)
		rec := loadResource(hLib, hResInfo)
		memIcon := lockResource(rec)
		if memIcon == 0 {
			return nil, errors.New("could not lock resource")
		}

		data := (*[1 << 30]byte)(unsafe.Pointer(memIcon))[0:size]

		hIconRet, _ := createIconFromResourceEx(&data[0], size, true, 0x00030000, 0, 0, 0x00000000)
		var iconInfo win.ICONINFO
		getIconInfo(hIconRet, uintptr(unsafe.Pointer(&iconInfo)))

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

		screenDevice := win.GetDC(0)
		hdc := win.CreateCompatibleDC(screenDevice)
		win.ReleaseDC(0, screenDevice)

		var bits unsafe.Pointer
		winBitmap := win.CreateDIBSection(hdc, &bitmapInfo, win.DIB_RGB_COLORS, &bits, 0, 0)

		var pixels = (*[1 << 30]byte)(bits)[0:bitmapInfo.BiSizeImage]

		win.SelectObject(hdc, win.HGDIOBJ(winBitmap))
		win.DrawIconEx(hdc, 0, 0, win.HICON(hIconRet), w, h, 0, 0, win.DI_NORMAL)

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
					A: pixels[((y*w+x)*4)+3],
					R: uint8(pixels[((y*w+x)*4)+2]),
					G: uint8(pixels[((y*w+x)*4)+1]),
					B: uint8(pixels[((y*w+x)*4)+0]),
				})
			}
		}

		images = append(images, rgba)

		win.DestroyIcon(win.HICON(hIconRet))
	}

	syscall.FreeLibrary(syscall.Handle(hLib))

	return images, nil
}
