package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccAlicloudApiGatewayDomainDataSource_basic creates a domain and queries it
// through the data source using filter interpolation (no depends_on on the data block,
// per the established convention).
func TestAccAliCloudApiGatewayDomainDataSource_basic(t *testing.T) {
	resourceId := "alicloud_api_gateway_domain.default"
	dataSourceId := "data.alicloud_api_gateway_domain.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAccDomainDs_%d", rand)
	// SetDomain requires an ICP-filed domain; see resource test for rationale.
	domainName := os.Getenv("ALICLOUD_API_GATEWAY_DOMAIN_NAME")
	if domainName == "" {
		domainName = fmt.Sprintf("tf-testacc-domain-ds-%d.example.com", rand)
	}
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceApigatewayDomainConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithEnvVariable(t, "ALICLOUD_API_GATEWAY_DOMAIN_NAME")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"group_id":    "${alicloud_api_gateway_group.default.id}",
					"domain_name": domainName,
				}) + fmt.Sprintf(`
data "alicloud_api_gateway_domain" "default" {
  group_id    = alicloud_api_gateway_domain.default.group_id
  domain_name = alicloud_api_gateway_domain.default.domain_name
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceId, "domain_name", domainName),
					resource.TestCheckResourceAttrSet(dataSourceId, "id"),
				),
			},
		},
	})
}
