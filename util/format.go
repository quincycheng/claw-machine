package util

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/TylerBrock/colorjson"
	"github.com/briandowns/spinner"
	lipgloss "github.com/charmbracelet/lipgloss"
	textinput "github.com/erikgeiser/promptkit/textinput"
	colorgrad "github.com/mazznoer/colorgrad"
)

func TrimQuotes(s string) string {
	if len(s) >= 2 {
		if c := s[len(s)-1]; s[0] == c && (c == '"' || c == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func PrintPrettyJson(str string) {
	var obj map[string]interface{}
	json.Unmarshal([]byte(str), &obj)

	// Make a custom formatter with indent set
	f := colorjson.NewFormatter()
	f.Indent = 4

	// Marshall the Colorized JSON
	s, _ := f.Marshal(obj)
	fmt.Println(string(s))

}

var highlightStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#3490CE"))

func PrintHeader(header string) {
	grad, _ := colorgrad.NewGradient().HtmlColors("#2BB3EE", "grey").Build()

	for i, c := range header {
		fmt.Print(lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(grad.At(float64(i) / float64(len(header)-1)).Hex())).
			Render(string(c)))
	}
	fmt.Println()
}

func StartSpinner(loadingMsg string, finsihMsg string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[1], 100*time.Millisecond) // Build our new spinner
	s.Suffix = "  " + loadingMsg + "..."
	s.FinalMSG = "✅ " + finsihMsg

	s.Start()
	return s
}

func InputPassword(prompt string) string {
	input := textinput.New(fmt.Sprintf("🔑 %s:", prompt))
	input.Placeholder = "Password will NOT be stored"
	input.Hidden = true

	password, err := input.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	return password
}
