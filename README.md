# `slicesutils` – Utility Functions for Slices in Go (Legacy)

A lightweight package providing utility functions for working with slices in Go, designed for versions prior to the introduction of the [`slices`](https://pkg.go.dev/slices) package in the standard library.

This package is considered **legacy** and is intended for projects that require backwards compatibility with earlier Go versions.

Some functions are still usefull in newer Go versions.

## Examples

### Max
```go
import "bitbucket.org/yourteam/slicesutils"

max := slicesutils.Max(3, 1, 9, 5)
// max == 9
doubled := slicesutils.Map([]int{1, 2, 3}, func(n int) int {
    return n * 2
})
// doubled == []int{2, 4, 6}
```

### Filter
```go
evens := slicesutils.Filter([]int{1, 2, 3, 4}, func(n int) bool {
    return n%2 == 0
})
// evens == []int{2, 4}
```

### Reduce
```go
sum := slicesutils.Reduce([]int{1, 2, 3, 4}, func(acc, n int) int {
    return acc + n
}, 0)
// sum == 10
```

### ParallelForEach
```go
slicesutils.ParallelForEach([]string{"a", "b", "c"}, func(s string) {
    fmt.Println(s) // Some expensive operations
})
```

### Chunk
```go
chunks := slicesutils.Chunk([]int{1, 2, 3, 4, 5}, 2)
// chunks == [][]int{{1, 2}, {3, 4}, {5}}
```

### Distinct
```go
unique := slicesutils.Distinct([]int{1, 2, 2, 3, 3, 3})
// unique == []int{1, 2, 3}
```

## Open Source and Contributions
This package is open source and maintained with the intention of being useful for the Go community, especially when working with legacy code or older Go versions. While some functions may overlap with newer standard packages, this library serves as a practical alternative and is open to improvement. Contributions, suggestions, and issue reports are all welcome!