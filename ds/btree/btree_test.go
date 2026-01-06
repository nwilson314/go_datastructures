package btree

import "testing"

func TestBTree_NewTree(t *testing.T) {
	tree := New(3) // order 3 = max 2 keys per node, max 3 children
	if tree == nil {
		t.Fatal("expected non-nil tree")
	}
}

func TestBTree_InsertAndSearch_Single(t *testing.T) {
	tree := New(3)

	tree.Insert(10, "ten")

	val, found := tree.Search(10)
	if !found {
		t.Fatal("expected to find key 10")
	}
	if val != "ten" {
		t.Errorf("expected 'ten', got '%v'", val)
	}
}

func TestBTree_Search_NotFound(t *testing.T) {
	tree := New(3)

	tree.Insert(10, "ten")

	_, found := tree.Search(99)
	if found {
		t.Error("expected not to find key 99")
	}
}

func TestBTree_InsertAndSearch_Multiple_NoSplit(t *testing.T) {
	tree := New(3) // max 2 keys per node

	tree.Insert(10, "ten")
	tree.Insert(20, "twenty")

	// Both should be in root, no split needed
	val, found := tree.Search(10)
	if !found || val != "ten" {
		t.Errorf("key 10: found=%v, val=%v", found, val)
	}

	val, found = tree.Search(20)
	if !found || val != "twenty" {
		t.Errorf("key 20: found=%v, val=%v", found, val)
	}
}

func TestBTree_Insert_MaintainsOrder(t *testing.T) {
	tree := New(3)

	// Insert out of order
	tree.Insert(30, "thirty")
	tree.Insert(10, "ten")
	tree.Insert(20, "twenty")

	// All should be findable
	for _, key := range []int{10, 20, 30} {
		_, found := tree.Search(key)
		if !found {
			t.Errorf("expected to find key %d", key)
		}
	}
}

func TestBTree_Insert_Split_Root(t *testing.T) {
	tree := New(3) // order 3 = max 2 keys per node

	// Insert 3 keys - this MUST trigger a split
	tree.Insert(10, "ten")
	tree.Insert(20, "twenty")
	tree.Insert(30, "thirty") // this causes split

	// After split of order-3 tree with [10, 20, 30]:
	//       [20]         <- new root with middle key
	//      /    \
	//   [10]    [30]     <- two children

	// Verify structure: root has 1 key
	if len(tree.root.entries) != 1 {
		t.Errorf("expected root to have 1 entry after split, got %d", len(tree.root.entries))
	}
	if tree.root.entries[0].key != 20 {
		t.Errorf("expected root key to be 20, got %d", tree.root.entries[0].key)
	}

	// Root should have 2 children
	if len(tree.root.children) != 2 {
		t.Fatalf("expected root to have 2 children, got %d", len(tree.root.children))
	}

	// Left child should have key 10
	if len(tree.root.children[0].entries) != 1 || tree.root.children[0].entries[0].key != 10 {
		t.Errorf("expected left child to have key 10")
	}

	// Right child should have key 30
	if len(tree.root.children[1].entries) != 1 || tree.root.children[1].entries[0].key != 30 {
		t.Errorf("expected right child to have key 30")
	}

	// All keys should still be findable
	val, found := tree.Search(10)
	if !found || val != "ten" {
		t.Errorf("key 10: found=%v, val=%v", found, val)
	}

	val, found = tree.Search(20)
	if !found || val != "twenty" {
		t.Errorf("key 20: found=%v, val=%v", found, val)
	}

	val, found = tree.Search(30)
	if !found || val != "thirty" {
		t.Errorf("key 30: found=%v, val=%v", found, val)
	}
}

func TestBTree_Insert_Split_HigherOrder(t *testing.T) {
	tree := New(5) // order 5 = max 4 keys per node, max 5 children

	// Insert 5 keys - triggers split on 5th insert
	tree.Insert(10, "ten")
	tree.Insert(20, "twenty")
	tree.Insert(30, "thirty")
	tree.Insert(40, "forty")
	tree.Insert(50, "fifty") // split!

	// After split with [10, 20, 30, 40, 50]:
	//        [30]            <- middle key promoted
	//       /    \
	//  [10,20]  [40,50]      <- 2 keys each side

	// Verify structure
	if len(tree.root.entries) != 1 {
		t.Fatalf("expected root to have 1 entry, got %d", len(tree.root.entries))
	}
	if tree.root.entries[0].key != 30 {
		t.Errorf("expected root key 30, got %d", tree.root.entries[0].key)
	}

	if len(tree.root.children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(tree.root.children))
	}

	// Left child: [10, 20]
	if len(tree.root.children[0].entries) != 2 {
		t.Errorf("expected left child to have 2 entries, got %d", len(tree.root.children[0].entries))
	}

	// Right child: [40, 50]
	if len(tree.root.children[1].entries) != 2 {
		t.Errorf("expected right child to have 2 entries, got %d", len(tree.root.children[1].entries))
	}

	// All keys findable
	for _, key := range []int{10, 20, 30, 40, 50} {
		_, found := tree.Search(key)
		if !found {
			t.Errorf("expected to find key %d", key)
		}
	}
}

func TestBTree_Insert_MultipleSplits(t *testing.T) {
	tree := New(3) // order 3 = small nodes, frequent splits

	// Insert 7 keys - will cause multiple splits
	keys := []int{10, 20, 30, 40, 50, 60, 70}
	for _, k := range keys {
		tree.Insert(k, k*10)
	}

	// Tree should have grown in height
	// With order 3, after 7 inserts we should have height 3:
	//
	//           [40]
	//          /    \
	//       [20]    [60]
	//       /  \    /  \
	//    [10] [30] [50] [70]

	// All keys must be findable
	for _, k := range keys {
		val, found := tree.Search(k)
		if !found {
			t.Errorf("expected to find key %d", k)
		}
		if val != k*10 {
			t.Errorf("key %d: expected value %d, got %v", k, k*10, val)
		}
	}

	// Root should not be a leaf (should have children)
	if len(tree.root.children) == 0 {
		t.Error("expected root to have children after multiple splits")
	}
}

func TestBTree_Insert_ManyKeys(t *testing.T) {
	tree := New(4) // order 4

	// Insert 100 keys
	for i := 0; i < 100; i++ {
		tree.Insert(i, i*2)
	}

	// All keys must be findable
	for i := 0; i < 100; i++ {
		val, found := tree.Search(i)
		if !found {
			t.Errorf("expected to find key %d", i)
		}
		if val != i*2 {
			t.Errorf("key %d: expected %d, got %v", i, i*2, val)
		}
	}

	// Search for non-existent keys
	_, found := tree.Search(-1)
	if found {
		t.Error("should not find key -1")
	}
	_, found = tree.Search(100)
	if found {
		t.Error("should not find key 100")
	}
}

func TestBTree_Insert_NonSequential(t *testing.T) {
	tree := New(3)

	// Insert in an order that forces splitting a non-last child
	// After inserting 50, 30, 70: tree is root=[50], children=[[30], [70]]
	// Insert 20: left child becomes [20, 30]
	// Insert 40: left child is full, split it — this exposes the bug
	for _, k := range []int{50, 30, 70, 20, 40} {
		tree.Insert(k, k)
	}

	// All keys should be findable
	for _, k := range []int{20, 30, 40, 50, 70} {
		_, found := tree.Search(k)
		if !found {
			t.Errorf("expected to find key %d", k)
		}
	}
}
