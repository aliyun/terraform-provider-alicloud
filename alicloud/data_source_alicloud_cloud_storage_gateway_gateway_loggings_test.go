package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudStorageGatewayGatewayLoggingsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudStorageGatewayGatewayLoggingsDataSourceName(rand, map[string]string{
			"gateway_id": `"${data.alicloud_cloud_storage_gateway_gateways.default.ids.0}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudStorageGatewayGatewayLoggingsDataSourceName(rand, map[string]string{
			"gateway_id": "",
		}),
	}
	var existAlicloudCloudStorageGatewayGatewayLoggingsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"loggings.#":              "1",
			"loggings.0.id":           CHECKSET,
			"loggings.0.status":       CHECKSET,
			"loggings.0.gateway_id":   CHECKSET,
			"loggings.0.sls_logstore": CHECKSET,
			"loggings.0.sls_project":  CHECKSET,
		}
	}
	var fakeAlicloudCloudStorageGatewayGatewayLoggingsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"loggings.#": "0",
		}
	}
	var alicloudCloudStorageGatewayGatewayLoggingsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cloud_storage_gateway_gateway_loggings.default",
		existMapFunc: existAlicloudCloudStorageGatewayGatewayLoggingsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCloudStorageGatewayGatewayLoggingsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudCloudStorageGatewayGatewayLoggingsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, statusConf)
}
func testAccCheckAlicloudCloudStorageGatewayGatewayLoggingsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testacccloudstoragegatewaygatewaylogging%d"
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

resource "alicloud_cloud_storage_gateway_gateway_logging" "default" {
  gateway_id   = data.alicloud_cloud_storage_gateway_gateways.default.ids.0
  sls_logstore = alicloud_log_store.default.name
  sls_project  = alicloud_log_project.default.name
  status       = "Enabled"
}

data "alicloud_cloud_storage_gateway_gateway_loggings" "default" {
  %s
}
`, rand, defaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
