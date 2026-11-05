// Package acctestcredentials selects test identities without depending on the provider.
package acctestcredentials

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
)

const (
	Domestic      = "Domestic"
	International = "International"
	ModeEnv       = "ALICLOUD_ACCTEST_CREDENTIAL_MODE"
)

var siteKeys = []string{"ALICLOUD_ACCOUNT_SITE", "ALIBABA_CLOUD_ACCOUNT_TYPE", "ALICLOUD_ACCOUNT_TYPE"}
var fallbackKeys = []string{
	"ALIBABA_CLOUD_ACCESS_KEY_ID", "ALIBABA_CLOUD_ACCESS_KEY_SECRET", "ALIBABA_CLOUD_SECURITY_TOKEN",
	"ALIBABACLOUD_ACCESS_KEY_ID", "ALIBABACLOUD_ACCESS_KEY_SECRET", "ALIBABACLOUD_SECURITY_TOKEN",
	"ALICLOUD_ECS_ROLE_NAME", "ALIBABA_CLOUD_ECS_METADATA", "ALICLOUD_PROFILE", "ALIBABA_CLOUD_PROFILE",
	"ALICLOUD_SHARED_CREDENTIALS_FILE", "ALIBABA_CLOUD_CREDENTIALS_FILE", "ALICLOUD_CREDENTIALS_URI", "ALIBABA_CLOUD_CREDENTIALS_URI",
	"ALICLOUD_ASSUME_ROLE_ARN", "ALIBABA_CLOUD_ROLE_ARN", "ALICLOUD_ASSUME_ROLE_SESSION_NAME", "ALIBABA_CLOUD_ROLE_SESSION_NAME",
	"ALIBABA_CLOUD_OIDC_PROVIDER_ARN", "ALIBABA_CLOUD_OIDC_TOKEN_FILE", "ALIBABA_CLOUD_OIDC_TOKEN",
	"ALIBABA_CLOUD_ACCOUNT_ID",
}

// Selection contains environment changes only. Never log or serialize its values.
type Selection struct {
	Values map[string]string
	Unset  []string
	Skip   string
	Source string
}

func currentSite(env map[string]string) string {
	for _, key := range siteKeys {
		if value := strings.TrimSpace(env[key]); value != "" {
			return value
		}
	}
	return Domestic
}

// Resolve is pure: it does not mutate the process or contact cloud APIs.
func Resolve(env map[string]string, site string) (Selection, error) {
	var result Selection
	if site != Domestic && site != International {
		return result, fmt.Errorf("unsupported acceptance account site")
	}
	mode := strings.TrimSpace(env[ModeEnv])
	if mode != "" && mode != "managed" && mode != "custom" {
		return result, fmt.Errorf("invalid acceptance credential mode")
	}
	if site == Domestic {
		if currentSite(env) != Domestic {
			result.Skip = "requested Domestic account; current account site is different"
		}
		return result, nil
	}
	// A caller-supplied identity must never silently become a managed identity.
	if mode == "custom" || currentSite(env) == International {
		if currentSite(env) != International {
			return result, fmt.Errorf("international acceptance test cannot replace custom credentials; an explicitly International test environment is required")
		}
		if strings.TrimSpace(env["ALICLOUD_ACCESS_KEY"]) == "" || strings.TrimSpace(env["ALICLOUD_SECRET_KEY"]) == "" {
			return result, fmt.Errorf("International account is selected but default test credentials are incomplete")
		}
		result.Source = "current"
		return result, nil
	}
	prefix := "ALICLOUD_INTL_"
	result.Source = "international"
	configured := false
	for _, key := range []string{"ACCESS_KEY", "SECRET_KEY", "SECURITY_TOKEN", "ACCOUNT_ID"} {
		if _, ok := env[prefix+key]; ok {
			configured = true
		}
	}
	key, secret, token, account := "", "", "", ""
	if configured {
		key, secret, token, account = env[prefix+"ACCESS_KEY"], env[prefix+"SECRET_KEY"], env[prefix+"SECURITY_TOKEN"], env[prefix+"ACCOUNT_ID"]
	} else {
		key, secret, token, account = env["ALICLOUD_ACCESS_KEY_1"], env["ALICLOUD_SECRET_KEY_1"], env["ALICLOUD_SECURITY_TOKEN_1"], env["ALICLOUD_ACCOUNT_ID_1"]
		result.Source = "legacy-international"
	}
	if key == "" && secret == "" && token == "" && account == "" && !configured {
		if mode == "managed" {
			return result, fmt.Errorf("international test credentials are not configured: ALICLOUD_INTL_ACCESS_KEY / ALICLOUD_INTL_SECRET_KEY")
		}
		result.Skip = "international test credentials are unavailable: configure ALICLOUD_INTL_ACCESS_KEY / ALICLOUD_INTL_SECRET_KEY"
		return result, nil
	}
	if strings.TrimSpace(key) == "" || strings.TrimSpace(secret) == "" {
		return result, fmt.Errorf("international test credentials are incomplete: ALICLOUD_INTL_ACCESS_KEY / ALICLOUD_INTL_SECRET_KEY must be a pair")
	}
	if key != strings.TrimSpace(key) || secret != strings.TrimSpace(secret) || token != strings.TrimSpace(token) {
		return result, fmt.Errorf("international test credentials contain surrounding whitespace")
	}
	result.Values = map[string]string{"ALICLOUD_ACCESS_KEY": key, "ALICLOUD_SECRET_KEY": secret}
	result.Unset = append(result.Unset, fallbackKeys...)
	for _, field := range []struct{ key, value string }{{"ALICLOUD_SECURITY_TOKEN", token}, {"ALICLOUD_ACCOUNT_ID", account}} {
		if field.value != "" {
			result.Values[field.key] = field.value
		} else {
			result.Unset = append(result.Unset, field.key)
		}
	}
	for _, key := range siteKeys {
		result.Values[key] = International
	}
	return result, nil
}

var scopes = struct {
	sync.Mutex
	selected map[*testing.T]string
}{selected: make(map[*testing.T]string)}

// Select lasts through the entire test, including import and destroy. It must
// run before constructing clients; environment mutation is not parallel-safe.
func Select(t *testing.T, site string) {
	t.Helper()
	scopes.Lock()
	previous, exists := scopes.selected[t]
	scopes.Unlock()
	if exists {
		if previous != site {
			t.Fatalf("conflicting account sites in one acceptance test")
		}
		return
	}
	env := make(map[string]string)
	for _, entry := range os.Environ() {
		key, value, _ := strings.Cut(entry, "=")
		env[key] = value
	}
	selection, err := Resolve(env, site)
	if err != nil {
		t.Fatal(err)
	}
	if selection.Skip != "" {
		t.Skip(selection.Skip)
	}
	// Only an actual identity switch mutates the environment. Preserve existing
	// parallel DomesticSite checks that do not change credentials.
	if len(selection.Values) != 0 || len(selection.Unset) != 0 {
		t.Setenv(ModeEnv, env[ModeEnv])
	}
	scopes.Lock()
	scopes.selected[t] = site
	scopes.Unlock()
	t.Cleanup(func() { scopes.Lock(); delete(scopes.selected, t); scopes.Unlock() })
	for _, key := range selection.Unset {
		// Setenv registers exact restoration (including originally-unset variables).
		t.Setenv(key, "")
		if err := os.Unsetenv(key); err != nil {
			t.Fatal("cannot clear test credential fallback")
		}
	}
	for key, value := range selection.Values {
		t.Setenv(key, value)
	}
	t.Logf("Acceptance account site=%s credentialSource=%s", site, selection.Source)
}
