package data_test

import (
	. "fun/pkg/data"
	"strconv"
	"sync"
	"testing"
)

// Data is the type of data we're storing in the data structure.
type Data int

// String converts Data to a string
func (data Data) String() string {
	return strconv.Itoa(int(data))
}

// listAssert verifies state of the list
func listAssert[T ListData](t *testing.T, list *List[T], values []T) {
	if list == nil {
		t.Error("no list provided to compare data to")
		return
	}

	if list.Length() != len(values) {
		t.Error("expected list.Length ==", len(values), "got", list.Length())
	}

	failed := false
	i := 0
	var currentParentNode *ListNode[T]
	currentNode := list.Head()
	for currentNode != nil && i < len(values) {
		value, _ := currentNode.Value()
		if value != values[i] {
			t.Error("expected value", i, "does not match, value ", value)
			failed = true
		}
		currentParentNode = currentNode
		currentNode = currentNode.Next()
		i++
	}
	if currentParentNode != nil {
		if list.Tail() == nil {
			t.Errorf("tail is incorrect, should not be nil")
		} else if currentParentNode != list.Tail() {
			t.Errorf("tail is incorrect, should be %p not %p", currentParentNode, list.Tail())
		}
	}
	if failed {
		t.Log("expected: ", values)
		t.Log("got: ", list.String())
	}
}

func Test_New(t *testing.T) {
	list := NewList[Data]()
	listAssert(t, list, nil)
}

func Test_Insert(t *testing.T) {
	list := NewList[Data]()
	values := []Data{10}
	list.Insert(values[0])
	listAssert(t, list, values)

	values2 := []Data{20, 10}
	list.Insert(values2[0])
	listAssert(t, list, values2)
}

func Test_Delete(t *testing.T) {
	list := NewList[Data]()
	list.Insert(4)
	list.Insert(3)
	list.Insert(2)
	list.Insert(1)

	before := []Data{1, 2, 3, 4}
	listAssert(t, list, before)

	if ok := list.Delete(2); !ok {
		t.Error("failed to delete 2")
		t.Log(list.String())
	}
	afterDelete2 := []Data{1, 3, 4}
	listAssert(t, list, afterDelete2)

	deleted, ok := list.DeleteHead()
	if !ok {
		t.Error("failed result trying to delete head from non-empty list")
	}
	if deleted != 1 {
		t.Errorf("expect deleted 1, got %d", deleted)
	}
	afterDeleteHead := []Data{3, 4}
	listAssert(t, list, afterDeleteHead)

	deleted, ok = list.DeleteTail()
	if !ok {
		t.Error("failed result trying to delete tail from non-empty list")
	}
	if deleted != 4 {
		t.Errorf("expect deleted 4, got %d", deleted)
	}
	afterDeleteTail := []Data{3}
	listAssert(t, list, afterDeleteTail)

	list.Delete(3)
	emptyList := []Data{}
	listAssert(t, list, emptyList)
	_, ok = list.DeleteTail()
	if ok {
		t.Error("deleting tail from empty list returned", ok, "expecting", false)
	}
	listAssert(t, list, emptyList)

	list.Insert(7)
	list.DeleteHead()
	listAssert(t, list, emptyList)
	_, ok = list.DeleteHead()
	if ok {
		t.Error("deleting head from empty list returned", ok, "expecting", false)
	}
	listAssert(t, list, emptyList)
}

func Test_Concurrency(t *testing.T) {
	const threads = 10
	const iterations = 100

	list := NewList[Data]()
	listData := []Data{777}
	list.Append(listData[0])

	listAssert(t, list, listData)

	var wg sync.WaitGroup
	wg.Add(threads)

	for i := 0; i < threads; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				list.Append(Data(j))
				list.Delete(Data(j))
				list.Insert(Data(j))
				list.Delete(Data(j))
			}
		}()
	}

	wg.Wait()

	listAssert(t, list, listData)
}
