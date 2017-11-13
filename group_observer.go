package entitas

type ObserverEvent uint

type GroupObserver interface {
	CollectedEntities() []Entity
	Activate()
	Deactivate()
	ClearCollectedEntities()
}

type groupObserver struct {
	entities map[Entity]bool
	active   bool
}

func NewGroupObserver(group Group, event EventType) GroupObserver {
	observer := &groupObserver{
		entities: make(map[Entity]bool),
		active:   true,
	}

	callback := func(group Group, entity Entity) {
		addEntity(observer, group, entity)
	}

	if event == EventAddedOrRemoved {
		group.AddEvent(EventAdded, callback)
		group.AddEvent(EventRemoved, callback)
	} else {
		group.AddEvent(event, callback)
	}

	return observer
}

func (observer *groupObserver) CollectedEntities() []Entity {
	entities := make([]Entity, 0, len(observer.entities))

	for entity := range observer.entities {
		entities = append(entities, entity)
	}
	return entities
}

func (observer *groupObserver) Activate() {
	observer.active = true
}

func (observer *groupObserver) Deactivate() {
	observer.active = false
}

func (observer *groupObserver) ClearCollectedEntities() {
	observer.entities = make(map[Entity]bool)
}

func addEntity(observer *groupObserver, group Group, entity Entity) {
	if observer.active {
		observer.entities[entity] = true
	}
}
