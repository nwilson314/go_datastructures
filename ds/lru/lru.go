package lru

import "container/list"

// Cache represents an LRU Cache.
// It is not thread-safe.
type Cache[K comparable, V any] struct {
	// cap is the maximum number of items.
	cap int

	// lruList is a doubly linked list.
	// Front = Most Recently Used.
	// Back = Least Recently Used.
	// Stores *Entry[K, V].
	lruList *list.List

	// cache is a map pointing to the list element for O(1) access.
	cache map[K]*list.Element
}

type Entry[K comparable, V any] struct {
	key   K
	value V
}

func New[K comparable, V any](capacity int) *Cache[K, V] {
	return &Cache[K, V]{
		cap:     capacity,
		lruList: list.New(),
		cache:   make(map[K]*list.Element),
	}
}

func (c *Cache[K, V]) Len() int {
	return c.lruList.Len()
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	var zero V
	el, ok := c.cache[key]
	if !ok {
		return zero, false
	}

	c.lruList.MoveToFront(el)

	return el.Value.(Entry[K, V]).value, true

}

func (c *Cache[K, V]) Put(key K, value V) {
	el, ok := c.cache[key]

	if ok {
		// Value has been moved to front, but we still need to update the value.
		el.Value = Entry[K, V]{key, value}
		c.lruList.MoveToFront(el)
		return
	}

	if c.Len() == c.cap {
		backEl := c.lruList.Back()
		c.lruList.Remove(c.lruList.Back())
		delete(c.cache, backEl.Value.(Entry[K, V]).key)
	}

	c.lruList.PushFront(Entry[K, V]{
		key, value,
	})
	c.cache[key] = c.lruList.Front()
}
