package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudStorageGatewayExpressSyncsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudStorageGatewayExpressSyncsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cloud_storage_gateway_express_sync.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCloudStorageGatewayExpressSyncsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cloud_storage_gateway_express_sync.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudStorageGatewayExpressSyncsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cloud_storage_gateway_express_sync.default.express_sync_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudStorageGatewayExpressSyncsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cloud_storage_gateway_express_sync.default.express_sync_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudStorageGatewayExpressSyncsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cloud_storage_gateway_express_sync.default.id}"]`,
			"name_regex": `"${alicloud_cloud_storage_gateway_express_sync.default.express_sync_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudStorageGatewayExpressSyncsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cloud_storage_gateway_express_sync.default.id}_fake"]`,
			"name_regex": `"${alicloud_cloud_storage_gateway_express_sync.default.express_sync_name}_fake"`,
		}),
	}
	var existAlicloudCloudStorageGatewayExpressSyncsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"names.#":                   "1",
			"syncs.#":                   "1",
			"syncs.0.bucket_name":       CHECKSET,
			"syncs.0.bucket_prefix":     "",
			"syncs.0.bucket_region":     CHECKSET,
			"syncs.0.description":       fmt.Sprintf("tf-testaccexpresssync-%d", rand),
			"syncs.0.id":                CHECKSET,
			"syncs.0.express_sync_id":   CHECKSET,
			"syncs.0.express_sync_name": fmt.Sprintf("tf-testaccexpresssync-%d", rand),
			"syncs.0.mns_topic":         CHECKSET,
		}
	}
	var fakeAlicloudCloudStorageGatewayExpressSyncsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudCloudStorageGatewayExpressSyncsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cloud_storage_gateway_express_syncs.default",
		existMapFunc: existAlicloudCloudStorageGatewayExpressSyncsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCloudStorageGatewayExpressSyncsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudCloudStorageGatewayExpressSyncsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudCloudStorageGatewayExpressSyncsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testaccexpresssync-%d"
}

variable "region" {	
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
	zone_id      = data.alicloud_cloud_storage_gateway_stocks.default.stocks.0.zone_id
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
  type                     = "File"
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

resource "alicloud_cloud_storage_gateway_gateway_file_share" "default" {
  gateway_file_share_name = var.name
  gateway_id              = alicloud_cloud_storage_gateway_gateway.default.id
  local_path              = alicloud_cloud_storage_gateway_gateway_cache_disk.default.local_file_path
  oss_bucket_name         = alicloud_oss_bucket.default.bucket
  oss_endpoint            = alicloud_oss_bucket.default.extranet_endpoint
  protocol                = "NFS"
  remote_sync             = true
  polling_interval        = 4500
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

data "alicloud_cloud_storage_gateway_express_syncs" "default" {	
	%s
}
`, rand, defaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
