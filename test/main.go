package main

import (
	"fmt"
	"github.com/jangsky215/go-entitas"
	"math"
)

type Component1 struct {
	AA int
}

func (c *Component1) Type() int {
	return 1
}

func NewC1() entitas.Component {
	return &Component1{99}
}

func main() {
	fmt.Println(math.MaxUint32 / 30000)

	entitas.TotalComponents = 99
	pool := entitas.NewPool(0)
	fmt.Println(pool)

	matcher := entitas.AllOf(1)
	group := pool.Group(matcher)
	observer := entitas.NewGroupObserver(group, entitas.EventAddedOrRemoved)

	pool.RegisterComponent(&Component1{})

	c1 := pool.CreateComponent(1)
	entity1 := pool.CreateEntity(c1)
	fmt.Println(entity1)
	fmt.Println(pool)

	c, err := entity1.Component(1)
	if err == nil {
		cc := c.(*Component1)
		fmt.Println((cc).AA)
	}
	fmt.Println(observer.CollectedEntities(), len(observer.CollectedEntities()))

	entity2 := pool.CreateEntity()
	fmt.Println(entity2)

	fmt.Println(observer.CollectedEntities(), len(observer.CollectedEntities()))

	observer.ClearCollectedEntities()
	entity2.Destroy()

	fmt.Println(observer.CollectedEntities())

}
