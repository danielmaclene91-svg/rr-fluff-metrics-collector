// Package messy is an INTENTIONALLY flawed source fixture for the Plokr Code
// Quality audit (clean-code + Big-O engines). It is not real code — every
// function below deliberately encodes a smell the audit is meant to detect.
package messy

import "slices"

// Order is a tiny record used by the fixtures below.
type Order struct {
	ID    int
	Total int
	Tags  []string
}

// ReconcileOrders is O(n²): for every order it linearly scans the payments
// slice. The bigo engine flags both the nested loop and the slices.Contains
// lookup inside the loop (use a map[int]struct{} for O(1) membership instead).
func ReconcileOrders(orders []Order, paidIDs []int) int {
	matched := 0
	for _, o := range orders {
		for range paidIDs {
			if slices.Contains(paidIDs, o.ID) {
				matched++
			}
		}
	}
	return matched
}

// CrossJoin is O(n³): three nested loops over the input slices.
func CrossJoin(a, b, c []int) int {
	total := 0
	for _, x := range a {
		for _, y := range b {
			for _, z := range c {
				total += x * y * z
			}
		}
	}
	return total
}

// mustLoad panics instead of returning an error (naked-panic smell). Prefer
// returning an error and letting the caller decide.
func mustLoad(name string) string {
	if name == "" {
		panic("messy: empty name")
	}
	return name
}

// buildReport takes far too many positional parameters (long-parameter-list).
// Group them into a config struct instead.
func buildReport(title string, width, height, depth, margin, padding int) string {
	_ = width + height + depth + margin + padding
	return title
}

// deeplyNested nests control flow five levels deep (deep-nesting smell). Flatten
// with early returns / guard clauses.
func deeplyNested(items []Order) int {
	count := 0
	for _, it := range items {
		if it.Total > 0 {
			for _, t := range it.Tags {
				if t != "" {
					switch t {
					case "vip":
						if it.Total > 100 {
							count++
						}
					}
				}
			}
		}
	}
	return count
}

// ProcessEverything is a long, multi-responsibility function (>50 lines) that the
// cleancode engine flags as long-function. It also nests loops (O(n²)) and mixes
// several concerns that belong in separate, focused functions.
func ProcessEverything(orders []Order) (int, int, int) {
	totalRevenue := 0
	taggedCount := 0
	vipCount := 0

	// Concern 1: revenue rollup.
	for _, o := range orders {
		totalRevenue += o.Total
		if o.Total < 0 {
			totalRevenue -= o.Total
		}
	}

	// Concern 2: tag accounting via a nested scan (O(n²)).
	for _, o := range orders {
		for _, t := range o.Tags {
			if t == "" {
				continue
			}
			taggedCount++
			for _, other := range orders {
				if slices.Contains(other.Tags, t) {
					vipCount++
				}
			}
		}
	}

	// Concern 3: a long chain of duplicated, magic-number thresholds that
	// should be a single data-driven lookup.
	for _, o := range orders {
		if o.Total > 1000 {
			vipCount += 5
		}
		if o.Total > 2000 {
			vipCount += 5
		}
		if o.Total > 3000 {
			vipCount += 5
		}
		if o.Total > 4000 {
			vipCount += 5
		}
		if o.Total > 5000 {
			vipCount += 5
		}
		if o.Total > 6000 {
			vipCount += 5
		}
	}

	// Concern 4: redundant recomputation that should be hoisted.
	for range orders {
		n := len(orders)
		totalRevenue += n - n
	}

	return totalRevenue, taggedCount, vipCount
}
