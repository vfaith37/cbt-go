package security

import (
	"net/http"
	"regexp"
)

type WAF struct {
	rules []*regexp.Regexp
}

func NewWAF() *WAF {
	return &WAF{
		rules: []*regexp.Regexp{
			regexp.MustCompile(`(?i)(<script|javascript:|onload=|onerror=)`),
			regexp.MustCompile(`(?i)(UNION.*SELECT|INSERT.*INTO|DELETE.*FROM)`),
			regexp.MustCompile(`(?i)(../../|\.\.\/\.\.\/|/etc/passwd|/etc/shadow)`),
		},
	}
}

func (w *WAF) CheckRequest(req *http.Request) bool {
	// Check headers
	for header, values := range req.Header {
		for _, value := range values {
			if w.containsThreats(value) {
				return false
			}
		}
	}

	// Check URL path
	if w.containsThreats(req.URL.Path) {
		return false
	}

	// Check query parameters
	for _, values := range req.URL.Query() {
		for _, value := range values {
			if w.containsThreats(value) {
				return false
			}
		}
	}

	return true
}
