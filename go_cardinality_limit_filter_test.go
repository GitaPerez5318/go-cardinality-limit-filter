package cardinalitylimitfilter

import "testing"

func TestAllowWithinLimit(t *testing.T) {
	f := NewCardinalityLimitFilter(3)

	for _, v := range []string{"a", "b", "c"} {
		if !f.Allow(v) {
			t.Fatalf("Allow(%q) = false, want true", v)
		}
	}
	if f.Count() != 3 {
		t.Fatalf("Count() = %d, want 3", f.Count())
	}
	if f.Exceeded() {
		t.Fatal("Exceeded() = true, want false")
	}
}

func TestAllowRepeatedValue(t *testing.T) {
	f := NewCardinalityLimitFilter(2)

	if !f.Allow("x") {
		t.Fatal("first Allow(x) = false, want true")
	}
	if !f.Allow("x") {
		t.Fatal("second Allow(x) = false, want true")
	}
	if f.Count() != 1 {
		t.Fatalf("Count() = %d, want 1", f.Count())
	}
}

func TestAllowExceedsLimit(t *testing.T) {
	f := NewCardinalityLimitFilter(2)

	if !f.Allow("a") {
		t.Fatal("Allow(a) = false, want true")
	}
	if !f.Allow("b") {
		t.Fatal("Allow(b) = false, want true")
	}
	if f.Allow("c") {
		t.Fatal("Allow(c) = true, want false (limit reached)")
	}
	if !f.Exceeded() {
		t.Fatal("Exceeded() = false, want true")
	}
	if f.Count() != 2 {
		t.Fatalf("Count() = %d, want 2", f.Count())
	}
}

func TestAllowRepeatedAfterExceeded(t *testing.T) {
	f := NewCardinalityLimitFilter(1)

	if !f.Allow("a") {
		t.Fatal("Allow(a) = false, want true")
	}
	if f.Allow("b") {
		t.Fatal("Allow(b) = true, want false")
	}
	if !f.Allow("a") {
		t.Fatal("Allow(a) = false, want true (already seen)")
	}
}

func TestZeroLimitAllowsEverything(t *testing.T) {
	f := NewCardinalityLimitFilter(0)

	for i := 0; i < 100; i++ {
		v := string(rune('a' + i%26))
		if !f.Allow(v) {
			t.Fatalf("Allow(%q) = false, want true with zero limit", v)
		}
	}
	if f.Exceeded() {
		t.Fatal("Exceeded() = true, want false with zero limit")
	}
	if f.Count() != 0 {
		t.Fatalf("Count() = %d, want 0 (no tracking with zero limit)", f.Count())
	}
}

func TestNegativeLimitAllowsEverything(t *testing.T) {
	f := NewCardinalityLimitFilter(-5)

	if !f.Allow("anything") {
		t.Fatal("Allow(anything) = false, want true with negative limit")
	}
	if f.Exceeded() {
		t.Fatal("Exceeded() = true, want false with negative limit")
	}
}

func TestEmptyValue(t *testing.T) {
	f := NewCardinalityLimitFilter(1)

	if !f.Allow("") {
		t.Fatal("Allow(\"\") = false, want true")
	}
	if f.Allow("other") {
		t.Fatal("Allow(other) = true, want false (limit reached)")
	}
}
