package slicesutils

import (
	"cmp"
	"runtime"
	"sort"
	"sync"
)

// Max returns the maximum value in the provided slice.
// If no elements are provided, it panics with "No element provided to Max".
func Max(elements ...int) int {
	if len(elements) == 0 {
		panic("No element provided to Max")
	}

	maxValue := elements[0]
	for _, num := range elements {
		if num > maxValue {
			maxValue = num
		}
	}
	return maxValue
}

// ParallelMap applies the given map function concurrently to each element in the input slice.
// It creates a fixed number of worker goroutines to process the elements in parallel.
// The input slice is divided into chunks and each chunk is processed by a worker goroutine.
// The results are collected and returned as a new slice in the same order as the input.
// The map function takes an element of type T as input and returns an element of type U.
// The number of worker goroutines is determined by the number of available CPU cores.
// This function blocks until all worker goroutines have completed their tasks.
func ParallelMap[I any, O any, S ~[]I](inputSlice S, mapFunc func(I) O) []O {
	if inputSlice == nil {
		return []O{}
	}

	outputSlice := make([]O, len(inputSlice))
	numWorkers := runtime.NumCPU()
	if len(inputSlice) < numWorkers {
		numWorkers = len(inputSlice)
	}

	var wg sync.WaitGroup

	inputChan := make(chan int, len(inputSlice))

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range inputChan {
				outputSlice[idx] = mapFunc(inputSlice[idx])
			}
		}()
	}

	// Send index to workers
	for i := range inputSlice {
		inputChan <- i
	}
	close(inputChan)

	wg.Wait()

	return outputSlice
}

// Map applies a mapping function to each element of the input slice and returns
// a new slice containing the results.
func Map[I any, O any, S ~[]I](inputSlice S, mapFunc func(I) O) []O {
	outputSlice := make([]O, len(inputSlice))

	for i, input := range inputSlice {
		outputSlice[i] = mapFunc(input)
	}

	return outputSlice
}

// SafeMap applies a mapping function to each element of an input slice, returning a new slice
// with the results. If the mapping function returns an error for any element or panics, SafeMap will
// return that error and halt further processing.
func SafeMap[I any, O any, S ~[]I](inputSlice S, mappingFunc func(I) (O, error)) ([]O, error) {
	outputSlice := make([]O, len(inputSlice))

	for i, input := range inputSlice {
		output, err := SafeExcecute(func() (out O, errAux error) {
			out, errAux = mappingFunc(input)
			return
		})

		if err != nil {
			return nil, err
		}
		outputSlice[i] = output
	}

	return outputSlice, nil
}

// Filter applies a filter function to each element in the inputSlice and returns a new slice
// containing only the elements for which the filter function returns true.
// The filter function takes an element of type T as input and returns a boolean value.
// The inputSlice is not modified.
func Filter[I any, S ~[]I](inputSlice S, filterFunc func(I) bool) S {
	newSliceLen := 0

	for _, input := range inputSlice {
		if filterFunc(input) {
			inputSlice[newSliceLen] = input
			newSliceLen++
		}
	}

	return inputSlice[:newSliceLen]
}

// Reduce applies a function to each element of the input slice and returns a single value.
// The reduceFunc function takes two arguments: an accumulator value of type U and an element of the input slice of type I.
// It returns a new accumulator value of type O.
// The initialValue is the initial value of the accumulator.
// Reduce iterates over the inputSlice, applying the reduceFunc function to each element and the current accumulator value.
// The result of each iteration becomes the new accumulator value for the next iteration.
// Finally, the function returns the final accumulator value.
func Reduce[I any, O any, S ~[]I](inputSlice S, reduceFunc func(O, I) O, initialValue O) O {
	accumulator := initialValue

	for _, input := range inputSlice {
		accumulator = reduceFunc(accumulator, input)
	}

	return accumulator
}

// SafeReduce is a generic function that safely reduces a slice of input elements
// into a single output value by applying a user-defined reduce function. It ensures
// that if an error is encountered during the reduction process, the reduce stops and returns the error.
func SafeReduce[I any, O any, S ~[]I](inputSlice S, reduceFunc func(O, I) (O, error), initialValue O) (O, error) {
	accumulator := initialValue

	for _, input := range inputSlice {
		accumAux, err := SafeExcecute(func() (out O, errAux error) {
			out, errAux = reduceFunc(accumulator, input)
			return
		})

		if err != nil {
			return accumAux, err
		}
		accumulator = accumAux
	}

	return accumulator, nil
}

// Sort sorts a slice of any type in place based on the provided less function.
// The less function should return true if the first argument is considered to be less than the second.
func Sort[I any, S ~[]I](slice S, less func(i, j I) bool) S {
	sort.Slice(slice, func(i, j int) bool {
		return less(slice[i], slice[j])
	})
	return slice
}

