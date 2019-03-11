package servemux

import (
	"net/http"
	"strings"
)

// Value represents a trie value.
type Value = http.Handler

// Trie is a prefix search tree.
type Trie struct {
	value    Value
	param    string
	children map[string]*Trie // TODO: strange to use map within a trie :-|
}

// NewTrie allocates and returns a new *Trie.
func NewTrie() *Trie {
	return &Trie{
		children: make(map[string]*Trie),
	}
}

// Put inserts a new value into the tree.
func (t *Trie) Put(key string, val Value) bool {
	node := t

	for part, i := splitter(key, 0); ; part, i = splitter(key, i) {
		if isParam(part) {
			node.param = part
		}

		child, _ := node.children[part]
		if child == nil {
			child = NewTrie()
			node.children[part] = child
		}
		node = child
		if i == -1 {
			break
		}
	}

	isNewVal := node.value == nil
	node.value = val
	return isNewVal
}

// Get returns the value associated with the given key.
func (t *Trie) Get(key string) Value {
	node := t
	for part, i := splitter(key, 0); ; part, i = splitter(key, i) {
		node, _ = selectChild(node, part)
		if node == nil {
			return nil
		}
		if i == -1 {
			break
		}
	}

	return node.value
}

// GetWithParams returns the value associated with the given key.
func (t *Trie) GetWithParams(key string) (Value, map[string]string) {
	var params map[string]string
	node := t
	for part, i := splitter(key, 0); ; part, i = splitter(key, i) {
		n, isParamMatch := selectChild(node, part)
		if n == nil {
			return nil, params
		}
		if isParamMatch {
			if params == nil {
				params = map[string]string{
					node.param[1:]: part,
				}
			} else {
				params[node.param[1:]] = part
			}
		}
		node = n
		if i == -1 {
			break
		}
	}

	return node.value, params
}

func splitter(path string, start int) (segment string, next int) {
	if /* len(path) == 0 || start < 0 || */ start > len(path)-1 {
		return "", -1
	}

	end := strings.IndexRune(path[start+1:], '/')
	if end == -1 {
		return path[start+1:], -1
	}

	return path[start+1 : start+end+1], start + end + 1
}

func selectChild(node *Trie, key string) (*Trie, bool) {
	c, found := node.children[key]
	if found {
		return c, false
	}

	if node.param != "" {
		return node.children[node.param], true
	}

	return nil, false
}

func isParam(key string) bool {
	// return strings.HasPrefix(key, ":")

	if len(key) == 0 {
		return false
	}
	return key[0] == ':'
}
