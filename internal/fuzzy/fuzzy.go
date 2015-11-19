package fuzzy

import (
	"bytes"
	"strings"
)

type MatchResult struct {
	Success bool
	Score   int
	Value   string
}

func Match(query, str string) *MatchResult {
	if len(str) == 0 || len(query) == 0 {
		return &MatchResult{
			Success: false,
		}
	}

	strLen := len(str)
	compareString := strings.ToLower(str)
	pattern := strings.ToLower(query)

	sb := bytes.Buffer{}
	patternIdx := 0
	firstMatchIndex := -1
	lastMatchIndex := 0
	var ch byte

	for idx := 0; idx < strLen; idx++ {
		ch = str[idx]
		if compareString[idx] == pattern[patternIdx] {
			if firstMatchIndex < 0 {
				firstMatchIndex = idx
			}
			lastMatchIndex = idx + 1

			sb.WriteByte( /*opt.Prefix + */ ch /* + opt.Suffix*/)
			patternIdx += 1
		} else {
			sb.WriteByte(ch)
		}

		if patternIdx == len(pattern) && (idx+1) != len(compareString) {
			sb.WriteString(str[idx+1:])
			break
		}
	}

	if patternIdx == len(pattern) {
		return &MatchResult{
			Success: true,
			Value:   sb.String(),
			Score:   CalcScore(query, str, firstMatchIndex, lastMatchIndex-firstMatchIndex),
		}
	}

	return &MatchResult{
		Success: false,
	}
}

func CalcScore(query, str string, firstIndex, matchLen int) int {
	// a match found near the beginning of a string is scored more than a match found near the end
	// a match is scored more if the characters in the patterns are closer to each other, while the score is lower if they are more spread out
	score := 100 * (len(query) + 1) / ((1 + firstIndex) + (matchLen + 1))
	// a match with less characters assigning more weights
	if len(str)-len(query) < 5 {
		score = score + 20
	} else if len(str)-len(query) < 10 {
		score = score + 10
	}

	return score
}
