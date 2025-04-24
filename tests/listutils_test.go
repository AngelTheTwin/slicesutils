package tests

import (
	"testing"

	"github.com/AngelTheTwin/slicesutils"
)

var items []int = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

func TestFindIndeex(t *testing.T) {
	index := slicesutils.FindIndex(items, func(item int) bool {
		return item == 5
	})

	if index != 4 {
		t.Errorf("Expected index 4, but got %d", index)
	}

	index = slicesutils.FindIndex(items, func(item int) bool {
		return item == 11
	})

	if index != -1 {
		t.Errorf("Expected index -1, but got %d", index)
	}
}

func TestFind(t *testing.T) {
	item, ok := slicesutils.Find(items, func(item int) bool {
		return item == 5
	})

	if !ok {
		t.Errorf("Expected to find item 5")
	}

	if item != 5 {
		t.Errorf("Expected item 5, but got %d", item)
	}

	item, ok = slicesutils.Find(items, func(item int) bool {
		return item == 11
	})

	if ok {
		t.Errorf("Expected not to find item 11")
	}

	if item != 0 {
		t.Errorf("Expected item 0, but got %d", item)
	}
}

func TestParallelMap(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expected := []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}

	result := slicesutils.ParallelMap(items, func(item int) int {
		return item * 2
	})

	for i, item := range result {
		if item != expected[i] {
			t.Errorf("Expected %d, but got %d", expected[i], item)
		}
	}
}

func TestRemoveElement_OneOcurrence(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expected := []int{1, 2, 3, 4, 6, 7, 8, 9, 10}

	result := slicesutils.RemoveElement(items, 5, 1)

	if ok := slicesutils.Compare(expected, result); !ok {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestRemoveElement_MultipleOcurrences(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 5, 5, 5}
	expected := []int{1, 2, 3, 4, 6, 7, 8, 9, 10, 5, 5}

	result := slicesutils.RemoveElement(items, 5, 2)

	if ok := slicesutils.Compare(expected, result); !ok {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestRemoveElement_AllOcurrences(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 5, 5, 5}
	expected := []int{1, 2, 3, 4, 6, 7, 8, 9, 10}

	result := slicesutils.RemoveElement(items, 5, -1)

	if ok := slicesutils.Compare(expected, result); !ok {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestRemoveElements(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expected := []int{1, 2, 3, 6, 7, 8, 9, 10}

	result := slicesutils.RemoveElements(items, 5, 4)

	if ok := slicesutils.Compare(expected, result); !ok {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestChunkSlice(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	expected := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}

	result := slicesutils.Chunk(input, 3)

	if len(result) != len(expected) {
		t.Errorf("Expected slice of length %v, but got %v", expected, result)
	}

	for i, chunk := range result {
		if ok := slicesutils.Compare(chunk, expected[i]); !ok {
			t.Errorf("Expected %v, but got %v", expected[i], chunk)
		}
	}
}

func TestChunkSlice_ChunkSizeBiggerThanSliceLength(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	expected := [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}}

	result := slicesutils.Chunk(input, 10)

	if len(result) != len(expected) {
		t.Errorf("Expected slice of length %v, but got %v", expected, result)
	}

	for i, chunk := range result {
		if ok := slicesutils.Compare(chunk, expected[i]); !ok {
			t.Errorf("Expected %v, but got %v", expected[i], chunk)
		}
	}
}

func TestParallelForEach(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	output := make([]int, len(items))
	expected := []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}

	slicesutils.ParallelForEach(items, func(item int) {
		output[item-1] = item * 2
	})

	if ok := slicesutils.Compare(expected, output); !ok {
		t.Errorf("Expected %v, but got %v", expected, output)
	}
}
func TestParallelForEach_AppendingItems(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	output := 0
	expected := 55

	slicesutils.ParallelForEach(items, func(item int) {
		output += item
	})

	// compare the slices as sets to avoid order issues as the order is not guaranteed when using ParallelForEach
	if output != expected {
		t.Errorf("Expected %v, but got %v", expected, output)
	}
}

func TestAll(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	result := slicesutils.All(input, func(item int) bool {
		return item < 10
	})

	if !result {
		t.Errorf("Expected true, but got false")
	}

	result = slicesutils.All(input, func(item int) bool {
		return item < 5
	})

	if result {
		t.Errorf("Expected false, but got true")
	}
}

func TestAny(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	result := slicesutils.Any(input, func(item int) bool {
		return item == 5
	})

	if !result {
		t.Errorf("Expected true, but got false")
	}

	result = slicesutils.Any(input, func(item int) bool {
		return item == 10
	})

	if result {
		t.Errorf("Expected false, but got true")
	}
}

func TestFilter(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	expected := []int{2, 4, 6, 8}

	result := slicesutils.Filter(input, func(item int) bool {
		return item%2 == 0
	})

	if ok := slicesutils.Compare(expected, result); !ok {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestDistinct(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 3}
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	result := slicesutils.Distinct(input)

	if ok := slicesutils.Compare(expected, result); !ok {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

type IdentifiableItem struct {
	ID   int
	Type string
}

func (i IdentifiableItem) Id() int {
	return i.ID
}
func TestUniqueItemsById(t *testing.T) {
	input := []IdentifiableItem{
		{ID: 1},
		{ID: 2},
		{ID: 3},
		{ID: 4},
		{ID: 5},
		{ID: 6},
		{ID: 7},
		{ID: 8},
		{ID: 9},
		{ID: 1},
		{ID: 2},
		{ID: 3},
	}

	expected := []IdentifiableItem{
		{ID: 1},
		{ID: 2},
		{ID: 3},
		{ID: 4},
		{ID: 5},
		{ID: 6},
		{ID: 7},
		{ID: 8},
		{ID: 9},
	}

	result := slicesutils.UniqueItemsById(input)

	if ok := slicesutils.Compare(expected, result); !ok {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestDifference(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	other := []int{1, 2, 3, 4, 5}
	expected := []int{6, 7, 8, 9}

	result := slicesutils.Difference(input, other)

	if ok := slicesutils.Compare(expected, result); !ok {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

type StringWeigh string

const (
	Light StringWeigh = "Light"
	Heavy StringWeigh = "Heavy"
)

func TestWeightedSort(t *testing.T) {
	input := []IdentifiableItem{
		{ID: 1, Type: "A"},
		{ID: 2, Type: "B"},
		{ID: 3, Type: "A"},
		{ID: 4, Type: "B"},
		{ID: 5, Type: "A"},
		{ID: 6, Type: "B"},
	}

	expected := []IdentifiableItem{
		{ID: 1, Type: "A"},
		{ID: 3, Type: "A"},
		{ID: 5, Type: "A"},
		{ID: 2, Type: "B"},
		{ID: 4, Type: "B"},
		{ID: 6, Type: "B"},
	}

	weightsMap := map[string]StringWeigh{
		"A": Heavy,
		"B": Light,
	}
	getWeight := func(item IdentifiableItem) StringWeigh {
		return weightsMap[item.Type]
	}

	result := slicesutils.WeightedSort(input, getWeight, func(a, b IdentifiableItem) bool {
		return a.ID < b.ID
	})

	if ok := slicesutils.Compare(expected, result); !ok {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
