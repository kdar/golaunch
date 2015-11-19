package system

import (
	"sync"
	"testing"

	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

func TestResolveLink(t *testing.T) {
	a := assertions.New(t)

	sys := NewSystem()
	path := sys.ResolveLink("testdata/Start Tor Browser.lnk")

	a.So(path, should.Equal, "D:\\Program Files\\Tor Browser\\Browser\\firefox.exe")
}

func TestResolveLinkConcurrent(t *testing.T) {
	a := assertions.New(t)

	sys := NewSystem()

	var wg sync.WaitGroup
	for x := 0; x < 10; x++ {
		wg.Add(1)
		go func() {
			path := sys.ResolveLink("testdata/Start Tor Browser.lnk")
			a.So(path, should.Equal, "D:\\Program Files\\Tor Browser\\Browser\\firefox.exe")
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestGetAppIcon(t *testing.T) {
	a := assertions.New(t)

	sys := NewSystem()
	img, err := sys.GetAppIcon("testdata/Smite.lnk")

	a.So(img, should.NotBeEmpty)
	a.So(err, should.BeNil)
}
