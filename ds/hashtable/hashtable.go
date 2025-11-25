package hashtable

import (
	"fmt"
	"hash/fnv"
)

type Entry[K comparable, V any] struct {
	key       K
	value     V
	occupied  bool
	tombstone bool
}

type HashTable[K comparable, V any] struct {
	data          []Entry[K, V]
	size          int
	tombstoneSize int
	cap           int
}

func New[K comparable, V any](initialCap int) *HashTable[K, V] {
	if initialCap <= 0 {
		initialCap = 16
	}
	return &HashTable[K, V]{
		data: make([]Entry[K, V], initialCap),
		cap:  initialCap,
	}
}

func (h *HashTable[K, V]) Len() int {
	return h.size
}

func (h *HashTable[K, V]) Put(k K, v V) {
	key := hashKey(k) % uint64(h.cap)

	if float64((h.size+h.tombstoneSize)/h.cap) > 0.7 {
		h.resize()
	}

	firstTombstone := -1
	for i := 0; i < h.cap; i++ {
		idx := (key + uint64(i)) % uint64(h.cap)

		if h.data[idx].occupied {
			if h.data[idx].key == k {
				h.data[idx].value = v
				return
			}
			continue
		}
		if h.data[idx].tombstone {
			if firstTombstone < 0 {
				firstTombstone = int(idx)
			}
			continue
		}

		if firstTombstone < 0 {
			// Insert at first unoccupied position
			h.data[idx] = Entry[K, V]{
				key:       k,
				value:     v,
				occupied:  true,
				tombstone: false,
			}
		} else {
			// Insert at firstTombstone
			h.data[firstTombstone] = Entry[K, V]{
				key:       k,
				value:     v,
				occupied:  true,
				tombstone: false,
			}
			h.tombstoneSize -= 1
		}
		h.size += 1
		return
	}

	if firstTombstone >= 0 {
		// Found a tombstone even though we completed loop
		h.data[firstTombstone] = Entry[K, V]{
			key:       k,
			value:     v,
			occupied:  true,
			tombstone: false,
		}
		h.size += 1
		h.tombstoneSize -= 1
		return
	}
	panic("Hashtable is full!")
}

func (h *HashTable[K, V]) Get(k K) (V, bool) {
	key := hashKey(k) % uint64(h.cap)
	var zero V

	for i := 0; i < h.cap; i++ {
		idx := (key + uint64(i)) % uint64(h.cap)

		if h.data[idx].occupied {
			if h.data[idx].key == k {
				return h.data[idx].value, true
			}
		} else if !h.data[idx].occupied && !h.data[idx].tombstone {
			return zero, false
		}
	}
	return zero, false
}

func (h *HashTable[K, V]) Delete(k K) bool {
	key := hashKey(k)
	var zeroKey K
	var zeroVal V

	for i := 0; i < h.cap; i++ {
		idx := (key + uint64(i)) % uint64(h.cap)

		if h.data[idx].occupied {
			if h.data[idx].key == k {
				h.data[idx].occupied = false
				h.data[idx].tombstone = true
				h.data[idx].key = zeroKey
				h.data[idx].value = zeroVal
				h.size -= 1
				h.tombstoneSize += 1
				return true
			}
		} else if !h.data[idx].occupied && !h.data[idx].tombstone {
			return false
		}
	}

	return false
}

func (h *HashTable[K, V]) resize() {
	newData := make([]Entry[K, V], h.cap*2)

	for i := 0; i < h.cap; i++ {
		d := h.data[i]
		if d.occupied {
			newKey := hashKey(d.key)
			newIdx := newKey % uint64(h.cap*2)
			ent := newData[newIdx]
			for ent.occupied {
				newIdx := (newKey + 1) % uint64(h.cap*2)
				ent = newData[newIdx]
			}

			newData[newIdx].occupied = true
			newData[newIdx].key = d.key
			newData[newIdx].value = d.value
		}
	}
	h.data = newData
	h.cap *= 2
	h.tombstoneSize = 0
}

func hashKey[K comparable](key K) uint64 {
	h := fnv.New64a()

	switch k := any(key).(type) {
	case string:
		h.Write([]byte(k))
	case int:
		fmt.Fprintf(h, "%d", k)
	default:
		fmt.Fprintf(h, "%v", k)
	}

	return h.Sum64()
}
