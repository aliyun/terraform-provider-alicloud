package cloudapi_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudApiGatewayDomainFramework_basic exercises the framework-native
// alicloud_api_gateway_domain across create -> import -> update wss_enable ->
// update the three SetDomainCertificate attributes. SetDomain requires a real
// ICP-filed domain: the case is skipped unless ALICLOUD_API_GATEWAY_DOMAIN_NAME
// is set, so it never sends a fake .example.com domain that ICP would reject.
func TestAccAliCloudApiGatewayDomainFramework_basic(t *testing.T) {
	domainName := os.Getenv("ALICLOUD_API_GATEWAY_DOMAIN_NAME")
	if domainName == "" {
		t.Skip("Skipping: ALICLOUD_API_GATEWAY_DOMAIN_NAME is not set; SetDomain requires an ICP-filed domain")
	}

	r := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-domain-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccApiGatewayDomainConfigBasic(name, domainName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_api_gateway_domain.default", "domain_name", domainName),
					resource.TestCheckResourceAttrSet("alicloud_api_gateway_domain.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_api_gateway_domain.default", "sub_domain"),
				),
			},
			{
				ResourceName:      "alicloud_api_gateway_domain.default",
				ImportState:       true,
				ImportStateVerify: true,
				// These four attributes are write-only: the read API does not
				// return them, so they cannot be verified on import.
				ImportStateVerifyIgnore: []string{
					"ssl_ocsp_enable",
					"ssl_ocsp_cache_enable",
					"client_cert_s_dn_pass_through",
					"wss_enable",
				},
			},
			{
				Config: testAccApiGatewayDomainConfigWss(name, domainName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_api_gateway_domain.default", "wss_enable", "ON"),
				),
			},
			{
				Config: testAccApiGatewayDomainConfigCert(name, domainName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_api_gateway_domain.default", "ssl_ocsp_enable", "true"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_domain.default", "ssl_ocsp_cache_enable", "true"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_domain.default", "client_cert_s_dn_pass_through", "true"),
				),
			},
		},
	})
}

func testAccApiGatewayDomainGroup(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_api_gateway_group" "default" {
  name        = var.name
  description = "tf_testAcc api gateway group for domain"
}
`, name)
}

func testAccApiGatewayDomainConfigBasic(name, domainName string) string {
	return fmt.Sprintf(`
%[1]s

resource "alicloud_api_gateway_domain" "default" {
  group_id    = alicloud_api_gateway_group.default.id
  domain_name = "%[2]s"
}
`, testAccApiGatewayDomainGroup(name), domainName)
}

func testAccApiGatewayDomainConfigWss(name, domainName string) string {
	return fmt.Sprintf(`
%[1]s

resource "alicloud_api_gateway_domain" "default" {
  group_id    = alicloud_api_gateway_group.default.id
  domain_name = "%[2]s"
  wss_enable  = "ON"
}
`, testAccApiGatewayDomainGroup(name), domainName)
}

func testAccApiGatewayDomainConfigCert(name, domainName string) string {
	return fmt.Sprintf(`
%[1]s

resource "alicloud_api_gateway_domain" "default" {
  group_id                       = alicloud_api_gateway_group.default.id
  domain_name                    = "%[2]s"
  ssl_ocsp_enable                = true
  ssl_ocsp_cache_enable          = true
  client_cert_s_dn_pass_through  = true
}
`, testAccApiGatewayDomainGroup(name), domainName)
}
