package drill

import (
	"math/rand"
	"time"
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

type KeyStat struct {
	Correct   int
	Incorrect int
}

type DrillState struct {
	Prompt          []rune
	Position        int
	TotalKeystrokes int
	Errors          int
	HasError        bool
	Complete        bool
	StartTime       time.Time
	keyCorrect      map[rune]int
	keyIncorrect    map[rune]int
}

func NewDrillState(prompt []rune) *DrillState {
	return &DrillState{
		Prompt:       prompt,
		StartTime:    time.Now(),
		keyCorrect:   make(map[rune]int),
		keyIncorrect: make(map[rune]int),
	}
}

func (d *DrillState) HandleKey(ch rune) {
	if d.Complete || d.HasError {
		return
	}

	d.TotalKeystrokes++
	expected := d.Prompt[d.Position]

	if ch == expected {
		d.keyCorrect[expected]++
		d.Position++
		if d.Position >= len(d.Prompt) {
			d.Complete = true
		}
	} else {
		d.Errors++
		d.keyIncorrect[ch]++
		d.HasError = true
	}
}

func (d *DrillState) HandleBackspace() {
	if d.HasError {
		d.HasError = false
	}
}

func (d *DrillState) Accuracy() float64 {
	if d.TotalKeystrokes == 0 {
		return 100.0
	}
	correct := d.TotalKeystrokes - d.Errors
	return float64(correct) / float64(d.TotalKeystrokes) * 100
}

func (d *DrillState) WPM() float64 {
	elapsed := time.Since(d.StartTime).Minutes()
	if elapsed < 0.001 {
		return 0
	}
	return (float64(d.Position) / 5.0) / elapsed
}

func (d *DrillState) KeyStats() map[rune]KeyStat {
	stats := make(map[rune]KeyStat)
	for k, v := range d.keyCorrect {
		s := stats[k]
		s.Correct = v
		stats[k] = s
	}
	for k, v := range d.keyIncorrect {
		s := stats[k]
		s.Incorrect = v
		stats[k] = s
	}
	return stats
}
