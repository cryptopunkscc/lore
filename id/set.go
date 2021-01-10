package id

type IDFunc func(ID)

type Set map[ID]struct{}

func NewSet() Set {
	return make(map[ID]struct{}, 0)
}

func (set Set) Add(id ID) {
	set[id] = struct{}{}
}

func (set Set) Remove(id ID) {
	delete(set, id)
}

func (set Set) Contains(id ID) bool {
	_, ok := set[id]
	return ok
}

func (set Set) Each(f IDFunc) {
	for k, _ := range set {
		f(k)
	}
}
