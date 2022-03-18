package alicloud

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"fmt"
)

func TestAccAlicloudECSInstancesDataSourceBasic(t *testing.T) {
	rand := acctest.RandInt()

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudInstancesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_instance.default.instance_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudInstancesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_instance.default.instance_name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudInstancesDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alicloud_instance.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlicloudInstancesDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alicloud_instance.default.id}_fake" ]`,
		}),
	}

	imageIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudInstancesDataSourceConfig(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccCheckAlicloudInstancesDataSource%d"`, rand),
			"image_id":   `"${data.alicloud_images.default.images.0.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudInstancesDataSourceConfig(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccCheckAlicloudInstancesDataSource%d"`, rand),
			"image_id":   `"${data.alicloud_images.default.images.0.id}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudInstancesDataSourceConfig(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccCheckAlicloudInstancesDataSource%d"`, rand),
			"status":     `"Running"`,
		}),
		fakeConfig: testAccCheckAlicloudInstancesDataSourceConfig(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccCheckAlicloudInstancesDataSource%d"`, rand),
			"status":     `"Stopped"`,
		}),
	}

	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudInstancesDataSourceConfig(rand, map[string]string{
			"vpc_id": `"${alicloud_vpc.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudInstancesDataSourceConfig(rand, map[string]string{
			"vpc_id": `"${alicloud_vpc.default.id}_fake"`,
		}),
	}

	vSwitchConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudInstancesDataSourceConfig(rand, map[string]string{
			"vswitch_id": `"${alicloud_vswitch.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudInstancesDataSourceConfig(rand, map[string]string{
			"vswitch_id": `"${alicloud_vswitch.default.id}_fake"`,
		}),
	}

	availabilityZoneConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudInstancesDataSourceConfig(rand, map[string]string{
			"name_regex":        fmt.Sprintf(`"tf-testAccCheckAlicloudInstancesDataSource%d"`, rand),
			"availability_zone": `"${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"`,
		}),
		fakeConfig: testAccCheckAlicloudInstancesDataSourceConfig(rand, map[string]string{
			"name_regex":        fmt.Sprintf(`"tf-testAccCheckAlicloudInstancesDataSource%d"`, rand),
			"availability_zone": `"${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}_fake"`,
		}),
	}

	ramRoleNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudInstancesDataSourceConfig(rand, map[string]string{
			"ram_role_name": `"${alicloud_instance.default.role_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudInstancesDataSourceConfig(rand, map[string]string{
			"ram_role_name": `"${alicloud_instance.default.role_name}_fake"`,
		}),
	}

	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudInstancesDataSourceConfig(rand, map[string]string{
			"name_regex":        fmt.Sprintf(`"tf-testAccCheckAlicloudInstancesDataSource%d"`, rand),
			"resource_group_id": `"${alicloud_instance.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudInstancesDataSourceConfig(rand, map[string]string{
			"name_regex":        fmt.Sprintf(`"tf-testAccCheckAlicloudInstancesDataSource%d"`, rand),
			"resource_group_id": `"${alicloud_instance.default.resource_group_id}_fake"`,
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudInstancesDataSourceConfigWithTag(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccCheckAlicloudInstancesDataSource%d"`, rand),
		},
			`tags = {
				from = "datasource"
				usage1 = "test"
				usage2 = "test"
				usage3 = "test"
				usage4 = "test"
				usage5 = "test"
				usage6 = "test"
			}`,
		),
		fakeConfig: testAccCheckAlicloudInstancesDataSourceConfigWithTag(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccCheckAlicloudInstancesDataSource%d"`, rand),
		},
			`tags = {
				from = "datasource_fake"
				usage1 = "test"
				usage2 = "test"
				usage3 = "test"
				usage4 = "test"
				usage5 = "test"
				usage6 = "test"
			}`,
		),
	}

	pagingConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudInstancesDataSourceConfig(rand, map[string]string{
			"ids":         `[ "${alicloud_instance.default.id}" ]`,
			"page_number": `1`,
		}),
		fakeConfig: testAccCheckAlicloudInstancesDataSourceConfig(rand, map[string]string{
			"ids":         `[ "${alicloud_instance.default.id}" ]`,
			"page_number": `2`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudInstancesDataSourceConfigWithTag(rand, map[string]string{
			"ids":               `[ "${alicloud_instance.default.id}" ]`,
			"name_regex":        fmt.Sprintf(`"tf-testAccCheckAlicloudInstancesDataSource%d"`, rand),
			"image_id":          `"${data.alicloud_images.default.images.0.id}"`,
			"status":            `"Running"`,
			"vpc_id":            `"${alicloud_vpc.default.id}"`,
			"vswitch_id":        `"${alicloud_vswitch.default.id}"`,
			"availability_zone": `"${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"`,
			"resource_group_id": `"${alicloud_instance.default.resource_group_id}"`,
			"page_number":       `1`,
		},
			`tags = {
				from = "datasource"
				usage1 = "test"
				usage2 = "test"
				usage3 = "test"
				usage4 = "test"
				usage5 = "test"
				usage6 = "test"
			}`,
		),
		fakeConfig: testAccCheckAlicloudInstancesDataSourceConfigWithTag(rand, map[string]string{
			"ids":               `[ "${alicloud_instance.default.id}_fake" ]`,
			"name_regex":        fmt.Sprintf(`"tf-testAccCheckAlicloudInstancesDataSource%d"`, rand),
			"image_id":          `"${data.alicloud_images.default.images.0.id}"`,
			"status":            `"Running"`,
			"vpc_id":            `"${alicloud_vpc.default.id}"`,
			"vswitch_id":        `"${alicloud_vswitch.default.id}"`,
			"availability_zone": `"${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"`,
			"resource_group_id": `"${alicloud_instance.default.resource_group_id}"`,
			"page_number":       `2`,
		},
			`tags = {
				from = "datasource_fake"
				usage1 = "test"
				usage2 = "test"
				usage3 = "test"
				usage4 = "test"
				usage5 = "test"
				usage6 = "test"
			}`,
		),
	}

	instancesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, imageIdConf, statusConf,
		vpcIdConf, vSwitchConf, availabilityZoneConf, ramRoleNameConf, resourceGroupIdConf, tagsConf, pagingConf, allConf)
}

func testAccCheckAlicloudInstancesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	%s

	data "alicloud_resource_manager_resource_groups" "default" {
	  name_regex = "default"
	}

	variable "name" {
		default = "tf-testAccCheckAlicloudInstancesDataSource%d"
	}

	resource "alicloud_instance" "default" {
		availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
		vswitch_id = "${alicloud_vswitch.default.id}"
		private_ip = "172.16.0.10"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		system_disk_category = "cloud_efficiency"
		security_groups = ["${alicloud_security_group.default.id}"]
		resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups[0].id
		role_name = "${alicloud_ram_role.default.name}"
		data_disks {
				name  = "${var.name}-disk1"
				size =        "20"
				category =  "cloud_efficiency"
				description = "disk1"
		}
		data_disks {
				name  = "${var.name}-disk2"
				size =        "20"
				category =  "cloud_efficiency"
				description = "disk2"
		}
        tags = {
			from = "datasource"
			usage1 = "test"
			usage2 = "test"
			usage3 = "test"
			usage4 = "test"
			usage5 = "test"
			usage6 = "test"

		}
	}
	
	resource "alicloud_ram_role" "default" {
	  name = "${var.name}"
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": "sts:AssumeRole",
			  "Effect": "Allow",
			  "Principal": {
				"Service": [
				  "ecs.aliyuncs.com"
				]
			  }
			}
		  ],
		  "Version": "1"
		}
	  EOF
	  description = "this is a test"
	  force = true
	}

	data "alicloud_instances" "default" {
		%s
	}`, EcsInstanceCommonNoZonesTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}

func testAccCheckAlicloudInstancesDataSourceConfigWithTag(rand int, attrMap map[string]string, tags string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	%s

	variable "name" {
		default = "tf-testAccCheckAlicloudInstancesDataSource%d"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	  name_regex = "default"
	}

	resource "alicloud_instance" "default" {
		availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
		vswitch_id = "${alicloud_vswitch.default.id}"
		private_ip = "172.16.0.10"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		system_disk_category = "cloud_efficiency"
		security_groups = ["${alicloud_security_group.default.id}"]
		resource_group_id = "${data.alicloud_resource_manager_resource_groups.default.groups[0].id}"
		data_disks {
				name  = "${var.name}-disk1"
				size =        "20"
				category =  "cloud_efficiency"
				description = "disk1"
		}
		data_disks {
				name  = "${var.name}-disk2"
				size =        "20"
				category =  "cloud_efficiency"
				description = "disk2"
		}
        tags = {
			from = "datasource"
			usage1 = "test"
			usage2 = "test"
			usage3 = "test"
			usage4 = "test"
			usage5 = "test"
			usage6 = "test"

		}
	}

	data "alicloud_instances" "default" {
		%s
		%s
	}`, EcsInstanceCommonNoZonesTestCase, rand, strings.Join(pairs, "\n  "), tags)
	return config
}

var existInstancesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                                  "1",
		"names.#":                                "1",
		"instances.#":                            "1",
		"total_count":                            CHECKSET,
		"instances.0.id":                         CHECKSET,
		"instances.0.region_id":                  CHECKSET,
		"instances.0.availability_zone":          CHECKSET,
		"instances.0.private_ip":                 "172.16.0.10",
		"instances.0.status":                     string(Running),
		"instances.0.name":                       fmt.Sprintf("tf-testAccCheckAlicloudInstancesDataSource%d", rand),
		"instances.0.instance_type":              CHECKSET,
		"instances.0.vpc_id":                     CHECKSET,
		"instances.0.vswitch_id":                 CHECKSET,
		"instances.0.image_id":                   CHECKSET,
		"instances.0.resource_group_id":          CHECKSET,
		"instances.0.public_ip":                  "",
		"instances.0.eip":                        "",
		"instances.0.description":                "",
		"instances.0.security_groups.#":          "1",
		"instances.0.key_name":                   "",
		"instances.0.creation_time":              CHECKSET,
		"instances.0.instance_charge_type":       string(PostPaid),
		"instances.0.internet_max_bandwidth_out": "0",
		"instances.0.spot_strategy":              string(NoSpot),
		"instances.0.disk_device_mappings.#":     "3",
	}
}

var fakeInstancesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":       "0",
		"names.#":     "0",
		"instances.#": "0",
	}
}

var instancesCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_instances.default",
	existMapFunc: existInstancesMapFunc,
	fakeMapFunc:  fakeInstancesMapFunc,
}
