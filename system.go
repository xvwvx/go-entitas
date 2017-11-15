package entitas

type System interface {
	Initialize(context Context)
	Execute()
}

type Systems struct {
	systems []System
}

func (ss Systems) Add(s System) {
	ss.systems = append(ss.systems, s)
}

func (ss Systems) Initialize(context Context) {
	for _, system := range ss.systems {
		system.Initialize(context)
	}
}

func (ss Systems) Execute() {
	for _, system := range ss.systems {
		system.Execute()
	}
}
