package policy

import (
	"regexp"
	"sync"
)

// Redactor 敏感資訊遮罩器
type Redactor struct {
	enabled  bool
	patterns []compiledPattern
	mu       sync.RWMutex
}

type compiledPattern struct {
	name        string
	regex       *regexp.Regexp
	replacement string
}

// DefaultRedactPatterns 預設遮罩規則
var DefaultRedactPatterns = []RedactPattern{
	{
		Name:        "ipv4",
		Pattern:     `\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`,
		Replacement: "[REDACTED_IP]",
	},
	{
		Name:        "auth_header",
		Pattern:     `(?i)(authorization|bearer|token)[:\s]+[^\s"]+`,
		Replacement: "[REDACTED_AUTH]",
	},
	{
		Name:        "email",
		Pattern:     `\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`,
		Replacement: "[REDACTED_EMAIL]",
	},
	{
		Name:        "cookie",
		Pattern:     `(?i)cookie[:\s]+[^\s;]+`,
		Replacement: "[REDACTED_COOKIE]",
	},
	{
		Name:        "api_key",
		Pattern:     `(?i)(api[_-]?key|apikey)[:\s=]+[^\s"]+`,
		Replacement: "[REDACTED_API_KEY]",
	},
	{
		Name:        "password",
		Pattern:     `(?i)(password|passwd|pwd)[:\s=]+[^\s"]+`,
		Replacement: "[REDACTED_PASSWORD]",
	},
}

// NewRedactor 建立遮罩器
func NewRedactor(cfg RedactConfig) *Redactor {
	r := &Redactor{
		enabled:  cfg.Enabled,
		patterns: make([]compiledPattern, 0),
	}

	patterns := cfg.Patterns
	if len(patterns) == 0 {
		patterns = DefaultRedactPatterns
	}

	for _, p := range patterns {
		regex, err := regexp.Compile(p.Pattern)
		if err != nil {
			continue // 跳過無效的正規表示式
		}
		r.patterns = append(r.patterns, compiledPattern{
			name:        p.Name,
			regex:       regex,
			replacement: p.Replacement,
		})
	}

	return r
}

// Apply 套用遮罩規則
func (r *Redactor) Apply(data string) string {
	if !r.enabled || len(r.patterns) == 0 {
		return data
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	result := data
	for _, p := range r.patterns {
		result = p.regex.ReplaceAllString(result, p.replacement)
	}

	return result
}

// ApplyToMap 套用遮罩規則到 map
func (r *Redactor) ApplyToMap(data map[string]interface{}) map[string]interface{} {
	if !r.enabled {
		return data
	}

	result := make(map[string]interface{})
	for k, v := range data {
		switch val := v.(type) {
		case string:
			result[k] = r.Apply(val)
		case map[string]interface{}:
			result[k] = r.ApplyToMap(val)
		default:
			result[k] = v
		}
	}
	return result
}

// AddPattern 動態新增遮罩規則
func (r *Redactor) AddPattern(pattern RedactPattern) error {
	regex, err := regexp.Compile(pattern.Pattern)
	if err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.patterns = append(r.patterns, compiledPattern{
		name:        pattern.Name,
		regex:       regex,
		replacement: pattern.Replacement,
	})

	return nil
}

// GetPatternNames 取得所有規則名稱
func (r *Redactor) GetPatternNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, len(r.patterns))
	for i, p := range r.patterns {
		names[i] = p.name
	}
	return names
}
