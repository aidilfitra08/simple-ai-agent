package utils

import (
	"regexp"
	"strings"
)

// RemoveThinkBlocks removes <think>...</think> blocks from text
func RemoveThinkBlocks(input string) string {
	// Regex to match <think>...</think> including newlines
	re := regexp.MustCompile(`(?s)<think>.*?</think>`)

	// Replace with empty string
	result := re.ReplaceAllString(input, "")

	// Remove only leading newlines (at the beginning of the string)
	result = strings.TrimLeft(result, "\n")

	return result
}
