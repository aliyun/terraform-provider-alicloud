package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSAEApplicationDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1, 1000)
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
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
			"ids.#":                                           "1",
			"applications.#":                                  "1",
			"applications.0.app_name":                         fmt.Sprintf("tftestaccsaenames%d", rand),
			"applications.0.app_description":                  fmt.Sprintf("tftestaccsaenames%d", rand),
			"applications.0.id":                               CHECKSET,
			"applications.0.application_id":                   CHECKSET,
			"applications.0.command":                          "",
			"applications.0.command_args":                     "",
			"applications.0.config_map_mount_desc":            "[]",
			"applications.0.create_time":                      CHECKSET,
			"applications.0.region_id":                        CHECKSET,
			"applications.0.repo_name":                        CHECKSET,
			"applications.0.repo_namespace":                   CHECKSET,
			"applications.0.repo_origin_type":                 CHECKSET,
			"applications.0.envs":                             "[]",
			"applications.0.custom_host_alias":                "[]",
			"applications.0.jar_start_args":                   "",
			"applications.0.jar_start_options":                "",
			"applications.0.jdk":                              CHECKSET,
			"applications.0.liveness":                         "",
			"applications.0.min_ready_instances":              CHECKSET,
			"applications.0.mount_desc.#":                     CHECKSET,
			"applications.0.mount_host":                       "",
			"applications.0.nas_id":                           "",
			"applications.0.oss_ak_id":                        "",
			"applications.0.oss_ak_secret":                    "",
			"applications.0.oss_mount_descs":                  "[]",
			"applications.0.oss_mount_details.#":              CHECKSET,
			"applications.0.php_arms_config_location":         "",
			"applications.0.package_url":                      "",
			"applications.0.php_config":                       "",
			"applications.0.php_config_location":              "",
			"applications.0.post_start":                       "",
			"applications.0.pre_stop":                         "",
			"applications.0.readiness":                        "",
			"applications.0.security_group_id":                "",
			"applications.0.sls_configs":                      "",
			"applications.0.status":                           CHECKSET,
			"applications.0.termination_grace_period_seconds": CHECKSET,
			"applications.0.acr_assume_role_arn":              "",
			"applications.0.timezone":                         CHECKSET,
			"applications.0.tomcat_config":                    "",
			"applications.0.vpc_id":                           CHECKSET,
			"applications.0.war_start_options":                "",
			"applications.0.web_container":                    "",
			"applications.0.namespace_id":                     fmt.Sprintf("%s:tftestacc%d", os.Getenv("ALICLOUD_REGION"), rand),
			"applications.0.package_type":                     "Image",
			"applications.0.vswitch_id":                       CHECKSET,
			"applications.0.image_url":                        CHECKSET,
			"applications.0.replicas":                         "5",
			"applications.0.cpu":                              "500",
			"applications.0.memory":                           "2048",
			"applications.0.tags.%":                           "2",
			"applications.0.tags.Created":                     "tfTestAcc7",
			"applications.0.tags.For":                         "Tftestacc7",
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
	alicloudSaeNamespaceCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameConf, fieldConf, namespaceIdConf, statusConf, allConf)
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
  vpc_id          = data.alicloud_vpcs.default.ids.0
  timezone = "Asia/Shanghai"
  replicas=        "5"
  cpu=             "500"
  memory =          "2048"
  tags  = {
	Created = "tfTestAcc7"
	For =  "Tftestacc7"
  }
}
data "alicloud_sae_applications" "default" {
	%s
}
`, rand, defaultRegionToTest, rand, strings.Join(pairs, " \n "))
	return config
}
