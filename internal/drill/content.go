package drill

import (
	"math/rand"
	"strings"

	"github.com/nmelo/typer/data"
)

func LoadWordList() []string {
	var words []string
	for _, w := range strings.Split(data.WordList, "\n") {
		w = strings.TrimSpace(w)
		if w != "" {
			words = append(words, w)
		}
	}
	return words
}

func FilterWords(words []string, allowedKeys map[rune]bool) []string {
	var filtered []string
	for _, w := range words {
		ok := true
		for _, ch := range w {
			if !allowedKeys[ch] {
				ok = false
				break
			}
		}
		if ok && len(w) > 0 {
			filtered = append(filtered, w)
		}
	}
	return filtered
}

func GenerateWordDrill(words []string, n int) []string {
	if len(words) == 0 {
		return nil
	}
	drill := make([]string, n)
	for i := range drill {
		drill[i] = words[rand.Intn(len(words))]
	}
	return drill
}
