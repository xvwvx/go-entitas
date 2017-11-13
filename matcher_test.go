package entitas

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMatcher(t *testing.T) {
	Convey("Subject: Matchers", t, func() {
		Convey("Given a few entities", func() {
			all1 := AllOf(Type(1))
			all2 := AllOf(Type(1))

			So(all1.Hash() == all2.Hash(), ShouldBeTrue)

		})

		Convey("", func() {
			any1 := AnyOf(Type(1), Type(2), Type(3))
			any2 := AnyOf(Type(2), Type(1), Type(3), Type(3))
			So(any1.Equals(any2), ShouldBeTrue)
		})
	})
}
