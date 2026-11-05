package apigateway_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/apigateway"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudApiGatewayDomain_basic(t *testing.T) {
	resourceId := "alicloud_api_gateway_domain.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-beijing"
	}
	name := fmt.Sprintf("tf-testAcc-domain%d", rand)
	domain := fmt.Sprintf("tf-testAcc-domain%d.example.com", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckApiGatewayDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudApiGatewayDomainBasic(name, domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "domain_name", domain),
					resource.TestCheckResourceAttrSet(resourceId, "group_id"),
					resource.TestCheckResourceAttr(resourceId, "custom_domain_type", "INTERNET"),
					resource.TestCheckResourceAttr(resourceId, "bind_stage_name", "RELEASE"),
					resource.TestCheckResourceAttrSet(resourceId, "sub_domain"),
					resource.TestCheckResourceAttrSet(resourceId, "domain_binding_status"),
					resource.TestCheckResourceAttrSet(resourceId, "domain_legal_status"),
				),
			},
			{
				Config: testAccAliCloudApiGatewayDomainUpdate(name, domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "domain_name", domain),
					resource.TestCheckResourceAttr(resourceId, "bind_stage_name", ""),
					resource.TestCheckResourceAttr(resourceId, "is_http_redirect_to_https", "true"),
					resource.TestCheckResourceAttr(resourceId, "is_force", "true"),
					resource.TestCheckResourceAttrSet(resourceId, "sub_domain"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bind_stage_name", "custom_domain_type", "is_http_redirect_to_https", "is_force"},
			},
		},
	})
}

func testAccCheckApiGatewayDomainDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_api_gateway_domain" {
			continue
		}
		_, err := apigateway.NewApiGatewayServiceV2(client).DescribeApiGatewayDomain(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing api gateway domain %s failed: %s", rs.Primary.ID, err)
		}
		return fmt.Errorf("api gateway domain %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccAliCloudApiGatewayDomainBasic(name, domain string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%[1]s"
}

resource "alicloud_api_gateway_group" "default" {
  name = var.name
}

resource "alicloud_api_gateway_domain" "default" {
  group_id          = alicloud_api_gateway_group.default.id
  domain_name       = "%[2]s"
  custom_domain_type = "INTERNET"
  bind_stage_name   = "RELEASE"
}
`, name, domain)
}

func testAccAliCloudApiGatewayDomainUpdate(name, domain string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%[1]s"
}

resource "alicloud_api_gateway_group" "default" {
  name = var.name
}

resource "alicloud_api_gateway_domain" "default" {
  group_id            = alicloud_api_gateway_group.default.id
  domain_name         = "%[2]s"
  custom_domain_type  = "INTERNET"
  bind_stage_name     = ""
  is_http_redirect_to_https = true
  is_force            = true
}
`, name, domain)
}
