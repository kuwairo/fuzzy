package main

import "fuzzy/search"

func main() {
	search.Search(
		"The house you bought is very nice",
		[]string{"houze", "very"},
		&search.Options{},
	)
}
