package cardinalitylimitfilter

import (
	"sync"
)

// CardinalityLimitFilter tracks label values for a given metric and label name.
// It enforces a maximum cardinality: once the limit is reached, any new label
// value is considered "exceeded" and the caller can drop or aggregate it.
//
type CardinalityLimitFilter struct {
	mu       sync.Mutex
	limit    int
	seen     map[string]struct{}
	exceeded bool
}

// NewCardinalityLimitFilter creates a filter with the given maximum number of
// distinct label values. A limit of zero or less means no cardinality limit is
// enforced (all values are allowed).
func NewCardinalityLimitFilter(limit int) *CardinalityLimitFilter {
	return &CardinalityLimitFilter{
		limit: limit,
		seen:  make(map[string]struct{}),
	}
}

// Allow reports whether the given label value should be allowed under the
// cardinality limit. The first time a value is seen it is counted; subsequent
// calls with the same value always return true. Once the number of distinct
// values reaches the limit, any new value returns false and the filter enters an
// "exceeded" state.
func (f *CardinalityLimitFilter) Allow(value string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.limit <= 0 {
		return true
	}

	if _, ok := f.seen[value]; ok {
		return true
	}

	if len(f.seen) >= f.limit {
		f.exceeded = true
		return false
	}

	f.seen[value] = struct{}{}
	return true
}

// Exceeded reports whether any label value has been rejected because the
// cardinality limit was reached.
func (f *CardinalityLimitFilter) Exceeded() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.exceeded
}

// Count returns the number of distinct label values currently tracked. It never
// exceeds the configured limit.
func (f *CardinalityLimitFilter) Count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.seen)
}
