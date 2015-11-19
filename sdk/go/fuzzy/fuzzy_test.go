package fuzzy

import "testing"

var matchTests = []struct {
	query string
	str   string
	score int
}{
	{"S", "Shutdown", 76},
	{"h", "Shutdown", 60},
	{"u", "Shutdown", 50},
	{"t", "Shutdown", 43},
	{"d", "Shutdown", 38},
	{"o", "Shutdown", 35},
	{"w", "Shutdown", 32},
	{"n", "Shutdown", 30},
	{"chrome", "chrome", 107},
	{"chrom", "chrome", 105},
	{"chrom", "chrom", 105},
	{"nix", "i love me some unix", 19},
	{"nix", "unix i love", 76},
	{"abc", "def", 0},
	{"abc", "zxyabdef", 0},
	{"ad", "zxyabdef", 47},
	{"", "hey", 0},
	{"hey", "", 0},
}

func TestMatch(t *testing.T) {
	for i, tt := range matchTests {
		result := Match(tt.query, tt.str)
		if result.Score != tt.score {
			t.Fatalf("%d: match failed. got: %d, wanted: %d", i, result.Score, tt.score)
		}
	}
}
