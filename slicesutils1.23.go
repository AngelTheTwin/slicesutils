//go:build go1.23
// +build go1.23

package slicesutils

import (
	"cmp"
	"iter"
	"math"
)

func MaxSeq[I cmp.Ordered](inputSeq iter.Seq[I]) I {
	next, stop := iter.Pull(inputSeq)

	defer stop()

	first, ok := next()
	if !ok {
		panic("MaxSeq: empty sequence")
	}
	mx := first
	for nextItem, ok := next(); ok; nextItem, ok = next() {
		mx = max(mx, nextItem)
	}

	return mx
}

func MaxSeqFunc[I any](inputSeq iter.Seq[I], maxFunc func(I, I) I) I {
	next, stop := iter.Pull(inputSeq)

	defer stop()

	first, ok := next()
	if !ok {
		panic("MaxSeq: empty sequence")
	}
	mx := first
	for nextItem, ok := next(); ok; nextItem, ok = next() {
		mx = maxFunc(mx, nextItem)
	}

	return mx
}

func MapSeq[I any, O any](inputSeq iter.Seq[I], mapFunc func(I) O) iter.Seq[O] {
	return func(yield func(O) bool) {
		for input := range inputSeq {
			if !yield(mapFunc(input)) {
				return
			}
		}
	}
}

func SafeMapSeq[I any, O any](inputSeq iter.Seq[I], mapFunc func(I) (O, error)) iter.Seq[O] {
	return func(yield func(O) bool) {
		for input := range inputSeq {
			out, errAux := SafeExcecute(func() (O, error) {
				return mapFunc(input)
			})
			if errAux != nil {
				return
			}
			if !yield(out) {
				return
			}
		}
	}
}

func FilterSeq[I any](inputSeq iter.Seq[I], filterFunc func(I) bool) iter.Seq[I] {
	return func(yield func(I) bool) {
		for input := range inputSeq {
			if filterFunc(input) && !yield(input) {
				return
			}
		}
	}
}

func ReduceSeq[I any, O any](inputSeq iter.Seq[I], reduceFunc func(O, I) O, initialValue O) O {
	result := initialValue
	for input := range inputSeq {
		result = reduceFunc(result, input)
	}
	return result
}

func SafeReduceSeq[I any, O any](inputSeq iter.Seq[I], reduceFunc func(O, I) (O, error), initialValue O) (O, error) {
	result := initialValue
	for input := range inputSeq {
		accumAux, err := SafeExcecute(func() (O, error) {
			return reduceFunc(result, input)
		})

		if err != nil {
			return result, err
		}
		result = accumAux
	}
	return result, nil
}

// RemoveElementSeq returns a sequence that yields the elements of inputSeq
// with at most n occurrences of element removed.
//
//	n == 0    → nothing is removed
//	n  > 0    → remove the first n matches, then pass the rest through
//	n == -1   → remove ALL matches
//
// The consumer can stop early by returning false from yield.
func RemoveElementSeq[I comparable](inputSeq iter.Seq[I], element I, occurrencesToDelete int) iter.Seq[I] {
	limit := occurrencesToDelete
	if limit == -1 {
		limit = math.MaxInt
	}
	return func(yield func(I) bool) {
		for input := range inputSeq {
			if input == element {
				if limit > 0 {
					limit--
					continue
				}
			}
			if !yield(input) {
				return
			}
		}
	}
}

func RemoveElementsSeq[I comparable](inputSeq iter.Seq[I], elements ...I) iter.Seq[I] {
	itemsToRemoveMap := Reduce(elements, func(accum map[I]bool, curr I) map[I]bool {
		accum[curr] = true
		return accum
	}, map[I]bool{})

	return func(yield func(I) bool) {
		for input := range inputSeq {
			if _, ok := itemsToRemoveMap[input]; ok {
				continue
			}
			if !yield(input) {
				return
			}
		}
	}
}

func FindSeq[I any](inputSeq iter.Seq[I], findFunc func(I) bool) (foundItem I, didFind bool) {
	for input := range inputSeq {
		if findFunc(input) {
			foundItem = input
			didFind = true
			break
		}
	}

	return foundItem, didFind
}

