package entitas

import (
	. "github.com/smartystreets/goconvey/convey"
	"sort"
	"testing"
)

const (
	ComponentA int = iota
	ComponentB
	ComponentC
	ComponentD
	ComponentE
	ComponentF
	NumComponents
)

// ------- A

type componentA struct {
	value int
}

func NewComponentA(value int) Component {
	return &componentA{value}
}

func (c1 *componentA) int() int {
	return ComponentA
}

func (c1 *componentA) String() string {
	return "A"
}

// ------- B

type componentB struct {
	value float32
}

func NewComponentB(value float32) Component {
	return &componentB{value}
}

func (c1 *componentB) int() int {
	return ComponentB
}

func (c1 *componentB) String() string {
	return "B"
}

// ------- C

type componentC struct {
}

func NewComponentC() Component {
	return &componentC{}
}

func (c1 *componentC) int() int {
	return ComponentC
}

func (c1 *componentC) String() string {
	return "C"
}

// ------- D

type componentD struct {
}

func NewComponentD() Component {
	return &componentC{}
}

func (c1 *componentD) int() int {
	return ComponentD
}

func (c1 *componentD) String() string {
	return "D"
}

// ------- E

type componentE struct {
}

func NewComponentE() Component {
	return &componentE{}
}

func (c1 *componentE) int() int {
	return ComponentE
}

func (c1 *componentE) String() string {
	return "E"
}

// ------- F

type componentF struct {
}

func NewComponentF() Component {
	return &componentF{}
}

func (c1 *componentF) int() int {
	return ComponentF
}

func (c1 *componentF) String() string {
	return "F"
}

func TestComponentSorting(t *testing.T) {
	Convey("Given components and a component list", t, func() {
		c1 := NewComponentA(1)
		c2 := NewComponentB(0.0)

		components := []Component{c2, c1}

		Convey("It should be sortable by type", func() {
			sort.Sort(Components(components))
			So(components, ShouldResemble, []Component{c1, c2})
		})
	})
}
