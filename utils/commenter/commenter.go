package commenter

import "strings"

// Comment adds a comment character to the beginning of the given line.
func Comment(line string, char string) string {
	// just for html
	if char == "<!-- -->" {
		return "<!--" + " " + line + " " + "-->"
	}
	return char + " " + line
}

// Uncomment removes a comment character or string from the beginning of the given line, if present.
func Uncomment(line string, char string) string {
	trimmedLine := strings.TrimSpace(line)

	//just for html
	if char == "<!-- -->" {
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

	if strings.HasPrefix(trimmedLine, char) {
		// Check for both `//` and `// ` prefixes.
		if strings.HasPrefix(trimmedLine, char+" ") {
			return strings.Replace(line, char+" ", "", 1)
		}
		return strings.Replace(line, char, "", 1)
	}
	return line
}

// ToggleComments toggles comments on or off for the given line based on its current state.
func ToggleComments(line string, char string) string {
	trimmedLine := strings.TrimSpace(line)

	//just for html
	if char == "<!-- -->" {
		if strings.HasPrefix(trimmedLine, "<!--") && strings.HasSuffix(trimmedLine, "-->") {
			return Uncomment(line, char)
		}
		return Comment(line, char)
	}

	if strings.HasPrefix(trimmedLine, char) {
		return Uncomment(line, char)
	} else {
		return Comment(line, char)
	}
}
