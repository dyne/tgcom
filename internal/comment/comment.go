package comment

import "strings"

// Comment adds comments to the specified content based on comment characters.
func Comment(line string, commentChars map[string]string) string {
	return commentChars["singleLine"] + " " + line
}

// Uncomment removes comments from the specified content based on comment characters.
func Uncomment(line string, commentChars map[string]string) string {
	return strings.TrimPrefix(strings.TrimSpace(line), commentChars["singleLine"]+" ")
}

// ToggleComments toggles comments for the specified content based on comment characters.
func ToggleComments(line string, commentChars map[string]string) string {

	if strings.HasPrefix(strings.TrimSpace(line), commentChars["singleLine"]) {
		return strings.TrimPrefix(strings.TrimSpace(line), commentChars["singleLine"]+" ")
	} else {
		return commentChars["singleLine"] + " " + line
	}
}
