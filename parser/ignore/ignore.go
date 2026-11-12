package ignore

import (
	"strings"

	"github.com/open-policy-agent/conftest/parser/internal/convert"
)

// Parser is a ignore (dockerignore, gitignore) parser.
type Parser struct{}

// Entry is a single line of an ignore file.
type Entry struct {
	Kind     string
	Value    string
	Original string
}

// Unmarshal unmarshals ignore files.
func (pp *Parser) Unmarshal(p []byte, v any) error {
	ignoreEntries := parseLines(strings.Split(string(p), "\n"))

	// Wrap the entry list in another list, to ensure it's
	// treated as a single file.
	entryListList := [][]Entry{ignoreEntries}

	return convert.Remarshal("ignore", entryListList, v)
}

func parseLines(lines []string) []Entry {
	entries := make([]Entry, 0, len(lines))
	for _, l := range lines {
		stripped := trimRightSpace(l)

		switch {
		case stripped == "":
			entries = append(entries, Entry{"Empty", l, l})
		case strings.HasPrefix(l, "#"):
			entries = append(entries, Entry{"Comment", strings.TrimSpace(stripped[1:]), l})
		case strings.HasPrefix(l, "!"):
			entries = append(entries, Entry{"NegatedPath", unescape(stripped[1:]), l})
		default:
			entries = append(entries, Entry{"Path", unescape(stripped), l})
		}
	}

	return entries
}

// trimRightSpace removes trailing spaces, keeping one space that is escaped
// with a backslash. A line may end in a backslash.
func trimRightSpace(s string) string {
	stop := len(s)
	for stop > 0 && s[stop-1] == ' ' {
		stop--
	}

	if stop > 0 && stop < len(s) && s[stop-1] == '\\' {
		stop++
	}

	return s[:stop]
}

// unescape drops single backslashes and turns a double backslash into one.
func unescape(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] != '\\' {
			b.WriteByte(s[i])
		} else if i < len(s)-1 && s[i+1] == '\\' {
			b.WriteByte(s[i])
			i++
		}
	}

	return b.String()
}
