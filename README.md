# `slicesutils` – Utility Functions for Slices in Go (Legacy) (New functions for iter.Seq go >=1.23 just added!)


A lightweight package providing utility functions for working with slices in Go, designed for versions prior to the introduction of the [`slices`](https://pkg.go.dev/slices) package in the standard library.

This package uses **generics** (introduced in Go 1.18) to provide type-safe and flexible utilities. While it's considered **legacy** and intended for projects that require backwards compatibility with earlier standard library features, some functions remain useful even in newer versions of Go.

## Examples

### Max
```go
import "bitbucket.org/yourteam/slicesutils"

max := slicesutils.Max(3, 1, 9, 5)
// max == 9
```

### Map
```go
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

## Sequence Utilities (Go 1.23+)

### MaxSeq
```go
import "bitbucket.org/yourteam/slicesutils"
import "iter"

seq := iter.Seq([]int{3, 1, 9, 5})
max := slicesutils.MaxSeq(seq)
// max == 9
```

### MapSeq
```go
import "bitbucket.org/yourteam/slicesutils"
import "iter"

seq := iter.Seq([]int{1, 2, 3})
doubledSeq := slicesutils.MapSeq(seq, func(n int) int {
    return n * 2
})
// doubledSeq yields 2, 4, 6
```

### FilterSeq
```go
import "bitbucket.org/yourteam/slicesutils"
import "iter"

seq := iter.Seq([]int{1, 2, 3, 4})
evensSeq := slicesutils.FilterSeq(seq, func(n int) bool {
    return n%2 == 0
})
// evensSeq yields 2, 4
```

### ReduceSeq
```go
import "bitbucket.org/yourteam/slicesutils"
import "iter"

seq := iter.Seq([]int{1, 2, 3, 4})
sum := slicesutils.ReduceSeq(seq, func(acc, n int) int {
    return acc + n
}, 0)
// sum == 10
```

## Open Source and Contributions
This package is open source and maintained with the intention of being useful for the Go community, especially when working with legacy code or older Go versions. While some functions may overlap with newer standard packages, this library serves as a practical alternative and is open to improvement. Contributions, suggestions, and issue reports are all welcome!
