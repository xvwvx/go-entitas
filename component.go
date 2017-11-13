package entitas

type Type uint16

type Types []Type

func (ts Types) Len() int {
	return len(ts)
}

func (ts Types) Less(i, j int) bool {
	return ts[i] < ts[j]
}

func (ts Types) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

//
type Component interface {
	Type() Type
}

//
type Components []Component

func (t Components) Len() int {
	return len(t)
}

func (t Components) Less(i, j int) bool {
	return t[i].Type() < t[j].Type()
}

func (t Components) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
