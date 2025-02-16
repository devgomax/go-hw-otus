package hw03frequencyanalysis

import (
	"cmp"
	"slices"
	"strings"
)

// Top10 возвращает до 10 наиболее часто встречаемых в строке слов
func Top10(raw string) []string {
	counter := make(map[string]int)

	for _, w := range strings.Fields(raw) {
		counter[w]++
	}

	words := make([]string, 0, len(counter))
	for k := range counter {
		words = append(words, k)
	}

	slices.SortFunc(words, func(a, b string) int {
		return cmp.Or(
			cmp.Compare(counter[b], counter[a]),
			strings.Compare(a, b),
		)
	})

	topCount := min(10, len(words))

	return words[:topCount]
}
