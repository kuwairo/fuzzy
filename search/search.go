package search

import (
	"os"
	"strings"

	"golang.org/x/exp/slices"
)

type Options struct {
	InFile          bool
	CaseInsensitive bool
	Reverse         bool
	MatchLimit      int
	DistThreshold   int
}

func DefaultOptions() *Options {
	return &Options{
		MatchLimit:    10,
		DistThreshold: 2,
	}
}

// TODO: concurrent search (lim = runtime.NumCPU())

func Search(src string, patterns []string, options *Options) (map[string][]int, error) {
	if options == nil {
		options = DefaultOptions()
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

	matches := make(map[string][]int)
	textRunes := []rune(text)

	for _, pattern := range patterns {
		tmp := [][]int{}
		p := []rune(pattern)

		for i := range textRunes {
			left, right := textRunes[:i], textRunes[i:]

			tmp = newDistance(tmp, left, p, options.DistThreshold, 0)
			tmp = newDistance(tmp, right, p, options.DistThreshold, i)

			if len(p)+i <= len(textRunes) {
				window := textRunes[i : len(p)+i]
				tmp = newDistance(tmp, window, p, options.DistThreshold, i)
			}
		}

		matched := matches[pattern]
		for _, pair := range tmp {
			idx, dist := pair[0], pair[1]

			if dist <= options.DistThreshold {
				if !slices.Contains(matched, dist) {
					matched = append(matched, idx)
				}
			}
		}

		if options.Reverse {
			for i, j := 0, len(matched)-1; i < j; i, j = i+1, j+1 {
				matched[i], matched[j] = matched[j], matched[i]
			}
		}
		matches[pattern] = matched
	}

	for pattern, matched := range matches {
		limit, count := options.MatchLimit, len(matched)
		if limit > count {
			limit = count
		}
		matches[pattern] = matched[:limit]
	}
	return matches, nil
}

func newDistance(store [][]int, r1 []rune, r2 []rune, max int, idx int) [][]int {
	dist := LevenshteinDistance(r1, r2)
	if dist <= max {
		store = append(store, []int{idx, dist})
	}
	return store
}
