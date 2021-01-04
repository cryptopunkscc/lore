package id

type Set struct {
	items map[string]struct{}
}

func NewSet() *Set {
	return &Set{
		items: make(map[string]struct{}, 0),
	}
}

func (set *Set) Add(item string) {
	set.items[item] = struct{}{}
}

func (set *Set) Remove(item string) {
	delete(set.items, item)
}

func (set *Set) Contains(item string) bool {
	_, ok := set.items[item]
	return ok
}

func (set *Set) Each(f func(string)) {
	for k, _ := range set.items {
		f(k)
	}
}

func (set *Set) Union(set2 *Set) *Set {
	set2.Each(func(item string) {
		set.Add(item)
	})
	return set
}

func (set *Set) Intersection(set2 *Set) *Set {
	s := NewSet()
	set.Each(func(item string) {
		if set2.Contains(item) {
			s.Add(item)
		}
	})
	return s
}

func (set *Set) Subtract(set2 *Set) *Set {
	set2.Each(func(item string) {
		set.Remove(item)
	})
	return set
}