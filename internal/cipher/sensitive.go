package cipher

import "strings"

// sensitivePatterns mirrors the patterns used across the project to identify
// keys that contain secret material and should be encrypted at rest.
var sensitivePatterns = []string{
	"PASSWORD",
	"PASSWD",
	"SECRET",
	"TOKEN",
	"API_KEY",
	"APIKEY",
	"PRIVATE_KEY",
	"PRIVATE",
	"CREDENTIAL",
	"AUTH",
	"DSN",
	"DATABASE_URL",
	"ENCRYPTION_KEY",
}

// isSensitiveKey reports whether the given key name matches any of the known
// sensitive patterns (case-insensitive substring match).
func isSensitiveKey(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}
