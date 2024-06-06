package comment

import (
	"strings"
)

// Comment adds comments to the specified content based on comment characters.
func Comment(line string, commentChars map[string]string) string {
	return commentChars["singleLine"] + " " + line
}

// Uncomment removes comments from the specified content based on comment characters.
func Uncomment(line string, commentChars map[string]string) string {
	trimmedLine := strings.TrimSpace(line)
	if strings.HasPrefix(trimmedLine, commentChars["singleLine"]) {
		// Check for both `//` and `// ` prefixes.
		if strings.HasPrefix(trimmedLine, commentChars["singleLine"]+" ") {
			return strings.Replace(line, commentChars["singleLine"]+" ", "", 1)
		}
		return strings.Replace(line, commentChars["singleLine"], "", 1)
	}
	return line
}

// ToggleComments toggles comments for the specified content based on comment characters.
func ToggleComments(line string, commentChars map[string]string) string {
	trimmedLine := strings.TrimSpace(line)
	if strings.HasPrefix(trimmedLine, commentChars["singleLine"]) {
		return Uncomment(line, commentChars)
	} else {
		return Comment(line, commentChars)
	}
}
