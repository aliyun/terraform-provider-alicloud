---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_launch_templates"
sidebar_current: "docs-alicloud-datasource-ecs-launch-templates"
description: |-
  Provides a list of Ecs Launch Templates to the user.
---

# alicloud_ecs_launch_templates

This data source provides the Ecs Launch Templates of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.120.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
}

data "alicloud_images" "default" {
  name_regex = "^ubuntu_18.*64"
  owners     = "system"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = "terraform-example"
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name   = "terraform-example"
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_ecs_launch_template" "default" {
  launch_template_name          = "terraform-example"
  description                   = "terraform-example"
  image_id                      = data.alicloud_images.default.images.0.id
  host_name                     = "terraform-example"
  instance_charge_type          = "PrePaid"
  instance_name                 = "terraform-example"
  instance_type                 = data.alicloud_instance_types.default.instance_types.0.id
  internet_charge_type          = "PayByBandwidth"
  internet_max_bandwidth_in     = "5"
  internet_max_bandwidth_out    = "5"
  io_optimized                  = "optimized"
  key_pair_name                 = "key_pair_name"
  ram_role_name                 = "ram_role_name"
  network_type                  = "vpc"
  security_enhancement_strategy = "Active"
  spot_price_limit              = "5"
  spot_strategy                 = "SpotWithPriceLimit"
  security_group_ids            = [alicloud_security_group.default.id]
  system_disk {
    category             = "cloud_ssd"
    description          = "Test For Terraform"
    name                 = "terraform-example"
    size                 = "40"
    delete_with_instance = "false"
  }

  user_data  = "xxxxxxx"
  vswitch_id = alicloud_vswitch.default.id
  vpc_id     = alicloud_vpc.default.id
  zone_id    = data.alicloud_zones.default.zones.0.id

  template_tags = {
    Create = "Terraform"
    For    = "example"
  }

  network_interfaces {
    name              = "eth0"
    description       = "hello1"
    primary_ip        = "10.0.0.2"
    security_group_id = alicloud_security_group.default.id
    vswitch_id        = alicloud_vswitch.default.id
  }

  data_disks {
    name                 = "disk1"
    description          = "description"
    delete_with_instance = "true"
    category             = "cloud"
    encrypted            = "false"
    performance_level    = "PL0"
    size                 = "20"
  }
  data_disks {
    name                 = "disk2"
    description          = "description2"
    delete_with_instance = "true"
    category             = "cloud"
    encrypted            = "false"
    performance_level    = "PL0"
    size                 = "20"
  }
}

data "alicloud_ecs_launch_templates" "example" {
  ids = ["${alicloud_ecs_launch_template.default.id}"]
}

output "first_ecs_launch_template_id" {
  value = data.alicloud_ecs_launch_templates.example.templates.0.id
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Launch Template IDs.
* `launch_template_name` - (Optional, ForceNew) The Launch Template Name.
* `template_tags` - (Optional) The template tags.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Launch Template name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `template_resource_group_id` - (Optional, ForceNew) The template resource group id.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Launch Template names.
* `templates` - A list of Ecs Launch Templates. Each element contains the following attributes:
  * `auto_release_time` - (Optional) Instance auto release time.
  * `created_by` - CreatedBy.
  * `data_disks` - The list of data disks created with instance.
    * `encrypted` - Encrypted the data in this disk.
    * `name` - The name of the data disk.
    * `performance_level` - PerformanceLevel.
    * `size` - The performance level of the ESSD used as the data disk.
    * `snapshot_id` - The snapshot ID used to initialize the data disk. If the size specified by snapshot is greater that the size of the disk, use the size specified by snapshot as the size of the data disk.
    * `category` - The category of the disk.
    * `delete_with_instance` - Indicates whether the data disk is released with the instance.
    * `description` - The description of the data disk.
  * `default_version_number` - The Default Version Number.
  * `deployment_set_id` - The Deployment Set Id.
  * `description` - The Description of Template.
  * `enable_vm_os_config` - Whether to enable the instance operating system configuration.
  * `host_name` - Instance host name.
  * `id` - The ID of the Launch Template.
  * `image_id` - The Image Id.
  * `image_owner_alias` - Mirror source.
  * `instance_charge_type` - Internet bandwidth billing method.
  * `instance_name` - The Instance Name.
  * `instance_type` - Instance type.
  * `internet_charge_type` - Internet bandwidth billing method.
  * `internet_max_bandwidth_in` - The maximum inbound bandwidth from the Internet network, measured in Mbit/s.
  * `internet_max_bandwidth_out` - Maximum outbound bandwidth from the Internet, its unit of measurement is Mbit/s.
  * `io_optimized` - Whether it is an I/O-optimized instance or not.
  * `key_pair_name` - The name of the key pair.
  * `latest_version_number` - The Latest Version Number.
  * `launch_template_id` - The ID of the Launch Template.
  * `launch_template_name` - The Launch Template Name.
  * `modified_time` - The Modified Time.
  * `network_interfaces` - The list of network interfaces created with instance.
    * `description` - The ENI description.
    * `name` - The ENI name.
    * `primary_ip` - The primary private IP address of the ENI.
    * `security_group_id` - The security group ID must be one in the same VPC.
    * `vswitch_id` - The vSwitch ID for ENI. The instance must be in the same zone of the same VPC network as the ENI, but they may belong to different VSwitches.
  * `network_type` - Network type of the instance.
  * `password_inherit` - Whether to use the password preset by the mirror.
  * `period` - The subscription period of the instance.
  * `private_ip_address` - The private IP address of the instance.
  * `ram_role_name` - The RAM role name of the instance.
  * `resource_group_id` - The ID of the resource group to which to assign the instance, Elastic Block Storage (EBS) device, and ENI.
  * `security_enhancement_strategy` - Whether or not to activate the security enhancement feature and install network security software free of charge.
  * `security_group_id` - The security group ID.
  * `security_group_ids` - The security group IDs.
  * `spot_duration` - The protection period of the preemptible instance.
  * `spot_price_limit` - Sets the maximum hourly instance price.
  * `spot_strategy` - The spot strategy for a Pay-As-You-Go instance.
  * `system_disk` - The System Disk.
    * `category` - The category of the system disk.
    * `delete_with_instance` - Specifies whether to release the system disk when the instance is released.
    * `description` - System disk description.
    * `iops` - The Iops.
    * `name` - System disk name.
    * `performance_level` - The performance level of the ESSD used as the system disk.
    * `size` - Size of the system disk, measured in GB.
  * `template_tags` - The template tags.
  * `user_data` - The User Data.
  * `version_description` - The Version Description.
  * `vpc_id` - VpcId.
  * `vswitch_id` - The vswitch id.
  * `zone_id` - The Zone Id.
  * `http_endpoint` - Whether to enable access to instance metadata.
  * `http_tokens` - Whether to use the hardened mode (IMDSv2) when accessing instance metadata.
  * `http_put_response_hop_limit` - The HTTP PUT response hop limit required for instance metadata requests.