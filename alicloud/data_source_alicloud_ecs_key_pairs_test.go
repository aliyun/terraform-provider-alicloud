package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudEcsKeyPairsDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_ecs_key_pairs.default"
	name := fmt.Sprintf("tf-testAcc-EcsKeyPair%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEcsKeyPairsConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecs_key_pair_attachment.default.key_pair_name}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecs_key_pair.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ecs_key_pair_attachment.default.key_pair_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ecs_key_pair.default.key_pair_name}_fake",
		}),
	}

	fingerPrintConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"finger_print": "${alicloud_ecs_key_pair.default.finger_print}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"finger_print": "${alicloud_ecs_key_pair.default.finger_print}_fake",
		}),
	}

	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"resource_group_id": "${alicloud_ecs_key_pair.default.resource_group_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "KeyPair",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "KeyPair_Fake",
			},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_ecs_key_pair.default.id}"},
			"name_regex":        "${alicloud_ecs_key_pair_attachment.default.key_pair_name}",
			"finger_print":      "${alicloud_ecs_key_pair.default.finger_print}",
			"resource_group_id": "${alicloud_ecs_key_pair.default.resource_group_id}",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "KeyPair",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_ecs_key_pair.default.id}_fake"},
			"name_regex":        "${alicloud_ecs_key_pair.default.key_pair_name}_fake",
			"finger_print":      "${alicloud_ecs_key_pair.default.finger_print}_fake",
			"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "KeyPair_Fake",
			},
		}),
	}

	var existAliCloudEcsKeyPairsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                     "1",
			"names.#":                                   "1",
			"pairs.#":                                   "1",
			"pairs.0.id":                                CHECKSET,
			"pairs.0.key_pair_name":                     CHECKSET,
			"pairs.0.key_name":                          CHECKSET,
			"pairs.0.finger_print":                      CHECKSET,
			"pairs.0.resource_group_id":                 CHECKSET,
			"pairs.0.tags.%":                            "2",
			"pairs.0.tags.Created":                      "TF",
			"pairs.0.tags.For":                          "KeyPair",
			"pairs.0.instances.#":                       "1",
			"pairs.0.instances.0.instance_id":           CHECKSET,
			"pairs.0.instances.0.instance_name":         CHECKSET,
			"pairs.0.instances.0.description":           CHECKSET,
			"pairs.0.instances.0.image_id":              CHECKSET,
			"pairs.0.instances.0.region_id":             CHECKSET,
			"pairs.0.instances.0.availability_zone":     CHECKSET,
			"pairs.0.instances.0.instance_type":         CHECKSET,
			"pairs.0.instances.0.vswitch_id":            CHECKSET,
			"pairs.0.instances.0.public_ip":             CHECKSET,
			"pairs.0.instances.0.private_ip":            CHECKSET,
			"pairs.0.instances.0.key_name":              CHECKSET,
			"pairs.0.instances.0.status":                CHECKSET,
			"key_pairs.#":                               "1",
			"key_pairs.0.id":                            CHECKSET,
			"key_pairs.0.key_pair_name":                 CHECKSET,
			"key_pairs.0.key_name":                      CHECKSET,
			"key_pairs.0.finger_print":                  CHECKSET,
			"key_pairs.0.resource_group_id":             CHECKSET,
			"key_pairs.0.tags.%":                        "2",
			"key_pairs.0.tags.Created":                  "TF",
			"key_pairs.0.tags.For":                      "KeyPair",
			"key_pairs.0.instances.#":                   "1",
			"key_pairs.0.instances.0.instance_id":       CHECKSET,
			"key_pairs.0.instances.0.instance_name":     CHECKSET,
			"key_pairs.0.instances.0.description":       CHECKSET,
			"key_pairs.0.instances.0.image_id":          CHECKSET,
			"key_pairs.0.instances.0.region_id":         CHECKSET,
			"key_pairs.0.instances.0.availability_zone": CHECKSET,
			"key_pairs.0.instances.0.instance_type":     CHECKSET,
			"key_pairs.0.instances.0.vswitch_id":        CHECKSET,
			"key_pairs.0.instances.0.public_ip":         CHECKSET,
			"key_pairs.0.instances.0.private_ip":        CHECKSET,
			"key_pairs.0.instances.0.key_name":          CHECKSET,
			"key_pairs.0.instances.0.status":            CHECKSET,
		}
	}

	var fakeAliCloudEcsKeyPairsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"pairs.#":     "0",
			"key_pairs.#": "0",
		}
	}

	var aliCloudEcsKeyPairsInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecs_key_pairs.default",
		existMapFunc: existAliCloudEcsKeyPairsMapFunc,
		fakeMapFunc:  fakeAliCloudEcsKeyPairsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudEcsKeyPairsInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, fingerPrintConf, resourceGroupIdConf, tagsConf, allConf)
}

func dataSourceEcsKeyPairsConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}
	
	data "alicloud_zones" "default" {
  		available_disk_category     = "cloud_efficiency"
  		available_resource_creation = "VSwitch"
	}
	
	data "alicloud_images" "default" {
		name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		most_recent = true
  		owners      = "system"
	}
	
	data "alicloud_instance_types" "default" {
  		availability_zone    = data.alicloud_zones.default.zones.0.id
  		image_id             = data.alicloud_images.default.images.0.id
  		system_disk_category = "cloud_efficiency"
	}
	
	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}
	
	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}
	
	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}
	
	resource "alicloud_instance" "default" {
  		image_id                   = data.alicloud_images.default.images.0.id
  		instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  		security_groups            = alicloud_security_group.default.*.id
  		internet_charge_type       = "PayByTraffic"
  		internet_max_bandwidth_out = "10"
  		availability_zone          = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  		instance_charge_type       = "PostPaid"
  		system_disk_category       = "cloud_efficiency"
  		vswitch_id                 = alicloud_vswitch.default.id
  		instance_name              = var.name
		description                = var.name
	}
	
	resource "alicloud_ecs_key_pair" "default" {
  		key_pair_name     = var.name
  		public_key        = "ssh-rsa AAAAB3Nza12345678qwertyuudsfsg"
  		resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.1
  		tags = {
    		Created = "TF"
    		For     = "KeyPair",
  		}
	}
	
	resource "alicloud_ecs_key_pair_attachment" "default" {
  		key_pair_name = alicloud_ecs_key_pair.default.key_pair_name
  		instance_ids  = [alicloud_instance.default.id]
	}
`, name)
}