func Reverse[I any, S ~[]I](slice S) S {
	for i := 0; i <= len(slice)/2; i++ {
		j := len(slice) - i - 1
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

// WeightedSort sorts a slice of any type based on a weight function and a less function.
// The weight function determines the primary sorting order by returning an integer weight for each element.
// The less function is used as a secondary sorting order when two elements have the same weight.
func WeightedSort[I any, W cmp.Ordered, S ~[]I](slice S, getWeighfn func(I) W, less func(i, j I) bool) S {
	sort.Slice(slice, func(i, j int) bool {
		weightI := getWeighfn(slice[i])
		weightJ := getWeighfn(slice[j])

		if weightI != weightJ {
			return weightI < weightJ
		}

		return less(slice[i], slice[j])
	})
	return slice
}

// RemoveElement removes the first n occurrences of the specified element from the given slice.
// Passing nil to occurrencesToDelete will remove all occurrences of the element.
// It returns a new slice with the element removed, or the original slice if the element is not found.
func RemoveElement[I comparable, S ~[]I](slice S, element I, occurrencesToDelete *int) S {
	if len(slice) == 0 {
		return slice
	}

	limit := -1 // Default to removing all occurrences
	if occurrencesToDelete != nil {
		if *occurrencesToDelete <= 0 {
			return slice
		}
		limit = *occurrencesToDelete
	}

	count := 0

	newSliceLen := 0
	for _, e := range slice {
		if e == element && (limit == -1 || count < limit) {
			count++
			continue
		}
		slice[newSliceLen] = e
		newSliceLen++
	}

	return slice[:newSliceLen]
}

// RemoveFirstOccurrence removes the first occurrence of the specified element from the given slice.
// It's a shorthand for calling RemoveElement with occurrencesToDelete set to 1.
func RemoveFirstOccurrence[I comparable, S ~[]I](slice S, element I) S {
	ocurrencesToDelete := 1
	return RemoveElement(slice, element, &ocurrencesToDelete)
}

// RemoveElements removes all occurrences of the specified elements from the given slice.
// It returns a new slice with the elements removed.
//
// The function uses a map to keep track of the elements to be removed, ensuring efficient lookups.
func RemoveElements[I comparable, S ~[]I](slice S, elements ...I) S {
	itemsToRemoveMap := make(map[I]struct{}, len(slice))
	for _, e := range elements {
		itemsToRemoveMap[e] = struct{}{}
	}

	newSliceLen := 0
	for _, e := range slice {
		if _, found := itemsToRemoveMap[e]; found {
			continue
		}
		slice[newSliceLen] = e
		newSliceLen++
	}

	return slice[:newSliceLen]
}

// ParallelForEach applies a given function to each element of the input slice in parallel.
// The number of parallel workers is determined by the minimum of the number of CPU cores
// and the length of the input slice.
//
// Type Parameters:
//
//	T: The type of elements in the input slice.
//
// Parameters:
//
//	inputSlice: The slice of elements to process. If nil, the function returns immediately.
//	forEachFunc: The function to apply to each element of the input slice.
//
// Example usage:
//
//	ParallelForEach([]int{1, 2, 3, 4}, func(n int) {
//	    fmt.Println(n)
//	})
func ParallelForEach[I any, S ~[]I](inputSlice S, forEachFunc func(I)) {
	if inputSlice == nil {
		return
	}

	numWorkers := runtime.NumCPU()
	if len(inputSlice) < numWorkers {
		numWorkers = len(inputSlice)
	}

	var wg sync.WaitGroup

	inputChan := make(chan I, len(inputSlice))

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for input := range inputChan {
				forEachFunc(input)
			}
		}()
	}

	// Send input to workers
	for _, input := range inputSlice {
		inputChan <- input
	}
	close(inputChan)

	wg.Wait()
}

// Find searches for an element in the inputSlice that satisfies the given findFunc.
// It returns the first element that matches the condition or the zero value of type T if no match is found.
func Find[I any, S ~[]I](inputSlice S, findFunc func(I) bool) (foundItem I, didFind bool) {
	for _, input := range inputSlice {
		if findFunc(input) {
			return input, true
		}
	}
	var zero I
	return zero, false
}

func SafeFind[I any, S ~[]I](inputSlice S, findFunc func(I) (bool, error)) (foundItem I, didFind bool, err error) {
	for _, input := range inputSlice {

		didFind, err := SafeExcecute(func() (out bool, errAux error) {
			out, errAux = findFunc(input)
			return
		})

		if err != nil {
			return foundItem, false, err
		}

		if didFind {
			return input, true, nil
		}
	}
	return foundItem, false, nil
}

// FindIndex returns the index of the first element in the inputSlice that satisfies the findFunc condition.
// If no element satisfies the condition, it returns -1.
func FindIndex[I any, S ~[]I](inputSlice S, findFunc func(I) bool) int {
	for i, input := range inputSlice {
		if findFunc(input) {
			return i
		}
	}
	return -1
}

// Contains checks if the given element is present in the slice.
// It returns true if the element is found, otherwise it returns false.
func Contains[I comparable, S ~[]I](slice S, element I) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}

