package tst_test

import (
	"github.com/mlposey/tst"
	"testing"
)

type mockItem struct {
	key   string
	Value int
}

func (mock mockItem) Key() string {
	return mock.key
}

// TST.Add should succeed in adding string or Item types.
func TestTST_Add_succeed(t *testing.T) {
	tree := tst.New()

	tree.Add("cat")
	tree.Add(mockItem{"can", 3})
}

// TST.Add should panic if receiving non string or Item types.
func TestTST_Add_fail(t *testing.T) {
	tree := tst.New()

	defer func() {
		if recover() == nil {
			t.Error("Did not panic when adding unsupported type")
		}
	}()
	tree.Add(4)
}

// TST.Get should retrieve an item if its key exists.
func TestTST_Get_exists(t *testing.T) {
	tree := tst.New()

	items := []string{"cat", "car", "tom", "c"}
	for _, item := range items {
		tree.Add(item)
	}

	// Search forwards.
	for _, item := range items {
		if tree.Get(item) != item {
			t.Error("Failed to retrieve existing item from tree")
		}
	}

	// Search backwards.
	for i := len(items) - 1; i >= 0; i-- {
		if tree.Get(items[i]) != items[i] {
			t.Error("Failed to retrieve existing item from tree")
		}
	}
}

// TST.Get should return nil if the key does not exist.
func TestTST_Get_noexist(t *testing.T) {
	tree := tst.New()

	if tree.Get("test") != nil {
		t.Error("A ghost added items to the tree")
		// ...but we really just hope this doesn't fail.
	}

	tree.Add("tom")

	for _, query := range []string{"t", "to", "too", "tomo"} {
		if tree.Get(query) != nil {
			t.Error("Retrieved wrong item from tree")
		}
	}

	// An index panic could happen here.
	tree.Get("")
}

// TST.MatchPrefix should return a slice of only items whose keys
// either equal or are prefixed with a prefix argument.
func TestTST_MatchPrefix_exists(t *testing.T) {
	tree := tst.New()

	prefix := "c"

	// The number of times each key should appear in the MatchPrefix result
	expected := map[string]int{
		"cat":  1,
		"car":  1,
		"tom":  0,
		prefix: 1,
	}
	for key := range expected {
		tree.Add(key)
	}

	actual := make(map[string]int)
	for _, item := range tree.MatchPrefix(prefix) {
		actual[item.(string)]++
	}

	for key := range expected {
		if expected[key] != actual[key] {
			t.Error("Found bad match in MatchPrefix result")
		}
	}
}

// TST.MatchPrefix should return an empty slice if there are no keys having
// the provided prefix.
func TestTST_MatchPrefix_noexist(t *testing.T) {
	tree := tst.New()

	items := []string{"cat", "car", "tom", "c"}
	for _, item := range items {
		tree.Add(item)
	}

	msg := "Found bad match in MatchPrefix result"
	if len(tree.MatchPrefix("can")) != 0 {
		t.Error(msg)
	}
	if len(tree.MatchPrefix("")) != 0 {
		t.Error(msg)
	}
}
