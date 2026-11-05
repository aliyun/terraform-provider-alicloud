package ddoscoo_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAliCloudDdosCooWebCacheCustomRule covers the Framework-native
// alicloud_ddoscoo_web_cache_custom_rule resource: create, update the mutable attributes
// (uri / mode / cache_ttl) and import. The rule depends on a website domain resource that is
// created from an existing DdosCoo instance, mirroring alicloud_ddoscoo_web_cc_rule.
func TestAccAliCloudDdosCooWebCacheCustomRule(t *testing.T) {
	r := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf-acc-cache-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckDdosCooWebCacheCustomRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDdosCooWebCacheCustomRuleConfig(name, "/blog/", "standard", 3600),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_ddoscoo_web_cache_custom_rule.default", "id"),
					resource.TestCheckResourceAttr("alicloud_ddoscoo_web_cache_custom_rule.default", "name", name),
					resource.TestCheckResourceAttrSet("alicloud_ddoscoo_web_cache_custom_rule.default", "domain"),
					resource.TestCheckResourceAttr("alicloud_ddoscoo_web_cache_custom_rule.default", "uri", "/blog/"),
					resource.TestCheckResourceAttr("alicloud_ddoscoo_web_cache_custom_rule.default", "mode", "standard"),
					resource.TestCheckResourceAttr("alicloud_ddoscoo_web_cache_custom_rule.default", "cache_ttl", "3600"),
					resource.TestCheckResourceAttrSet("alicloud_ddoscoo_web_cache_custom_rule.default", "request_id"),
				),
			},
			{
				Config: testAccDdosCooWebCacheCustomRuleConfig(name, "/news/", "aggressive", 7200),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_ddoscoo_web_cache_custom_rule.default", "uri", "/news/"),
					resource.TestCheckResourceAttr("alicloud_ddoscoo_web_cache_custom_rule.default", "mode", "aggressive"),
					resource.TestCheckResourceAttr("alicloud_ddoscoo_web_cache_custom_rule.default", "cache_ttl", "7200"),
					resource.TestCheckResourceAttr("alicloud_ddoscoo_web_cache_custom_rule.default", "name", name),
				),
			},
			{
				Config: testAccDdosCooWebCacheCustomRuleConfig(name, "/api/", "bypass", 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_ddoscoo_web_cache_custom_rule.default", "uri", "/api/"),
					resource.TestCheckResourceAttr("alicloud_ddoscoo_web_cache_custom_rule.default", "mode", "bypass"),
				),
			},
			{
				ResourceName:      "alicloud_ddoscoo_web_cache_custom_rule.default",
				ImportState:       true,
				ImportStateVerify: true,
				// request_id is refreshed on every apply and is not stable across reads, so it is
				// ignored during import verification.
				ImportStateVerifyIgnore: []string{"request_id"},
			},
		},
	})
}

func testAccDdosCooWebCacheCustomRuleConfig(name, uri, mode string, cacheTtl int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%[1]s"
}

data "alicloud_ddoscoo_instances" "default" {
}

resource "alicloud_ddoscoo_domain_resource" "default" {
  domain       = "tf-testacc%[1]s.alibaba.com"
  instance_ids = [data.alicloud_ddoscoo_instances.default.ids.0]
  proxy_types {
    proxy_ports = [443]
    proxy_type  = "https"
  }
  real_servers = ["177.167.32.11"]
  rs_type      = 0
}

resource "alicloud_ddoscoo_web_cache_custom_rule" "default" {
  domain    = alicloud_ddoscoo_domain_resource.default.id
  name      = var.name
  uri       = "%[2]s"
  mode      = "%[3]s"
  cache_ttl = %[4]d
}
`, name, uri, mode, cacheTtl)
}

// testAccCheckDdosCooWebCacheCustomRuleDestroy verifies that every rule left in state has been
// removed by re-reading DescribeWebCacheConfigs for its domain and asserting the named rule is
// gone from CustomRules.
func testAccCheckDdosCooWebCacheCustomRuleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ddoscoo_web_cache_custom_rule" {
			continue
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		if err := ddosCooWebCacheCustomRuleAssertAbsent(client, rs.Primary.ID); err != nil {
			return err
		}
	}
	return nil
}

// ddosCooWebCacheCustomRuleAssertAbsent splits the `<domain>:<name>` ID and asks DescribeWebCacheConfigs
// whether the named rule still exists; a missing domain or rule counts as destroyed.
func ddosCooWebCacheCustomRuleAssertAbsent(client *connectivity.AliyunClient, id string) error {
	domain, name, err := ddosCooWebCacheCustomRuleSplitID(id)
	if err != nil {
		return fmt.Errorf("checking destroy of alicloud_ddoscoo_web_cache_custom_rule %s: %w", id, err)
	}
	request := map[string]interface{}{
		"Domains.1": domain,
		"RegionId":  client.RegionId,
	}
	response, err := client.RpcPost("ddoscoo", "2020-01-01", "DescribeWebCacheConfigs", nil, request, true)
	if err != nil {
		return fmt.Errorf("checking destroy of alicloud_ddoscoo_web_cache_custom_rule %s: %w", id, err)
	}
	configs, _ := response["DomainCacheConfigs"].([]interface{})
	for _, c := range configs {
		cfg, _ := c.(map[string]interface{})
		rules, _ := cfg["CustomRules"].([]interface{})
		for _, r := range rules {
			rule, _ := r.(map[string]interface{})
			if fmt.Sprint(rule["Name"]) == name {
				return fmt.Errorf("the resource alicloud_ddoscoo_web_cache_custom_rule %s was not destroyed: rule %q still exists", id, name)
			}
		}
	}
	return nil
}

func ddosCooWebCacheCustomRuleSplitID(id string) (string, string, error) {
	for i := len(id) - 1; i >= 0; i-- {
		if id[i] == ':' {
			domain, name := id[:i], id[i+1:]
			if domain != "" && name != "" {
				return domain, name, nil
			}
			break
		}
	}
	return "", "", fmt.Errorf("invalid resource id %q: expected <domain>:<name>", id)
}
