package bpstudio_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudBpStudioApplicationsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.PreCheckWithRegions(t, true, connectivity.BpStudioApplicationSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudBpStudioApplicationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_bp_studio_application.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudBpStudioApplicationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_bp_studio_application.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudBpStudioApplicationsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_bp_studio_application.default.application_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudBpStudioApplicationsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_bp_studio_application.default.application_name}_fake"`,
		}),
	}
	keywordConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudBpStudioApplicationsDataSourceName(rand, map[string]string{
			"keyword": `"${alicloud_bp_studio_application.default.application_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudBpStudioApplicationsDataSourceName(rand, map[string]string{
			"keyword": `"${alicloud_bp_studio_application.default.application_name}_fake"`,
		}),
	}
	resourceGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudBpStudioApplicationsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_bp_studio_application.default.id}"]`,
			"resource_group_id": `"${alicloud_bp_studio_application.default.resource_group_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudBpStudioApplicationsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_bp_studio_application.default.id}"]`,
			"resource_group_id": `"${alicloud_bp_studio_application.default.resource_group_id}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudBpStudioApplicationsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_bp_studio_application.default.id}"]`,
			"status": `"success"`,
		}),
		FakeConfig: testAccCheckAlicloudBpStudioApplicationsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_bp_studio_application.default.id}"]`,
			"status": `"release"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudBpStudioApplicationsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_bp_studio_application.default.id}"]`,
			"name_regex":        `"${alicloud_bp_studio_application.default.application_name}"`,
			"keyword":           `"${alicloud_bp_studio_application.default.application_name}"`,
			"resource_group_id": `"${alicloud_bp_studio_application.default.resource_group_id}"`,
			"status":            `"success"`,
		}),
		FakeConfig: testAccCheckAlicloudBpStudioApplicationsDataSourceName(rand, map[string]string{
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
			"applications.0.id":                acctest.CHECKSET,
			"applications.0.application_id":    acctest.CHECKSET,
			"applications.0.application_name":  acctest.CHECKSET,
			"applications.0.resource_group_id": acctest.CHECKSET,
			"applications.0.image_url":         acctest.CHECKSET,
			"applications.0.create_time":       acctest.CHECKSET,
			"applications.0.status":            acctest.CHECKSET,
		}
	}
	var fakeAlicloudBpStudioApplicationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"names.#":        "0",
			"applications.#": "0",
		}
	}
	var alicloudBpStudioApplicationsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_bp_studio_applications.default",
		ExistMapFunc: existAlicloudBpStudioApplicationsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudBpStudioApplicationsDataSourceNameMapFunc,
	}
	alicloudBpStudioApplicationsCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, keywordConf, resourceGroupIdConf, statusConf, allConf)
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
`, rand, acctest.DefaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
