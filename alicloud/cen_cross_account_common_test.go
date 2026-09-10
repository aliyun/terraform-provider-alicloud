package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// cenCrossAccountCreds holds AK/SK pairs loaded from the aliyun CLI profiles
// used by all CEN cross-account tests. It is populated at the top of each
// such test (before Steps are evaluated) so HCL provider blocks can embed
// credentials directly — the in-tree provider caches its profile lookup in a
// package global, which makes `profile = "X"` unreliable once a second alias
// is present, so we bypass it entirely.
type cenCrossAccountCreds struct {
	// utAK / utSK come from profile "TerraformUT" (account A — CEN owner).
	utAK, utSK string
	// testAK / testSK come from profile "TerraformTest" (account B — resource owner).
	testAK, testSK string
}

var sharedCENCrossAccountCreds cenCrossAccountCreds

// testAccPreCheckCENCrossAccount prepares the environment for any CEN
// cross-account test: it loads AK/SK for profiles "TerraformUT" and
// "TerraformTest" from ~/.aliyun/config.json, then unsets conflicting env
// credentials (restored via t.Cleanup) so the provider uses the injected
// AK/SK from the HCL rather than the shell environment. The test is skipped
// when the CLI config or either profile is missing.
func testAccPreCheckCENCrossAccount(t *testing.T) {
	u, err := user.Current()
	if err != nil {
		t.Fatalf("resolve current user: %s", err)
	}
	cfgPath := filepath.Join(u.HomeDir, ".aliyun", "config.json")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Skipf("aliyun CLI config not found at %s: %s", cfgPath, err)
	}
	var parsed struct {
		Profiles []map[string]interface{} `json:"profiles"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("parse %s: %s", cfgPath, err)
	}
	find := func(name string) (ak, sk string, ok bool) {
		for _, p := range parsed.Profiles {
			if n, _ := p["name"].(string); n == name {
				ak, _ = p["access_key_id"].(string)
				sk, _ = p["access_key_secret"].(string)
				mode, _ := p["mode"].(string)
				if mode != "" && mode != "AK" {
					return "", "", false
				}
				return ak, sk, ak != "" && sk != ""
			}
		}
		return "", "", false
	}
	utAK, utSK, ok := find("TerraformUT")
	if !ok {
		t.Skipf("profile TerraformUT not usable (AK mode) in %s", cfgPath)
	}
	testAK, testSK, ok := find("TerraformTest")
	if !ok {
		t.Skipf("profile TerraformTest not usable (AK mode) in %s", cfgPath)
	}
	sharedCENCrossAccountCreds = cenCrossAccountCreds{
		utAK: utAK, utSK: utSK,
		testAK: testAK, testSK: testSK,
	}
	// Clear any env-based credentials for the duration of this single test;
	// t.Cleanup restores originals so neighboring tests / the shell are not
	// affected.
	for _, envVar := range []string{
		"ALICLOUD_ACCESS_KEY", "ALICLOUD_SECRET_KEY",
		"ALIBABA_CLOUD_ACCESS_KEY_ID", "ALIBABA_CLOUD_ACCESS_KEY_SECRET",
		"ALICLOUD_CLI_ACCESS_KEY", "ALICLOUD_CLI_SECRET_KEY",
	} {
		name := envVar
		original, present := os.LookupEnv(name)
		if present {
			if err := os.Unsetenv(name); err != nil {
				t.Fatalf("unset %s: %s", name, err)
			}
			t.Cleanup(func() {
				_ = os.Setenv(name, original)
			})
		}
	}
	if v := os.Getenv("ALICLOUD_REGION"); v == "" {
		log.Println("[INFO] Test: Using cn-beijing as test region")
		os.Setenv("ALICLOUD_REGION", "cn-beijing")
	} else {
		defaultRegionToTest = v
	}
}

// cenCrossAccountProviderFactories returns a ProviderFactories map plus a
// pointer-to-slice the factory populates with each created *schema.Provider.
// Callers can iterate the slice from a describe closure to pick the
// provider configured with the AK that owns the resource under test.
func cenCrossAccountProviderFactories() (map[string]terraform.ResourceProviderFactory, *[]*schema.Provider) {
	var providers []*schema.Provider
	factories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	return factories, &providers
}

// cenCrossAccountClientByAK returns the *AliyunClient from the provider whose
// AK matches wantAK, falling back to the first configured provider. The
// returned client is nil only if no providers have been configured yet.
func cenCrossAccountClientByAK(providers []*schema.Provider, wantAK string) *connectivity.AliyunClient {
	for _, p := range providers {
		meta := p.Meta()
		if meta == nil {
			continue
		}
		client := meta.(*connectivity.AliyunClient)
		if client.AccessKey == wantAK {
			return client
		}
	}
	for _, p := range providers {
		if meta := p.Meta(); meta != nil {
			return meta.(*connectivity.AliyunClient)
		}
	}
	return nil
}

// cenCrossAccountProviderBlocks renders two provider blocks — the default
// (account B / TerraformTest) and alias "a" (account A / TerraformUT) — with
// the injected AK/SK. All CEN cross-account test configs include this as the
// first fragment of their HCL.
func cenCrossAccountProviderBlocks() string {
	creds := sharedCENCrossAccountCreds
	return fmt.Sprintf(`
# Default provider = Account B (TerraformTest): owns the resource that gets
# granted to account A's CEN.
provider "alicloud" {
  access_key = %q
  secret_key = %q
}

# Alias "a" = Account A (TerraformUT): owns the CEN instance and the
# transit_router_*_attachment that consumes the grant.
provider "alicloud" {
  alias      = "a"
  access_key = %q
  secret_key = %q
}
`, creds.testAK, creds.testSK, creds.utAK, creds.utSK)
}
