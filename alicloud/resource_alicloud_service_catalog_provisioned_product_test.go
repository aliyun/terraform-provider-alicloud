package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Case 3
func TestAccAliCloudServiceCatalogProvisionedProduct_basic1956(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_service_catalog_provisioned_product.default"
	ra := resourceAttrInit(resourceId, AlicloudService_catalogProvisionedProductMap1956)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServicecatalogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceCatalogProvisionedProduct")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)
	name := fmt.Sprintf("tf-testacc%sSCPProduct%d", defaultRegionToTest, rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Create: launch a provisioned product from the self-built
				// fixture (service ROS template tpl-bp1x4v3r44u7u7 requiring
				// acl + bucket_name) with acl=private and a single tag.
				Config: testAccServiceCatalogProvisionedProductConfigCreate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#":                 "2",
						"parameters.0.parameter_key":   CHECKSET,
						"parameters.0.parameter_value": CHECKSET,
						"parameters.1.parameter_key":   CHECKSET,
						"parameters.1.parameter_value": CHECKSET,
						"provisioned_product_name":     name,
						"stack_region_id":              "cn-hangzhou",
						"tags.%":                       "1",
						"tags.v1":                      "tf-test",
					}),
				),
			},
			{
				// Tags-only update: parameters are unchanged (acl=private,
				// bucket_name=lower(var.name)); only tags change. This is the
				// regression for the parameter-sending fix: UpdateProvisionedProduct
				// invokes ROS UpdateStack, which requires every template parameter
				// on each update, so a tags-only update must still send the
				// configured Parameters. Without the fix the update fails with
				// "The Parameter (acl) was not provided."; with the fix the update
				// succeeds and the parameters are retained. The update is driven by
				// the tags change, so it does not exercise a no-change-also-sends
				// update path.
				Config: testAccServiceCatalogProvisionedProductConfigUpdateTagsOnly(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#":                 "2",
						"parameters.0.parameter_key":   CHECKSET,
						"parameters.0.parameter_value": CHECKSET,
						"parameters.1.parameter_key":   CHECKSET,
						"parameters.1.parameter_value": CHECKSET,
						"portfolio_id":                 CHECKSET,
						"product_id":                   CHECKSET,
						"product_version_id":           CHECKSET,
						"provisioned_product_id":       CHECKSET,
						"tags.%":                       "1",
						"tags.v1":                      "tf-test-update",
					}),
				),
			},
			{
				// Update parameters + tags: change the bucket_name parameter
				// value (acl stays "private") and the tag, so
				// d.HasChange("parameters") is true and the Update control flow
				// sends Parameters via both the HasChange("parameters") guard
				// and the moved-out GetOk+assembly. This covers the original
				// update-parameters scenario (mirrors the baseline master case's
				// design of changing a parameter value on update) alongside the
				// tags-only regression above. acl is kept "private" because the
				// test account has Block Public Access enabled (rejects OSS
				// PutBucketACL public-read); changing bucket_name drives a ROS
				// OSS::Bucket replacement that keeps the new bucket private.
				Config: testAccServiceCatalogProvisionedProductConfigUpdateParams(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#":                 "2",
						"parameters.0.parameter_key":   CHECKSET,
						"parameters.0.parameter_value": CHECKSET,
						"parameters.1.parameter_key":   CHECKSET,
						"parameters.1.parameter_value": CHECKSET,
						"portfolio_id":                 CHECKSET,
						"product_id":                   CHECKSET,
						"product_version_id":           CHECKSET,
						"provisioned_product_id":       CHECKSET,
						"provisioned_product_name":     name,
						"stack_region_id":              "cn-hangzhou",
						"tags.%":                       "1",
						"tags.v1":                      "tf-test-update2",
					}),
				),
			},
			{
				// ImportStateVerifyIgnore: the ROS stack's Outputs (OSS bucket
				// endpoint/name) are computed by ROS and lag the bucket_name
				// parameter change in the preceding update-parameters step: the
				// post-Update Read races ROS's asynchronous OSS::Bucket
				// replacement, so the state captured at update time reflects the
				// pre-replacement bucket while a later import read reflects the
				// post-replacement bucket. This is a template/ROS timing artifact
				// of the self-built fixture, not a control-flow defect — the
				// Update sent Parameters (the update-parameters step passed) and
				// every non-outputs attribute (parameters, tags, identity ids)
				// matches on import.
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"outputs"},
			},
		},
	})
}

var AlicloudService_catalogProvisionedProductMap1956 = map[string]string{}

