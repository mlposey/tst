// Package tst implements a ternary search tree.
package tst

// An Item is a value that can be placed in a TST.
// Those who wish to use TST to store a set of strings can
// avoid this and use `string` for TST methods instead.
type Item interface {
	Key() string
}

// TST is a ternary search tree.
//
// Ternary search trees (or tries) are data structures that facilitate
// string storage and searching. The data, which can be strings or
// structs that have string keys, is efficiently stored according to
// character prefixes. This storage method provides more complex
// search capabilities than other key/value stores like hash tables.
type TST struct {
	root *node
}

// New creates and returns a TST instance.
func New() *TST {
	// We may need to initialize stuff later and this would reduce
	// refactoring steps in usage code.
	return &TST{}
}

// Add places item into the tree.
// This function will panic if item is not a string or tst.Item type.
func (tree *TST) Add(item interface{}) {
	tree.root = tree.add(tree.root, tree.getKey(item), item, 0)
}

func (tree *TST) add(n *node, key string, value interface{}, i int) *node {
	c := key[i]
	if n == nil {
		n = &node{C: c}
	}

	if c < n.C {
		n.Left = tree.add(n.Left, key, value, i)
	} else if c > n.C {
		n.Right = tree.add(n.Right, key, value, i)
	} else if i < len(key)-1 {
		n.Mid = tree.add(n.Mid, key, value, i+1)
	} else {
		n.Value = value
	}
	return n
}

// Get returns the item from the tree whose key matches or nil if no
// item with that key exists.
func (tree *TST) Get(key string) interface{} {
	n := tree.get(tree.root, key, 0)
	if n == nil {
		return nil
	}
	return n.Value // This may also be nil, but it's the var we really want.
}

func (tree *TST) get(n *node, key string, i int) *node {
	if n == nil || len(key) == 0 {
		return n
	}
	c := key[i]

	if c < n.C {
		return tree.get(n.Left, key, i)
	} else if c > n.C {
		return tree.get(n.Right, key, i)
	} else if i < len(key)-1 {
		return tree.get(n.Mid, key, i+1)
	} else {
		return n
	}
}

// MatchPrefix returns a slice of values whose keys are prefixed with prefix.
// The prefix itself is also among the keys checked.
func (tree *TST) MatchPrefix(prefix string) []interface{} {
	var queue []interface{}

	if prefix == "" {
		return queue
	}

	// The node that full matches the prefix. We're interested in all
	// of its children.
	start := tree.get(tree.root, prefix, 0)

	if start == nil {
		return queue
	}
	if start.Value != nil {
		// The prefix itself could be a match.
		queue = append(queue, start.Value)
	}

	tree.matchPrefix(start.Mid, prefix, &queue)
	return queue
}

func (tree *TST) matchPrefix(n *node, prefix string, queue *[]interface{}) {
	if n == nil {
		return
	}
	if n.Value != nil {
		*queue = append(*queue, n.Value)
	}

	tree.matchPrefix(n.Left, prefix, queue)
	tree.matchPrefix(n.Mid, prefix+string(n.C), queue)
	tree.matchPrefix(n.Right, prefix, queue)
}

func (tree *TST) getKey(item interface{}) string {
	switch val := item.(type) {
	case string:
		return val
	case Item:
		return val.Key()
	default:
		panic("TST.Add must take either a string or an Item type")
	}
}

// nodes join together to create a TST.
// A node should always have a value for C. Other fields are optional
// depending on the context within the TST.
type node struct {
	C     byte
	Value interface{}
	Left  *node
	Mid   *node
	Right *node
}
