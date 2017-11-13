package entitas

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMatcher(t *testing.T) {
	Convey("Subject: Matchers", t, func() {
		Convey("Given a few entities", func() {
			all1 := AllOf(ComponentType(1))
			all2 := AllOf(ComponentType(1))

			So(all1.Hash() == all2.Hash(), ShouldBeTrue)

		})

		Convey("", func() {
			any1 := AnyOf(ComponentType(1), ComponentType(2), ComponentType(3))
			any2 := AnyOf(ComponentType(2), ComponentType(1), ComponentType(3), ComponentType(3))
			So(any1.Equals(any2), ShouldBeTrue)
		})
	})
}
