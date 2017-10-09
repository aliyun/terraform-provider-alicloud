---
layout: "alicloud"
page_title: "Alicloud: alicloud_instance"
sidebar_current: "docs-alicloud-resource-instance"
description: |-
  Provides a ECS instance resource.
---

# alicloud\_instance

Provides a ECS instance resource.

~> **NOTE:** You can launch an ECS instance for a VPC network via specifying parameter `vswitch_id`. One instance can only belong to one VSwitch.

~> **NOTE:** If a VSwitchId is specified for creating an instance, SecurityGroupId and VSwitchId must belong to one VPC.

~> **NOTE:** Several instance types have outdated in some regions and availability zones, such as `ecs.t1.*`, `ecs.s2.*`, `ecs.n1.*` and so on. If you want to keep them, you should set `is_outdated` to true. For more about the upgraded instance type, refer to `alicloud_instance_types` datasource.

## Example Usage

```
# Create a new ECS instance for a VPC
resource "alicloud_security_group" "group" {
  name        = "tf_test_foo"
  description = "foo"
  vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_instance" "instance" {
  # cn-beijing
  availability_zone = "cn-beijing-b"
  security_groups = ["${alicloud_security_group.group.*.id}"]

  allocate_public_ip = true

  # series III
  instance_type        = "ecs.n4.large"
  system_disk_category = "cloud_efficiency"
  image_id             = "ubuntu_140405_64_40G_cloudinit_20161115.vhd"
  instance_name        = "test_foo"
  vswitch_id = "${alicloud_vswitch.vswitch.id}"
}

# Create a new ECS instance for VPC
resource "alicloud_vpc" "vpc" {
  # Other parameters...
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id = "${alicloud_vpc.vpc.id}"
  # Other parameters...
}

resource "alicloud_slb" "slb" {
  name       = "test-slb-tf"
  vpc_id     = "${alicloud_vpc.vpc.id}"
  vswitch_id = "${alicloud_vswitch.vswitch.id}"
}
```

## Argument Reference

The following arguments are supported:

* `image_id` - (Required) The Image to use for the instance. ECS instance's image can be replaced via changing 'image_id'.
* `instance_type` - (Required) The type of instance to start.
* `io_optimized` - (Deprecated) It has been deprecated on instance resource. All the launched alicloud instances will be I/O optimized.
* `is_outdated` - (Optional) Whether to use outdated instance type. Default to false.
* `security_groups` - (Required)  A list of security group ids to associate with.
* `availability_zone` - (Optional) The Zone to start the instance in.
* `instance_name` - (Optional) The name of the ECS. This instance_name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin or end with a hyphen, and must not begin with http:// or https://. If not specified, 
Terraform will autogenerate a default name is `ECS-Instance`.
* `allocate_public_ip` - (Optional) Associate a public ip address with an instance in a VPC or Classic. Boolean value, Default is false.
* `system_disk_category` - (Optional) Valid values are `cloud_efficiency`, `cloud_ssd` and `cloud`. `cloud` only is used to some no I/O optimized instance. Default to `cloud_efficiency`.
* `system_disk_size` - (Optional) Size of the system disk, value range: 40GB ~ 500GB. Default is 40GB. ECS instance's system disk can be reset when replacing system disk.
* `description` - (Optional) Description of the instance, This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Default value is null.
* `internet_charge_type` - (Optional) Internet charge type of the instance, Valid values are `PayByBandwidth`, `PayByTraffic`. Default is `PayByTraffic`.
* `internet_max_bandwidth_in` - (Optional) Maximum incoming bandwidth from the public network, measured in Mbps (Mega bit per second). Value range: [1, 200]. If this value is not specified, then automatically sets it to 200 Mbps.
* `internet_max_bandwidth_out` - (Optional) Maximum outgoing bandwidth to the public network, measured in Mbps (Mega bit per second). Value range:  [0, 100], If this value is not specified, then automatically sets it to 0 Mbps.
* `host_name` - (Optional) Host name of the ECS, which is a string of at least two characters. “hostname” cannot start or end with “.” or “-“. In addition, two or more consecutive “.” or “-“ symbols are not allowed. On Windows, the host name can contain a maximum of 15 characters, which can be a combination of uppercase/lowercase letters, numerals, and “-“. The host name cannot contain dots (“.”) or contain only numeric characters.
On other OSs such as Linux, the host name can contain a maximum of 30 characters, which can be segments separated by dots (“.”), where each segment can contain uppercase/lowercase letters, numerals, or “_“.
* `password` - (Optional) Password to an instance is a string of 8 to 30 characters. It must contain uppercase/lowercase letters and numerals, but cannot contain special symbols. In order to take effect new password, the instance will be restarted after modifying the password.
* `vswitch_id` - (Optional) The virtual switch ID to launch in VPC. If you want to create instances in VPC network, this parameter must be set.
* `instance_charge_type` - (Optional) Valid values are `PrePaid`, `PostPaid`, The default is `PostPaid`.
* `period` - (Optional) The time that you have bought the resource, in month. Only valid when instance_charge_type is set as `PrePaid`. Value range [1, 12].
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `user_data` - (Optional) The user data to provide when launching the instance.

## Attributes Reference

The following attributes are exported:

* `id` - The instance ID.
* `availability_zone` - The Zone to start the instance in.
* `instance_name` - The instance name.
* `host_name` - The instance host name.
* `description` - The instance description.
* `status` - The instance status.
* `image_id` - The instance Image Id.
* `instance_type` - The instance type.
* `private_ip` - The instance private ip.
* `public_ip` - The instance public ip.
* `vswitch_id` - If the instance created in VPC, then this value is  virtual switch ID.
* `tags` - The instance tags, use jsonencode(item) to display the value.
