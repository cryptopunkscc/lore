package id

type Counter map[ID]uint

func NewCounter() Counter {
	return make(map[ID]uint, 0)
}

func (c Counter) Increment(id ID) uint {
	c[id] = c[id] + 1
	return c[id]
}

func (c Counter) Decrement(id ID) uint {
	if c[id] == 0 {
		return 0
	}
	if c[id] == 1 {
		delete(c, id)
		return 0
	}
	c[id] = c[id] - 1
	return c[id]
}

func (c Counter) Count(id ID) uint {
	return c[id]
}

func (c Counter) Set() Set {
	var set = NewSet()
	for id, _ := range c {
		set.Add(id)
	}
	return set
}
