package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 3
func TestAccAlicloudService_catalogProvisionedProduct_basic1956(t *testing.T) {
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudService_catalogProvisionedProductBasicDependence1956)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"parameters": []map[string]interface{}{
						{
							"parameter_key":   "role_name",
							"parameter_value": "${var.name}",
						},
					},
					"provisioned_product_name": "${var.name}",
					"stack_region_id":          "${data.alicloud_ros_regions.all.regions.5.region_id}",
					"product_version_id":       "${data.alicloud_service_catalog_product_versions.default.versions.0.id}",
					"product_id":               "${data.alicloud_service_catalog_product_as_end_users.default.users.0.id}",
					"tags": map[string]string{
						"k4": "v4",
					},
					"portfolio_id": "${data.alicloud_service_catalog_launch_options.default.options.0.portfolio_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#":           "1",
						"product_version_id":     CHECKSET,
						"provisioned_product_id": CHECKSET,
						"product_id":             CHECKSET,
						"tags.%":                 "1",
						"portfolio_id":           CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parameters": []map[string]interface{}{
						{
							"parameter_key":   "role_name",
							"parameter_value": "${var.name}-update",
						},
					},
					"tags": map[string]string{
						"k4": "v4_update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#": "1",
						"tags.%":       "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudService_catalogProvisionedProductMap1956 = map[string]string{}

func AlicloudService_catalogProvisionedProductBasicDependence1956(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_service_catalog_product_as_end_users" "default" {
  name_regex = "ram模板创建"
}

data "alicloud_service_catalog_product_versions" "default" {
  name_regex = "1.0.0"
  product_id = data.alicloud_service_catalog_product_as_end_users.default.users.0.id
}

data "alicloud_service_catalog_launch_options" "default" {
  product_id = data.alicloud_service_catalog_product_as_end_users.default.users.0.id
}

data "alicloud_ros_regions" "all" {}
`, name)
}
