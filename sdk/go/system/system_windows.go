package system

import (
	"bytes"
	"encoding/base64"
	"golaunch/sdk/go/extracticon"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unicode/utf16"
	"runtime"

	"github.com/mattn/go-ole"
	"github.com/mattn/go-ole/oleutil"
	"golang.org/x/sys/windows/registry"
)

// Possibly get filename description from exe
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms647003(v=vs.85).aspx

func init() {
	ole.CoInitializeEx(0, COINIT_SPEED_OVER_MEMORY)
}

type System struct {
	oleShellObject *ole.IUnknown
	wshell         *ole.IDispatch
	extract        *extracticon.Extract
}

func NewSystem() *System {
	//ole.CoInitializeEx(0, COINIT_APARTMENTTHREADED | COINIT_SPEED_OVER_MEMORY)
	oleShellObject, _ := oleutil.CreateObject("WScript.Shell")
	wshell, _ := oleShellObject.QueryInterface(ole.IID_IDispatch)
	return &System{
		oleShellObject: oleShellObject,
		wshell:         wshell,
		extract:        extracticon.New(),
	}
}

func (s *System) Close() {
	s.wshell.Release()
}

func (s *System) AppIcon(path string) (image.Image, error) {
	if filepath.Ext(path) == ".lnk" {
		path = s.ResolveLink(path)
	}

	icon, err := s.extract.From(path)
	if err != nil {
		icon, err = s.extract.FromExt(filepath.Ext(path))
		if err != nil {
			return nil, err
		}
	}

	return icon, nil
}

func (s *System) EmbeddedAppIcon(path string) (string, error) {
	icon, err := s.AppIcon(path)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer([]byte("data:image/png;base64,"))

	b64encoder := base64.NewEncoder(base64.StdEncoding, buf)
	err = png.Encode(b64encoder, icon)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (s *System) ResolveLink(path string) string {
	cs, err := oleutil.CallMethod(s.wshell, "CreateShortcut", path)
	if err != nil {
		if v, ok := err.(*ole.OleError); ok {
			log.Println(v.SubError())
		}
		return path
	}

	idispatch := cs.ToIDispatch()

	// FIXME: this gives us the icon id too, but we currently don't use it in
	// the extracticon api.
	iconLocation := strings.Split(oleutil.MustGetProperty(idispatch, "IconLocation").ToString(), ",")
	if len(iconLocation) > 0 && strings.Contains(iconLocation[0], ".") {
		return iconLocation[0]
	}

	return oleutil.MustGetProperty(idispatch, "TargetPath").ToString()
}

// ResolveMSILink finds the target for what are called Advertisment
// Shortcuts. These are special shortcuts installed with windows
// installer MSI that don't follow the conventions of regular shortcuts.
// I had to do some funky things with CoInitialize because this code
// seems to only work with apartment threads, while the rest of the
// code in here fails with apartment threads.
func (s *System) ResolveMSILink(path string) string {
	ch := make(chan string)
  go func() {
		// lock the thread so we can use apartment threads
		runtime.LockOSThread()
	  ole.CoInitializeEx(0, COINIT_APARTMENTTHREADED | COINIT_SPEED_OVER_MEMORY)
		defer ole.CoUninitialize()

		path, _ = filepath.Abs(path)

		product := make([]uint16, 39)
		feature := make([]uint16, 39)
		component := make([]uint16, 39)
		_, err := MsiGetShortcutTarget(syscall.StringToUTF16Ptr(path), &product[0], &feature[0], &component[0])
		if err != nil {
			log.Println(err)
			ch <- ""
			return
		}

		dwlen := syscall.MAX_PATH
		target := make([]uint16, syscall.MAX_PATH)
		_, err = MsiGetComponentPath(&product[0], &component[0], &target[0], &dwlen)
		if err != nil {
			log.Println(err)
			ch <- ""
			return
		}

		ch <- string(utf16.Decode(target[:dwlen]))
	}()

	return <-ch
}

func (s *System) RunProgram(path string, args string, dir string, user string) error {
	if filepath.Ext(path) == ".lnk" {
		lcs := oleutil.MustCallMethod(s.wshell, "CreateShortcut", path).ToIDispatch()

		// can possible call the method Count and Item to retrieve individual arg items
		// https://msdn.microsoft.com/en-us/library/ss1ysb2a(v=vs.84).aspx
		largs := oleutil.MustGetProperty(lcs, "Arguments").ToString()
		lwd := oleutil.MustGetProperty(lcs, "WorkingDirectory").ToString()

    // attempt to see if it's a MSI shortcut first
		path = s.ResolveMSILink(path)
		if path == "" {
			// just try to get the target path if not a MSI shortcut
			path = oleutil.MustGetProperty(lcs, "TargetPath").ToString()
			if path == "" {
				// I'm not sure if this is the correct thing to do, but games
				// like Minesweeper aren't MSI shortcuts and they have no
				// TargetPath. The only thing I can figure out is the exe
				// is located in the IconLocation.
				path = oleutil.MustGetProperty(lcs, "IconLocation").ToString()
				index := strings.LastIndex(path, ",")
				if index >= 0 {
					path = path[:index]
				}
			}
		}

		lwdstat, err := os.Stat(lwd)
		if dir == "" && !os.IsNotExist(err) && lwdstat.IsDir() {
			dir = lwd
		}

		if args == "" {
			args = largs
		}
	}

	action := ""
	if user == "administrator" {
		action = "runas"
	}

	if dir == "" {
		dir = filepath.Dir(path)
	}

	// resolve environment variables
	path, _ = registry.ExpandString(path)
	dir, _ = registry.ExpandString(dir)

	return ShellExecute(action, path, args, dir)

	// if filepath.Ext(path) == ".lnk" {
	// 	cs := oleutil.MustCallMethod(s.wshell, "CreateShortcut", path).ToIDispatch()
	// 	path = oleutil.MustGetProperty(cs, "TargetPath").ToString()
	// 	// can possible call the method Count and Item to retrieve individual arg items
	// 	// https://msdn.microsoft.com/en-us/library/ss1ysb2a(v=vs.84).aspx
	// 	args := oleutil.MustGetProperty(cs, "Arguments").ToString()
	// 	wd := oleutil.MustGetProperty(cs, "WorkingDirectory").ToString()
	// 	cmd := exec.Command(path)
	//
	// 	wdstat, err := os.Stat(wd)
	// 	if !os.IsNotExist(err) && wdstat.IsDir() {
	// 		cmd.Dir = wd
	// 	}
	//
	// 	if len(args) > 0 {
	// 		args, _ := shellwords.Parse("cmd " + args)
	// 		cmd.Args = args
	// 	}
	//
	// 	return cmd.Start()
	// }
	//
	// cmd := exec.Command(path)
	// return cmd.Start()
}

func (s *System) OpenFolder(path string) error {
	return ShellExecute("open", path, "", "")
}
