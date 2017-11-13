package entitas

import (
	"fmt"
	"reflect"
	"sort"
)

const (
	componentHashFactor uint = 647
	allHashFactor            = 653
	anyHashFactor            = 659
	noneHashFactor           = 661
	arrayHashFactor          = 673
)

type Matcher interface {
	Matches(entity Entity) bool
	Hash() uint
	Types() []Type
	Equals(m Matcher) bool
	String() string
}

// BaseMatcher
type BaseMatcher struct {
	types []Type
	hash  uint
}

func newBaseMatcher(types ...Type) BaseMatcher {
	mtype := make(map[Type]bool)
	for _, t := range types {
		mtype[t] = true
	}

	types = make([]Type, 0, len(mtype))
	for t := range mtype {
		types = append(types, t)
	}
	sort.Sort(Types(types))

	return BaseMatcher{types: types}
}

func (b *BaseMatcher) Hash() uint {
	return b.hash
}

func (b *BaseMatcher) Types() []Type {
	return b.types
}

func (a *BaseMatcher) Equals(m Matcher) bool {
	return reflect.DeepEqual(a.Types(), m.Types())
}

// AllOf
type AllMatcher struct {
	BaseMatcher
}

func AllOf(types ...Type) Matcher {
	b := newBaseMatcher(types...)
	b.hash = Hash(allHashFactor, b.Types()...)
	return &AllMatcher{b}
}

func (a *AllMatcher) Matches(e Entity) bool {
	return e.HasComponent(a.Types()...)
}

func (a *AllMatcher) String() string {
	return fmt.Sprintf("AllOf(%v)", print(a.Types()...))
}

// AnyOf
type AnyMatcher struct {
	BaseMatcher
}

func AnyOf(types ...Type) Matcher {
	b := newBaseMatcher(types...)
	b.hash = Hash(anyHashFactor, b.Types()...)
	return &AnyMatcher{b}
}

func (a *AnyMatcher) Matches(e Entity) bool {
	return e.HasAnyComponent(a.Types()...)
}

func (a *AnyMatcher) String() string {
	return fmt.Sprintf("AnyOf(%v)", print(a.Types()...))
}

// NonoOf
type NoneMatcher struct {
	BaseMatcher
}

func NoneOf(types ...Type) Matcher {
	b := newBaseMatcher(types...)
	b.hash = Hash(noneHashFactor, b.Types()...)
	return &NoneMatcher{b}
}

func (n *NoneMatcher) Matches(e Entity) bool {
	return !e.HasAnyComponent(n.Types()...)
}

func (n *NoneMatcher) String() string {
	return fmt.Sprintf("NonoOf(%v)", print(n.Types()...))
}

// Utilities
func Hash(factor uint, types ...Type) uint {
	var hash uint
	for _, t := range types {
		hash ^= uint(t) * componentHashFactor
	}
	hash ^= uint(len(types)) * factor
	return uint(hash)
}

func HashMatcher(matchers ...Matcher) uint {
	if len(matchers) == 1 {
		return matchers[0].Hash()
	}

	hash := uint(0)
	for _, m := range matchers {
		hash ^= m.Hash()
	}
	hash ^= uint(len(matchers)) * arrayHashFactor
	return hash
}

func print(types ...Type) string {
	return fmt.Sprintf("%v", types)
}
