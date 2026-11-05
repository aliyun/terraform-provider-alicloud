package alicloud

import (
	"fmt"
	"os"
)

// Local compat copies of ECS instance test helpers referenced by
// TestAccAlicloudProviderEcs in provider_test.go. The originals moved to
// service/ecs (package ecs_test) with the instance test migration; root
// package alicloud cannot import ecs_test, so the subset TestAccAlicloudProviderEcs
// uses is kept here verbatim, behavior-preserving.

func resourceInstanceVpcConfigDependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_instance_types" "default" {
  		instance_type_family = "ecs.sn1ne"
	}

	data "alicloud_images" "default" {
  		name_regex    = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		most_recent   = true
  		owners        = "system"
  		instance_type = data.alicloud_instance_types.default.instance_types.0.id
	}

resource "alicloud_vpc" "default" {
    vpc_name = var.name
}

resource "alicloud_vswitch" "default" {
  vpc_id  = alicloud_vpc.default.id
  zone_id = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)
  vswitch_name   = var.name
}

resource "alicloud_security_group" "default" {
  count = "2"
  name   = "${var.name}"
  vpc_id = alicloud_vpc.default.id
}
resource "alicloud_security_group_rule" "default" {
	count = 2
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${element(alicloud_security_group.default.*.id,count.index)}"
  	cidr_ip = "172.16.0.0/24"
}

variable "name" {
	default = "%s"
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
		  force = "true"
}

resource "alicloud_key_pair" "default" {
	key_pair_name = "${var.name}"
}

resource "alicloud_kms_key" "key" {
  description            = var.name
  pending_window_in_days = "7"
  key_state              = "Enabled"
}

resource "alicloud_ecs_network_interface" "primary" {
  network_interface_name = "${var.name}primary"
  vswitch_id             = alicloud_vswitch.default.id
  security_group_ids     = [alicloud_security_group.default.0.id]
}

resource "alicloud_ecs_network_interface" "secondary" {
  network_interface_name = "${var.name}secondary"
  vswitch_id             = alicloud_vswitch.default.id
  security_group_ids     = [alicloud_security_group.default.1.id]
}
`, name)
}

func resourceECSInstanceVpcDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_images" "default" {
  		name_regex = "^ubuntu_18.*64"
  		owners     = "system"
	}

	data "alicloud_instance_types" "default" {
  		image_id          = data.alicloud_images.default.images.0.id
	}

	resource "alicloud_vpc" "default" {
  		cidr_block = "192.168.0.0/16"
  		vpc_name   = var.name
	}

	resource "alicloud_vswitch" "default" {
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 2)
  		zone_id      = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  		vswitch_name = var.name
	}

	resource "alicloud_security_group" "default" {
  		count  = 5
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_security_group" "update" {
  		count  = 2
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_vpc" "vpcUpdate" {
  		cidr_block = "10.0.0.0/8"
	}

	resource "alicloud_vswitch" "vswitchUpdate" {
  		cidr_block = "10.1.0.0/16"
  		vpc_id     = alicloud_vpc.vpcUpdate.id
  		zone_id    = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
	}

	resource "alicloud_security_group" "vpcUpdateSg" {
  		count  = 5
  		vpc_id = alicloud_vpc.vpcUpdate.id
	}

	resource "alicloud_security_group" "vpcUpdateSg2" {
  		count  = 2
  		vpc_id = alicloud_vpc.vpcUpdate.id
	}

	resource "alicloud_key_pair" "default" {
  		key_pair_name = var.name
	}

	resource "alicloud_ram_role" "default" {
  		name     = var.name
  		force    = "true"
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
	}
`, name)
}

func resourceInstancePrePaidConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_instance_types" "default" {
  		instance_type_family = "ecs.sn1ne"
	}

	data "alicloud_images" "default" {
  		name_regex    = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		most_recent   = true
  		owners        = "system"
  		instance_type = data.alicloud_instance_types.default.instance_types.0.id
	}

	resource "alicloud_vpc" "default" {
  		vpc_name = var.name
	}

	resource "alicloud_vswitch" "default" {
  		vpc_id       = alicloud_vpc.default.id
  		zone_id      = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  		cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)
  		vswitch_name = var.name
	}

	resource "alicloud_security_group" "default" {
  		count  = 2
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_ram_role" "default" {
  		name     = var.name
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
  		force    = "true"
	}

	resource "alicloud_key_pair" "default" {
  		key_pair_name = var.name
	}
`, name)
}

func resourceInstanceBasicConfigDependence(name string) string {
	return fmt.Sprintf(`

data "alicloud_instance_types" "default" {
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  owners      = "system"
}

variable "resource_group_id" {
		default = "%s"
}

resource "alicloud_security_group" "default" {
  count = "2"
  name   = "${var.name}"
}
resource "alicloud_security_group_rule" "default" {
	count = 2
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${element(alicloud_security_group.default.*.id,count.index)}"
  	cidr_ip = "172.16.0.0/24"
}

variable "name" {
	default = "%s"
}

resource "alicloud_key_pair" "default" {
	key_name = "${var.name}"
}

`, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"), name)
}

func testAccCheckSpotInstanceDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_images" "default" {
  		name_regex    = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		most_recent   = true
  		owners        = "system"
  		instance_type = data.alicloud_instance_types.default.instance_types.0.id
	}

	data "alicloud_instance_types" "default" {
  		spot_strategy        = "SpotWithPriceLimit"
  		cpu_core_count       = 2
  		memory_size          = 8
  		instance_type_family = "ecs.g6"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name = var.name
	}

	resource "alicloud_vswitch" "default" {
  		vpc_id       = alicloud_vpc.default.id
  		zone_id      = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  		cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)
  		vswitch_name = var.name
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}
	
	`, name)
}

var testAccInstanceCheckMap = map[string]string{
	"image_id":          CHECKSET,
	"instance_type":     CHECKSET,
	"security_groups.#": "1",

	"availability_zone":    CHECKSET,
	"system_disk_category": "cloud_efficiency",
	//"credit_specification":          "",
	"spot_strategy":    "NoSpot",
	"spot_price_limit": "0",
	// "security_enhancement_strategy": "Active",
	"vswitch_id": CHECKSET,
	"user_data":  "I_am_user_data",

	"description":   "",
	"host_name":     CHECKSET,
	"password":      "",
	"is_outdated":   NOSET,
	"volume_tags.%": "0",
	"tags.%":        NOSET,

	"private_ip":                 CHECKSET,
	"public_ip":                  "",
	"status":                     "Running",
	"internet_charge_type":       "PayByTraffic",
	"internet_max_bandwidth_out": "0",
	"instance_charge_type":       "PostPaid",
	// the attributes of below are suppressed  when the value of instance_charge_type is `PostPaid`
	"period":             NOSET,
	"force_delete":       NOSET,
	"include_data_disks": NOSET,
	"dry_run":            "false",
	"system_disk_id":     CHECKSET,
	"create_time":        CHECKSET,
	"start_time":         CHECKSET,
	"expired_time":       CHECKSET,
}
