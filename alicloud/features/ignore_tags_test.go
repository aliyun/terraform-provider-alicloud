package features

import "testing"

func TestUnitIgnoreTagsEmpty(t *testing.T) {
	cases := []struct {
		name   string
		ignore IgnoreTags
		want   bool
	}{
		{"zero value", IgnoreTags{}, true},
		{"nil slices", IgnoreTags{Keys: nil, KeyPrefixes: nil}, true},
		{"empty slices", IgnoreTags{Keys: []string{}, KeyPrefixes: []string{}}, true},
		{"keys only", IgnoreTags{Keys: []string{"env"}}, false},
		{"prefixes only", IgnoreTags{KeyPrefixes: []string{"tf-"}}, false},
		{"both", IgnoreTags{Keys: []string{"env"}, KeyPrefixes: []string{"tf-"}}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.ignore.Empty(); got != c.want {
				t.Fatalf("Empty() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestUnitIgnoreTagsMatchExactKeys(t *testing.T) {
	ignore := IgnoreTags{Keys: []string{"env", "owner"}}
	cases := map[string]bool{
		"env":         true,
		"owner":       true,
		"environment": false,
		"Env":         false,
		"OWNER":       false,
		"":            false,
		"team":        false,
	}
	for key, want := range cases {
		if got := ignore.Match(key); got != want {
			t.Errorf("Match(%q) = %v, want %v", key, got, want)
		}
	}
}

func TestUnitIgnoreTagsMatchKeyPrefixes(t *testing.T) {
	ignore := IgnoreTags{KeyPrefixes: []string{"tf-", "cost:"}}
	cases := map[string]bool{
		"tf-":                    true,
		"tf-created-by":          true,
		"cost:center":            true,
		"cost:":                  true,
		"tf":                     false,
		"cost":                   false,
		"x-tf-created-by":        false,
		"Cost:center":            false,
		"TF-created-by":          false,
		"":                       false,
		"created-by-tf-suffixes": false,
	}
	for key, want := range cases {
		if got := ignore.Match(key); got != want {
			t.Errorf("Match(%q) = %v, want %v", key, got, want)
		}
	}
}

func TestUnitIgnoreTagsMatchExactAndPrefixOverlap(t *testing.T) {
	ignore := IgnoreTags{Keys: []string{"tf-env"}, KeyPrefixes: []string{"tf-"}}
	for _, key := range []string{"tf-env", "tf-owner", "tf-"} {
		if !ignore.Match(key) {
			t.Errorf("Match(%q) = false, want true", key)
		}
	}
	for _, key := range []string{"env", "tf", "tfenv"} {
		if ignore.Match(key) {
			t.Errorf("Match(%q) = true, want false", key)
		}
	}
}

func TestUnitIgnoreTagsMatchIsKeyOnly(t *testing.T) {
	// Match takes only the key, so a tag value can never make a key ignored.
	ignore := IgnoreTags{Keys: []string{"secret"}, KeyPrefixes: []string{"sensitive-"}}
	if ignore.Match("plain") {
		t.Fatal("a key that matches nothing must stay visible regardless of its value")
	}
	if !ignore.Match("secret") {
		t.Fatal("exact key must match")
	}
}

func TestUnitIgnoreTagsMatchEmptyConfiguration(t *testing.T) {
	var ignore IgnoreTags
	for _, key := range []string{"", "env", "acs:created"} {
		if ignore.Match(key) {
			t.Errorf("Match(%q) = true with empty configuration, want false", key)
		}
	}
}
