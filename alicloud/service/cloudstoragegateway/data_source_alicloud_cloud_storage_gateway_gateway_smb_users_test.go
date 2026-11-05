package cloudstoragegateway_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudCloudStorageGatewaySmbUsersDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudStorageGatewaySmbUsersDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cloud_storage_gateway_gateway_smb_user.default.id}"]`,
			"gateway_id": `"${alicloud_cloud_storage_gateway_gateway.default.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudCloudStorageGatewaySmbUsersDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cloud_storage_gateway_gateway_smb_user.default.id}_fake"]`,
			"gateway_id": `"${alicloud_cloud_storage_gateway_gateway.default.id}"`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudStorageGatewaySmbUsersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cloud_storage_gateway_gateway_smb_user.default.username}"`,
			"gateway_id": `"${alicloud_cloud_storage_gateway_gateway.default.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudCloudStorageGatewaySmbUsersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cloud_storage_gateway_gateway_smb_user.default.username}_fake"`,
			"gateway_id": `"${alicloud_cloud_storage_gateway_gateway.default.id}"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudStorageGatewaySmbUsersDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cloud_storage_gateway_gateway_smb_user.default.id}"]`,
			"gateway_id": `"${alicloud_cloud_storage_gateway_gateway.default.id}"`,
			"name_regex": `"${alicloud_cloud_storage_gateway_gateway_smb_user.default.username}"`,
		}),
		FakeConfig: testAccCheckAlicloudCloudStorageGatewaySmbUsersDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cloud_storage_gateway_gateway_smb_user.default.id}_fake"]`,
			"gateway_id": `"${alicloud_cloud_storage_gateway_gateway.default.id}"`,
			"name_regex": `"${alicloud_cloud_storage_gateway_gateway_smb_user.default.username}_fake"`,
		}),
	}

	var existAlicloudCloudStorageGatewaySmbUsersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":              "1",
			"users.#":            "1",
			"users.0.id":         acctest.CHECKSET,
			"users.0.username":   fmt.Sprintf("tf-testacccsguser%d", rand),
			"users.0.gateway_id": acctest.CHECKSET,
		}
	}
	var fakeAlicloudCloudStorageGatewaySmbUsersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"users.#": "0",
		}
	}
	var alicloudCsgSmbUsersCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cloud_storage_gateway_gateway_smb_users.default",
		ExistMapFunc: existAlicloudCloudStorageGatewaySmbUsersDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudCloudStorageGatewaySmbUsersDataSourceNameMapFunc,
	}

	alicloudCsgSmbUsersCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudCloudStorageGatewaySmbUsersDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testacccsguser%d"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_cloud_storage_gateway_storage_bundle" "example" {
  storage_bundle_name = var.name
}
resource "alicloud_cloud_storage_gateway_gateway" "default" {
  description              = "tf-acctestDesalone"
  gateway_class            = "Standard"
  type                     = "File"
  payment_type             = "PayAsYouGo"
  vswitch_id               = data.alicloud_vswitches.default.ids.0
  release_after_expiration = false
  public_network_bandwidth = 40
  storage_bundle_id        = alicloud_cloud_storage_gateway_storage_bundle.example.id
  location                 = "Cloud"
  gateway_name             = var.name
}

resource "alicloud_cloud_storage_gateway_gateway_smb_user" "default" {
	username = var.name
    password = "%d"
	gateway_id = alicloud_cloud_storage_gateway_gateway.default.id
}

data "alicloud_cloud_storage_gateway_gateway_smb_users" "default" {	
	%s
}
`, rand, rand, strings.Join(pairs, " \n "))
	return config
}
