package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSAEApplicationInstancesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1, 1000)
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	appIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSaeApplicationInstancesDataSourceName(rand, map[string]string{
			"application_id": `"${alicloud_sae_application.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudSaeApplicationInstancesDataSourceName(rand, map[string]string{
			"application_id": `"${alicloud_sae_application.default.id}"`,
			"ids":            `["fake"]`,
		}),
	}
	var existAlicloudSaeApplicationInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              "2",
			"instances.#":                        "2",
			"instances.0.instance_id":            CHECKSET,
			"instances.0.group_id":               CHECKSET,
			"instances.0.instance_container_ip":  CHECKSET,
			"instances.0.instance_health_status": CHECKSET,
			"instances.0.image_url":              CHECKSET,
			"instances.0.package_version":        CHECKSET,
		}
	}
	var fakeAlicloudSaeApplicationInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}
	var alicloudSaeApplicationInstancesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_sae_application_instances.default",
		existMapFunc: existAlicloudSaeApplicationInstancesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudSaeApplicationInstancesDataSourceNameMapFunc,
	}
	alicloudSaeApplicationInstancesCheckInfo.dataSourceTestCheck(t, rand, appIdConf)
}

func testAccCheckAlicloudSaeApplicationInstancesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tftestaccsaeinstance%d"
}
data "alicloud_vpcs" "default"	{
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
}
resource "alicloud_sae_namespace" "default" {
	namespace_description = var.name
	namespace_id = "%s:tftestacc%d"
	namespace_name = var.name
}
resource "alicloud_sae_application" "default" {
  app_description= var.name
  app_name=        var.name
  namespace_id=    alicloud_sae_namespace.default.namespace_id
  image_url=     "registry-vpc.cn-hangzhou.aliyuncs.com/lxepoo/apache-php5"
  package_type=    "Image"
  jdk=             "Open JDK 8"
  vswitch_id=      data.alicloud_vswitches.default.ids.0
  vpc_id          = data.alicloud_vpcs.default.ids.0
  timezone = "Asia/Shanghai"
  replicas=        "2"
  cpu=             "500"
  memory =          "2048"
}
data "alicloud_sae_application_instances" "default" {
	%s
}
`, rand, defaultRegionToTest, rand, strings.Join(pairs, " \n "))
	return config
}
