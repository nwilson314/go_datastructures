package btree

import "slices"

type Entry struct {
	key   int
	value any
}

type Node struct {
	entries  []Entry
	children []*Node
}

type BTree struct {
	root  *Node
	order int
}

func New(order int) *BTree {
	return &BTree{
		root:  nil,
		order: order,
	}
}

func (b *BTree) Insert(key int, value any) {
	entry := Entry{
		key:   key,
		value: value,
	}
	if b.root == nil {
		b.root = &Node{
			entries:  make([]Entry, 0, b.order-1),
			children: make([]*Node, 0, b.order),
		}
		b.root.entries = append(b.root.entries, entry)
		return
	}

	if len(b.root.entries) >= b.order-1 {
		// split root and make tree taller
		newRoot := &Node{
			entries:  make([]Entry, 0, b.order-1),
			children: make([]*Node, 0, b.order),
		}

		// Move old root to new root's children
		newRoot.children = append(newRoot.children, b.root)
		b.root = newRoot

		b.splitChild(b.root, 0)
	}

	b.insertNonFull(b.root, entry)
}

func (b *BTree) Search(key int) (any, bool) {
	if b.root == nil {
		return nil, false
	}
	return b.searchNode(key, b.root)
}

func (b *BTree) searchNode(key int, node *Node) (any, bool) {
	// Recursive search through node and its children
	index, found := b.searchEntries(key, node.entries)
	if found {
		return node.entries[index].value, true
	}

	if len(node.children) == 0 {
		return nil, false
	}
	return b.searchNode(key, node.children[index])

}

func (b *BTree) searchEntries(key int, entries []Entry) (int, bool) {
	// Binary search for the index of the child pointer
	// If the key is found, return the index and true
	// If the key is not found, return the index of the child pointer to the left of the key
	// and false

	left := 0
	right := len(entries) - 1

	for left <= right {
		mid := (left + right) / 2
		if entries[mid].key == key {
			return mid, true
		} else if entries[mid].key < key {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left, false
}

func (b *BTree) splitChild(node *Node, index int) {
	// Only split if there's a single child and it is full
	mid := (b.order - 1) / 2

	rightSib := &Node{
		entries:  make([]Entry, 0, b.order-1),
		children: make([]*Node, 0, b.order),
	}

	leftSib := node.children[index]

	promotedEntry := leftSib.entries[mid]
	rightSib.entries = append(rightSib.entries, leftSib.entries[mid+1:]...)
	leftSib.entries = leftSib.entries[:mid]

	if len(leftSib.children) > 0 {
		rightSib.children = append(rightSib.children, leftSib.children[mid+1:]...)
		leftSib.children = leftSib.children[:mid+1]
	}

	// node.entries = append(node.entries, promotedEntry)
	// node.children = append(node.children, rightSib)
	node.entries = slices.Insert(node.entries, index, promotedEntry)
	node.children = slices.Insert(node.children, index+1, rightSib)
}

func (b *BTree) insertNonFull(node *Node, entry Entry) {
	if len(node.children) == 0 {
		// Leaf node
		index, _ := b.searchEntries(entry.key, node.entries)
		node.entries = slices.Insert(node.entries, index, entry)
	} else {
		// Internal node

		// Find the child to insert into
		index, _ := b.searchEntries(entry.key, node.entries)
		if len(node.children[index].entries) >= b.order-1 {
			b.splitChild(node, index)
			index, _ = b.searchEntries(entry.key, node.entries)
		}
		b.insertNonFull(node.children[index], entry)
	}
}
