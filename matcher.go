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
	ComponentTypes() []ComponentType
	Equals(m Matcher) bool
	String() string
}

// baseMatcher
type baseMatcher struct {
	types []ComponentType
	hash  uint
}

func newBaseMatcher(types ...ComponentType) baseMatcher {
	mtype := make(map[ComponentType]bool)
	for _, t := range types {
		mtype[t] = true
	}

	types = make([]ComponentType, 0, len(mtype))
	for t := range mtype {
		types = append(types, t)
	}
	sort.Sort(ComponentTypes(types))

	return baseMatcher{types: types}
}

func (b *baseMatcher) Hash() uint {
	return b.hash
}

func (b *baseMatcher) ComponentTypes() []ComponentType {
	return b.types
}

// AllOf
type AllMatcher struct {
	baseMatcher
}

func AllOf(types ...ComponentType) Matcher {
	b := newBaseMatcher(types...)
	b.hash = Hash(allHashFactor, b.ComponentTypes()...)
	return &AllMatcher{b}
}

func (a *AllMatcher) Matches(e Entity) bool {
	return e.HasComponent(a.ComponentTypes()...)
}

func (a *AllMatcher) String() string {
	return fmt.Sprintf("AllOf(%v)", print(a.ComponentTypes()...))
}

func (a *baseMatcher) Equals(m Matcher) bool {
	return reflect.DeepEqual(a, m)
}

// AnyOf
type AnyMatcher struct {
	baseMatcher
}

func AnyOf(types ...ComponentType) Matcher {
	b := newBaseMatcher(types...)
	b.hash = Hash(anyHashFactor, b.ComponentTypes()...)
	return &AnyMatcher{b}
}

func (a *AnyMatcher) Matches(e Entity) bool {
	return e.HasAnyComponent(a.ComponentTypes()...)
}

func (a *AnyMatcher) String() string {
	return fmt.Sprintf("AnyOf(%v)", print(a.ComponentTypes()...))
}

func (a *AnyMatcher) Equals(m Matcher) bool {
	return reflect.DeepEqual(a, m)
}

// NonoOf
type NoneMatcher struct {
	baseMatcher
}

func NoneOf(types ...ComponentType) Matcher {
	b := newBaseMatcher(types...)
	b.hash = Hash(noneHashFactor, b.ComponentTypes()...)
	return &NoneMatcher{b}
}

func (n *NoneMatcher) Matches(e Entity) bool {
	return !e.HasAnyComponent(n.ComponentTypes()...)
}

func (n *NoneMatcher) String() string {
	return fmt.Sprintf("NonoOf(%v)", print(n.ComponentTypes()...))
}

func (n *NoneMatcher) Equals(m Matcher) bool {
	return reflect.DeepEqual(n, m)
}

// Utilities
func Hash(factor uint, types ...ComponentType) uint {
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

func print(types ...ComponentType) string {
	return fmt.Sprintf("%v", types)
}
