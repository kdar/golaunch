package idletime

import (
	"testing"

	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

func TestGet(t *testing.T) {
	a := assertions.New(t)

	i, err := Get()
	a.So(err, should.BeNil)
	a.So(i, should.BeGreaterThanOrEqualTo, 0)
}
