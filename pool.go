package entitas

import "fmt"

var (
	TotalComponents int
)

type ComponentNewFunc func() Component

type PoolEntityChanged func(Pool, Entity)

type PoolGroupChanged func(Pool, Group)

type PoolEntityEvent uint

const (
	PoolEntityCreated PoolEntityEvent = iota
	PoolEntityWillBeDestroyed
	PoolEntityDestroyed
)

type Pool interface {
	CreateComponent(ts ComponentType) Component
	AddComponentNewFunc(ts ComponentType, f ComponentNewFunc)

	CreateEntity(cs ...Component) Entity
	Entities() []Entity
	Count() int
	HasEntity(e Entity) bool
	destroyEntity(e Entity)
	DestroyAllEntities()
	Group(matcher ...Matcher) Group

	AddEvent(PoolEntityEvent, PoolEntityChanged)
	AddGroupCreatedEvent(changed PoolGroupChanged)
}

type pool struct {
	index       EntityID
	entities    map[EntityID]Entity
	entitiesCache       []Entity
	unused      []Entity

	groups      map[uint]Group
	groupsIndex map[ComponentType][]Group

	cacheComponents  [][]Component
	componentNewFunc []ComponentNewFunc

	entityChanged map[PoolEntityEvent][]PoolEntityChanged
	groupChanged []PoolGroupChanged
}

func NewPool(index EntityID) Pool {
	if TotalComponents == 0 {
		panic("please set entitas.TotalComponents")
	}
	return &pool{
		index:            index,
		entities:         make(map[EntityID]Entity),
		groups:           make(map[uint]Group),
		groupsIndex:      make(map[ComponentType][]Group),
		unused:           make([]Entity, 0),
		cacheComponents:  make([][]Component, TotalComponents),
		componentNewFunc: make([]ComponentNewFunc, TotalComponents),
		entityChanged:    make(map[PoolEntityEvent][]PoolEntityChanged),
	}
}

func (p *pool) CreateComponent(ts ComponentType) (component Component) {
	cache := p.cacheComponents[ts]
	length := len(cache)
	if length > 0 {
		last := length - 1
		component = cache[last]
		p.cacheComponents[ts] = cache[:last]
	} else {
		component = p.componentNewFunc[ts]()
	}
	return
}

func (p *pool) AddComponentNewFunc(ts ComponentType, f ComponentNewFunc) {
	p.componentNewFunc[ts] = f
}

func (p *pool) CreateEntity(cs ...Component) Entity {
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

func (p *pool) Entities() []Entity {
	if p.entitiesCache == nil {
		entities := make([]Entity, 0, len(p.entities))

		for _, e := range p.entities {
			entities = append(entities, e)
		}
		p.entitiesCache = entities
	}
	return p.entitiesCache
}

func (p *pool) Count() int {
	return len(p.entities)
}

func (p *pool) HasEntity(e Entity) bool {
	entity, exist := p.entities[e.ID()]
	return exist && entity == e
}

func (p *pool) destroyEntity(e Entity) {
	if p.HasEntity(e) {
		p.onEntityChanged(PoolEntityWillBeDestroyed, e)
		e.RemoveAllComponents()
		e.RemoveAllEvents()
		p.onEntityChanged(PoolEntityDestroyed, e)

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

func (p *pool) DestroyAllEntities() {
	for _, e := range p.entities {
		p.onEntityChanged(PoolEntityWillBeDestroyed, e)
		e.internalDestroy()
		p.onEntityChanged(PoolEntityDestroyed, e)
	}
	p.entities = make(map[EntityID]Entity)
	p.entitiesCache = nil
}

func (p *pool) Group(matchers ...Matcher) Group {
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

func (p *pool) AddEvent(event PoolEntityEvent, action PoolEntityChanged) {
	actions := p.entityChanged[event]
	p.entityChanged[event] = append(actions, action)
}

func (p *pool) AddGroupCreatedEvent(changed PoolGroupChanged){

}

func (p *pool) String() string {
	return fmt.Sprintf("Pool(%d entities, %d reusable, %d groups)",
		len(p.entities), len(p.unused), len(p.groups))
}

// private
func (p *pool) onEntityChanged(t PoolEntityEvent, entity Entity) {
	events := p.entityChanged[t]
	for _, event := range events {
		event(p, entity)
	}
}

func (p *pool) onGroupChanged(group Group) {
	for _, event := range p.groupChanged {
		event(p, group)
	}
}

func (p *pool) componentAdded(e Entity, c Component) {
	p.forMatchingGroup(e, c, func(g Group) {
		g.HandleEntity(e)
	})
}

func (p *pool) componentUpdated(e Entity, c Component) {
	p.forMatchingGroup(e, c, func(g Group) {
		g.UpdateEntity(e)
	})
}

func (p *pool) componentRemoved(e Entity, c Component) {
	p.cacheComponents[c.ComponentType()] = append(p.cacheComponents[c.ComponentType()], c)

	p.forMatchingGroup(e, c, func(g Group) {
		g.HandleEntity(e)
	})
}

func (p *pool) getEntity() (entity Entity) {
	
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

	p.onEntityChanged(PoolEntityCreated, entity)

	return entity
}

func (p *pool) forMatchingGroup(e Entity, c Component, f func(g Group)) {
	if p.HasEntity(e) {
		for _, g := range p.groupsIndex[c.ComponentType()] {
			f(g)
		}
	}
}
