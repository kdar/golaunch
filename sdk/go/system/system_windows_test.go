package system

import (
	"fmt"
	"sync"
	"testing"

	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

func TestResolveLink(t *testing.T) {
	a := assertions.New(t)

	sys := NewSystem()
	path := sys.ResolveLink("testdata/Pidgin.lnk")

	a.So(path, should.Equal, "D:\\Program Files (x86)\\Pidgin\\pidgin.exe")
}

func TestResolveMSILink(t *testing.T) {
	//a := assertions.New(t)

	sys := NewSystem()
	path := sys.ResolveMSILink("testdata\\Minesweeper.lnk")
	fmt.Println(path)
}

func TestResolveLinkConcurrent(t *testing.T) {
	a := assertions.New(t)

	sys := NewSystem()

	var wg sync.WaitGroup
	for x := 0; x < 10; x++ {
		wg.Add(1)
		go func() {
			path := sys.ResolveLink("testdata/Pidgin.lnk")
			a.So(path, should.Equal, "D:\\Program Files (x86)\\Pidgin\\pidgin.exe")
			wg.Done()
		}()
	}

	wg.Wait()
}

// func TestGetAppIcon(t *testing.T) {
// 	a := assertions.New(t)
//
// 	sys := NewSystem()
// 	img, err := sys.GetAppIcon("testdata/Pidgin.lnk")
//
// 	a.So(img, should.NotBeEmpty)
// 	a.So(err, should.BeNil)
// }
