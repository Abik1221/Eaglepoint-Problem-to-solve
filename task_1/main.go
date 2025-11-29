package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// this struct holds the results of text processing
type TextAnalysis struct {
	WordCount         int            `json:"word_count"`
	AverageWordLength float64        `json:"average_word_length"`
	LongestWords      []string       `json:"longest_words"`
	WordFrequency     map[string]int `json:"word_frequency"`
}

// analyzeText takes a string and returns detailed statistics about it
func AnalyzeText(text string) TextAnalysis {
	words := strings.Fields(text)

	analysis := TextAnalysis{
		WordCount:     0,
		LongestWords:  []string{},
		WordFrequency: make(map[string]int),
	}

	// if the string was empty we are done here.
	if len(words) == 0 {
		return analysis
	}

	totalLetters := 0
	maxWordLength := 0

	// now we loop through every word to gather our stats.
	for _, word := range words {
		// we want "The" and "the" to count as the same word, so lowercase everything.
		lowerWord := strings.ToLower(word)

		// update our frequency map
		analysis.WordFrequency[lowerWord]++

		// keep track of total letters so we can calculate the average later
		wordLength := len(lowerWord)
		totalLetters += wordLength

		// check if this is the longest word we've seen so far
		if wordLength > maxWordLength {
			// if new winner find replace it
			maxWordLength = wordLength
			analysis.LongestWords = []string{lowerWord}
		} else if wordLength == maxWordLength {
			// it's a tie so add it to the list if we haven't seen it already.
			isDuplicate := false
			for _, existingWord := range analysis.LongestWords {
				if existingWord == lowerWord {
					isDuplicate = true
					break
				}
			}
			if !isDuplicate {
				analysis.LongestWords = append(analysis.LongestWords, lowerWord)
			}
		}
	}

	// final calculations
	analysis.WordCount = len(words)

	// calculate average and round it to 2 decimal places for cleaner output
	analysis.AverageWordLength = float64(totalLetters) / float64(analysis.WordCount)
	analysis.AverageWordLength = float64(int(analysis.AverageWordLength*100)) / 100

	return analysis
}

func main() {
	text := "The quick brown fox jumps over the lazy dog the fox"

	result := AnalyzeText(text)

	// pretty print the JSON output
	jsonOutput, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("Error formatting output:", err)
		return
	}

	fmt.Println("Input:")
	fmt.Printf("\"%s\"\n\n", text)
	fmt.Println("Output:")
	fmt.Println(string(jsonOutput))
}
