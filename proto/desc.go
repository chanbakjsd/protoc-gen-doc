package proto

import (
	"go/doc"
	"strings"
	"unicode"

	"google.golang.org/protobuf/compiler/protogen"
)

// ConvertCommentSet parses the comment set as a description.
func ConvertCommentSet(c protogen.CommentSet) Desc {
	return ParseDesc(string(c.Leading))
}

// ParseDesc parses the string provided as a description, a comment describing
// a type, method or value.
func ParseDesc(s string) Desc {
	s = clean(s)
	var deprecated bool
	if strings.Contains(s, "Deprecated: ") {
		deprecated = true
	}
	return Desc{
		Text:       s,
		Deprecated: deprecated,
	}
}

// Desc is a struct containing information retrieved from the description.
type Desc struct {
	Text       string
	Deprecated bool
}

// Long returns the description with the name removed.
func (d Desc) Long(name string) string {
	text := d.Text
	if d.Deprecated {
		text = strings.ReplaceAll(text, "Deprecated: ", "")
	}
	long := strings.TrimPrefix(text, name+" is ")
	long = strings.TrimPrefix(long, name+" are ")
	long = strings.TrimPrefix(long, name+" ")
	if long == "" {
		return ""
	}
	r := []rune(long)
	r[0] = unicode.ToTitle(r[0])
	return string(r)
}

// Short returns the first sentence of the description with the name removed.
func (d Desc) Short(name string) string {
	text := d.Long(name)
	synopsis := doc.Synopsis(text)
	return strings.TrimSuffix(synopsis, ".")
}

// clean rewrites the provided string's whitespaces by removing leading,
// trailing, and excessive newlines in addition to collapsing multiple word
// separators into one character.
func clean(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	var nlCount int
	var whitespace bool
	for _, v := range s {
		// Line separator.
		if unicode.Is(unicode.Zl, v) {
			nlCount++
			continue
		}
		// Whitespace.
		if unicode.IsSpace(v) {
			whitespace = true
			continue
		}
		// Ignore whitespaces surrounding line separators.
		if nlCount > 0 {
			whitespace = false
		}
		// Only print newlines if it is not at the start.
		if b.Len() > 0 {
			switch nlCount {
			case 0:
				if whitespace {
					b.WriteByte(' ')
				}
			case 1:
				// Turn a single new line into space.
				b.WriteByte(' ')
			default:
				// Double new line starts a new paragraph.
				b.WriteString("\n\n")
			}
		}
		whitespace = false
		nlCount = 0
		b.WriteRune(v)
	}
	return b.String()
}
