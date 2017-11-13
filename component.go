package entitas

type ComponentType uint16

type ComponentTypes []ComponentType

func (ts ComponentTypes) Len() int {
	return len(ts)
}

func (ts ComponentTypes) Less(i, j int) bool {
	return ts[i] < ts[j]
}

func (ts ComponentTypes) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

//
type Component interface {
	ComponentType() ComponentType
}

//
type Components []Component

func (t Components) Len() int {
	return len(t)
}

func (t Components) Less(i, j int) bool {
	return t[i].ComponentType() < t[j].ComponentType()
}

func (t Components) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
