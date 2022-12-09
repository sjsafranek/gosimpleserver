package cache

import (
	"container/list"
	"sync"
)

type LRUCache[K comparable, V any] struct {
	capacity int                 // capacity
	ll       *list.List          // doubly linked list
	table    map[K]*list.Element // hash table for checking if list node exists
	lock     sync.RWMutex
}

// Entry is the value of a list node.
type Entry[K comparable, V any] struct {
	key   K
	value *V
}

// Get a list node from the hash map.
func (self *LRUCache[K, V]) Get(key K) *V {
	// set lock for thread safety
	self.lock.Lock()
	defer self.lock.Unlock()

	// check if list node exists
	if node, ok := self.table[key]; ok {
		// move node to front
		self.ll.MoveToFront(node)

		// return node value
		return node.Value.(*list.Element).Value.(Entry[K, V]).value
	}
	return nil
}

// Put key and value in the LRUCache
func (self *LRUCache[K, V]) Set(key K, value V) {
	// set lock for thread safety
	self.lock.Lock()
	defer self.lock.Unlock()

	// check if list node exists
	if node, ok := self.table[key]; ok {
		// move the node to front
		self.ll.MoveToFront(node)

		// update the value of a list node
		node.Value.(*list.Element).Value = Entry[K, V]{key: key, value: &value}
		return
	}

	// delete the last list node if the list is full
	if self.ll.Len() == self.capacity {
		// get the key that we want to delete
		idx := self.ll.Back().Value.(*list.Element).Value.(Entry[K, V]).key

		// delete the node pointer in the hash map by key
		delete(self.table, idx)

		// remove the last list node
		self.ll.Remove(self.ll.Back())
	}

	// initialize a list node
	node := &list.Element{
		Value: Entry[K, V]{
			key:   key,
			value: &value,
		},
	}

	// push the new list node into the list.
	// save the node pointer in the hash map.
	self.table[key] = self.ll.PushFront(node)
}

func (self *LRUCache[K, V]) Del(key K) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if node, ok := self.table[key]; ok {
		delete(self.table, key)
		self.ll.Remove(node)
	}
}

func (self *LRUCache[K, V]) Has(key K) bool {
	self.lock.Lock()
	defer self.lock.Unlock()
	_, ok := self.table[key]
	return ok
}

func (self *LRUCache[K, V]) Len() int {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.ll.Len()
}

// func (self *LRUCache[K, V]) Flush() {
//     self.lock.Lock()
//     defer self.lock.Unlock()
//     for key, node := range self.table {
//         delete(self.table, key)
//         self.ll.Remove(node)
//     }
// }

// New creates LRU cache with size capacity.
// func New[K comparable, V any](capacity int) *LRUCache[K, V] {
func New[K comparable, V any](capacity int) ICache[K, V] {
	return &LRUCache[K, V]{
		capacity: capacity,
		ll:       new(list.List),
		table:    make(map[K]*list.Element, capacity),
	}
}
