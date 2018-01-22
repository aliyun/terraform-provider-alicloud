---
layout: "alicloud"
page_title: "Alicloud: alicloud_instances"
sidebar_current: "docs-alicloud-datasource-instances"
description: |-
    Provides a list of ECS instances to the user.
---

# alicloud\_instances

The Instances data source list ECS instance resource accoring to its ID, name regex, image id, status and other fields.

## Example Usage

```
data "alicloud_instances" "instances" {
	name_regex = "web_server"
	status = "Running"
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of ECS instance ID.
* `name_regex` - (Optional) A regex string to apply to the instance list returned by Alicloud.
* `image_id` - (Optional) The image ID of some ECS instance used.
* `status` - (Optional) List specified status instances. Valid values: "Creating", "Starting", "Running", "Stopping" and "Stopped". Default to list all status.
* `vpc_id` - (Optional) List several instances in the specified VPC.
* `vswitch_id` - (Optional) List several instances in the specified VSwitch.
* `availability_zone` - (Optional) List several instances in the specified availability zone.
* `tags` - (Optional) A mapping of tags marked ECS instanes.
* `output_file` - (Optional) The name of file that can save instances data source after running `terraform plan`.

## Attributes Reference

* `instances` A list of instnaces. It contains several attributes to `Block Instances`.

### Block Instances

Attributes for instanes:

* `id` - ID of the instance.
* `region_id` - Region Id the instance belongs.
* `availability_zone` - Availability zone the instance belongs.
* `status` - Instance current status.
* `name` - Instance name.
* `description` - Instance description.
* `instance_type` - Instance type.
* `vpc_id` - VPC ID the instance belongs.
* `vswitch_id` - VSwitch ID the instance belongs.
* `image_id` - Image id the instance used.
* `private_ip` - Instance private IP address.
* `public_ip` - Instance public IP address.
* `eip` - EIP address the VPC instance used.
* `security_groups` - List security group ID the instance belongs.
* `key_name` - Key pair the instance used.
* `creation_time` - Instance creation time.
* `instance_charge_type` - Instance charge type.
* `internet_charge_type` - Instance network charge type.
* `internet_max_bandwidth_out` - Instance internet out max bandwidth
* `spot_strategy` - Spot strategy the instance used.
* `disk_device_mappings` - Description of the disk the instance attached.
  * `device` - Device information of the created disk: such as /dev/xvdb.
  * `size` - Size of the created disk.
  * `category` - Cloud disk category.
  * `type` - Cloud disk type. System disk or data disk.
* `tags` - A mapping of tags marked ECS instanes.