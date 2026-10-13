package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlicloudBpStudioApplicationsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.BpStudioApplicationSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudBpStudioApplicationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_bp_studio_application.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudBpStudioApplicationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_bp_studio_application.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudBpStudioApplicationsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_bp_studio_application.default.application_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudBpStudioApplicationsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_bp_studio_application.default.application_name}_fake"`,
		}),
	}
	keywordConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudBpStudioApplicationsDataSourceName(rand, map[string]string{
			"keyword": `"${alicloud_bp_studio_application.default.application_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudBpStudioApplicationsDataSourceName(rand, map[string]string{
			"keyword": `"${alicloud_bp_studio_application.default.application_name}_fake"`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudBpStudioApplicationsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_bp_studio_application.default.id}"]`,
			"resource_group_id": `"${alicloud_bp_studio_application.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudBpStudioApplicationsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_bp_studio_application.default.id}"]`,
			"resource_group_id": `"${alicloud_bp_studio_application.default.resource_group_id}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudBpStudioApplicationsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_bp_studio_application.default.id}"]`,
			"status": `"success"`,
		}),
		fakeConfig: testAccCheckAlicloudBpStudioApplicationsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_bp_studio_application.default.id}"]`,
			"status": `"release"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudBpStudioApplicationsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_bp_studio_application.default.id}"]`,
			"name_regex":        `"${alicloud_bp_studio_application.default.application_name}"`,
			"keyword":           `"${alicloud_bp_studio_application.default.application_name}"`,
			"resource_group_id": `"${alicloud_bp_studio_application.default.resource_group_id}"`,
			"status":            `"success"`,
		}),
		fakeConfig: testAccCheckAlicloudBpStudioApplicationsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_bp_studio_application.default.id}_fake"]`,
			"name_regex":        `"${alicloud_bp_studio_application.default.application_name}_fake"`,
			"keyword":           `"${alicloud_bp_studio_application.default.application_name}_fake"`,
			"resource_group_id": `"${alicloud_bp_studio_application.default.resource_group_id}_fake"`,
			"status":            `"release"`,
		}),
	}
	var existAlicloudBpStudioApplicationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"names.#":                          "1",
			"applications.#":                   "1",
			"applications.0.id":                CHECKSET,
			"applications.0.application_id":    CHECKSET,
			"applications.0.application_name":  CHECKSET,
			"applications.0.resource_group_id": CHECKSET,
			"applications.0.image_url":         CHECKSET,
			"applications.0.create_time":       CHECKSET,
			"applications.0.status":            CHECKSET,
		}
	}
	var fakeAlicloudBpStudioApplicationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"names.#":        "0",
			"applications.#": "0",
		}
	}
	var alicloudBpStudioApplicationsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_bp_studio_applications.default",
		existMapFunc: existAlicloudBpStudioApplicationsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudBpStudioApplicationsDataSourceNameMapFunc,
	}
	alicloudBpStudioApplicationsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, keywordConf, resourceGroupIdConf, statusConf, allConf)
}

func testAccCheckAlicloudBpStudioApplicationsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccBpStudioApplication-%d"
	}

	variable "area_id" {
		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
		name_regex = "tf"
	}

	data "alicloud_zones" "default" {
		available_disk_category     = "cloud_efficiency"
		available_resource_creation = "VSwitch"
	}

	data "alicloud_images" "default" {
		name_regex    = "^ubuntu_[0-9]+_[0-9]+_x64*"
		most_recent   = true
		owners        = "system"
		instance_type = data.alicloud_instance_types.default.instance_types.0.id
	}

	data "alicloud_instance_types" "default" {
		availability_zone    = data.alicloud_zones.default.zones.0.id
		instance_type_family = "ecs.sn1ne"
	}

	data "alicloud_vpcs" "default" {
		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
		vpc_id  = data.alicloud_vpcs.default.ids.0
		zone_id = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
		name   = var.name
		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_instance" "default" {
		image_id             = data.alicloud_images.default.images.0.id
		instance_type        = data.alicloud_instance_types.default.instance_types.0.id
		instance_name        = var.name
		security_groups      = alicloud_security_group.default.*.id
		availability_zone    = data.alicloud_zones.default.zones.0.id
		instance_charge_type = "PostPaid"
		system_disk_category = "cloud_efficiency"
		vswitch_id           = data.alicloud_vswitches.default.ids.0
	}

	resource "alicloud_bp_studio_application" "default" {
		application_name  = var.name
		template_id       = "YAUUQIYRSV1CMFGX"
		resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
		area_id           = var.area_id
		instances {
			id        = alicloud_instance.default.id
			node_name = alicloud_instance.default.instance_name
			node_type = "ecs"
  		}
  		configuration = {
			enableMonitor = "1"
  		}
  		variables = {
			test = "1"
		}
	}
	
	data "alicloud_bp_studio_applications" "default" {
		%s
	}
`, rand, defaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
