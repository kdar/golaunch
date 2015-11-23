package idletime

//go:generate go run $GOROOT/src/syscall/mksyscall_windows.go -output zsyscall_windows.go syscall_windows.go

//sys   getLastInputInfo(plii uintptr) (b bool, err error) [failretval==false] = user32.GetLastInputInfo
//sys   getTickCount() (i uint32) [failretval==false] = kernel32.GetTickCount
