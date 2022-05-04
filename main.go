package main

import (
	"flag"
	"fmt"
	"fuzzy/search"
	"log"
	"strings"
)

const (
	defaultSrc      = "The quick brown fox jumps over the lazy dog"
	defaultPatterns = "teh doug brown"
)

func main() {
	src := flag.String("s", defaultSrc, "search source")
	patterns := flag.String("p", defaultPatterns, "search patterns")
	inFile := flag.Bool("f", false, "in-file search")
	caseInsensitive := flag.Bool("c", false, "case-insensitive search")
	reverse := flag.Bool("r", false, "reverse search")
	matchLimit := flag.Int("m", 10, "displayed matches limit")
	distThreshold := flag.Int("t", 1, "Levenshtein distance threshold")
	flag.Parse()

	// TODO: results output to a file

	patternsSlice := strings.Split(*patterns, " ")

	results, err := search.Search(
		*src,
		patternsSlice,
		&search.Options{
			InFile:          *inFile,
			CaseInsensitive: *caseInsensitive,
			Reverse:         *reverse,
			MatchLimit:      *matchLimit,
			DistThreshold:   *distThreshold,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(results)
}
