//go:build go1.23
// +build go1.23

package tests

import (
	"slices"
	"testing"

	"github.com/AngelTheTwin/slicesutils"
)

var itemsSeq = slices.Values(items)

func TestFindIndexSeq(t *testing.T) {
	index, didFind := slicesutils.FindIndexSeq(itemsSeq, func(item int) bool {
		return item == 5
	})

	if !didFind {
		t.Errorf("Expected to find index")
	}

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

func TestFindSeq(t *testing.T) {
	item, ok := slicesutils.FindSeq(itemsSeq, func(item int) bool {
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

func TestRemoveElement_OneOcurrence_Seq(t *testing.T) {
	expected := []int{1, 2, 3, 4, 6, 7, 8, 9, 10}

	elementsToRemove := 1
	result := slicesutils.RemoveElementSeq(itemsSeq, 5, elementsToRemove)

	if ok := slicesutils.CompareSeq(slices.Values(expected), result); !ok {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestRemoveElement_MultipleOcurrencesSeq(t *testing.T) {
	itemsSlice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 5, 5, 5}
	items := slices.Values(itemsSlice)
	expectedSlice := []int{1, 2, 3, 4, 6, 7, 8, 9, 10, 5, 5}
	expected := slices.Values(expectedSlice)

	elementsToRemove := 2
	result := slicesutils.RemoveElementSeq(items, 5, elementsToRemove)

	if ok := slicesutils.CompareSeq(expected, result); !ok {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestRemoveElement_AllOcurrencesSeq(t *testing.T) {
	itemsSlice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 5, 5, 5}
	expectedSlice := []int{1, 2, 3, 4, 6, 7, 8, 9, 10}
	items := slices.Values(itemsSlice)
	expected := slices.Values(expectedSlice)

	result := slicesutils.RemoveElementSeq(items, 5, -1)

	if ok := slicesutils.CompareSeq(expected, result); !ok {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestRemoveElementsSeq(t *testing.T) {
	itemsSlice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expectedSlice := []int{1, 2, 3, 6, 7, 8, 9, 10}
	items := slices.Values(itemsSlice)
	expected := slices.Values(expectedSlice)

	result := slicesutils.RemoveElementsSeq(items, 5, 4)

	if ok := slicesutils.CompareSeq(expected, result); !ok {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestAllSeq(t *testing.T) {
	result := slicesutils.AllSeq(itemsSeq, func(item int) bool {
		return item <= 10
	})

	if !result {
		t.Errorf("Expected true, but got false")
	}

	result = slicesutils.AllSeq(itemsSeq, func(item int) bool {
		return item < 5
	})

	if result {
		t.Errorf("Expected false, but got true")
	}
}

func TestAnySeq(t *testing.T) {
	result := slicesutils.AnySeq(itemsSeq, func(item int) bool {
		return item == 5
	})

	if !result {
		t.Errorf("Expected true, but got false")
	}

	result = slicesutils.AnySeq(itemsSeq, func(item int) bool {
		return item == 20
	})

	if result {
		t.Errorf("Expected false, but got true")
	}
}

func TestFilterSeq(t *testing.T) {
	expected := slices.Values([]int{2, 4, 6, 8, 10})

	result := slicesutils.FilterSeq(itemsSeq, func(item int) bool {
		return item%2 == 0
	})

	if ok := slicesutils.CompareSeq(expected, result); !ok {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestDistinctSeq(t *testing.T) {
	input := slices.Values([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 3})
	expected := slices.Values([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})

	result := slicesutils.DistinctSeq(input)

	if ok := slicesutils.CompareSeq(expected, result); !ok {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestGroupBySeq(t *testing.T) {
	result := slicesutils.GroupBySeq(
		itemsSeq,
		func(item int) string {
			if item%2 == 0 {
				return "even"
			}
			return "odd"
		},
	)

	for key, group := range result {
		for item := range group {
			if key == "even" && item%2 != 0 {
				t.Errorf("Expected even group, but got odd item %d", item)
			}
			if key == "odd" && item%2 == 0 {
				t.Errorf("Expected odd group, but got even item %d", item)
			}
		}
	}
}
