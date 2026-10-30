package cloudstoragegateway_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudCloudStorageGatewayGatewayBlockVolumesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudStorageGatewayGatewayBlockVolumesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cloud_storage_gateway_gateway_block_volume.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudCloudStorageGatewayGatewayBlockVolumesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cloud_storage_gateway_gateway_block_volume.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudStorageGatewayGatewayBlockVolumesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cloud_storage_gateway_gateway_block_volume.default.gateway_block_volume_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudCloudStorageGatewayGatewayBlockVolumesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cloud_storage_gateway_gateway_block_volume.default.gateway_block_volume_name}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudStorageGatewayGatewayBlockVolumesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cloud_storage_gateway_gateway_block_volume.default.id}"]`,
			"status": `"0"`,
		}),
		FakeConfig: testAccCheckAlicloudCloudStorageGatewayGatewayBlockVolumesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cloud_storage_gateway_gateway_block_volume.default.id}"]`,
			"status": `"1"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudStorageGatewayGatewayBlockVolumesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cloud_storage_gateway_gateway_block_volume.default.id}"]`,
			"name_regex": `"${alicloud_cloud_storage_gateway_gateway_block_volume.default.gateway_block_volume_name}"`,
			"status":     `"0"`,
		}),
		FakeConfig: testAccCheckAlicloudCloudStorageGatewayGatewayBlockVolumesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cloud_storage_gateway_gateway_block_volume.default.id}_fake"]`,
			"name_regex": `"${alicloud_cloud_storage_gateway_gateway_block_volume.default.gateway_block_volume_name}_fake"`,
			"status":     `"1"`,
		}),
	}

	var existAlicloudCloudStorageGatewayGatewayBlockVolumesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"names.#":                             "1",
			"volumes.#":                           "1",
			"volumes.0.address":                   acctest.CHECKSET,
			"volumes.0.cache_mode":                "Cache",
			"volumes.0.chap_enabled":              "true",
			"volumes.0.chap_in_user":              acctest.CHECKSET,
			"volumes.0.chunk_size":                "8192",
			"volumes.0.disk_id":                   acctest.CHECKSET,
			"volumes.0.disk_type":                 acctest.CHECKSET,
			"volumes.0.enabled":                   acctest.CHECKSET,
			"volumes.0.gateway_block_volume_name": fmt.Sprintf("tftestacc%d", rand),
			"volumes.0.gateway_id":                acctest.CHECKSET,
			"volumes.0.id":                        acctest.CHECKSET,
			"volumes.0.index_id":                  acctest.CHECKSET,
			"volumes.0.local_path":                acctest.CHECKSET,
			"volumes.0.lun_id":                    acctest.CHECKSET,
			"volumes.0.oss_bucket_name":           acctest.CHECKSET,
			"volumes.0.oss_bucket_ssl":            "true",
			"volumes.0.oss_endpoint":              acctest.CHECKSET,
			"volumes.0.port":                      acctest.CHECKSET,
			"volumes.0.protocol":                  acctest.CHECKSET,
			"volumes.0.size":                      acctest.CHECKSET,
			"volumes.0.state":                     acctest.CHECKSET,
			"volumes.0.status":                    acctest.CHECKSET,
			"volumes.0.target":                    acctest.CHECKSET,
			"volumes.0.total_download":            acctest.CHECKSET,
			"volumes.0.total_upload":              acctest.CHECKSET,
			"volumes.0.volume_state":              acctest.CHECKSET,
		}
	}
	var fakeAlicloudCloudStorageGatewayGatewayBlockVolumesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudCloudStorageGatewayGatewayBlockVolumesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cloud_storage_gateway_gateway_block_volumes.default",
		ExistMapFunc: existAlicloudCloudStorageGatewayGatewayBlockVolumesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudCloudStorageGatewayGatewayBlockVolumesDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudCloudStorageGatewayGatewayBlockVolumesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudCloudStorageGatewayGatewayBlockVolumesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tftestacc%d"
}

data "alicloud_cloud_storage_gateway_stocks" "default" {
  gateway_class = "Standard"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
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

resource "alicloud_cloud_storage_gateway_gateway_block_volume" "default" {
  cache_mode                = "Cache"
  chap_enabled              = true
  chap_in_user              = var.name
  chap_in_password          = var.name
  chunk_size                = "8192"
  gateway_block_volume_name = var.name
  gateway_id                = alicloud_cloud_storage_gateway_gateway.default.id
  local_path                = alicloud_cloud_storage_gateway_gateway_cache_disk.default.local_file_path
  oss_bucket_name           = alicloud_oss_bucket.default.bucket
  oss_bucket_ssl            = true
  oss_endpoint              = alicloud_oss_bucket.default.extranet_endpoint
  protocol                  = "iSCSI"
  size                      = 100
}

data "alicloud_cloud_storage_gateway_gateway_block_volumes" "default" {
  gateway_id = alicloud_cloud_storage_gateway_gateway.default.id
  %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
