package commenter

import (
	"testing"
)

func TestComment(t *testing.T) {
	commentChars := "//"
	tests := []struct {
		line     string
		expected string
	}{
		{"I want to test!", "// I want to test!"},
		{"// This is already comment", "// // This is already comment"},
		{"   This has leading spaces", "//    This has leading spaces"},
		{"		with tab", "// 		with tab"},
		{"", "// "},
	}

	for _, test := range tests {
		result := Comment(test.line, commentChars)
		if result != test.expected {
			t.Errorf("Expected Comment(%q) to be %q, but got %q", test.line, test.expected, result)
		}
	}
}

func TestUncomment(t *testing.T) {
	commentChars := "//"
	tests := []struct {
		line     string
		expected string
	}{
		{"// I want to test!", "I want to test!"},
		{"     // This has leading spaces", "     This has leading spaces"},
		{"//		with tab", "		with tab"},
		{"This does not have a comment", "This does not have a comment"},
	}

	for _, test := range tests {
		result := Uncomment(test.line, commentChars)
		if result != test.expected {
			t.Errorf("Expected Uncomment(%q) to be %q, but got %q", test.line, test.expected, result)
		}
	}
}

func TestToggleComments(t *testing.T) {
	commentChars := "//"
	tests := []struct {
		line     string
		expected string
	}{
		{"Hello, world!", "// Hello, world!"},
		{"// This is a comment", "This is a comment"},
		{"//     This has leading spaces", "    This has leading spaces"},
		{"", "// "},
	}

	for _, test := range tests {
		result := ToggleComments(test.line, commentChars)
		if result != test.expected {
			t.Errorf("Expected ToggleComments(%q) to be %q, but got %q", test.line, test.expected, result)
		}
	}
}
