package main

import (
	"io"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

func configureProcessOutput() {
	configurePlatformConsole()
	log.SetOutput(&textFixWriter{w: os.Stderr})
}

type textFixWriter struct {
	w io.Writer
}

func (w *textFixWriter) Write(p []byte) (int, error) {
	return w.w.Write([]byte(fixLegacyUTF8String(string(p))))
}

func fixLegacyUTF8String(s string) string {
	current := s
	for i := 0; i < 3; i++ {
		if !strings.ContainsAny(current, "ÃÄÅâÂ") {
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
