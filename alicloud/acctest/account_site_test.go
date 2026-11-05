package acctest

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAccountSiteIsolation checks the real muxed provider, not just environment
// variables. Each subtest configures the same SDK provider and must leave the
// next test with the original identity.
func TestAccAccountSiteIsolation(t *testing.T) {
	var original, international string
	for _, tc := range []struct {
		name string
		site AccountSite
	}{
		{"DomesticBefore", DomesticSite}, {"International", IntlSite}, {"DomesticAfter", DomesticSite},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { PreCheckWithAccountSiteType(t, tc.site); PreCheck(t) },
				ProtoV5ProviderFactories: ProtoV5ProviderFactories(),
				Steps: []resource.TestStep{{Config: `data "alicloud_caller_identity" "current" {}`,
					Check: func(state *terraform.State) error {
						current := state.RootModule().Resources["data.alicloud_caller_identity.current"]
						if current == nil || current.Primary.Attributes["account_id"] == "" {
							return fmt.Errorf("caller identity was not returned")
						}
						account := current.Primary.Attributes["account_id"]
						t.Logf("Verified acceptance account site=%s accountId=%s", tc.site, account)
						switch tc.name {
						case "DomesticBefore":
							original = account
						case "International":
							international = account
							if original != "" && account == original {
								return fmt.Errorf("international test reused the domestic account")
							}
							if expected := os.Getenv("ALICLOUD_ACCOUNT_ID"); expected != "" && account != expected {
								return fmt.Errorf("selected account does not match configured account metadata")
							}
						case "DomesticAfter":
							if original != "" && account != original {
								return fmt.Errorf("domestic identity was not restored")
							}
						}
						return nil
					},
				}},
			})
		})
	}
	if os.Getenv("TF_ACC") != "" && original != "" && international == "" {
		t.Fatal("international identity was not exercised")
	}
}
