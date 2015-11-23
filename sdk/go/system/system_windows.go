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

	"github.com/mattn/go-ole"
	"github.com/mattn/go-ole/oleutil"
)

// Possibly get filename description from exe
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms647003(v=vs.85).aspx

type System struct {
	oleShellObject *ole.IUnknown
	wshell         *ole.IDispatch
	extract        *extracticon.Extract
}

func NewSystem() *System {
	ole.CoInitializeEx(0, 0)
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
	ole.CoUninitialize()
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

func (s *System) RunProgram(path string, args string, dir string, user string) error {
	if filepath.Ext(path) == ".lnk" {
		lcs := oleutil.MustCallMethod(s.wshell, "CreateShortcut", path).ToIDispatch()
		lpath := oleutil.MustGetProperty(lcs, "TargetPath").ToString()
		// can possible call the method Count and Item to retrieve individual arg items
		// https://msdn.microsoft.com/en-us/library/ss1ysb2a(v=vs.84).aspx
		largs := oleutil.MustGetProperty(lcs, "Arguments").ToString()
		lwd := oleutil.MustGetProperty(lcs, "WorkingDirectory").ToString()

		path = lpath

		lwdstat, err := os.Stat(lwd)
		if dir == "" && !os.IsNotExist(err) && lwdstat.IsDir() {
			dir = lwd
		}

		if args == "" && len(largs) > 0 {
			args = largs
		}
	}

	action := ""
	if user == "administrator" {
		action = "runas"
	}

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
