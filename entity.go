package entitas

import (
	"errors"
	"fmt"
)

var (
	ErrComponentExists       = errors.New("component exists")
	ErrComponentDoesNotExist = errors.New("component does not exist")
)

type EntityComponentChanged func(Entity, Component)

type EntityID uint64

type Entity interface {
	ID() EntityID

	CreateComponent(ts int) Component
	AddComponent(cs ...Component) error
	UpdateComponent(cs ...Component)
	RemoveComponent(ts ...int) error
	RemoveAllComponents()

	HasComponent(ts ...int) bool
	HasAnyComponent(ts ...int) bool
	Component(t int) (Component, error)
	Components() []Component
	ComponentTypes() []int

	AddEvent(ev EventType, action EntityComponentChanged)
	RemoveAllEvents()
	HasEvents() bool

	Destroy()
	internalDestroy()
}

type entity struct {
	id               EntityID
	components       []Component
	componentChanged map[EventType][]EntityComponentChanged

	componentsCache     []Component
	componentTypesCache []int

	context Context
}

func newEntity(context Context, id EntityID) Entity {
	return &entity{
		id:               id,
		components:       make([]Component, TotalComponents),
		componentChanged: make(map[EventType][]EntityComponentChanged),
		context:          context,
	}
}

//private
func (e *entity) onComponentChanged(ev EventType, c Component) {
	if actions, ok := e.componentChanged[ev]; ok {
		for _, action := range actions {
			action(e, c)
		}
	}
}

//public
func (e *entity) CreateComponent(ts int) Component {
	return e.context.CreateComponent(ts)
}

func (e *entity) HasComponent(ts ...int) bool {
	for _, t := range ts {
		if e.components[t] == nil {
			return false
		}
	}
	return true
}

func (e *entity) HasAnyComponent(ts ...int) bool {
	for _, t := range ts {
		if e.components[t] != nil {
			return true
		}
	}
	return false
}

func (e *entity) Component(t int) (Component, error) {
	c := e.components[t]
	if c == nil {
		return nil, ErrComponentDoesNotExist
	}
	return c, nil
}

func (e *entity) Components() []Component {
	components := e.componentsCache
	if components == nil {
		components = make([]Component, 0, len(e.components))

		for _, c := range e.components {
			components = append(components, c)
		}

		e.componentsCache = components
	}
	return components
}

func (e *entity) ComponentTypes() []int {
	types := e.componentTypesCache
	if types == nil {
		types = make([]int, 0, len(e.components))
		for t, c := range e.components {
			if c != nil {
				types = append(types, int(t))
			}
		}
		e.componentTypesCache = types
	}
	return types
}

func (e *entity) AddComponent(cs ...Component) error {
	for _, c := range cs {
		t := c.Type()
		if e.HasComponent(t) {
			return ErrComponentExists
		}
		e.components[t] = c
		e.onComponentChanged(EventAdded, c)
	}

	if len(cs) > 0 {
		e.componentsCache = nil
		e.componentTypesCache = nil
	}

	return nil
}

func (e *entity) UpdateComponent(cs ...Component) {
	for _, c := range cs {
		t := c.Type()
		old := e.components[t]
		e.components[t] = c
		if old != nil {
			if old != c {
				e.onComponentChanged(EventRemoved, old)
			}
			e.onComponentChanged(EventUpdated, c)
		} else {
			e.onComponentChanged(EventAdded, c)
		}
	}

	if len(cs) > 0 {
		e.componentsCache = nil
		e.componentTypesCache = nil
	}
}

func (e *entity) RemoveComponent(ts ...int) error {
	for _, t := range ts {
		c, err := e.Component(t)
		if err != nil {
			return err
		}
		e.components[t] = nil
		e.onComponentChanged(EventRemoved, c)
	}

	if len(ts) > 0 {
		e.componentsCache = nil
		e.componentTypesCache = nil
	}

	return nil
}

func (e *entity) RemoveAllComponents() {
	components := e.components

	e.components = make([]Component, TotalComponents)
	e.componentsCache = nil
	e.componentTypesCache = nil

	for _, c := range components {
		if c != nil {
			e.onComponentChanged(EventRemoved, c)
		}
	}

}

func (e *entity) ID() EntityID {
	return e.id
}

func (e *entity) AddEvent(ev EventType, action EntityComponentChanged) {
	actions := e.componentChanged[ev]
	e.componentChanged[ev] = append(actions, action)
}

func (e *entity) HasEvents() bool {
	return len(e.componentChanged) > 0
}

func (e *entity) RemoveAllEvents() {
	e.componentChanged = make(map[EventType][]EntityComponentChanged)
}

func (e *entity) Destroy() {
	e.context.destroyEntity(e)
}

func (e *entity) internalDestroy() {
	e.RemoveAllComponents()
	e.RemoveAllEvents()
}

func (e *entity) String() string {
	return fmt.Sprintf("Entity_%d(types %v)", e.id, e.ComponentTypes())
}
