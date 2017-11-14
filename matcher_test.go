package entitas

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMatcher(t *testing.T) {
	Convey("Subject: Matchers", t, func() {
		Convey("Given a few entities", func() {
			all1 := AllOf(int(1))
			all2 := AllOf(int(1))

			So(all1.Hash() == all2.Hash(), ShouldBeTrue)

		})

		Convey("", func() {
			any1 := AnyOf(int(1), int(2), int(3))
			any2 := AnyOf(int(2), int(1), int(3), int(3))
			So(any1.Equals(any2), ShouldBeTrue)
		})
	})
}
