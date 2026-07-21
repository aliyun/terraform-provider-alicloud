// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudSsoUserProvisioningEventDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudSsoUserProvisioningEventSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_cloud_sso_directory.defaultQSrGmc.id}"]`,
			"directory_id": `"${alicloud_cloud_sso_directory.defaultQSrGmc.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudSsoUserProvisioningEventSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_cloud_sso_directory.defaultQSrGmc.id}_fake"]`,
			"directory_id": `"${alicloud_cloud_sso_directory.defaultQSrGmc.id}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudSsoUserProvisioningEventSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_cloud_sso_directory.defaultQSrGmc.id}"]`,
			"directory_id": `"${alicloud_cloud_sso_directory.defaultQSrGmc.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudSsoUserProvisioningEventSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_cloud_sso_directory.defaultQSrGmc.id}_fake"]`,
			"directory_id": `"${alicloud_cloud_sso_directory.defaultQSrGmc.id}"`,
		}),
	}

	CloudSsoUserProvisioningEventCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}

var existCloudSsoUserProvisioningEventMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"events.#": "0",
	}
}

var fakeCloudSsoUserProvisioningEventMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"events.#": "0",
	}
}

var CloudSsoUserProvisioningEventCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cloud_sso_user_provisioning_events.default",
	existMapFunc: existCloudSsoUserProvisioningEventMapFunc,
	fakeMapFunc:  fakeCloudSsoUserProvisioningEventMapFunc,
}

func testAccCheckAlicloudCloudSsoUserProvisioningEventSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCloudSsoUserProvisioningEvent%d"
}
resource "alicloud_cloud_sso_directory" "defaultQSrGmc" {
  directory_global_access_status = "Disabled"
  password_policy {
    min_password_length          = "8"
    min_password_different_chars = "8"
    max_password_age             = "90"
    password_reuse_prevention    = "1"
    max_login_attempts           = "5"
  }
  mfa_authentication_setting_info {
    mfa_authentication_advance_settings = "OnlyRiskyLogin"
    operation_for_risk_login            = "EnforceVerify"
  }
  directory_name = "tftest"
}


data "alicloud_cloud_sso_user_provisioning_events" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
