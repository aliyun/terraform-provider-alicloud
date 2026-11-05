package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	esa "github.com/aliyun/terraform-provider-alicloud/alicloud/service/esa"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAliCloudEsaApiSecToken_basic0 exercises the create / update / import
// lifecycle of alicloud_esa_api_sec_token.
//
// The ApiSecToken APIs (AddApiSecToken / GetApiSecToken / SetApiSecToken /
// DeleteApiSecToken) are Private per the frozen CloudSpec evidence
// (cloudspecCommit 5ff43bdc). The inner members of TokenSources (array<struct>)
// and Credentials (array<JWK>) are not enumerated in that evidence, so those
// attributes are opaque JSON strings and the example payloads below are
// illustrative pending ESA product-owner confirmation (浩德, emp 229705).
// ACC therefore validates the resource CRUD plumbing; any API-level rejection
// of the payload shape is a product-contract gap to confirm with ESA, not a
// code defect.
func TestAccAliCloudEsaApiSecToken_basic0(t *testing.T) {
	resourceId := "alicloud_esa_api_sec_token.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sEsaApiSecToken%d", acctest.DefaultRegionToTest, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEsaApiSecTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudEsaApiSecTokenBasic0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttrSet(resourceId, "token_id"),
					resource.TestCheckResourceAttr(resourceId, "token_name", name),
					resource.TestCheckResourceAttrSet(resourceId, "token_sources"),
					resource.TestCheckResourceAttrSet(resourceId, "credentials"),
				),
			},
			{
				Config: testAccAliCloudEsaApiSecTokenBasic0Update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "token_name", name+"-upd"),
					resource.TestCheckResourceAttrSet(resourceId, "token_sources"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token_sources", "credentials"},
			},
		},
	})
}

func testAccCheckEsaApiSecTokenDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_esa_api_sec_token" {
			continue
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		_, err := esa.NewService(client).DescribeApiSecToken(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing api sec token %s failed: %s", rs.Primary.ID, err)
		}
		return fmt.Errorf("api sec token %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccAliCloudEsaApiSecTokenBasic0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%[1]s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "default" {
  site_name   = "${lower(var.name)}.com"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_api_sec_token" "default" {
  site_id       = alicloud_esa_site.default.id
  token_name    = var.name
  token_type    = "jwt"
  token_sources = jsonencode([{ source_type = "api" }])
  credentials   = jsonencode([{ kty = "oct", k = "dGVzdC1rZXk" }])
}
`, name)
}

func testAccAliCloudEsaApiSecTokenBasic0Update(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%[1]s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "default" {
  site_name   = "${lower(var.name)}.com"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_api_sec_token" "default" {
  site_id       = alicloud_esa_site.default.id
  token_name    = "${var.name}-upd"
  token_type    = "jwt"
  token_sources = jsonencode([{ source_type = "api" }, { source_type = "api-group" }])
  credentials   = jsonencode([{ kty = "oct", k = "dGVzdC1rZXk" }])
}
`, name)
}
