package main

import (
	"fmt"
	"fuzzy/search"
)

func main() {
	results, _ := search.Search(
		"The house you bought is very nice",
		[]string{"houze", "very"},
		nil,
	)

	fmt.Println("The house you bought is very nice")
	fmt.Println(results)
}
