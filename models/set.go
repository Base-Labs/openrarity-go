package models

// Set implements a set based on generics, and all elements in it are not duplicate.
type Set[T comparable] struct {
	dedup map[T]struct{}
	data  []T
}

// NewSet is the constructor of Set
func NewSet[T comparable](size int) *Set[T] {
	return &Set[T]{
		dedup: make(map[T]struct{}, size),
		data:  make([]T, 0, size),
	}
}

// Add is used to add element
func (c *Set[T]) Add(item T) {
	if _, exists := c.dedup[item]; !exists {
		c.dedup[item] = struct{}{}
		c.data = append(c.data, item)
	}
}

// List is used to get all the elements in the set
func (c *Set[T]) List() []T {
	return c.data
}

// IsSubset is used to determine whether subSet is a subset of mainSet
func IsSubset[V comparable](mainSet []V, subSet []V) bool {
	mainSetMap := make(map[V]struct{}, len(mainSet))
	for _, item := range mainSet {
		mainSetMap[item] = struct{}{}
	}
	for _, item := range subSet {
		if _, exists := mainSetMap[item]; !exists {
			return false
		}
	}
	return true
}
