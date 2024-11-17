package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type OrderedMap struct {
	left, right *OrderedMap
	key         *int
	value       int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{} // need to implement
}

func (m *OrderedMap) Insert(key int, value int) {
	if m.key == nil {
		m.key = &key
		m.value = value
	} else if *m.key == key {
		m.value = value
	} else if *m.key > key {
		if m.left == nil {
			m.left = &OrderedMap{
				left:  nil,
				right: nil,
				key:   &key,
				value: value,
			}
		} else {
			m.left.Insert(key, value)
		}
	} else {
		if m.right == nil {
			m.right = &OrderedMap{
				left:  nil,
				right: nil,
				key:   &key,
				value: value,
			}
		} else {
			m.right.Insert(key, value)
		}
	}
}

func (m *OrderedMap) Erase(key int) {
	if m.key == nil {
		return
	}

	if *m.key == key {
		//TODO
		panic("test")

	}

	if key > *m.key && m.right != nil {
		if *m.right.key == key {
			if m.right.left != nil {
				m.right = m.right.left
			} else if m.right.right != nil {
				m.right = m.right.right
			} else {
				m.right = nil
			}
		} else {
			m.right.Erase(key)
		}
	} else if key < *m.key && m.left != nil {
		if *m.left.key == key {
			if m.left.left != nil {
				m.left = m.left.left
			} else if m.left.right != nil {
				m.left = m.left.right
			} else {
				m.left = nil
			}

		} else {
			m.left.Erase(key)
		}

	}

}

func (m *OrderedMap) Contains(key int) bool {
	if m.key == nil {
		return false
	} else if *m.key == key {
		return true
	} else if *m.key > key {
		if m.left == nil {
			return false
		}
		return m.left.Contains(key)
	} else if *m.key < key {
		if m.right == nil {
			return false
		}
		return m.right.Contains(key)
	} else {
		return *m.key == key
	}
}

func (m *OrderedMap) Size() int {
	size := 0
	if m.left != nil {
		size += m.left.Size()
	}
	if m.right != nil {
		size += m.right.Size()
	}
	if m.key != nil {
		size++
	}
	return size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	if m.left != nil {
		m.left.ForEach(action)
	}

	action(*m.key, m.value)

	if m.right != nil {
		m.right.ForEach(action)
	}
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
