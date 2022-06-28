package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudStorageGatewayGatewayFileSharesDataSource(t *testing.T) {
	resourceId := "data.alicloud_cloud_storage_gateway_gateway_file_shares.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cloudstoragegatewaygatewayfileshare-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceCloudStorageGatewayGatewayFileSharesDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"gateway_id": "${alicloud_cloud_storage_gateway_gateway_file_share.default.gateway_id}",
			"name_regex": "${alicloud_cloud_storage_gateway_gateway_file_share.default.gateway_file_share_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"gateway_id": "${alicloud_cloud_storage_gateway_gateway_file_share.default.gateway_id}",
			"name_regex": "${alicloud_cloud_storage_gateway_gateway_file_share.default.gateway_file_share_name}-fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"gateway_id": "${alicloud_cloud_storage_gateway_gateway_file_share.default.gateway_id}",
			"ids":        []string{"${alicloud_cloud_storage_gateway_gateway_file_share.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"gateway_id": "${alicloud_cloud_storage_gateway_gateway_file_share.default.gateway_id}",
			"ids":        []string{"${alicloud_cloud_storage_gateway_gateway_file_share.default.id}-fake"},
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"gateway_id": "${alicloud_cloud_storage_gateway_gateway_file_share.default.gateway_id}",
			"name_regex": "${alicloud_cloud_storage_gateway_gateway_file_share.default.gateway_file_share_name}",
			"ids":        []string{"${alicloud_cloud_storage_gateway_gateway_file_share.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"gateway_id": "${alicloud_cloud_storage_gateway_gateway_file_share.default.gateway_id}",
			"name_regex": "${alicloud_cloud_storage_gateway_gateway_file_share.default.gateway_file_share_name}-fake",
			"ids":        []string{"${alicloud_cloud_storage_gateway_gateway_file_share.default.id}"},
		}),
	}
	var existCloudStorageGatewayGatewayFileShareMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"ids.0":                             CHECKSET,
			"shares.#":                          "1",
			"shares.0.id":                       CHECKSET,
			"shares.0.gateway_file_share_name":  fmt.Sprintf("tf-testacc-cloudstoragegatewaygatewayfileshare-%d", rand),
			"shares.0.access_based_enumeration": "false",
			"shares.0.address":                  CHECKSET,
			"shares.0.backend_limit":            "0",
			"shares.0.browsable":                "false",
			"shares.0.bucket_infos":             CHECKSET,
			"shares.0.buckets_stub":             "false",
			"shares.0.cache_mode":               "Cache",
			"shares.0.client_side_cmk":          "",
			"shares.0.client_side_encryption":   "false",
			"shares.0.direct_io":                "false",
			"shares.0.disk_id":                  CHECKSET,
			"shares.0.disk_type":                CHECKSET,
			"shares.0.download_limit":           "0",
			"shares.0.enabled":                  "true",
			"shares.0.express_sync_id":          "",
			"shares.0.fast_reclaim":             "false",
			"shares.0.fe_limit":                 "0",
			"shares.0.file_num_limit":           CHECKSET,
			"shares.0.fs_size_limit":            CHECKSET,
			"shares.0.gateway_id":               CHECKSET,
			"shares.0.ignore_delete":            "false",
			"shares.0.in_place":                 "false",
			"shares.0.in_rate":                  "0",
			"shares.0.index_id":                 CHECKSET,
			"shares.0.kms_rotate_period":        "0",
			"shares.0.lag_period":               "5",
			"shares.0.local_path":               CHECKSET,
			"shares.0.mns_health":               "MNSNotEnabled",
			"shares.0.nfs_v4_optimization":      "false",
			"shares.0.obsolete_buckets":         "",
			"shares.0.oss_bucket_name":          CHECKSET,
			"shares.0.oss_bucket_ssl":           "true",
			"shares.0.oss_endpoint":             CHECKSET,
			"shares.0.oss_health":               CHECKSET,
			"shares.0.oss_used":                 CHECKSET,
			"shares.0.out_rate":                 CHECKSET,
			"shares.0.partial_sync_paths":       "",
			"shares.0.path_prefix":              "",
			"shares.0.polling_interval":         "4500",
			"shares.0.protocol":                 "NFS",
			"shares.0.remaining_meta_space":     CHECKSET,
			"shares.0.remote_sync":              "true",
			"shares.0.remote_sync_download":     "false",
			"shares.0.ro_client_list":           "",
			"shares.0.ro_user_list":             "",
			"shares.0.rw_client_list":           "",
			"shares.0.rw_user_list":             "",
			"shares.0.server_side_cmk":          "",
			"shares.0.server_side_encryption":   "false",
			"shares.0.size":                     CHECKSET,
			"shares.0.squash":                   "none",
			"shares.0.state":                    "clean",
			"shares.0.support_archive":          "false",
			"shares.0.sync_progress":            "0",
			"shares.0.total_download":           "0",
			"shares.0.total_upload":             "0",
			"shares.0.transfer_acceleration":    "false",
			"shares.0.used":                     "0",
			"shares.0.windows_acl":              "false",
			"shares.0.bypass_cache_read":        "false",
		}
	}

	var fakeCloudStorageGatewayGatewayFileShareMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"shares.#": "0",
		}
	}

	var CloudStorageGatewayGatewayFileShareCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCloudStorageGatewayGatewayFileShareMapFunc,
		fakeMapFunc:  fakeCloudStorageGatewayGatewayFileShareMapFunc,
	}

	CloudStorageGatewayGatewayFileShareCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceCloudStorageGatewayGatewayFileSharesDependence(name string) string {
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
}`, name)
}
