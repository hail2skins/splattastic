package helpers

import (
	"html/template"
	"strings"

	"github.com/gin-gonic/gin"
)

func InitRouterWithFuncMap(r *gin.Engine) {
	funcMap := template.FuncMap{
		"mod":     func(i, j int) int { return i % j },
		"shorten": abbreviate, // Provides view method to shorten the names of dives in some views
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseGlob("../../templates/**/**")
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(tmpl)
}

// abbreviate takes a string and returns an abbreviated version of it
// TODO: This should be moved to a helpers package
// TODO: And more mappings should be added
func abbreviate(s string) string {
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
