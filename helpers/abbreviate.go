package helpers

import "strings"

func Abbreviate(s string) string {
	// Define the mapping of words to their abbreviations
	mapping := map[string]string{
		"Forward":     "FWD",
		"Springboard": "SB",
		"Somersault":  "SS",
		"Inward":      "IWD",
	}

	// Split the input string into words
	words := strings.Split(s, " ")

	// Iterate through the words and replace them with their abbreviations
	for i, word := range words {
		if abbreviation, ok := mapping[word]; ok {
			words[i] = abbreviation
		}
	}

	// Rejoin the words into a single string
	abbreviated := strings.Join(words, " ")

	return abbreviated
}
