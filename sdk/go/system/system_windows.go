package system

import (
	"bytes"
	"encoding/base64"
	"golaunch/sdk/go/extracticon"
	"image/png"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mattn/go-ole"
	"github.com/mattn/go-ole/oleutil"
	"github.com/mattn/go-shellwords"
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

func (s *System) GetAppIcon(path string) (string, error) {
	if filepath.Ext(path) == ".lnk" {
		path = s.GetRealImagePath(path)
	}

	icon, err := s.extract.From(path)
	if err != nil {
		icon, err = s.extract.FromExt(filepath.Ext(path))
		if err != nil {
			return "", err
		}
	}

	buf := bytes.NewBuffer([]byte("data:image/png;base64,"))

	b64encoder := base64.NewEncoder(base64.StdEncoding, buf)
	err = png.Encode(b64encoder, icon)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (s *System) GetRealImagePath(path string) string {
	cs, err := oleutil.CallMethod(s.wshell, "CreateShortcut", path)
	if err != nil {
		if v, ok := err.(*ole.OleError); ok {
			log.Println(v.SubError())
		}
		return path
	}

	// FIXME: this gives us the icon id too, but we currently don't use it in
	// the extracticon api.
	iconLocation := strings.Split(oleutil.MustGetProperty(cs.ToIDispatch(), "IconLocation").ToString(), ",")
	if len(iconLocation) > 0 && strings.Contains(iconLocation[0], ".") {
		return iconLocation[0]
	}

	return oleutil.MustGetProperty(cs.ToIDispatch(), "TargetPath").ToString()
}

func (s *System) RunProgram(path string) error {
	if filepath.Ext(path) == ".lnk" {
		cs := oleutil.MustCallMethod(s.wshell, "CreateShortcut", path).ToIDispatch()
		path = oleutil.MustGetProperty(cs, "TargetPath").ToString()
		// can possible call the method Count and Item to retrieve individual arg items
		// https://msdn.microsoft.com/en-us/library/ss1ysb2a(v=vs.84).aspx
		args := oleutil.MustGetProperty(cs, "Arguments").ToString()
		wd := oleutil.MustGetProperty(cs, "WorkingDirectory").ToString()
		cmd := exec.Command(path)

		wdstat, err := os.Stat(wd)
		if !os.IsNotExist(err) && wdstat.IsDir() {
			cmd.Dir = wd
		}

		if len(args) > 0 {
			args, _ := shellwords.Parse("cmd " + args)
			cmd.Args = args
		}

		return cmd.Start()
	}

	cmd := exec.Command(path)
	return cmd.Start()
}
