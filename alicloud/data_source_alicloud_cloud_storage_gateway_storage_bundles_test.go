package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudStorageGatewayStorageBundlesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudStorageGatewayStorageBundlesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cloud_storage_gateway_storage_bundle.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCloudStorageGatewayStorageBundlesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cloud_storage_gateway_storage_bundle.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudStorageGatewayStorageBundlesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cloud_storage_gateway_storage_bundle.default.storage_bundle_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudStorageGatewayStorageBundlesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cloud_storage_gateway_storage_bundle.default.storage_bundle_name}_fake"`,
		}),
	}
	pagingConf := dataSourceTestAccConfig{
		fakeConfig: testAccCheckAlicloudCloudStorageGatewayStorageBundlesDataSourceName(rand, map[string]string{
			"page_number": `2`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudStorageGatewayStorageBundlesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cloud_storage_gateway_storage_bundle.default.id}"]`,
			"name_regex": `"${alicloud_cloud_storage_gateway_storage_bundle.default.storage_bundle_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudStorageGatewayStorageBundlesDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_cloud_storage_gateway_storage_bundle.default.id}_fake"]`,
			"name_regex":  `"${alicloud_cloud_storage_gateway_storage_bundle.default.storage_bundle_name}_fake"`,
			"page_number": `2`,
		}),
	}
	var existAlicloudCloudStorageGatewayStorageBundlesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"total_count":                   CHECKSET,
			"bundles.#":                     "1",
			"bundles.0.description":         "",
			"bundles.0.location":            CHECKSET,
			"bundles.0.id":                  CHECKSET,
			"bundles.0.storage_bundle_id":   CHECKSET,
			"bundles.0.storage_bundle_name": fmt.Sprintf("tf-testAccStorageBundle-%d", rand),
			"bundles.0.create_time":         CHECKSET,
		}
	}
	var fakeAlicloudCloudStorageGatewayStorageBundlesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"names.#":   "0",
			"bundles.#": "0",
		}
	}
	var alicloudCloudStorageGatewayStorageBundlesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cloud_storage_gateway_storage_bundles.default",
		existMapFunc: existAlicloudCloudStorageGatewayStorageBundlesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCloudStorageGatewayStorageBundlesDataSourceNameMapFunc,
	}
	alicloudCloudStorageGatewayStorageBundlesCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, pagingConf, allConf)
}
func testAccCheckAlicloudCloudStorageGatewayStorageBundlesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccStorageBundle-%d"
}

resource "alicloud_cloud_storage_gateway_storage_bundle" "default" {
	storage_bundle_name = "${var.name}"
}

data "alicloud_cloud_storage_gateway_storage_bundles" "default" {	
	backend_bucket_region_id = "%s"
	%s
}
`, rand, defaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
