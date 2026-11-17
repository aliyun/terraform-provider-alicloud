package acctestcredentials

import (
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

func TestUnitResolveAccountSite(t *testing.T) {
	for _, tc := range []struct {
		name          string
		env           map[string]string
		site          string
		wantErr, skip bool
		source        string
	}{
		{"local missing", map[string]string{}, International, false, true, "legacy-international"},
		{"managed missing", map[string]string{ModeEnv: "managed"}, International, true, false, ""},
		{"legacy", map[string]string{"ALICLOUD_ACCESS_KEY_1": "key", "ALICLOUD_SECRET_KEY_1": "secret"}, International, false, false, "legacy-international"},
		{"named", map[string]string{"ALICLOUD_INTL_ACCESS_KEY": "key", "ALICLOUD_INTL_SECRET_KEY": "secret"}, International, false, false, "international"},
		{"partial named cannot fallback", map[string]string{"ALICLOUD_INTL_ACCESS_KEY": "key", "ALICLOUD_ACCESS_KEY_1": "old", "ALICLOUD_SECRET_KEY_1": "old-secret"}, International, true, false, ""},
		{"empty named cannot fallback", map[string]string{"ALICLOUD_INTL_ACCESS_KEY": "", "ALICLOUD_ACCESS_KEY_1": "old", "ALICLOUD_SECRET_KEY_1": "old-secret"}, International, true, false, ""},
		{"partial legacy", map[string]string{"ALICLOUD_SECRET_KEY_1": "secret"}, International, true, false, ""},
		{"whitespace", map[string]string{"ALICLOUD_INTL_ACCESS_KEY": " key", "ALICLOUD_INTL_SECRET_KEY": "secret"}, International, true, false, ""},
		{"domestic", map[string]string{}, Domestic, false, false, ""},
		{"domestic on international", map[string]string{"ALICLOUD_ACCOUNT_SITE": International}, Domestic, false, true, ""},
		{"unknown site", map[string]string{}, "unknown", true, false, ""},
		{"unknown mode", map[string]string{ModeEnv: "unknown"}, International, true, false, ""},
		{"custom cannot switch", map[string]string{ModeEnv: "custom", "ALICLOUD_ACCESS_KEY_1": "key", "ALICLOUD_SECRET_KEY_1": "secret"}, International, true, false, ""},
		{"custom international", map[string]string{ModeEnv: "custom", "ALICLOUD_ACCOUNT_SITE": International, "ALICLOUD_ACCESS_KEY": "current", "ALICLOUD_SECRET_KEY": "current-secret"}, International, false, false, "current"},
		{"incomplete current", map[string]string{"ALICLOUD_ACCOUNT_SITE": International, "ALICLOUD_ACCESS_KEY": "current"}, International, true, false, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			original := make(map[string]string)
			for k, v := range tc.env {
				original[k] = v
			}
			result, err := Resolve(tc.env, tc.site)
			if (err != nil) != tc.wantErr {
				t.Fatalf("unexpected error presence: %v", err)
			}
			if err == nil && ((result.Skip != "") != tc.skip || result.Source != tc.source) {
				t.Fatalf("unexpected decision (skip=%t source=%s)", result.Skip != "", result.Source)
			}
			if !reflect.DeepEqual(tc.env, original) {
				t.Fatal("resolver changed environment")
			}
			if err != nil && (strings.Contains(err.Error(), "current-secret") || strings.Contains(err.Error(), "old-secret")) {
				t.Fatal("secret leaked in error")
			}
		})
	}
}

func clearTestEnv(t *testing.T) {
	t.Helper()
	keys := append(append([]string{}, fallbackKeys...), siteKeys...)
	keys = append(keys, ModeEnv, "ALICLOUD_ACCESS_KEY", "ALICLOUD_SECRET_KEY", "ALICLOUD_SECURITY_TOKEN", "ALICLOUD_ACCOUNT_ID",
		"ALICLOUD_INTL_ACCESS_KEY", "ALICLOUD_INTL_SECRET_KEY", "ALICLOUD_INTL_SECURITY_TOKEN", "ALICLOUD_INTL_ACCOUNT_ID",
		"ALICLOUD_ACCESS_KEY_1", "ALICLOUD_SECRET_KEY_1", "ALICLOUD_SECURITY_TOKEN_1", "ALICLOUD_ACCOUNT_ID_1")
	for _, key := range keys {
		t.Setenv(key, "")
		os.Unsetenv(key)
	}
}

