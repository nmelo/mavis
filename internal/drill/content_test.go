package drill

import "testing"

func TestFilterWordsForKeys(t *testing.T) {
	words := []string{"fall", "salad", "flask", "hello", "dad", "ask"}
	keys := map[rune]bool{'f': true, 'j': true, 'd': true, 'k': true, 's': true, 'l': true, 'a': true, ';': true}

	filtered := FilterWords(words, keys)

	for _, w := range filtered {
		for _, ch := range w {
			if !keys[ch] {
				t.Errorf("word %q contains unlocked char %c", w, ch)
			}
		}
	}

	if len(filtered) < 4 {
		t.Errorf("expected at least 4 filtered words, got %d", len(filtered))
	}
}

func TestFilterWordsEmptyResult(t *testing.T) {
	words := []string{"hello", "world"}
	keys := map[rune]bool{'f': true, 'j': true}

	filtered := FilterWords(words, keys)
	if len(filtered) != 0 {
		t.Errorf("expected 0 words, got %d", len(filtered))
	}
}

func TestGenerateWordDrill(t *testing.T) {
	words := []string{"fall", "dad", "ask", "flask", "salad"}
	drill := GenerateWordDrill(words, 5)

	if len(drill) != 5 {
		t.Fatalf("expected 5 words in drill, got %d", len(drill))
	}

	for _, w := range drill {
		found := false
		for _, src := range words {
			if w == src {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("drill word %q not in source list", w)
		}
	}
}

func TestLoadWordList(t *testing.T) {
	words := LoadWordList()
	if len(words) < 100 {
		t.Errorf("word list too small: %d words", len(words))
	}
}

func TestFilterCodeSnippets(t *testing.T) {
	snippets := LoadCodeSnippets()
	if len(snippets) == 0 {
		t.Fatal("should have at least some code snippets")
	}

	keys := make(map[rune]bool)
	for _, r := range []rune{'f', 'j', 'd', 'k', 's', 'l', 'a', ';', 'g', 'h', 't', 'y', 'r', 'u', 'e', 'i', 'w', 'o', 'q', 'p'} {
		keys[r] = true
	}

	filtered := FilterCodeSnippets(snippets, keys)
	for _, s := range filtered {
		for _, ch := range s.Code {
			if ch != ' ' && ch != '\n' && !keys[ch] {
				t.Errorf("snippet %q contains locked char %c", s.Code[:20], ch)
			}
		}
	}
}
