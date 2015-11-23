package idletime

import (
	"time"
	"unsafe"
)

type lastInputInfo struct {
	cbSize uint32
	dwTime uint32
}

/// https://github.com/darwin/chromium-src-chrome-browser/blob/master/idle_win.cc
func Get() (time.Duration, error) {
	var lii lastInputInfo
	lii.cbSize = uint32(unsafe.Sizeof(lii))

	currentIdleTime := uint32(0)
	success, err := getLastInputInfo(uintptr(unsafe.Pointer(&lii)))
	if err != nil {
		return 0, err
	}

	if success {
		now := getTickCount()
		if now < lii.dwTime {
			// GetTickCount() wraps around every 49.7 days -- assume it wrapped just
			// once.
			kMaxDWORD := ^uint32(0)
			timeBeforeWrap := kMaxDWORD - lii.dwTime
			timeAfterWrap := now
			// The sum is always smaller than kMaxDWORD.
			currentIdleTime = timeBeforeWrap + timeAfterWrap
		} else {
			currentIdleTime = now - lii.dwTime
		}
	}

	return time.Duration(currentIdleTime) * time.Millisecond, nil
}
