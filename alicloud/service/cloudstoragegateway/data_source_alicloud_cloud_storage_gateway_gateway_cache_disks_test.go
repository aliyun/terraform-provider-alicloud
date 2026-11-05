package cloudstoragegateway_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudCloudStorageGatewayGatewayCacheDisksDataSource(t *testing.T) {
	resourceId := "data.alicloud_cloud_storage_gateway_gateway_cache_disks.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-cloudstoragegatewaygatewaycachedisk-%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceCloudStorageGatewayGatewayCacheDisksDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"gateway_id": "${alicloud_cloud_storage_gateway_gateway_cache_disk.default.gateway_id}",
			"ids":        []string{"${alicloud_cloud_storage_gateway_gateway_cache_disk.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"gateway_id": "${alicloud_cloud_storage_gateway_gateway_cache_disk.default.gateway_id}",
			"ids":        []string{"${alicloud_cloud_storage_gateway_gateway_cache_disk.default.id}-fake"},
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"gateway_id": "${alicloud_cloud_storage_gateway_gateway_cache_disk.default.gateway_id}",
			"ids":        []string{"${alicloud_cloud_storage_gateway_gateway_cache_disk.default.id}"},
			"status":     "0",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"gateway_id": "${alicloud_cloud_storage_gateway_gateway_cache_disk.default.gateway_id}",
			"ids":        []string{"${alicloud_cloud_storage_gateway_gateway_cache_disk.default.id}"},
			"status":     "1",
		}),
	}
	var existCloudStorageGatewayGatewayCacheDiskMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"ids.0":                         acctest.CHECKSET,
			"disks.#":                       "1",
			"disks.0.status":                "0",
			"disks.0.cache_disk_category":   "cloud_efficiency",
			"disks.0.cache_disk_size_in_gb": "50",
			"disks.0.cache_id":              acctest.CHECKSET,
			"disks.0.expired_time":          acctest.CHECKSET,
			"disks.0.gateway_id":            acctest.CHECKSET,
			"disks.0.iops":                  acctest.CHECKSET,
			"disks.0.is_used":               acctest.CHECKSET,
			"disks.0.id":                    acctest.CHECKSET,
			"disks.0.local_file_path":       acctest.CHECKSET,
			"disks.0.renew_url":             "",
		}
	}

	var fakeCloudStorageGatewayGatewayCacheDiskMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"disks.#": "0",
		}
	}

	var CloudStorageGatewayGatewayCacheDiskCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existCloudStorageGatewayGatewayCacheDiskMapFunc,
		FakeMapFunc:  fakeCloudStorageGatewayGatewayCacheDiskMapFunc,
	}

	CloudStorageGatewayGatewayCacheDiskCheckInfo.DataSourceTestCheck(t, rand, idsConf, statusConf)
}

func dataSourceCloudStorageGatewayGatewayCacheDisksDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
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
}`, name)
}
