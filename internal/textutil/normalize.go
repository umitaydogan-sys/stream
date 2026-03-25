package textutil

import (
	"strings"
	"unicode/utf8"
)

// FixLegacyUTF8String repairs common mojibake patterns that appear when
// UTF-8 text is accidentally reinterpreted through a legacy single-byte encoding.
func FixLegacyUTF8String(s string) string {
	current := s
	for i := 0; i < 3; i++ {
		if !strings.ContainsAny(current, "ÃƒÃ„Ã…Ã¢Ã‚") {
			return current
		}

		buf := make([]byte, 0, len(current))
		for _, r := range current {
			if r > 255 {
				return current
			}
			buf = append(buf, byte(r))
		}

		if !utf8.Valid(buf) {
			return current
		}

		next := string(buf)
		if next == current {
			return current
		}
		current = next
	}

	return current
}
