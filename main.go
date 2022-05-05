package main

import (
	"flag"
	"fmt"
	"fuzzy/search"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	defaultSrc      = "The quick brown fox jumps over the lazy dog"
	defaultPatterns = "teh doug brown"
)

var colors = []color.Attribute{
	color.FgRed,
	color.FgGreen,
	color.FgYellow,
	color.FgBlue,
	color.FgMagenta,
	color.FgCyan,
	color.FgWhite,
}

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

	var text []rune
	rand.Seed(time.Now().UnixNano())

	if *inFile {
		data, err := os.ReadFile(*src)
		if err != nil {
			log.Fatal(err)
		}
		text = []rune(string(data))

	} else {
		text = []rune(*src)
	}

	for pattern, matches := range results {
		fmt.Printf("`%s` matches:\n", pattern)
		if len(matches) < 1 {
			fmt.Printf("No matches found\n\n")
			continue
		}

		pLen := len([]rune(pattern))
		for _, idx := range matches {
			pEnd := pLen + idx
			if tLen := len(text); pEnd > tLen {
				pEnd = tLen
			}

			colorIdx := rand.Intn(len(colors) - 1)
			pColor := color.New(colors[colorIdx])

			fmt.Printf(
				"%s%s%s\n",
				string(text[:idx]),
				pColor.Sprint(string(text[idx:pEnd])),
				string(text[pEnd:]),
			)
		}
		fmt.Println()
	}
}