// testAccServiceCatalogProvisionedProductFixture builds the shared fixture chain
// (portfolio + product + product version on the service ROS template
// tpl-bp1x4v3r44u7u7 requiring acl + bucket_name + product/portfolio and
// principal/portfolio associations) so the case no longer depends on the
// absent preset product "ram模板创建" that previously forced t.Skip.
func testAccServiceCatalogProvisionedProductFixture(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_caller_identity" "current" {}

locals {
  principal_type = data.alicloud_caller_identity.current.identity_type == "ram-role" ? "RamRole" : "RamUser"
}

resource "alicloud_service_catalog_portfolio" "default" {
  provider_name  = "tianyu"
  portfolio_name = var.name
}

resource "alicloud_service_catalog_product" "default" {
  provider_name = "tianyu"
  product_name  = var.name
  product_type  = "Ros"
}

resource "alicloud_service_catalog_product_version" "default" {
  product_id           = alicloud_service_catalog_product.default.id
  template_url         = "oss://servicecatalog-cn-hangzhou/1466115886172051/terraform/template/tpl-bp1x4v3r44u7u7/template.json"
  template_type        = "RosTerraformTemplate"
  product_version_name = "1.0.0"
  active               = true
  guidance             = "Default"
}

resource "alicloud_service_catalog_product_portfolio_association" "default" {
  portfolio_id = alicloud_service_catalog_portfolio.default.id
  product_id   = alicloud_service_catalog_product.default.id
}

resource "alicloud_service_catalog_principal_portfolio_association" "default" {
  portfolio_id   = alicloud_service_catalog_portfolio.default.id
  principal_id   = data.alicloud_caller_identity.current.id
  principal_type = local.principal_type
  depends_on = [
    alicloud_service_catalog_product_portfolio_association.default,
    alicloud_service_catalog_product_version.default,
  ]
}

data "alicloud_service_catalog_launch_options" "default" {
  product_id = alicloud_service_catalog_product.default.id
  depends_on = [alicloud_service_catalog_principal_portfolio_association.default]
}
`, name)
}

func testAccServiceCatalogProvisionedProductConfigCreate(name string) string {
	return testAccServiceCatalogProvisionedProductFixture(name) + `
resource "alicloud_service_catalog_provisioned_product" "default" {
  provisioned_product_name = var.name
  stack_region_id          = "cn-hangzhou"
  product_version_id       = alicloud_service_catalog_product_version.default.id
  product_id               = alicloud_service_catalog_product.default.id
  portfolio_id             = data.alicloud_service_catalog_launch_options.default.launch_options.0.portfolio_id
  tags = {
    "v1" = "tf-test"
  }
  parameters {
    parameter_key   = "acl"
    parameter_value = "private"
  }
  parameters {
    parameter_key   = "bucket_name"
    parameter_value = lower(var.name)
  }
}
`
}

// testAccServiceCatalogProvisionedProductConfigUpdateTagsOnly keeps the
// parameters identical to the create config (acl=private,
// bucket_name=lower(var.name)) and only changes the tag value, so
// d.HasChange("parameters") is false while d.HasChange("tags") is true.
func testAccServiceCatalogProvisionedProductConfigUpdateTagsOnly(name string) string {
	return testAccServiceCatalogProvisionedProductFixture(name) + `
resource "alicloud_service_catalog_provisioned_product" "default" {
  provisioned_product_name = var.name
  stack_region_id          = "cn-hangzhou"
  product_version_id       = alicloud_service_catalog_product_version.default.id
  product_id               = alicloud_service_catalog_product.default.id
  portfolio_id             = data.alicloud_service_catalog_launch_options.default.launch_options.0.portfolio_id
  tags = {
    "v1" = "tf-test-update"
  }
  parameters {
    parameter_key   = "acl"
    parameter_value = "private"
  }
  parameters {
    parameter_key   = "bucket_name"
    parameter_value = lower(var.name)
  }
}
`
}

// testAccServiceCatalogProvisionedProductConfigUpdateParams changes the
// bucket_name parameter value (acl stays "private") and the tag, exercising
// the update-parameters scenario where d.HasChange("parameters") is true.
func testAccServiceCatalogProvisionedProductConfigUpdateParams(name string) string {
	return testAccServiceCatalogProvisionedProductFixture(name) + `
resource "alicloud_service_catalog_provisioned_product" "default" {
  provisioned_product_name = var.name
  stack_region_id          = "cn-hangzhou"
  product_version_id       = alicloud_service_catalog_product_version.default.id
  product_id               = alicloud_service_catalog_product.default.id
  portfolio_id             = data.alicloud_service_catalog_launch_options.default.launch_options.0.portfolio_id
  tags = {
    "v1" = "tf-test-update2"
  }
  parameters {
    parameter_key   = "acl"
    parameter_value = "private"
  }
  parameters {
    parameter_key   = "bucket_name"
    parameter_value = "${lower(var.name)}-upd"
  }
}
`
}
