package alicloud

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"fmt"
)

func TestAccAlicloudInstancesDataSourceBasic(t *testing.T) {
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
			"availability_zone": `"${data.alicloud_zones.default.zones.0.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudInstancesDataSourceConfig(rand, map[string]string{
			"name_regex":        fmt.Sprintf(`"tf-testAccCheckAlicloudInstancesDataSource%d"`, rand),
			"availability_zone": `"${data.alicloud_zones.default.zones.0.id}_fake"`,
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudInstancesDataSourceConfigWithTag(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccCheckAlicloudInstancesDataSource%d"`, rand),
		},
			`tags {
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
			`tags {
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

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudInstancesDataSourceConfigWithTag(rand, map[string]string{
			"ids":               `[ "${alicloud_instance.default.id}" ]`,
			"name_regex":        fmt.Sprintf(`"tf-testAccCheckAlicloudInstancesDataSource%d"`, rand),
			"image_id":          `"${data.alicloud_images.default.images.0.id}"`,
			"status":            `"Running"`,
			"vpc_id":            `"${alicloud_vpc.default.id}"`,
			"vswitch_id":        `"${alicloud_vswitch.default.id}"`,
			"availability_zone": `"${data.alicloud_zones.default.zones.0.id}"`,
		},
			`tags {
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
			"availability_zone": `"${data.alicloud_zones.default.zones.0.id}"`,
		},
			`tags {
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
		vpcIdConf, vSwitchConf, availabilityZoneConf, tagsConf, allConf)
}

func testAccCheckAlicloudInstancesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckAlicloudInstancesDataSource%d"
	}

	resource "alicloud_instance" "default" {
		vswitch_id = "${alicloud_vswitch.default.id}"
		private_ip = "172.16.0.10"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		system_disk_category = "cloud_efficiency"
		security_groups = ["${alicloud_security_group.default.id}"]
        tags {
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
	}`, EcsInstanceCommonTestCase, rand, strings.Join(pairs, "\n  "))
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

	resource "alicloud_instance" "default" {
		vswitch_id = "${alicloud_vswitch.default.id}"
		private_ip = "172.16.0.10"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		system_disk_category = "cloud_efficiency"
		security_groups = ["${alicloud_security_group.default.id}"]
        tags {
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
	}`, EcsInstanceCommonTestCase, rand, strings.Join(pairs, "\n  "), tags)
	return config
}

var existInstancesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                                  "1",
		"names.#":                                "1",
		"instances.#":                            "1",
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
		"instances.0.public_ip":                  "",
		"instances.0.eip":                        "",
		"instances.0.description":                "",
		"instances.0.security_groups.#":          "1",
		"instances.0.key_name":                   "",
		"instances.0.creation_time":              CHECKSET,
		"instances.0.instance_charge_type":       string(PostPaid),
		"instances.0.internet_max_bandwidth_out": "0",
		"instances.0.spot_strategy":              string(NoSpot),
		"instances.0.disk_device_mappings.#":     "1",
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
