// Package search contains functions for approximate string matching.
package search

import (
	"errors"
	"os"
	"runtime"
	"strings"

	"golang.org/x/exp/slices"
)

// Options contains the information about search configuration.
type Options struct {
	InFile          bool
	CaseInsensitive bool
	Reverse         bool
	MatchLimit      int
	DistThreshold   int
}

// DefaultOptions returns a pointer to an instance of Options with default
// search configuration.
func DefaultOptions() *Options {
	return &Options{
		MatchLimit:    10,
		DistThreshold: 1,
	}
}

// result is used internally as an intermediate container for search results.
type result struct {
	pattern string
	matches []int
}

// Search performs a fuzzy search for patterns in the specified source (a string
// or a filename with the InFile option turned on). It returns search results in
// a map, where the key is a pattern and the value is a slice of integers
// (indices of a pattern's first character in the source text).
func Search(src string, patterns []string, options *Options) (map[string][]int, error) {
	if options == nil {
		options = DefaultOptions()
	}

	if options.MatchLimit < 0 || options.DistThreshold < 0 {
		return nil, errors.New("options -m, -t must be >= 0")
	}

	text := src
	if options.InFile {
		data, err := os.ReadFile(src)
		if err != nil {
			return nil, err
		}
		text = string(data)
	}

	if options.CaseInsensitive {
		text = strings.ToLower(text)
		for i := range patterns {
			patterns[i] = strings.ToLower(patterns[i])
		}
	}

	textRunes, jobsLen := []rune(text), len(patterns)

	jobs := make(chan string, jobsLen)
	results := make(chan result, jobsLen)

	workers := runtime.NumCPU()
	if jobsLen < workers {
		workers = jobsLen
	}

	for i := 0; i < workers; i++ {
		go searchWorker(textRunes, jobs, results, options)
	}

	for i := 0; i < jobsLen; i++ {
		jobs <- patterns[i]
	}
	close(jobs)

	searchResults := make([]result, jobsLen)
	for i := 0; i < jobsLen; i++ {
		searchResults[i] = <-results
	}

	matches := make(map[string][]int)
	for _, res := range searchResults {
		limit, count := options.MatchLimit, len(res.matches)
		if limit > count {
			limit = count
		}
		matches[res.pattern] = res.matches[:limit]
	}
	return matches, nil
}

// searchWorker is used internally to search for patterns supplied through the
// string channel 'patterns'. The search is performed in the 'src' rune
// slice. Search results are sent to the 'matches' channel.
func searchWorker(src []rune, patterns <-chan string, matches chan<- result, options *Options) {
	for pattern := range patterns {
		p := []rune(pattern)
		tmp := [][]int{}

		for i := range src {
			left, right := src[:i], src[i:]

			tmp = newDistance(tmp, left, p, options.DistThreshold, 0)
			tmp = newDistance(tmp, right, p, options.DistThreshold, i)

			if len(p)+i <= len(src) {
				window := src[i : len(p)+i]
				tmp = newDistance(tmp, window, p, options.DistThreshold, i)
			}
		}

		matched := []int{}
		for _, pair := range tmp {
			idx, dist := pair[0], pair[1]

			if dist <= options.DistThreshold {
				if !slices.Contains(matched, idx) {
					matched = append(matched, idx)
				}
			}
		}

		if options.Reverse {
			for i, j := 0, len(matched)-1; i < j; i, j = i+1, j-1 {
				matched[i], matched[j] = matched[j], matched[i]
			}
		}
		matches <- result{pattern, matched}
	}
}

// newDistance calculates the Levenshtein distance between two rune slices and
// appends it with a supplied index to the 'store' slice, but only if the
// resulting distance is less or equal to the 'max' threshold. The 'store'
// slice is also a return value of the function.
func newDistance(store [][]int, r1 []rune, r2 []rune, max int, idx int) [][]int {
	dist := LevenshteinDistance(r1, r2)
	if dist <= max {
		store = append(store, []int{idx, dist})
	}
	return store
}
