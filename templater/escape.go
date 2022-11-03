package templater

import "strings"

func EscapeSlackText(text string) string {
	escapedCharacters := map[string]string{
		"&": "&amp;",
		"<": "&lt;",
		">": "&gt;",
	}

	// Since, maps in go are unordered, we need to maintain the
	// the order in which we escape these characters. The '&'
	// character must be escaped first since it is used in the
	// escaped character representation.
	for _, c := range []string{"&", "<", ">"} {
		text = strings.ReplaceAll(text, c, escapedCharacters[c])
	}

	return text
}
