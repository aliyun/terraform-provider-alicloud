package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_api_gateway_domain", &resource.Sweeper{
		Name:         "alicloud_api_gateway_domain",
		F:            testSweepApiGatewayDomain,
		Dependencies: []string{"alicloud_api_gateway_group"},
	})
}

// Domain has no dedicated list API; domains are removed when their parent
// API Gateway group is swept by testSweepApiGatewayGroup.
func testSweepApiGatewayDomain(region string) error {
	return nil
}

// TestAccAlicloudApiGatewayDomain_basic covers create -> import -> update wss_enable
// -> update ssl_ocsp / client_cert_s_dn_pass_through. All four update-only attributes
// appear in at least one Step Config to satisfy the TestingCoverageRate 100% gate.
func TestAccAliCloudApiGatewayDomain_basic(t *testing.T) {
	var v *cloudapi.DescribeDomainResponse
	resourceId := "alicloud_api_gateway_domain.default"
	ra := resourceAttrInit(resourceId, apiGatewayDomainBasicMap)
	serviceFunc := func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeApiGatewayDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAccDomain_%d", rand)
	// SetDomain requires an ICP-filed domain; a fake .example.com domain can
	// never pass ICP verification. The real domain is injected via env var; when
	// unset the env-guard preCheck skips the case so this placeholder is never used.
	domainName := os.Getenv("ALICLOUD_API_GATEWAY_DOMAIN_NAME")
	if domainName == "" {
		domainName = fmt.Sprintf("tf-testacc-domain-%d.example.com", rand)
	}
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceApigatewayDomainConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithEnvVariable(t, "ALICLOUD_API_GATEWAY_DOMAIN_NAME")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"group_id":    "${alicloud_api_gateway_group.default.id}",
					"domain_name": domainName,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name": domainName,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// These four attributes are write-only (@readonly): the read API
				// does not return them, so they cannot be verified on import.
				ImportStateVerifyIgnore: []string{
					"ssl_ocsp_enable",
					"ssl_ocsp_cache_enable",
					"client_cert_s_dn_pass_through",
					"wss_enable",
				},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_id":    "${alicloud_api_gateway_group.default.id}",
					"domain_name": domainName,
					"wss_enable":  "ON",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_id":                      "${alicloud_api_gateway_group.default.id}",
					"domain_name":                   domainName,
					"ssl_ocsp_enable":               "true",
					"ssl_ocsp_cache_enable":         "true",
					"client_cert_s_dn_pass_through": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceApigatewayDomainConfigDependence(name string) string {
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

var apiGatewayDomainBasicMap = map[string]string{
	"domain_name": CHECKSET,
}