func TestUnitSelectionRestoresIdentityAcrossTests(t *testing.T) {
	clearTestEnv(t)
	t.Setenv("ALICLOUD_ACCESS_KEY", "domestic-key")
	t.Setenv("ALICLOUD_SECRET_KEY", "domestic-secret")
	t.Setenv("ALICLOUD_SECURITY_TOKEN", "old-token")
	t.Setenv("ALICLOUD_ACCOUNT_ID", "old-account")
	t.Setenv("ALICLOUD_INTL_ACCESS_KEY", "intl-key")
	t.Setenv("ALICLOUD_INTL_SECRET_KEY", "intl-secret")
	t.Setenv("ALIBABA_CLOUD_PROFILE", "old-profile")
	t.Run("domestic before", func(t *testing.T) {
		Select(t, Domestic)
		if os.Getenv("ALICLOUD_ACCESS_KEY") != "domestic-key" {
			t.Fatal("wrong domestic identity")
		}
	})
	for i := 0; i < 2; i++ {
		t.Run("international", func(t *testing.T) {
			Select(t, International)
			Select(t, International)
			if os.Getenv("ALICLOUD_ACCESS_KEY") != "intl-key" || os.Getenv("ALICLOUD_SECRET_KEY") != "intl-secret" {
				t.Fatal("identity not switched")
			}
			for _, key := range []string{"ALICLOUD_SECURITY_TOKEN", "ALICLOUD_ACCOUNT_ID", "ALIBABA_CLOUD_PROFILE"} {
				if _, ok := os.LookupEnv(key); ok {
					t.Errorf("stale fallback %s", key)
				}
			}
			if os.Getenv("ALIBABA_CLOUD_ACCOUNT_TYPE") != International {
				t.Fatal("site not switched")
			}
		})
	}
	t.Run("domestic after", func(t *testing.T) {
		Select(t, Domestic)
		if os.Getenv("ALICLOUD_ACCESS_KEY") != "domestic-key" || os.Getenv("ALICLOUD_SECURITY_TOKEN") != "old-token" || os.Getenv("ALICLOUD_ACCOUNT_ID") != "old-account" || os.Getenv("ALIBABA_CLOUD_PROFILE") != "old-profile" {
			t.Fatal("identity was not restored")
		}
	})
}

func TestUnitSelectionKeepsOneCredentialBundle(t *testing.T) {
	result, err := Resolve(map[string]string{"ALICLOUD_INTL_ACCESS_KEY": "named", "ALICLOUD_INTL_SECRET_KEY": "named-secret", "ALICLOUD_INTL_SECURITY_TOKEN": "named-token", "ALICLOUD_INTL_ACCOUNT_ID": "named-account", "ALICLOUD_ACCESS_KEY_1": "legacy", "ALICLOUD_SECRET_KEY_1": "legacy-secret"}, International)
	if err != nil {
		t.Fatal(err)
	}
	if result.Values["ALICLOUD_ACCESS_KEY"] != "named" || result.Values["ALICLOUD_SECURITY_TOKEN"] != "named-token" || result.Values["ALICLOUD_ACCOUNT_ID"] != "named-account" {
		t.Fatal("credential bundles mixed")
	}
}

func TestUnitSelectionRejectsParallelAndConflictingSites(t *testing.T) {
	if scenario := os.Getenv("ACCTEST_SELECTION_CHILD"); scenario != "" {
		os.Setenv("ALICLOUD_INTL_ACCESS_KEY", "test-key")
		os.Setenv("ALICLOUD_INTL_SECRET_KEY", "test-secret")
		if scenario == "parallel" || scenario == "parallel-domestic" {
			requested := International
			if scenario == "parallel-domestic" {
				os.Setenv("ALICLOUD_ACCOUNT_SITE", Domestic)
				requested = Domestic
			}
			t.Run("parallel", func(t *testing.T) { t.Parallel(); Select(t, requested) })
			return
		}
		Select(t, International)
		Select(t, Domestic)
		return
	}
	for _, scenario := range []string{"parallel", "conflict", "parallel-domestic"} {
		command := exec.Command(os.Args[0], "-test.run=^TestUnitSelectionRejectsParallelAndConflictingSites$")
		command.Env = append(os.Environ(), "ACCTEST_SELECTION_CHILD="+scenario)
		output, err := command.CombinedOutput()
		if scenario == "parallel-domestic" {
			if err != nil {
				t.Fatal("unchanged domestic identity must remain parallel-compatible")
			}
			continue
		}
		if err == nil {
			t.Fatal("invalid selection was accepted")
		}
		expected := "conflicting account sites"
		if scenario == "parallel" {
			expected = "parallel"
		}
		if !strings.Contains(string(output), expected) {
			t.Fatal("failure was unrelated to the invalid selection")
		}
	}
}

func TestUnitSelectionRestoresAfterSkip(t *testing.T) {
	clearTestEnv(t)
	t.Setenv("ALICLOUD_ACCESS_KEY", "original")
	t.Setenv("ALICLOUD_INTL_ACCESS_KEY", "intl")
	t.Setenv("ALICLOUD_INTL_SECRET_KEY", "intl-secret")
	t.Run("skip after selection", func(t *testing.T) { Select(t, International); t.Skip("business prerequisite unavailable") })
	if os.Getenv("ALICLOUD_ACCESS_KEY") != "original" {
		t.Fatal("skip leaked identity")
	}
}
