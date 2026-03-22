package drill

import (
	"math/rand"
)

// GenerateCharDrill creates a character drill of the given length.
// newKeys are weighted at ~60%, review keys at ~40%.
func GenerateCharDrill(allKeys []rune, newKeys []rune, length int) []rune {
	if len(allKeys) == 0 {
		return nil
	}

	// Build weighted pool: new keys get 4x weight
	// Exclude space from char drills (space is for word/code drills only per spec)
	var pool []rune
	newSet := make(map[rune]bool)
	for _, k := range newKeys {
		newSet[k] = true
	}

	for _, k := range allKeys {
		if k == ' ' {
			continue
		}
		if newSet[k] {
			pool = append(pool, k, k, k, k) // 4x weight for new keys (~60% when 2 new vs 8 review)
		} else {
			pool = append(pool, k)
		}
	}

	if len(pool) == 0 {
		return nil
	}

	drill := make([]rune, length)
	for i := range drill {
		drill[i] = pool[rand.Intn(len(pool))]
	}
	return drill
}
