package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCloudStorageGatewayGatewayLogging_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_storage_gateway_gateway_logging.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudStorageGatewayGatewayLoggingMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SgwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudStorageGatewayGatewayLogging")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudstoragegatewaygatewaylogging%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudStorageGatewayGatewayLoggingBasicDependence0)
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
					"gateway_id":   "${data.alicloud_cloud_storage_gateway_gateways.default.ids.0}",
					"sls_logstore": "${alicloud_log_store.default.name}",
					"sls_project":  "${alicloud_log_project.default.name}",
					"status":       "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_id":   CHECKSET,
						"sls_logstore": CHECKSET,
						"sls_project":  CHECKSET,
						"status":       "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Disabled",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudCloudStorageGatewayGatewayLoggingMap0 = map[string]string{
	"gateway_id": CHECKSET,
	"status":     CHECKSET,
}

func AlicloudCloudStorageGatewayGatewayLoggingBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_cloud_storage_gateway_storage_bundles" "default" {
  backend_bucket_region_id = "%s"
  name_regex               = "default-NODELETING"
}

data "alicloud_cloud_storage_gateway_gateways" "default" {
  storage_bundle_id = data.alicloud_cloud_storage_gateway_storage_bundles.default.ids.0
  name_regex        = "default-NODELETING"
}

resource "alicloud_log_project" "default" {
  name        = var.name
  description = "created by terraform"
}

resource "alicloud_log_store" "default" {
  project               = alicloud_log_project.default.name
  name                  = var.name
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}
`, name, defaultRegionToTest)
}
