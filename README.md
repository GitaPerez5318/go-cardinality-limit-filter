# go-cardinality-limit-filter

A small Go library that tracks label values for a metric and drops new values once a configured cardinality limit is reached.

## Usage

```go
package main

import (
	"fmt"

	cardinalitylimitfilter "github.com/yourname/go-cardinality-limit-filter"
)

func main() {
	filter := cardinalitylimitfilter.NewCardinalityLimitFilter(3)

	values := []string{"a", "b", "c", "d", "a"}
	for _, v := range values {
		if filter.Allow(v) {
			fmt.Printf("allowed: %s\n", v)
		} else {
			fmt.Printf("dropped: %s\n", v)
		}
	}

	fmt.Println("count:", filter.Count())
	fmt.Println("exceeded:", filter.Exceeded())
}
```

Output:

```
allowed: a
allowed: b
allowed: c
dropped: d
allowed: a
count: 3
exceeded: true
```

## Why this exists

Metric systems often attach a label such as `user_id` or `request_path` to every sample. If the set of label values is unbounded, the number of time series grows without limit and memory usage follows. A cardinality limit caps the number of distinct values that are tracked; values beyond the limit are reported as rejected so the caller can drop them or aggregate them into an `other` bucket.

The trade-off is explicit: once the limit is reached, later values are rejected even if they are common. The library does not attempt to estimate which values are important; it simply counts distinct values in first-seen order. This keeps the implementation small and predictable.

## Edge cases

A limit of zero or less disables the limit: every value is allowed and nothing is tracked. A repeated value is always allowed, even after the limit has been reached, because it does not increase cardinality. The filter is safe for concurrent use.
