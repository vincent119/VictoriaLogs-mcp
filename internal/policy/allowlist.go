package policy

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Allowlist stream allowlist
type Allowlist struct {
	enabled bool
	allow   []string
	deny    []string
}

// ErrStreamNotAllowed stream not in allowlist
var ErrStreamNotAllowed = fmt.Errorf("stream not in allowlist")

// ErrStreamDenied stream is explicitly denied
var ErrStreamDenied = fmt.Errorf("stream access denied")

// NewAllowlist creates allowlist
func NewAllowlist(cfg AllowlistConfig) *Allowlist {
	return &Allowlist{
		enabled: cfg.Enabled,
		allow:   cfg.Streams,
		deny:    cfg.Deny,
	}
}

// Check checks if stream is allowed
func (a *Allowlist) Check(stream string) error {
	if !a.enabled {
		return nil
	}

	// Check deny list first
	for _, pattern := range a.deny {
		if matchPattern(stream, pattern) {
			return ErrStreamDenied
		}
	}

	// If allowlist is empty, allow all
	if len(a.allow) == 0 {
		return nil
	}

	// Check allowlist
	for _, pattern := range a.allow {
		if matchPattern(stream, pattern) {
			return nil
		}
	}

	return ErrStreamNotAllowed
}

// IsAllowed checks if stream is allowed (returns bool)
func (a *Allowlist) IsAllowed(stream string) bool {
	return a.Check(stream) == nil
}

// matchPattern uses glob pattern matching
func matchPattern(s, pattern string) bool {
	// Support simple glob patterns: * matches any character
	// Example: kubernetes/* matches kubernetes/pod, kubernetes/service, etc.

	// Direct match
	if s == pattern {
		return true
	}

	// Use filepath.Match for glob matching
	matched, err := filepath.Match(pattern, s)
	if err == nil && matched {
		return true
	}

	// Handle pattern/* case
	if strings.HasSuffix(pattern, "/*") {
		prefix := strings.TrimSuffix(pattern, "/*")
		if strings.HasPrefix(s, prefix+"/") {
			return true
		}
	}

	// Handle pattern/** case (match all subpaths)
	if strings.HasSuffix(pattern, "/**") {
		prefix := strings.TrimSuffix(pattern, "/**")
		if strings.HasPrefix(s, prefix+"/") || s == prefix {
			return true
		}
	}

	return false
}

// AddAllowPattern dynamically adds allow pattern
func (a *Allowlist) AddAllowPattern(pattern string) {
	a.allow = append(a.allow, pattern)
}

// AddDenyPattern dynamically adds deny pattern
func (a *Allowlist) AddDenyPattern(pattern string) {
	a.deny = append(a.deny, pattern)
}
