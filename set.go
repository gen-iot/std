package std

type Set struct {
	holder map[interface{}]bool
}

func NewSet() *Set {
	return &Set{
		holder: make(map[interface{}]bool, 16),
	}
}

func (this *Set) Len() int {
	return len(this.holder)
}

func (this *Set) Add(v interface{}) bool {
	if this.Contains(v) {
		return false
	}
	this.holder[v] = true
	return true
}

func (this *Set) Contains(v interface{}) bool {
	_, ok := this.holder[v]
	return ok
}

func (this *Set) Remove(v interface{}) {
	delete(this.holder, v)
}

func (this *Set) Collection() []interface{} {
	out := make([]interface{}, 0, len(this.holder))
	this.Foreach(func(it interface{}) {
		out = append(out, it)
	})
	return out
}

func (this *Set) Foreach(cb func(interface{})) {
	if cb == nil {
		return
	}
	for k := range this.holder {
		cb(k)
	}
}
