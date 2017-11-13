package entitas

type GroupChanged func(Group, Entity)

type Group interface {
	Entities() []Entity
	HandleEntity(e Entity)
	UpdateEntity(e Entity)
	Matches(e Entity) bool
	ContainsEntity(e Entity) bool

	AddEvent(EventType, GroupChanged)
	RemoveAllEvents()
}

type group struct {
	entities         map[EntityID]Entity
	cache            []Entity
	cacheInvalidated bool
	matchers         []Matcher

	groupChanged map[EventType][]GroupChanged
}

func newGroup(matchers ...Matcher) Group {
	return &group{
		entities:     make(map[EntityID]Entity),
		matchers:     matchers,
		groupChanged: make(map[EventType][]GroupChanged),
	}
}

func (g *group) Entities() []Entity {
	cache := g.cache
	if cache == nil {
		cache = make([]Entity, 0, len(g.entities))

		for _, e := range g.entities {
			cache = append(cache, e)
		}
		g.cache = cache
	}
	return cache
}

func (g *group) HandleEntity(e Entity) {
	if g.Matches(e) {
		g.addEntity(e)
	} else {
		g.removeEntity(e)
	}
}

func (g *group) UpdateEntity(e Entity) {
	if _, ok := g.entities[e.ID()]; ok {
		g.onGroupChanged(EventUpdated, e)
	}
}

func (g *group) Matches(e Entity) bool {
	for _, m := range g.matchers {
		if !m.Matches(e) {
			return false
		}
	}
	return true
}

func (g *group) ContainsEntity(e Entity) bool {
	_, exist := g.entities[e.ID()]
	return exist
}

func (g *group) AddEvent(event EventType, action GroupChanged) {
	actions := g.groupChanged[event]
	g.groupChanged[event] = append(actions, action)
}

func (g *group) RemoveAllEvents() {
	g.groupChanged = make(map[EventType][]GroupChanged)
}

// private
func (g *group) onGroupChanged(ev EventType, e Entity) {
	if events, ok := g.groupChanged[ev]; ok {
		for _, event := range events {
			event(g, e)
		}
	}
}

func (g *group) addEntity(e Entity) {
	if _, ok := g.entities[e.ID()]; !ok {
		g.entities[e.ID()] = e
		if g.cache != nil {
			g.cache = append(g.cache, e)
		}
		g.onGroupChanged(EventAdded, e)
	}
}

func (g *group) removeEntity(e Entity) {
	if _, ok := g.entities[e.ID()]; ok {
		delete(g.entities, e.ID())
		g.cache = nil
		g.onGroupChanged(EventRemoved, e)
	}
}
