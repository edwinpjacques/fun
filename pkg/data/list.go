// Package data implements various data structures.
package data

import (
	"errors"
	"fmt"
	"sync"
)

// ListData must be comparable, and can be converted to a string.
type ListData interface {
	comparable
	String() string
}

// ListNode is a singly-linked list data structure.
type ListNode[T ListData] struct {
	value T            // Value is storage for data in the list.
	next  *ListNode[T] // Pointer to the next element in the list.
}

// List data structure.
type List[T ListData] struct {
	head   *ListNode[T]  // Head of the list.
	tail   *ListNode[T]  // Tail of the list.
	length int           // Number of elements stored in the list.
	mux    *sync.RWMutex // Lock read and write operations.
}

// Create a new list.
func NewList[T ListData]() *List[T] {
	return &List[T]{mux: &sync.RWMutex{}}
}

// Length reports the number of elements in the list.
func (list *List[T]) Length() int {
	return list.length
}

// Head gets the head of the list.
func (list *List[T]) Head() *ListNode[T] {
	if list == nil {
		return nil
	}
	return list.head
}

// Head gets the head of the list.
func (list *List[T]) Tail() *ListNode[T] {
	if list == nil {
		return nil
	}
	return list.tail
}

// Value gets the value of a ListNode.
func (listNode *ListNode[T]) Value() (T, bool) {
	var unset T
	if listNode == nil {
		return unset, false
	}
	return listNode.value, true
}

// Next gets the next node of a ListNode.
func (listNode *ListNode[T]) Next() *ListNode[T] {
	if listNode == nil {
		return nil
	}
	return listNode.next
}

// Insert adds an element at the beginning of a list.
func (list *List[T]) Insert(value T) error {
	if list == nil {
		return errors.New("list is nil")
	}
	list.mux.Lock()
	defer list.mux.Unlock()
	listNode := &ListNode[T]{value, list.head}
	if list.tail == nil {
		list.tail = listNode
	}
	list.head = listNode
	list.length++
	return nil
}

// Append adds an element at the end of a list.
func (list *List[T]) Append(value T) error {
	if list == nil {
		return errors.New("list is nil")
	}
	list.mux.Lock()
	defer list.mux.Unlock()
	listNode := &ListNode[T]{value, nil}
	if list.tail == nil {
		list.tail = listNode
		list.head = listNode
	} else {
		list.tail.next = listNode
		list.tail = listNode
	}
	list.length++
	return nil
}

// findParent finds a node by its value and the parent.
func (list *List[T]) findParent(value T) (parent *ListNode[T], found *ListNode[T]) {
	if list == nil {
		return nil, nil
	}

	var lastNode *ListNode[T]
	currentNode := list.head
	for currentNode != nil {
		if currentNode.value == value {
			return lastNode, currentNode
		}
		lastNode = currentNode
		currentNode = currentNode.next
	}
	return nil, nil
}

// Find a value in the list.
func (list *List[T]) Find(value T) (listNode *ListNode[T]) {
	list.mux.RLock()
	defer func() {
		list.mux.RUnlock()
	}()
	_, found := list.findParent(value)
	return found
}

// Delete Data in the list.
func (list *List[T]) Delete(value T) bool {
	list.mux.Lock()
	defer list.mux.Unlock()
	parent, found := list.findParent(value)
	if found == nil {
		return false
	}
	if parent != nil {
		parent.next = found.next
	}
	if list.head == found {
		list.head = found.next
	}
	if list.tail == found {
		list.tail = parent
	}
	list.length--
	return true
}

// Delete the head node in the list.
func (list *List[T]) DeleteHead() (T, bool) {
	list.mux.Lock()
	defer list.mux.Unlock()
	var value T
	if list.head == nil {
		return value, false
	}
	value = list.head.value
	list.head = list.head.next
	if list.head == nil {
		list.tail = nil
	}
	list.length--
	return value, true
}

// Delete the tail node in the list.
func (list *List[T]) DeleteTail() (T, bool) {
	list.mux.Lock()
	defer list.mux.Unlock()
	var value T
	if list.tail == nil {
		return value, false
	} else {
		value = list.tail.value
	}
	// find the parent of list.Tail
	var parent *ListNode[T]
	for node := list.head; node != nil; node = node.next {
		if node.next == list.tail {
			parent = node
		}
	}
	if parent == nil {
		list.head = nil
		list.tail = nil
	} else {
		parent.next = nil
		list.tail = parent
	}
	list.length--
	return value, true
}

// For each value in the list, execute a method.
func (list *List[T]) ForEach(f func(T)) {
	list.mux.RLock()
	defer list.mux.RUnlock()
	currentNode := list.head
	for currentNode != nil {
		f(currentNode.value)
	}
}

// String converts List data into a string.
func (list *List[T]) String() string {
	if list == nil {
		return ""
	}
	list.mux.RLock()
	defer list.mux.RUnlock()
	currentNode := list.head
	values := make([]T, 0, list.length)
	for currentNode != nil {
		values = append(values, currentNode.value)
		currentNode = currentNode.next
	}

	s := fmt.Sprintf("Length: %d, Data:", list.length)
	for _, v := range values {
		s += " " + v.String()
	}
	return s
}
