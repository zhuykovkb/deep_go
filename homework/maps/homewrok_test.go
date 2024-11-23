package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Node struct {
	left, right *Node
	key         int
	value       int
}

type OrderedMap struct {
	rootNode *Node
	size     int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{
		rootNode: nil,
		size:     0,
	}
}

func (m *OrderedMap) Insert(key int, value int) {
	n := &Node{
		left:  nil,
		right: nil,
		key:   key,
		value: value,
	}

	if m.rootNode == nil {
		m.rootNode = n
		m.size++
		return
	}

	if key > m.rootNode.key {
		m.rootNode.right = insertIntoBst(m.rootNode.right, key, value)
	} else {
		m.rootNode.left = insertIntoBst(m.rootNode.left, key, value)
	}
	m.size++

}

func insertIntoBst(n *Node, key int, value int) *Node {
	if n == nil {
		return newNode(key, value)
	}

	if key < n.key {
		n.left = insertIntoBst(n.left, key, value)
	} else {
		n.right = insertIntoBst(n.right, key, value)
	}

	return n
}

func newNode(key int, value int) *Node {
	return &Node{
		left:  nil,
		right: nil,
		key:   key,
		value: value,
	}
}

func (m *OrderedMap) Erase(key int) {

	if m.rootNode == nil {
		return
	}

	if key < m.rootNode.key {
		m.rootNode.left = eraseNode(m.rootNode.left, key)
	} else {
		m.rootNode.right = eraseNode(m.rootNode.left, key)
	}

}

func eraseNode(n *Node, key int) *Node {
	if n == nil {
		return nil
	}

	if key < n.key {
		n.left = eraseNode(n.left, key)
	} else if key > n.key {
		n.right = eraseNode(n.right, key)
	}

	return nil
}

func (m *OrderedMap) Contains(key int) bool {
	return false
	//if m.key == nil {
	//	return false
	//} else if *m.key == key {
	//	return true
	//} else if *m.key > key {
	//	if m.left == nil {
	//		return false
	//	}
	//	return m.left.Contains(key)
	//} else if *m.key < key {
	//	if m.right == nil {
	//		return false
	//	}
	//	return m.right.Contains(key)
	//} else {
	//	return *m.key == key
	//}
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	//if m.left != nil {
	//	m.left.ForEach(action)
	//}
	//
	//action(*m.key, m.value)
	//
	//if m.right != nil {
	//	m.right.ForEach(action)
	//}
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
