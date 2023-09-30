package ansi

import (
	"regexp"
	"strings"
	"text/template"
)

// TemplateFuncMap contains a few useful template helpers
var TemplateFuncMap = template.FuncMap{
	"Left": func(s string, n int) string {
		if n > len(s) {
			n = len(s)
		}

		return s[:n]
	},
	"Matches": func(s, pattern string) bool {
		ok, _ := regexp.MatchString(pattern, s)
		return ok
	},
	"Mid": func(s string, l int, rs ...int) string {
		if l > len(s) {
			l = len(s)
		}

		if len(rs) > 0 {
			r := rs[0]
			if r > len(s) {
				r = len(s)
			}
			return s[l:r]
		}
		return s[l:]
	},
	"Right": func(s string, nn int) string {
		n := len(s) - nn
		if n < 0 {
			n = 0
		}

		return s[n:]
	},
	"Last": func(ss []string) string {
		return ss[len(ss)-1]
	},
	// strings functions
	"Compare":      strings.Compare, // 1.5+ only
	"Contains":     strings.Contains,
	"ContainsAny":  strings.ContainsAny,
	"Count":        strings.Count,
	"EqualFold":    strings.EqualFold,
	"HasPrefix":    strings.HasPrefix,
	"HasSuffix":    strings.HasSuffix,
	"Index":        strings.Index,
	"IndexAny":     strings.IndexAny,
	"Join":         strings.Join,
	"LastIndex":    strings.LastIndex,
	"LastIndexAny": strings.LastIndexAny,
	"Repeat":       strings.Repeat,
	"Replace":      strings.Replace,
	"Split":        strings.Split,
	"SplitAfter":   strings.SplitAfter,
	"SplitAfterN":  strings.SplitAfterN,
	"SplitN":       strings.SplitN,
	"Title":        strings.Title, //nolint:staticcheck
	"ToLower":      strings.ToLower,
	"ToTitle":      strings.ToTitle,
	"ToUpper":      strings.ToUpper,
	"Trim":         strings.Trim,
	"TrimLeft":     strings.TrimLeft,
	"TrimPrefix":   strings.TrimPrefix,
	"TrimRight":    strings.TrimRight,
	"TrimSpace":    strings.TrimSpace,
	"TrimSuffix":   strings.TrimSuffix,
}
