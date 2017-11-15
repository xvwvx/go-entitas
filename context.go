package entitas

import (
	"fmt"
	"reflect"
)

var (
	TotalComponents int
)

type ComponentNewFunc func() Component

type ContextEntityChanged func(Context, Entity)

type ContextGroupChanged func(Context, Group)

type ContextEntityEvent uint

const (
	ContextEntityCreated ContextEntityEvent = iota
	ContextEntityWillBeDestroyed
	ContextEntityDestroyed
)

type Context interface {
	CreateComponent(ts int) Component
	RegisterComponent(component Component)

	CreateEntity(cs ...Component) Entity
	Entities() []Entity
	Count() int
	HasEntity(e Entity) bool
	destroyEntity(e Entity)
	DestroyAllEntities()
	Group(matcher ...Matcher) Group

	AddEvent(ContextEntityEvent, ContextEntityChanged)
	AddGroupCreatedEvent(changed ContextGroupChanged)
}

type context struct {
	index         EntityID
	entities      map[EntityID]Entity
	entitiesCache []Entity
	unused        []Entity

	groups      map[uint]Group
	groupsIndex map[int][]Group

	cacheComponents   [][]Component
	registerComponent []reflect.Type

	entityChanged map[ContextEntityEvent][]ContextEntityChanged
	groupChanged  []ContextGroupChanged
}

func NewContext(index EntityID) Context {
	if TotalComponents == 0 {
		panic("please set entitas.TotalComponents")
	}
	return &context{
		index:             index,
		entities:          make(map[EntityID]Entity),
		groups:            make(map[uint]Group),
		groupsIndex:       make(map[int][]Group),
		unused:            make([]Entity, 0),
		cacheComponents:   make([][]Component, TotalComponents),
		registerComponent: make([]reflect.Type, TotalComponents),
		entityChanged:     make(map[ContextEntityEvent][]ContextEntityChanged),
	}
}

func (p *context) CreateComponent(ts int) (component Component) {
	cache := p.cacheComponents[ts]
	length := len(cache)
	if length > 0 {
		last := length - 1
		component = cache[last]
		p.cacheComponents[ts] = cache[:last]
	} else {
		value := reflect.New(p.registerComponent[ts])
		component = value.Interface().(Component)
	}
	return
}

func (p *context) RegisterComponent(component Component) {
	p.registerComponent[component.Type()] = reflect.TypeOf(component).Elem()
}

func (p *context) CreateEntity(cs ...Component) Entity {
	e := p.getEntity()
	e.AddComponent(cs...)
	p.entities[e.ID()] = e
	if p.entitiesCache != nil {
		p.entitiesCache = append(p.entitiesCache, e)
	}

	for _, g := range p.groups {
		g.HandleEntity(e)
	}
	return e
}

func (p *context) Entities() []Entity {
	if p.entitiesCache == nil {
		entities := make([]Entity, 0, len(p.entities))

		for _, e := range p.entities {
			entities = append(entities, e)
		}
		p.entitiesCache = entities
	}
	return p.entitiesCache
}

func (p *context) Count() int {
	return len(p.entities)
}

func (p *context) HasEntity(e Entity) bool {
	entity, exist := p.entities[e.ID()]
	return exist && entity == e
}

func (p *context) destroyEntity(e Entity) {
	if p.HasEntity(e) {
		p.onEntityChanged(ContextEntityWillBeDestroyed, e)
		e.RemoveAllComponents()
		e.RemoveAllEvents()
		p.onEntityChanged(ContextEntityDestroyed, e)

		delete(p.entities, e.ID())

		p.entitiesCache = nil
		for _, g := range p.groups {
			g.HandleEntity(e)
		}
		p.unused = append(p.unused, e)
	} else {
		panic("unknown entity")
	}
}

func (p *context) DestroyAllEntities() {
	for _, e := range p.entities {
		p.onEntityChanged(ContextEntityWillBeDestroyed, e)
		e.internalDestroy()
		p.onEntityChanged(ContextEntityDestroyed, e)
	}
	p.entities = make(map[EntityID]Entity)
	p.entitiesCache = nil
}

func (p *context) Group(matchers ...Matcher) Group {
	hash := HashMatcher(matchers...)
	if g, ok := p.groups[hash]; ok {
		return g
	}

	g := newGroup(matchers...)
	for _, e := range p.entities {
		g.HandleEntity(e)
	}
	p.groups[hash] = g

	for _, m := range matchers {
		for _, t := range m.ComponentTypes() {
			p.groupsIndex[t] = append(p.groupsIndex[t], g)
		}
	}

	p.onGroupChanged(g)

	return g
}

func (p *context) AddEvent(event ContextEntityEvent, action ContextEntityChanged) {
	actions := p.entityChanged[event]
	p.entityChanged[event] = append(actions, action)
}

func (p *context) AddGroupCreatedEvent(changed ContextGroupChanged) {

}

func (p *context) String() string {
	return fmt.Sprintf("Context(%d entities, %d reusable, %d groups)",
		len(p.entities), len(p.unused), len(p.groups))
}

// private
func (p *context) onEntityChanged(t ContextEntityEvent, entity Entity) {
	events := p.entityChanged[t]
	for _, event := range events {
		event(p, entity)
	}
}

func (p *context) onGroupChanged(group Group) {
	for _, event := range p.groupChanged {
		event(p, group)
	}
}

func (p *context) componentAdded(e Entity, c Component) {
	p.forMatchingGroup(e, c, func(g Group) {
		g.HandleEntity(e)
	})
}

func (p *context) componentUpdated(e Entity, c Component) {
	p.forMatchingGroup(e, c, func(g Group) {
		g.UpdateEntity(e)
	})
}

func (p *context) componentRemoved(e Entity, c Component) {
	t := c.Type()
	p.cacheComponents[t] = append(p.cacheComponents[t], c)

	p.forMatchingGroup(e, c, func(g Group) {
		g.HandleEntity(e)
	})
}

func (p *context) getEntity() (entity Entity) {

	length := len(p.unused)
	if length > 0 {
		last := length - 1
		entity = p.unused[last]
		p.unused = p.unused[:last]
	} else {
		entity = newEntity(p, p.index)
		p.index++
	}

	entity.AddEvent(EventAdded, p.componentAdded)
	entity.AddEvent(EventUpdated, p.componentUpdated)
	entity.AddEvent(EventRemoved, p.componentRemoved)

	p.onEntityChanged(ContextEntityCreated, entity)

	return entity
}

func (p *context) forMatchingGroup(e Entity, c Component, f func(g Group)) {
	if p.HasEntity(e) {
		for _, g := range p.groupsIndex[c.Type()] {
			f(g)
		}
	}
}
