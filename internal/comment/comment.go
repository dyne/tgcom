package comment

import (
	"strings"
)

// Comment adds comments to the specified content based on comment characters.
func Comment(line string, commentChars string) string {
	return commentChars + " " + line
}

// Uncomment removes comments from the specified content based on comment characters.
func Uncomment(line string, commentChars string) string {
	trimmedLine := strings.TrimSpace(line)
	if strings.HasPrefix(trimmedLine, commentChars) {
		// Check for both `//` and `// ` prefixes.
		if strings.HasPrefix(trimmedLine, commentChars+" ") {
			return strings.Replace(line, commentChars+" ", "", 1)
		}
		return strings.Replace(line, commentChars, "", 1)
	}
	return line
}

// ToggleComments toggles comments for the specified content based on comment characters.
func ToggleComments(line string, commentChars string) string {
	trimmedLine := strings.TrimSpace(line)
	if strings.HasPrefix(trimmedLine, commentChars) {
		return Uncomment(line, commentChars)
	} else {
		return Comment(line, commentChars)
	}
}
