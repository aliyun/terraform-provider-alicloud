package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSaeApplicationDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1, 1000)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSaeApplicationDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_sae_application.default.id}"]`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudSaeApplicationDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_sae_application.default.id}_fake"]`,
			"enable_details": "true",
		}),
	}
	nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSaeApplicationDataSourceName(rand, map[string]string{
			"app_name":       `"${alicloud_sae_application.default.app_name}"`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudSaeApplicationDataSourceName(rand, map[string]string{
			"app_name":       `"${alicloud_sae_application.default.app_name}_fake"`,
			"enable_details": "true",
		}),
	}

	fieldConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSaeApplicationDataSourceName(rand, map[string]string{
			"field_type":     `"appName"`,
			"field_value":    `"${alicloud_sae_application.default.app_name}"`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudSaeApplicationDataSourceName(rand, map[string]string{
			"field_type":     `"appName"`,
			"field_value":    `"${alicloud_sae_application.default.app_name}_fake"`,
			"enable_details": "true",
		}),
	}
	namespaceIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSaeApplicationDataSourceName(rand, map[string]string{
			"namespace_id":   `"${alicloud_sae_application.default.namespace_id}"`,
			"ids":            `["${alicloud_sae_application.default.id}"]`,
			"enable_details": "true",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSaeApplicationDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_sae_application.default.id}"]`,
			"status":         `"RUNNING"`,
			"enable_details": "true",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSaeApplicationDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_sae_application.default.id}"]`,
			"enable_details": "true",
			"app_name":       `"${alicloud_sae_application.default.app_name}"`,
			"field_type":     `"appName"`,
			"field_value":    `"${alicloud_sae_application.default.app_name}"`,
			"namespace_id":   `"${alicloud_sae_application.default.namespace_id}"`,
			"status":         `"RUNNING"`,
		}),
		fakeConfig: testAccCheckAlicloudSaeApplicationDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_sae_application.default.id}_fake"]`,
			"app_name":    `"${alicloud_sae_application.default.app_name}_fake"`,
			"field_type":  `"appName"`,
			"field_value": `"${alicloud_sae_application.default.app_name}_fake"`,
			"status":      `"UNKNOWN"`,
		}),
	}
	var existAlicloudSaeApplicationDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"applications.#":                 "1",
			"applications.0.app_name":        fmt.Sprintf("tftestaccsaenames%d", rand),
			"applications.0.app_description": fmt.Sprintf("tftestaccsaenames%d", rand),
			"applications.0.namespace_id":    fmt.Sprintf("%s:tftestacc%d", os.Getenv("ALICLOUD_REGION"), rand),
			"applications.0.package_type":    "Image",
			"applications.0.vswitch_id":      CHECKSET,
			"applications.0.image_url":       CHECKSET,
			"applications.0.replicas":        "5",
			"applications.0.cpu":             "500",
			"applications.0.memory":          "2048",
		}
	}
	var fakeAlicloudSaeApplicationDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"applications.#": "0",
		}
	}
	var alicloudSaeNamespaceCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_sae_applications.default",
		existMapFunc: existAlicloudSaeApplicationDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudSaeApplicationDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.SaeSupportRegions)
	}
	alicloudSaeNamespaceCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameConf, fieldConf, namespaceIdConf, statusConf, allConf)
}
func testAccCheckAlicloudSaeApplicationDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tftestaccsaenames%d"
}
data "alicloud_vpcs" "default"	{
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
}
data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
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

  timezone = "Asia/Shanghai"
  replicas=        "5"
  cpu=             "500"
  memory =          "2048"
}
data "alicloud_sae_applications" "default" {
	%s
}
`, rand, os.Getenv("ALICLOUD_REGION"), rand, strings.Join(pairs, " \n "))
	return config
}