// All checks if all elements in the given slice satisfy the provided predicate function.
// It returns true if all elements satisfy the predicate, otherwise it returns false.
func All[I any, S ~[]I](slice S, predicate func(I) bool) bool {
	for _, item := range slice {
		if !predicate(item) {
			return false
		}
	}
	return true
}

// Any checks if any element in the slice satisfies the given predicate function.
// It returns true if at least one element matches the predicate, otherwise false.
func Any[I any, S ~[]I](slice S, predicate func(I) bool) bool {
	for _, item := range slice {
		if predicate(item) {
			return true
		}
	}
	return false
}

// Chunk splits a slice into multiple smaller slices (chunks) of a specified size.
// If the chunkSize is less than or equal to 0, or if the input slice is empty,
// it returns an empty slice of slices.
//
// The function uses generics to work with slices of any type.
//
// Parameters:
//   - slice: The input slice to be split into chunks.
//   - chunkSize: The desired size of each chunk.
//
// Returns:
//
//	A slice of slices, where each inner slice is a chunk of the original slice.
func Chunk[I any, S ~[]I](slice S, chunkSize int) []S {
	if chunkSize <= 0 || len(slice) == 0 {
		return []S{}
	}

	// Pre-allocate the result slice with the expected number of chunks
	numChunks := (len(slice) + chunkSize - 1) / chunkSize // Ceiling of len(slice)/chunkSize
	chunks := make([]S, 0, numChunks)

	for i := 0; i < len(slice); i += chunkSize {
		// Safely calculate the end of the chunk
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

// Compare takes two slices of any comparable type and returns true if they are equal.
// Two slices are considered equal if they have the same length and all corresponding
// elements are equal.
//
// Returns true if the slices are equal, false otherwise.
func Compare[I comparable, S ~[]I](a, b S) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Distinct returns a new slice containing only the distinct elements from the input slice.
// The order of elements in the result slice is the same as their first occurrence in the input slice.
func Distinct[I comparable, S ~[]I](slice S) S {
	seenItems := make(map[I]struct{})

	newSliceLen := 0
	for _, item := range slice {
		if _, seen := seenItems[item]; seen {
			continue
		}
		seenItems[item] = struct{}{}
		slice[newSliceLen] = item
		newSliceLen++
	}

	return slice[:newSliceLen]
}

type identifiable[T any] interface {
	Id() T
}

// UniqueItemsById returns a slice containing only the unique items from the input slice,
// where uniqueness is determined by the item's Id. The function uses a map to track
// seen Ids and filters out duplicates. Items should implement the identifiable interface.
// Thus they must have a method Id() that returns a unique identifier.
func UniqueItemsById[Id comparable, I identifiable[Id], S ~[]I](slice S) S {
	seenItems := make(map[Id]struct{})

	newSliceLen := 0
	for _, item := range slice {
		id := item.Id()
		if _, seen := seenItems[id]; seen {
			continue
		}

		seenItems[id] = struct{}{}
		slice[newSliceLen] = item
		newSliceLen++
	}

	return slice[:newSliceLen]
}

// Intersection returns the common elements between two slices.
// It takes two slices of any comparable type and returns a slice containing
// the elements that are present in both input slices.
func Intersection[I comparable, S ~[]I](a, b S) S {
	set := make(map[I]struct{})
	for _, item := range a {
		set[item] = struct{}{}
	}

	var result S
	for _, item := range b {
		if _, ok := set[item]; ok {
			result = append(result, item)
		}
	}

	return result
}

// Union returns the union of two slices, removing duplicate elements.
// The function takes two slices of any comparable type and returns a new slice
// containing all unique elements from both input slices.
//
// Example usage:
//
//	a := []int{1, 2, 3}
//	b := []int{3, 4, 5}
//	result := Union(a, b) // result will be []int{1, 2, 3, 4, 5}
//
// The order of elements in the resulting slice is not guaranteed.
func Union[I comparable, S ~[]I](a, b S) S {
	set := make(map[I]struct{})
	for _, item := range a {
		set[item] = struct{}{}
	}
	for _, item := range b {
		set[item] = struct{}{}
	}

	var result S
	for item := range set {
		result = append(result, item)
	}

	return result
}

// Difference returns the elements in slice `a` that are not in slice `b`.
// It uses a map to track the elements in `b` for efficient lookups.
//
// Type Parameters:
//
//	T: a type that is comparable (supports the == and != operators).
//
// Parameters:
//
//	a: the first slice of elements.
//	b: the second slice of elements.
//
// Returns:
//
//	A slice containing the elements that are in `a` but not in `b`.
func Difference[I comparable, S ~[]I](a, b S) S {
	set := make(map[I]struct{})
	for _, item := range b {
		set[item] = struct{}{}
	}

	newSliceLen := 0
	for _, item := range a {
		if _, exists := set[item]; exists {
			continue
		}

		a[newSliceLen] = item
		newSliceLen++
	}

	return a[:newSliceLen]
}
