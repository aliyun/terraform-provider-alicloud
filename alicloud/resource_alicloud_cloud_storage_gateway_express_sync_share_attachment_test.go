package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCloudStorageGatewayExpressSyncShareAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_storage_gateway_express_sync_share_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudStorageGatewayExpressSyncShareAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SgwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressSyncShares")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccgatewayshare%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudStorageGatewayExpressSyncShareAttachmentBasicDependence0)
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
					"express_sync_id": "${alicloud_cloud_storage_gateway_express_sync.default.id}",
					"gateway_id":      "${alicloud_cloud_storage_gateway_gateway.default.id}",
					"share_name":      "${alicloud_cloud_storage_gateway_gateway_file_share.default.gateway_file_share_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"express_sync_id": CHECKSET,
						"gateway_id":      CHECKSET,
						"share_name":      CHECKSET,
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

var AlicloudCloudStorageGatewayExpressSyncShareAttachmentMap0 = map[string]string{}

func AlicloudCloudStorageGatewayExpressSyncShareAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

variable "region" {	
	default = "%s"
}

data "alicloud_cloud_storage_gateway_stocks" "default" {
  gateway_class = "Standard"
}

resource "alicloud_vpc" "vpc" {
  vpc_name   = var.name
  cidr_block = "192.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = "192.16.0.0/21"
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

resource "alicloud_cloud_storage_gateway_gateway_cache_disk" "default" {
  cache_disk_category   = "cloud_efficiency"
  gateway_id            = alicloud_cloud_storage_gateway_gateway.default.id
  cache_disk_size_in_gb = 50
}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
  acl    = "public-read-write"
}

resource "alicloud_cloud_storage_gateway_gateway_file_share" "default" {
  gateway_file_share_name = var.name
  gateway_id              = alicloud_cloud_storage_gateway_gateway.default.id
  local_path              = alicloud_cloud_storage_gateway_gateway_cache_disk.default.local_file_path
  oss_bucket_name         = alicloud_oss_bucket.default.bucket
  oss_endpoint            = alicloud_oss_bucket.default.extranet_endpoint
  protocol                = "NFS"
  remote_sync             = false
  fe_limit                = 0
  backend_limit           = 0
  cache_mode              = "Cache"
  squash                  = "none"
  lag_period              = 5
}

resource "alicloud_cloud_storage_gateway_express_sync" "default" {
  bucket_name       = alicloud_cloud_storage_gateway_gateway_file_share.default.oss_bucket_name
  bucket_region     = var.region
  description       = var.name
  express_sync_name = var.name
}
`, name, defaultRegionToTest)
}
