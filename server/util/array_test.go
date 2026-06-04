package util

import "testing"

func TestListToTreeCreatesChildrenWhenMissingKey(t *testing.T) {
	items := []map[string]interface{}{
		{"id": uint(1), "pid": uint(0), "name": "root"},
		{"id": uint(2), "pid": uint(1), "name": "child"},
		{"id": uint(3), "pid": uint(2), "name": "grandchild"},
	}

	tree := ArrayUtil.ListToTree(items, "id", "pid", "children")
	if len(tree) != 1 {
		t.Fatalf("tree length = %d, want 1", len(tree))
	}

	root := tree[0].(map[string]interface{})
	children := root["children"].([]interface{})
	if len(children) != 1 {
		t.Fatalf("root children length = %d, want 1", len(children))
	}

	child := children[0].(map[string]interface{})
	grandchildren := child["children"].([]interface{})
	if len(grandchildren) != 1 {
		t.Fatalf("child children length = %d, want 1", len(grandchildren))
	}
}

func TestListToTreeSupportsMixedNumericIDs(t *testing.T) {
	items := []map[string]interface{}{
		{"id": uint64(1), "pid": int(0), "name": "root"},
		{"id": int64(2), "pid": uint64(1), "name": "child"},
	}

	tree := ArrayUtil.ListToTree(items, "id", "pid", "children")
	root := tree[0].(map[string]interface{})
	children := root["children"].([]interface{})
	if children[0].(map[string]interface{})["name"] != "child" {
		t.Fatal("child node was not attached to root")
	}
}