func SafeFindSeq[I any](inputSeq iter.Seq[I], findFunc func(I) (bool, error)) (foundItem I, didFind bool, err error) {
	for input := range inputSeq {
		foundAux, errAux := SafeExcecute(func() (bool, error) {
			return findFunc(input)
		})

		if errAux != nil {
			return foundItem, didFind, errAux
		}

		if foundAux {
			foundItem = input
			didFind = true
			break
		}
	}

	return foundItem, didFind, nil
}

func FindIndexSeq[I any](inputSeq iter.Seq[I], findFunc func(I) bool) (foundIndex int, didFind bool) {
	index := 0
	for input := range inputSeq {
		if findFunc(input) {
			foundIndex = index
			didFind = true
			break
		}
		index++
	}

	return foundIndex, didFind
}

func ContainsSeq[I comparable](inputSeq iter.Seq[I], element I) bool {
	for input := range inputSeq {
		if input == element {
			return true
		}
	}
	return false
}

func AllSeq[I any](inputSeq iter.Seq[I], allFunc func(I) bool) bool {
	for input := range inputSeq {
		if !allFunc(input) {
			return false
		}
	}
	return true
}

func AnySeq[I any](inputSeq iter.Seq[I], anyFunc func(I) bool) bool {
	for input := range inputSeq {
		if anyFunc(input) {
			return true
		}
	}
	return false
}

func DistinctSeq[I comparable](inputSeq iter.Seq[I]) iter.Seq[I] {
	seen := make(map[I]bool)
	return func(yield func(I) bool) {
		for input := range inputSeq {
			if _, ok := seen[input]; !ok {
				seen[input] = true
				if !yield(input) {
					return
				}
			}
		}
	}
}

func Ennumerate[I any](inputSeq iter.Seq[I]) iter.Seq2[int, I] {
	return func(yield func(int, I) bool) {
		index := 0
		for input := range inputSeq {
			if !yield(index, input) {
				return
			}
			index++
		}
	}
}

func IntersectionSeq[I comparable](inputSeq1, inputSeq2 iter.Seq[I]) iter.Seq[I] {
	seen := make(map[I]bool)
	return func(yield func(I) bool) {
		for input := range inputSeq1 {
			seen[input] = true
		}
		for input := range inputSeq2 {
			if _, ok := seen[input]; ok {
				if !yield(input) {
					return
				}
			}
		}
	}
}

func UnionSeq[I comparable](inputSeq1, inputSeq2 iter.Seq[I]) iter.Seq[I] {
	seen := make(map[I]bool)
	return func(yield func(I) bool) {
		for input := range inputSeq1 {
			if _, ok := seen[input]; !ok {
				seen[input] = true
				if !yield(input) {
					return
				}
			}
		}
		for input := range inputSeq2 {
			if _, ok := seen[input]; !ok {
				seen[input] = true
				if !yield(input) {
					return
				}
			}
		}
	}
}

func DifferenceSeq[I comparable](a, b iter.Seq[I]) iter.Seq[I] {
	seen := make(map[I]bool)
	return func(yield func(I) bool) {
		for input := range b {
			seen[input] = true
		}
		for input := range a {
			if _, ok := seen[input]; !ok {
				if !yield(input) {
					return
				}
			}
		}
	}
}

func CompareSeq[I comparable](a, b iter.Seq[I]) bool {
	nextA, stopA := iter.Pull(a)
	nextB, stopB := iter.Pull(b)
	defer stopA()
	defer stopB()

	for {
		currA, okA := nextA()
		currB, okB := nextB()

		if okA != okB {
			return false
		}

		if okA {
			if currA != currB {
				return false
			}
		} else {
			return true
		}
	}
}

func GroupBySeq[I any, K comparable](inputSeq iter.Seq[I], keyFunc func(I) K) iter.Seq2[K, iter.Seq[I]] {
	groups := make(map[K][]I)

	for item := range inputSeq {
		key := keyFunc(item)
		groups[key] = append(groups[key], item)
	}

	// Step 2: Yield each group only once
	return func(yield func(K, iter.Seq[I]) bool) {
		for key, items := range groups {
			seq := func(yieldItem func(I) bool) {
				for _, item := range items {
					if !yieldItem(item) {
						return
					}
				}
			}
			if !yield(key, seq) {
				return
			}
		}
	}
}
