package sae_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudSAEApplicationDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1, 1000)
	acctest.CheckoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSaeApplicationDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_sae_application.default.id}"]`,
			"enable_details": "true",
		}),
		FakeConfig: testAccCheckAlicloudSaeApplicationDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_sae_application.default.id}_fake"]`,
			"enable_details": "true",
		}),
	}
	nameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSaeApplicationDataSourceName(rand, map[string]string{
			"app_name":       `"${alicloud_sae_application.default.app_name}"`,
			"enable_details": "true",
		}),
		FakeConfig: testAccCheckAlicloudSaeApplicationDataSourceName(rand, map[string]string{
			"app_name":       `"${alicloud_sae_application.default.app_name}_fake"`,
			"enable_details": "true",
		}),
	}

	fieldConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSaeApplicationDataSourceName(rand, map[string]string{
			"field_type":     `"appName"`,
			"field_value":    `"${alicloud_sae_application.default.app_name}"`,
			"enable_details": "true",
		}),
		FakeConfig: testAccCheckAlicloudSaeApplicationDataSourceName(rand, map[string]string{
			"field_type":     `"appName"`,
			"field_value":    `"${alicloud_sae_application.default.app_name}_fake"`,
			"enable_details": "true",
		}),
	}
	namespaceIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSaeApplicationDataSourceName(rand, map[string]string{
			"namespace_id":   `"${alicloud_sae_application.default.namespace_id}"`,
			"ids":            `["${alicloud_sae_application.default.id}"]`,
			"enable_details": "true",
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSaeApplicationDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_sae_application.default.id}"]`,
			"status":         `"RUNNING"`,
			"enable_details": "true",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSaeApplicationDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_sae_application.default.id}"]`,
			"enable_details": "true",
			"app_name":       `"${alicloud_sae_application.default.app_name}"`,
			"field_type":     `"appName"`,
			"field_value":    `"${alicloud_sae_application.default.app_name}"`,
			"namespace_id":   `"${alicloud_sae_application.default.namespace_id}"`,
			"status":         `"RUNNING"`,
		}),
		FakeConfig: testAccCheckAlicloudSaeApplicationDataSourceName(rand, map[string]string{
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
			"applications.0.id":                               acctest.CHECKSET,
			"applications.0.application_id":                   acctest.CHECKSET,
			"applications.0.command":                          "",
			"applications.0.command_args":                     "",
			"applications.0.config_map_mount_desc":            "[]",
			"applications.0.create_time":                      acctest.CHECKSET,
			"applications.0.region_id":                        acctest.CHECKSET,
			"applications.0.repo_name":                        acctest.CHECKSET,
			"applications.0.repo_namespace":                   acctest.CHECKSET,
			"applications.0.repo_origin_type":                 acctest.CHECKSET,
			"applications.0.envs":                             "[]",
			"applications.0.custom_host_alias":                "[]",
			"applications.0.jar_start_args":                   "",
			"applications.0.jar_start_options":                "",
			"applications.0.jdk":                              acctest.CHECKSET,
			"applications.0.liveness":                         "",
			"applications.0.min_ready_instances":              acctest.CHECKSET,
			"applications.0.mount_desc.#":                     acctest.CHECKSET,
			"applications.0.mount_host":                       "",
			"applications.0.nas_id":                           "",
			"applications.0.oss_ak_id":                        "",
			"applications.0.oss_ak_secret":                    "",
			"applications.0.oss_mount_descs":                  "[]",
			"applications.0.oss_mount_details.#":              acctest.CHECKSET,
			"applications.0.php_arms_config_location":         "",
			"applications.0.package_url":                      "",
			"applications.0.php_config":                       "",
			"applications.0.php_config_location":              "",
			"applications.0.post_start":                       "",
			"applications.0.pre_stop":                         "",
			"applications.0.readiness":                        "",
			"applications.0.security_group_id":                acctest.CHECKSET,
			"applications.0.sls_configs":                      "",
			"applications.0.status":                           acctest.CHECKSET,
			"applications.0.termination_grace_period_seconds": acctest.CHECKSET,
			"applications.0.acr_assume_role_arn":              "",
			"applications.0.timezone":                         acctest.CHECKSET,
			"applications.0.tomcat_config":                    "",
			"applications.0.vpc_id":                           acctest.CHECKSET,
			"applications.0.war_start_options":                "",
			"applications.0.web_container":                    "",
			"applications.0.namespace_id":                     fmt.Sprintf("%s:tftestacc%d", os.Getenv("ALICLOUD_REGION"), rand),
			"applications.0.package_type":                     "Image",
			"applications.0.vswitch_id":                       acctest.CHECKSET,
			"applications.0.image_url":                        acctest.CHECKSET,
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
	var alicloudSaeNamespaceCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_sae_applications.default",
		ExistMapFunc: existAlicloudSaeApplicationDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudSaeApplicationDataSourceNameMapFunc,
	}
	alicloudSaeNamespaceCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameConf, fieldConf, namespaceIdConf, statusConf, allConf)
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
`, rand, acctest.DefaultRegionToTest, rand, strings.Join(pairs, " \n "))
	return config
}
