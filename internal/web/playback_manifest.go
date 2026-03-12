package web

import (
	"html"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func passthroughPlaybackQuery(r *http.Request) url.Values {
	values := url.Values{}
	if r == nil {
		return values
	}
	for _, key := range []string{"token", "password"} {
		if value := strings.TrimSpace(r.URL.Query().Get(key)); value != "" {
			values.Set(key, value)
		}
	}
	return values
}

func serveManifestWithPassthrough(w http.ResponseWriter, r *http.Request, filePath, contentType string, rewrite func(string, url.Values) string) bool {
	if r.Method != http.MethodGet {
		return false
	}
	passthrough := passthroughPlaybackQuery(r)
	data, err := os.ReadFile(filePath)
	if err != nil {
		http.NotFound(w, r)
		return true
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "no-cache, no-store")
	_, _ = io.WriteString(w, rewrite(string(data), passthrough))
	return true
}

func rewriteHLSManifest(content string, passthrough url.Values) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "#") {
			lines[i] = rewriteQuotedURIAttribute(line, passthrough, false, "URI")
			continue
		}
		lines[i] = appendPlaybackQuery(line, passthrough, false)
	}
	return strings.Join(lines, "\n")
}

func rewriteDASHManifest(content string, passthrough url.Values) string {
	rewritten := content
	for _, attr := range []string{"media", "initialization", "sourceURL", "href", "xlink:href"} {
		rewritten = rewriteQuotedURIAttribute(rewritten, passthrough, true, attr)
	}
	return rewriteXMLTagText(rewritten, passthrough, "BaseURL")
}

func rewriteQuotedURIAttribute(content string, passthrough url.Values, xmlEscape bool, attr string) string {
	marker := attr + `="`
	searchFrom := 0
	for {
		offset := strings.Index(content[searchFrom:], marker)
		if offset == -1 {
			return content
		}
		start := searchFrom + offset + len(marker)
		endOffset := strings.Index(content[start:], `"`)
		if endOffset == -1 {
			return content
		}
		end := start + endOffset
		replaced := appendPlaybackQuery(content[start:end], passthrough, xmlEscape)
		content = content[:start] + replaced + content[end:]
		searchFrom = start + len(replaced)
	}
}

func rewriteXMLTagText(content string, passthrough url.Values, tag string) string {
	openTag := "<" + tag + ">"
	closeTag := "</" + tag + ">"
	searchFrom := 0
	for {
		openOffset := strings.Index(content[searchFrom:], openTag)
		if openOffset == -1 {
			return content
		}
		start := searchFrom + openOffset + len(openTag)
		closeOffset := strings.Index(content[start:], closeTag)
		if closeOffset == -1 {
			return content
		}
		end := start + closeOffset
		replaced := appendPlaybackQuery(content[start:end], passthrough, true)
		content = content[:start] + replaced + content[end:]
		searchFrom = start + len(replaced) + len(closeTag)
	}
}

func appendPlaybackQuery(raw string, passthrough url.Values, xmlEscape bool) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return raw
	}
	parseInput := raw
	if xmlEscape {
		parseInput = html.UnescapeString(parseInput)
	}
	parseInput = strings.ReplaceAll(parseInput, `\`, `/`)
	parsed, err := url.Parse(parseInput)
	if err != nil {
		return raw
	}
	query := parsed.Query()
	changed := false
	for key, values := range passthrough {
		if query.Get(key) != "" || len(values) == 0 {
			continue
		}
		query.Set(key, values[0])
		changed = true
	}
	if changed {
		parsed.RawQuery = query.Encode()
	}
	result := parsed.String()
	if xmlEscape {
		result = html.EscapeString(result)
	}
	if !changed && result == raw {
		return raw
	}
	return result
}
