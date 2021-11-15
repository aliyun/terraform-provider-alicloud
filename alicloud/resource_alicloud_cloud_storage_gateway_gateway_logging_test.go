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
					"gateway_id":   "${alicloud_cloud_storage_gateway_gateway.default.id}",
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

data "alicloud_cloud_storage_gateway_stocks" "default" {
  gateway_class = "Standard"
}

resource "alicloud_vpc" "vpc" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_cloud_storage_gateway_stocks.default.stocks.0.zone_id
  vswitch_name = var.name
}

resource "alicloud_cloud_storage_gateway_storage_bundle" "default" {
  storage_bundle_name = var.name
}

resource "alicloud_cloud_storage_gateway_gateway" "default" {
  description              = "tf-acctestDesalone"
  gateway_class            = "Standard"
  type                     = "File"
  payment_type             = "PayAsYouGo"
  vswitch_id               = alicloud_vswitch.default.id
  release_after_expiration = true
  public_network_bandwidth = 10
  storage_bundle_id        = alicloud_cloud_storage_gateway_storage_bundle.default.id
  location                 = "Cloud"
  gateway_name             = var.name
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
`, name)
}
