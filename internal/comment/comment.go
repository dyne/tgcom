package comment

import (
	"strings"
)

// Comment adds comments to the specified content based on comment characters.
func Comment(line string, commentChars string) string {
	// just for html
	if commentChars == "<!-- -->" {
		return "<!--" + " " + line + " " + "-->"
	}
	return commentChars + " " + line
}

// Uncomment removes comments from the specified content based on comment characters.
func Uncomment(line string, commentChars string) string {
	trimmedLine := strings.TrimSpace(line)

	//just for html
	if commentChars == "<!-- -->" {
		if strings.HasPrefix(trimmedLine, "<!--") && strings.HasSuffix(trimmedLine, "-->") {
			// Check for both `<!--` and `<!-- ` prefixes.
			if strings.HasPrefix(trimmedLine, "<!--"+" ") {
				line = strings.Replace(line, "<!-- ", "", 1)
				line = strings.Replace(line, "<!--", "", 1)
			}

			// Check for both '-->' and ' -->' suffixes
			if strings.HasSuffix(trimmedLine, " "+"-->") {
				line = strings.Replace(line, " "+"-->", "", 1)
				line = strings.Replace(line, "-->", "", 1)
			}

			return line
		}
	}

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
	//just for html
	if commentChars == "<!-- -->" {
		if strings.HasPrefix(trimmedLine, "<!--") && strings.HasSuffix(trimmedLine, "-->") {
			return Uncomment(line, commentChars)
		}
		return Comment(line, commentChars)
	}

	if strings.HasPrefix(trimmedLine, commentChars) {
		return Uncomment(line, commentChars)
	}
	return Comment(line, commentChars)
}
