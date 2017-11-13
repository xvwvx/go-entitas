package entitas

import (
	"errors"
	"fmt"
	"sort"
)

var (
	ErrComponentExists       = errors.New("component exists")
	ErrComponentDoesNotExist = errors.New("component does not exist")
)

type EntityComponentChanged func(Entity, Component)

type Entity interface {
	ID() uint64

	AddComponent(cs ...Component) error
	UpdateComponent(cs ...Component)
	RemoveComponent(ts ...Type) error
	RemoveAllComponents()

	HasComponent(ts ...Type) bool
	HasAnyComponent(ts ...Type) bool
	Component(t Type) (Component, error)
	Components() []Component
	ComponentTypes() []Type

	AddEvent(ev EventType, action EntityComponentChanged)
	RemoveAllEvents()
	HasEvents() bool

	Destroy()
}

type entity struct {
	id               uint64
	components       []Component
	componentChanged map[EventType][]EntityComponentChanged

	componentsCache     []Component
	componentTypesCache []Type

	pool Pool
}

func newEntity(pool Pool, id uint64) Entity {
	return &entity{
		id:               id,
		components:       make([]Component, TotalComponents),
		componentChanged: make(map[EventType][]EntityComponentChanged),
		pool:             pool,
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
func (e *entity) HasComponent(ts ...Type) bool {
	for _, t := range ts {
		if e.components[t] == nil {
			return false
		}
	}
	return true
}

func (e *entity) HasAnyComponent(ts ...Type) bool {
	for _, t := range ts {
		if e.components[t] != nil {
			return true
		}
	}
	return false
}

func (e *entity) Component(t Type) (Component, error) {
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

		sort.Sort(Components(components))
		e.componentsCache = components
	}
	return components
}

func (e *entity) ComponentTypes() []Type {
	types := e.componentTypesCache
	if types == nil {
		types = make([]Type, 0, len(e.components))
		for t, c := range e.components {
			if c != nil {
				types = append(types, Type(t))
			}
		}
		e.componentTypesCache = types
	}
	return types
}

func (e *entity) AddComponent(cs ...Component) error {
	for _, c := range cs {
		if e.HasComponent(c.Type()) {
			return ErrComponentExists
		}
		e.components[c.Type()] = c
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
		old := e.components[c.Type()]
		e.components[c.Type()] = c
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

func (e *entity) RemoveComponent(ts ...Type) error {
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

func (e *entity) ID() uint64 {
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
	e.pool.DestroyEntity(e)
}

func (e *entity) String() string {
	return fmt.Sprintf("Entity_%d(types %v)", e.id, e.ComponentTypes())
}
