package cloudapi_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudApiGatewayDomainDataSourceFramework_basic queries a bound domain
// through the data source after creating it. Skipped unless
// ALICLOUD_API_GATEWAY_DOMAIN_NAME is set, for the same ICP reason as the resource.
func TestAccAliCloudApiGatewayDomainDataSourceFramework_basic(t *testing.T) {
	domainName := os.Getenv("ALICLOUD_API_GATEWAY_DOMAIN_NAME")
	if domainName == "" {
		t.Skip("Skipping: ALICLOUD_API_GATEWAY_DOMAIN_NAME is not set; SetDomain requires an ICP-filed domain")
	}

	r := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-domain-ds-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccApiGatewayDomainDataSourceConfig(name, domainName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_api_gateway_domain.default", "domain_name", domainName),
					resource.TestCheckResourceAttrSet("data.alicloud_api_gateway_domain.default", "id"),
					resource.TestCheckResourceAttrSet("data.alicloud_api_gateway_domain.default", "sub_domain"),
				),
			},
		},
	})
}

func testAccApiGatewayDomainDataSourceConfig(name, domainName string) string {
	return fmt.Sprintf(`
%[1]s

resource "alicloud_api_gateway_domain" "default" {
  group_id    = alicloud_api_gateway_group.default.id
  domain_name = "%[2]s"
}

data "alicloud_api_gateway_domain" "default" {
  group_id    = alicloud_api_gateway_domain.default.group_id
  domain_name = alicloud_api_gateway_domain.default.domain_name
}
`, testAccApiGatewayDomainGroup(name), domainName)
}
