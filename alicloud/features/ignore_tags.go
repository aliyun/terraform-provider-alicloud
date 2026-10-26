package features

import "strings"

// IgnoreTags keeps tag keys out of a resource's state. Matching only ever looks at the key:
// tag values are neither read nor changed.
type IgnoreTags struct {
	// Keys are matched exactly, case sensitively.
	Keys []string
	// KeyPrefixes match when a tag key starts with them, case sensitively.
	KeyPrefixes []string
}

// Empty reports whether nothing is configured, in which case Read keeps the tags it read.
func (i IgnoreTags) Empty() bool {
	return len(i.Keys) == 0 && len(i.KeyPrefixes) == 0
}

// Match reports whether key is ignored, either by an exact key or by a prefix.
func (i IgnoreTags) Match(key string) bool {
	for _, ignored := range i.Keys {
		if key == ignored {
			return true
		}
	}
	for _, prefix := range i.KeyPrefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}
