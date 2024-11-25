package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homewrok_test.go

type Node struct {
	left, right *Node
	key, value  int
}

type OrderedMap struct {
	root *Node
	size int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{
		root: nil,
		size: 0,
	}
}

func (m *OrderedMap) Insert(key int, value int) {
	if m.root == nil {
		m.root = &Node{left: nil, right: nil, key: key, value: value}
		m.size++
		return
	}

	m.insert(m.root, key, value)
}

func (m *OrderedMap) insert(n *Node, key int, value int) *Node {
	if n == nil {
		m.size++
		return &Node{left: nil, right: nil, key: key, value: value}
	}

	if n.key == key {
		n.value = value
		return n
	}

	if key < n.key {
		n.left = m.insert(n.left, key, value)
	} else {
		n.right = m.insert(n.right, key, value)
	}

	return n
}

func (m *OrderedMap) Erase(key int) {
	if m.root == nil {
		return
	}

	m.root = m.erase(m.root, key)
}

func (m *OrderedMap) erase(n *Node, key int) *Node {
	if n == nil {
		return nil
	}

	if key < n.key {
		n.left = m.erase(n.left, key)
	} else if key > n.key {
		n.right = m.erase(n.right, key)
	} else if n.key == key {

		if n.left == nil {
			m.size--
			return n.right
		} else if n.right == nil {
			m.size--
			return n.left
		}

		minimum := findMinimum(n.right)
		n.key = minimum.key
		n.value = minimum.value
		n.right = m.erase(n.right, minimum.key)
	}

	return n
}

func findMinimum(n *Node) *Node {
	if n.left == nil {
		return n
	}
	return findMinimum(n.left)
}

func (m *OrderedMap) Contains(key int) bool {
	return m.contains(m.root, key)
}

func (m *OrderedMap) contains(n *Node, key int) bool {
	if n == nil {
		return false
	}

	if n.key == key {
		return true
	}

	if key < n.key {
		return m.contains(n.left, key)
	}

	return m.contains(n.right, key)
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	m.root.foreach(action)
}

func (n *Node) foreach(action func(int, int)) {
	if n == nil {
		return
	}

	n.left.foreach(action)
	action(n.key, n.value)
	n.right.foreach(action)
}

func TestOrderedMap(t *testing.T) {
	t.Run("Basic operations", func(t *testing.T) {
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
	})

	t.Run("Erase root with two children", func(t *testing.T) {
		data := NewOrderedMap()
		data.Insert(10, 10)
		data.Insert(5, 5)
		data.Insert(15, 15)

		data.Erase(10) // Удаляем корень

		assert.Equal(t, 2, data.Size())
		assert.False(t, data.Contains(10))
		assert.True(t, data.Contains(5))
		assert.True(t, data.Contains(15))
	})

	t.Run("Erase non-existent element", func(t *testing.T) {
		data := NewOrderedMap()
		data.Insert(10, 10)
		data.Insert(5, 5)

		data.Erase(20) // Пытаемся удалить отсутствующий элемент

		assert.Equal(t, 2, data.Size())
		assert.True(t, data.Contains(10))
		assert.True(t, data.Contains(5))
	})

	t.Run("ForEach on empty map", func(t *testing.T) {
		data := NewOrderedMap()
		var keys []int

		data.ForEach(func(key, _ int) {
			keys = append(keys, key)
		})

		assert.Empty(t, keys)
	})

	t.Run("Insert duplicate key", func(t *testing.T) {
		data := NewOrderedMap()
		data.Insert(10, 10)
		data.Insert(10, 20) // Обновляем значение

		assert.Equal(t, 1, data.Size())
		assert.True(t, data.Contains(10))

		var keys []int
		expectedKeys := []int{10}

		data.ForEach(func(key, value int) {
			keys = append(keys, key)
			if key == 10 {
				assert.Equal(t, 20, value)
			}
		})

		assert.True(t, reflect.DeepEqual(expectedKeys, keys))
	})

	t.Run("Erase all elements", func(t *testing.T) {
		data := NewOrderedMap()
		data.Insert(10, 10)
		data.Insert(5, 5)
		data.Insert(15, 15)

		data.Erase(10)
		data.Erase(5)
		data.Erase(15)

		assert.Zero(t, data.Size())
		assert.False(t, data.Contains(10))
		assert.False(t, data.Contains(5))
		assert.False(t, data.Contains(15))

		var keys []int
		data.ForEach(func(key, _ int) {
			keys = append(keys, key)
		})

		assert.Empty(t, keys)
	})

	t.Run("Insert and erase negative keys", func(t *testing.T) {
		data := NewOrderedMap()
		data.Insert(-10, 10)
		data.Insert(-5, 5)

		assert.True(t, data.Contains(-10))
		assert.True(t, data.Contains(-5))

		data.Erase(-10)
		assert.False(t, data.Contains(-10))
		assert.True(t, data.Contains(-5))

		data.Erase(-5)
		assert.False(t, data.Contains(-5))
		assert.Zero(t, data.Size())
	})
}
