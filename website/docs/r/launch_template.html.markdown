---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_launch_template"
sidebar_current: "docs-alicloud-resource-launch-tempate"
description: |-
  Provides an ECS Launch Template resource.
---

# alicloud\_launch\_template

Provides an ECS Launch Template resource.

For information about Launch Template and how to use it, see [Launch Template](https://www.alibabacloud.com/help/doc-detail/73916.html).

-> **DEPRECATED:**  This resource  has been deprecated from version `1.120.0`. Please use new resource [alicloud_ecs_launch_template](https://www.terraform.io/docs/providers/alicloud/r/ecs_launch_template).

## Example Usage

```terraform
data "alicloud_images" "images" {
  owners = "system"
}

data "alicloud_instances" "instances" {
}

resource "alicloud_launch_template" "template" {
  name                          = "tf-template"
  description                   = "test1"
  image_id                      = data.alicloud_images.images.images[0].id
  host_name                     = "tf-host"
  instance_charge_type          = "PrePaid"
  instance_name                 = "tf-instance-name"
  instance_type                 = data.alicloud_instances.instances.instances[0].instance_type
  internet_charge_type          = "PayByBandwidth"
  internet_max_bandwidth_in     = 5
  internet_max_bandwidth_out    = 0
  io_optimized                  = "none"
  key_pair_name                 = "tf-key-pair"
  ram_role_name                 = "xxxxx"
  network_type                  = "vpc"
  security_enhancement_strategy = "Active"
  spot_price_limit              = 5
  spot_strategy                 = "SpotWithPriceLimit"
  security_group_id             = "sg-zxcvj0lasdf102350asdf9a"
  system_disk_category          = "cloud_ssd"
  system_disk_description       = "test disk"
  system_disk_name              = "hello"
  system_disk_size              = 40
  resource_group_id             = "rg-zkdfjahg9zxncv0"
  userdata                      = "xxxxxxxxxxxxxx"
  vswitch_id                    = "sw-ljkngaksdjfj0nnasdf"
  vpc_id                        = "vpc-asdfnbg0as8dfk1nb2"
  zone_id                       = "beijing-a"

  tags = {
    tag1 = "hello"
    tag2 = "world"
  }
  network_interfaces {
    name              = "eth0"
    description       = "hello1"
    primary_ip        = "10.0.0.2"
    security_group_id = "xxxx"
    vswitch_id        = "xxxxxxx"
  }
  data_disks {
    name        = "disk1"
    description = "test1"
  }
  data_disks {
    name        = "disk2"
    description = "test2"
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_launch_template&spm=docs.r.launch_template.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `launch_template_name` - (Optional, ForceNew) The name of the launch template. It must be 2 to 128 characters in length. It must start with an English letter or Chinese, can contain numbers, periods (.), colons (:), underscores (_), and hyphens (-). It cannot start with "http://" or "https://".
* `name` - (Optional, ForceNew, Deprecated) It has been deprecated from version 1.120.0. New field `launch_template_name` instead.
* `description` - (Optional) Description of instance launch template version 1. It can be [2, 256] characters in length. It cannot start with "http://" or "https://". The default value is null.
* `host_name` - (Optional) Instance host name. It cannot start or end with a period (.) or a hyphen (-) and it cannot have two or more consecutive periods (.) or hyphens (-). For Windows: The host name can be [2, 15] characters in length. It can contain A-Z, a-z, numbers, periods (.), and hyphens (-). It cannot only contain numbers. For other operating systems: The host name can be [2, 64] characters in length. It can be segments separated by periods (.). It can contain A-Z, a-z, numbers, and hyphens (-).
* `image_id` - (Optional) Image ID.
* `image_owner_alias` - (Optional) The image owner alias. Valid values: `system`, `self`, `others`, `marketplace`.
* `instance_name` - (Optional) The name of the instance. The name is a string of 2 to 128 characters. It must begin with an English or a Chinese character. It can contain A-Z, a-z, Chinese characters, numbers, periods (.), colons (:), underscores (_), and hyphens (-).
* `instance_charge_type` - (Optional) Billing methods. Optional values:
  - PrePaid: Monthly, or annual subscription. Make sure that your registered credit card is invalid or you have insufficient balance in your PayPal account. Otherwise, InvalidPayMethod error may occur.
  - PostPaid: Pay-As-You-Go.

    Default value: PostPaid.
* `instance_type` - (Optional) Instance type. For more information, call resource_alicloud_instances to obtain the latest instance type list.
* `auto_release_time` - (Optional) Instance auto release time. The time is presented using the ISO8601 standard and in UTC time. The format is YYYY-MM-DDTHH:MM:SSZ.
* `auto_renew` - (Optional, ForceNew) Whether to enable auto-renewal for the instance.
* `auto_renew_period` - (Optional, ForceNew) The auto-renewal period of the instance.
* `internet_charge_type` - (Optional) Internet bandwidth billing method. Optional values: `PayByTraffic` | `PayByBandwidth`.
* `internet_max_bandwidth_in` - (Optional) The maximum inbound bandwidth from the Internet network, measured in Mbit/s. Value range: [1, 200].
* `internet_max_bandwidth_out` - (Optional) Maximum outbound bandwidth from the Internet, its unit of measurement is Mbit/s. Value range: [0, 100].
* `io_optimized` - (Optional) Whether it is an I/O-optimized instance or not. Optional values:
  - none
  - optimized
* `key_pair_name` - (Optional) The name of the key pair.
  - Ignore this parameter for Windows instances. It is null by default. Even if you enter this parameter, only the Password content is used.
  - The password logon method for Linux instances is set to forbidden upon initialization.
* `network_type` - (Optional) Network type of the instance. Value options: `classic` | `vpc`.
* `password_inherit` - (Optional) Whether to use the password inherited from the image.
* `period` - (Optional, ForceNew) The purchase duration of the instance.
* `period_unit` - (Optional, ForceNew) The unit of the purchase duration.
* `private_ip_address` - (Optional) The private IP address of the instance.
* `ram_role_name` - (Optional) The RAM role name of the instance. You can use the RAM API ListRoles to query instance RAM role names.
* `resource_group_id` - (Optional) The ID of the resource group to which the resource belongs.
* `security_enhancement_strategy` - (Optional) Whether or not to activate the security enhancement feature and install network security software free of charge. Optional values: Active | Deactive.
* `security_group_id` - (Optional) The security group ID.
* `security_group_ids` - (Optional) The list of security group IDs.
* `spot_duration` - (Optional) The protection period of the preemptible instance. Valid values: `0`, `1`, `2`, `3`, `4`, `5`, `6`. Default to `1`.
* `spot_price_limit` - (Optional) Sets the maximum hourly instance price. Supports up to three decimal places.
* `spot_strategy` - (Optional) The spot strategy for a Pay-As-You-Go instance. This parameter is valid and required only when InstanceChargeType is set to PostPaid. Value range:
  - NoSpot: Normal Pay-As-You-Go instance.
  - SpotWithPriceLimit: Sets the maximum price for a spot instance.
  - SpotAsPriceGo: The system automatically calculates the price. The maximum value is the Pay-As-You-Go price.
* `system_disk_category` - (Optional, Deprecated) It has been deprecated from version 1.120.0. New field `system_disk` instead.
* `system_disk_description` - (Optional, Deprecated) It has been deprecated from version 1.120.0. New field `system_disk` instead.
* `system_disk_name` - (Optional, Deprecated) It has been deprecated from version 1.120.0. New field `system_disk` instead.
* `system_disk_size` - (Optional, Deprecated) It has been deprecated from version 1.120.0. New field `system_disk` instead.
* `userdata` - (Optional, Deprecated) It has been deprecated from version 1.120.0. New field `user_data` instead.
* `user_data` - (Optional) User data of the instance, which is Base64-encoded. Size of the raw data cannot exceed 16 KB.
* `vswitch_id` - (Optional) When creating a VPC-Connected instance, you must specify its VSwitch ID.
* `version_description` - (Optional) The description of the launch template version.
* `vpc_id` - (Optional) The VPC ID.
* `zone_id` - (Optional) The zone ID of the instance.
* `http_endpoint` - (Optional) Whether to enable access to instance metadata.
* `http_tokens` - (Optional) Whether to use the hardened mode (IMDSv2) when accessing instance metadata.
* `http_put_response_hop_limit` - (Optional) The HTTP PUT response hop limit required for instance metadata requests.
* `deployment_set_id` - (Optional) The ID of the deployment set.
* `enable_vm_os_config` - (Optional) Whether to enable VM OS configuration.
* `tags` - (Optional) A mapping of tags to assign to the resource.
  - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
  - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.
* `template_tags` - (Optional) A mapping of tags to assign to the launch template.
* `template_resource_group_id` - (Optional, ForceNew) The ID of the resource group for the launch template.
* `default_version_number` - (Optional) The version number of the default launch template version.
* `update_default_version_number` - (Optional) Whether to update the default version of the launch template to the latest version automatically.
* `system_disk` - (Optional) The system disk. See [`system_disk`](#system_disk) below.
* `network_interfaces` - (Optional) The list of network interfaces created with instance. See [`network_interfaces`](#network_interfaces) below.
* `data_disks` - (Optional) The list of data disks created with instance. See [`data_disks`](#data_disks) below.
* `image_options` - (Optional) The options of images. See [`image_options`](#image_options) below.

### `system_disk`

The system_disk supports the following:

* `category` - (Optional) The category of the system disk.
* `delete_with_instance` - (Optional) Whether to release the system disk when the instance is released. Default to `true`.
* `description` - (Optional) System disk description.
* `iops` - (Optional) The Iops.
* `name` - (Optional) System disk name.
* `performance_level` - (Optional) The performance level of the ESSD used as the system disk.
* `size` - (Optional) Size of the system disk, measured in GB.
* `encrypted` - (Optional) Whether the system disk is encrypted.
* `kms_key_id` - (Optional) The ID of the KMS key to use for the system disk.

### `network_interfaces`

The network_interfaces supports the following:

* `description` - (Optional) The ENI description.
* `name` - (Optional) ENI name.
* `primary_ip` - (Optional) The primary private IP address of the ENI.
* `security_group_id` - (Optional) The security group ID must be one in the same VPC.
* `vswitch_id` - (Optional) The VSwitch ID for ENI. The instance must be in the same zone of the same VPC network as the ENI, but they may belong to different VSwitches.
* `delete_on_release` - (Optional) Whether to release the ENI when the instance is released.

### `data_disks`

The data_disks supports the following:

* `category` - (Optional) The category of the disk.
* `delete_with_instance` - (Optional) Delete this data disk when the instance is destroyed. Default to `true`.
* `description` - (Optional) The description of the data disk.
* `encrypted` - (Optional) Encrypted the data in this disk. Default to `false`.
* `name` - (Optional) The name of the data disk.
* `performance_level` - (Optional) The performance level of the ESSD used as the data disk.
* `size` - (Optional) The size of the data disk.
* `snapshot_id` - (Optional) The snapshot ID used to initialize the data disk.
* `device` - (Optional) The mount point of the data disk.
* `kms_key_id` - (Optional) The ID of the KMS key used for the data disk.

### `image_options`

The image_options supports the following:

* `login_as_non_root` - (Optional) Whether the instance that uses the image supports logons of the ecs-user user.

## Attributes Reference

The following attributes are exported:

* `id` - The Launch Template ID.
* `latest_version_number` - The latest version number of the launch template.

## Import

Launch Template can be imported using the id, e.g.

```shell
$ terraform import alicloud_launch_template.lt lt-abc1234567890000
```
