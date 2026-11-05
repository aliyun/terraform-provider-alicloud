package cloudstoragegateway_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudCloudStorageGatewayStorageBundlesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudStorageGatewayStorageBundlesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cloud_storage_gateway_storage_bundle.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudCloudStorageGatewayStorageBundlesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cloud_storage_gateway_storage_bundle.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudStorageGatewayStorageBundlesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cloud_storage_gateway_storage_bundle.default.storage_bundle_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudCloudStorageGatewayStorageBundlesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cloud_storage_gateway_storage_bundle.default.storage_bundle_name}_fake"`,
		}),
	}
	pagingConf := acctest.DataSourceTestAccConfig{
		FakeConfig: testAccCheckAlicloudCloudStorageGatewayStorageBundlesDataSourceName(rand, map[string]string{
			"page_number": `2`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudStorageGatewayStorageBundlesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cloud_storage_gateway_storage_bundle.default.id}"]`,
			"name_regex": `"${alicloud_cloud_storage_gateway_storage_bundle.default.storage_bundle_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudCloudStorageGatewayStorageBundlesDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_cloud_storage_gateway_storage_bundle.default.id}_fake"]`,
			"name_regex":  `"${alicloud_cloud_storage_gateway_storage_bundle.default.storage_bundle_name}_fake"`,
			"page_number": `2`,
		}),
	}
	var existAlicloudCloudStorageGatewayStorageBundlesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"total_count":                   acctest.CHECKSET,
			"bundles.#":                     "1",
			"bundles.0.description":         "",
			"bundles.0.location":            acctest.CHECKSET,
			"bundles.0.id":                  acctest.CHECKSET,
			"bundles.0.storage_bundle_id":   acctest.CHECKSET,
			"bundles.0.storage_bundle_name": fmt.Sprintf("tf-testAccStorageBundle-%d", rand),
			"bundles.0.create_time":         acctest.CHECKSET,
		}
	}
	var fakeAlicloudCloudStorageGatewayStorageBundlesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"names.#":   "0",
			"bundles.#": "0",
		}
	}
	var alicloudCloudStorageGatewayStorageBundlesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cloud_storage_gateway_storage_bundles.default",
		ExistMapFunc: existAlicloudCloudStorageGatewayStorageBundlesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudCloudStorageGatewayStorageBundlesDataSourceNameMapFunc,
	}
	alicloudCloudStorageGatewayStorageBundlesCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, pagingConf, allConf)
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
`, rand, acctest.DefaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
