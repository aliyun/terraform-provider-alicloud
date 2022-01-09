package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCloudStorageGatewayGatewayBlockVolume_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_storage_gateway_gateway_block_volume.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudStorageGatewayGatewayBlockVolumeMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SgwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudStorageGatewayGatewayBlockVolume")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tftestacccsvolume%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudStorageGatewayGatewayBlockVolumeBasicDependence0)
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
					"gateway_id":                "${alicloud_cloud_storage_gateway_gateway.default.id}",
					"gateway_block_volume_name": name,
					"chunk_size":                "8192",
					"chap_enabled":              "false",
					"oss_endpoint":              "${alicloud_oss_bucket.default.extranet_endpoint}",
					"oss_bucket_name":           "${alicloud_oss_bucket.default.bucket}",
					"cache_mode":                "Cache",
					"local_path":                "${alicloud_cloud_storage_gateway_gateway_cache_disk.default.local_file_path}",
					"protocol":                  "iSCSI",
					"oss_bucket_ssl":            "true",
					"size":                      "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_block_volume_name": name,
						"chunk_size":                "8192",
						"chap_enabled":              "false",
						"oss_endpoint":              CHECKSET,
						"oss_bucket_name":           CHECKSET,
						"cache_mode":                "Cache",
						"local_path":                CHECKSET,
						"protocol":                  "iSCSI",
						"oss_bucket_ssl":            "true",
						"gateway_id":                CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"chap_enabled":     "true",
					"chap_in_user":     "tftestAccnmSa123",
					"chap_in_password": "tftestAccnmSa456",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"chap_enabled": "true",
						"chap_in_user": "tftestAccnmSa123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"size": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"chap_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"chap_enabled": "false",
						"chap_in_user": "",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"is_source_deletion", "recovery", "size", "chap_in_password"},
			},
		},
	})
}

var AlicloudCloudStorageGatewayGatewayBlockVolumeMap0 = map[string]string{
	"is_source_deletion": NOSET,
	"recovery":           NOSET,
}

func AlicloudCloudStorageGatewayGatewayBlockVolumeBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_cloud_storage_gateway_stocks" "default" {
  gateway_class = "Standard"
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_cloud_storage_gateway_stocks.default.stocks.0.zone_id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_cloud_storage_gateway_stocks.default.stocks.0.zone_id
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_cloud_storage_gateway_storage_bundle" "default" {
  storage_bundle_name = var.name
}

resource "alicloud_cloud_storage_gateway_gateway" "default" {
  description              = "tf-acctestDesalone"
  gateway_class            = "Standard"
  type                     = "Iscsi"
  payment_type             = "PayAsYouGo"
  vswitch_id               = local.vswitch_id
  release_after_expiration = true
  public_network_bandwidth = 10
  storage_bundle_id        = alicloud_cloud_storage_gateway_storage_bundle.default.id
  location                 = "Cloud"
  gateway_name             = var.name
}


resource "alicloud_cloud_storage_gateway_gateway_cache_disk" "default" {
  cache_disk_category   = "cloud_efficiency"
  gateway_id            = alicloud_cloud_storage_gateway_gateway.default.id
  cache_disk_size_in_gb = 50
}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
  acl    = "public-read-write"
}
`, name)
}
