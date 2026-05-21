---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_custom"
description: |-
  Provides a Alicloud RDS Custom resource.
---

# alicloud_rds_custom

Provides a RDS Custom resource.

RDS dedicated host for users.

For information about RDS Custom and how to use it, see [What is Custom](https://next.api.alibabacloud.com/document/Rds/2014-08-15/RunRCInstances).

-> **NOTE:** Available since v1.247.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

variable "cluster_id" {
  default = "c18c40b2b336840e2b2bbf8ab291758e2"
}

variable "deploymentsetid" {
  default = "ds-2ze78ef5kyj9eveue92m"
}

variable "vswtich-id" {
  default = "example_vswitch"
}

variable "vpc_name" {
  default = "beijing111"
}

variable "example_region_id" {
  default = "cn-beijing"
}

variable "description" {
  default = "ran_1-08_rccreatenodepool_api"
}

variable "example_zone_id" {
  default = "cn-beijing-h"
}

variable "securitygroup_name" {
  default = "rds_custom_init_sg_cn_beijing"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "vpcId" {
  vpc_name = var.vpc_name
}

resource "alicloud_vswitch" "vSwitchId" {
  vpc_id       = alicloud_vpc.vpcId.id
  zone_id      = var.example_zone_id
  vswitch_name = var.vswtich-id
  cidr_block   = "172.16.5.0/24"
}

resource "alicloud_security_group" "securityGroupId" {
  vpc_id              = alicloud_vpc.vpcId.id
  security_group_name = var.securitygroup_name
}

resource "alicloud_ecs_deployment_set" "deploymentSet" {
}

resource "alicloud_ecs_key_pair" "KeyPairName" {
  key_pair_name = alicloud_vswitch.vSwitchId.id
}


resource "alicloud_rds_custom" "default" {
  amount        = "1"
  auto_renew    = false
  period        = "1"
  auto_pay      = true
  instance_type = "mysql.x2.xlarge.6cm"
  data_disk {
    category          = "cloud_essd"
    size              = "50"
    performance_level = "PL1"
  }
  status                        = "Running"
  security_group_ids            = ["${alicloud_security_group.securityGroupId.id}"]
  io_optimized                  = "optimized"
  description                   = var.description
  key_pair_name                 = alicloud_ecs_key_pair.KeyPairName.id
  zone_id                       = var.example_zone_id
  instance_charge_type          = "Prepaid"
  internet_max_bandwidth_out    = "0"
  image_id                      = "aliyun_2_1903_x64_20G_alibase_20240628.vhd"
  security_enhancement_strategy = "Active"
  period_unit                   = "Month"
  password                      = "jingyiTEST@123"
  system_disk {
    size     = "40"
    category = "cloud_essd"
  }
  host_name         = "1743386110"
  create_mode       = "0"
  spot_strategy     = "NoSpot"
  vswitch_id        = alicloud_vswitch.vSwitchId.id
  support_case      = "eni"
  deployment_set_id = var.deploymentsetid
  dry_run           = false
}
```

## Argument Reference

The following arguments are supported:

* `amount` - (Optional, Int) Specifies the number of RDS Custom instances to create. This parameter applies only when creating multiple RDS Custom instances at once.
Valid values: `1` to `5`. Default value: `1`.

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `auto_pay` - (Optional) Specifies whether to enable automatic payment. Valid values:
  - `true` (default): Enable automatic payment. You must ensure that your account balance is sufficient.
  - `false`: Generate an order without charging your account.

-> **NOTE:**  If your payment method has insufficient funds, you can set the AutoPay parameter to false. In this case, an unpaid order is generated, and you can log on to the RDS console to complete the payment manually.

-> **NOTE:** This parameter only takes effect when other resource properties are also modified. Changing this parameter alone will not trigger a resource update.

* `auto_renew` - (Optional) Specifies whether the instance is automatically renewed. This parameter applies only when you create a subscription instance. Valid values:
  - `true`: Enable auto-renewal.
  - `false`: Disable auto-renewal.

-> **NOTE:**  * If you purchase the instance on a monthly basis, the auto-renewal period is one month.

-> **NOTE:**  * If you purchase the instance on an annual basis, the auto-renewal period is one year.


-> **NOTE:** This parameter is only evaluated during resource creation and update. Modifying it in isolation will not trigger any action.

* `create_mode` - (Optional) Specifies whether the instance can be added to an ACK cluster. When this parameter is set to `1`, the created instance can be added to an ACK cluster by using the `AttachRCInstances` API operation, enabling efficient management of containerized applications.
  - `1`: Yes.
  - `0` (default): No.

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `data_disk` - (Optional, ForceNew, Computed, List) List of data disks.   See [`data_disk`](#data_disk) below.
* `deployment_set_id` - (Optional, ForceNew) Deployment set ID.
* `description` - (Optional, Computed) The instance description. It must be 2 to 256 characters in length and cannot start with http:// or https://.
* `direction` - (Optional) The instance specification change type. Valid values:

-> **NOTE:**  You do not need to specify this parameter because the system can automatically determine whether to upgrade or downgrade the instance. If you choose to specify it, follow the rules below:
  - `Up` (default): Upgrade the instance specification. Ensure that your account has sufficient balance.
  - `Down`: Downgrade the instance specification. Set Direction=Down when the instance type specified by InstanceType is lower than the current instance type.

-> **NOTE:** This parameter only takes effect when other resource properties are also modified. Changing this parameter alone will not trigger a resource update.

* `dry_run` - (Optional) Specifies whether to perform a dry run of the instance creation request. Valid values:
  - `true`: performs a dry run without creating the instance. The system checks the request parameters, request format, service limits, and available inventory.
  - `false` (default): sends the actual request and creates the instance if all checks pass.

-> **NOTE:** This parameter only takes effect when other resource properties are also modified. Changing this parameter alone will not trigger a resource update.

* `force` - (Optional) Specifies whether to forcibly release a running instance. Valid values:
  - `true`: Force release.
  - `false` (default): Do not force release.

-> **NOTE:** This parameter configures deletion behavior and is only evaluated when Terraform attempts to destroy the resource. Changes to this parameter during updates are stored but have no immediate effect.

* `force_stop` - (Optional) Specifies whether to force shut down the instance. Valid values:
  - `true`: Force shut down.
  - `false` (default): Shut down normally.

-> **NOTE:** This parameter only takes effect when other resource properties are also modified. Changing this parameter alone will not trigger a resource update.

* `host_name` - (Optional) The hostname of the instance.

-> **NOTE:** This parameter is only evaluated during resource creation and update. Modifying it in isolation will not trigger any action.

* `image_id` - (Optional) The ID of the image used by the instance.

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `instance_charge_type` - (Optional) The billing method. Valid values:
  - `Prepaid`: subscription.
  - `Postpaid`: pay-as-you-go.

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `instance_name` - (Optional, Computed, Available since v1.279.0) The name must be 2 to 128 characters in length, start with a letter or Chinese character, and can contain letters, Chinese characters, digits, periods (.), underscores (_), colons (:), or hyphens (-). By default, the instance name is the same as the InstanceId. When creating multiple RdsCustom instances, you can specify sequential instance names in batches by using square brackets ([]) and commas (,). For more information, see [Create an RDS Custom instance](https://help.aliyun.com/zh/rds/apsaradb-rds-for-mysql/create-an-rds-custom-instance?spm=a2c4g.11186623.0.0.36ef7288jg7aZD#00481f9ba381u).
* `instance_type` - (Required) The target instance type for configuration changes. For the list of instance types supported by RDS Custom instances, see [RDS Custom Instance Types](https://help.aliyun.com/document_detail/2844823.html).
* `internet_charge_type` - (Optional) Reserved parameter. Not supported currently.

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `internet_max_bandwidth_out` - (Optional, Int) The maximum outbound public bandwidth for Custom for SQL Server, measured in Mbit/s.
Valid values: 0 to 1024. Default value: 0.

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `io_optimized` - (Optional) This parameter is reserved and currently unsupported.

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `key_pair_name` - (Optional) The name of the key pair. Only a single key pair name is supported.

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `password` - (Optional) The account password for the instance. It must be 8 to 30 characters in length and contain at least three of the following character types: uppercase letters, lowercase letters, digits, and special characters. Supported special characters include: `()~!@#$%^&*-_+=|{}[]:;',.?/`.

-> **NOTE:** This parameter is only evaluated during resource creation and update. Modifying it in isolation will not trigger any action.

* `period` - (Optional, Int) The subscription duration of the resource. Default value: `1`.

-> **NOTE:** This parameter is only evaluated during resource creation and update. Modifying it in isolation will not trigger any action.

* `period_unit` - (Optional) The unit of subscription duration for the subscription billing method. Valid values:
  - `Year`: Year
  - `Month` (default): Month

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `private_ip_address` - (Optional, ForceNew, Computed, Available since v1.279.0) The private IP address of the instance. When assigning a private IP address to an ECS instance in a Virtual Private Cloud (VPC), you must select an available IP address from the CIDR block of the specified vSwitch (VSwitchId).
* `resource_group_id` - (Optional, Computed) The resource group ID. You can call ListResourceGroups to obtain it.
* `security_enhancement_strategy` - (Optional) This is a reserved parameter and is not currently supported.

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `security_group_ids` - (Optional, ForceNew, List) The ID of the security group to which the instance belongs. Instances in the same security group can access each other. The maximum number of instances that a security group can contain depends on the security group type. For more information, see the "Security groups" section in [Limits](https://help.aliyun.com/document_detail/25412.html).

-> **NOTE:**  The SecurityGroupId determines the network type of the instance. For example, if the specified security group uses the Virtual Private Cloud (VPC) network type, the instance is of the VPC type and you must also specify the VSwitchId parameter.

* `spot_strategy` - (Optional, Available since v1.252.0) The spot strategy for pay-as-you-go instances. This parameter takes effect only when `InstanceChargeType` is set to `PostPaid`. Valid values:
  - `NoSpot`: A regular pay-as-you-go instance.
  - `SpotAsPriceGo`: The system automatically bids based on the current market price.

Default value: `NoSpot`.

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `status` - (Optional, Computed) The status of the instance. Valid values:
  - `Pending`: The instance is being created.
  - `Running`: The instance is running.
  - `Starting`: The instance is starting.
  - `Stopping`: The instance is stopping.
  - `Stopped`: The instance is stopped.
* `support_case` - (Optional, Available since v1.252.0) The deployment type of RDS Custom. Valid values:
  - `eni`: Dual ENI.
  - `edge`: Edge node pool.
  - `share`: VPC.

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `system_disk` - (Optional, ForceNew, List) The system disk specification. See [`system_disk`](#system_disk) below.

-> **NOTE:** Since v1.279.0, `system_disk` is treated as a ForceNew field. Any change to this field, including its nested `category` and `size` values, will force replacement of the `alicloud_rds_custom` resource.

* `tags` - (Optional, Map) Details of the queried instances and their tags.
* `vswitch_id` - (Required, ForceNew) The virtual switch ID of the target instance. If you are creating a VPC-type RDS Custom instance, you must specify the virtual switch ID, and the security group and virtual switch must belong to the same Virtual Private Cloud (VPC).

-> **NOTE:**  If you specify the VSwitchId parameter, the ZoneId parameter you set must match the zone where the virtual switch is located. Alternatively, you can omit the ZoneId parameter, and the system will automatically select the zone of the specified virtual switch.

* `zone_id` - (Optional, ForceNew, Computed) The zone ID of the instance. You can call DescribeZones to obtain the list of available zones.

-> **NOTE:**  If you specify the VSwitchId parameter, the specified ZoneId must match the zone where the vSwitch is located. Alternatively, you can omit ZoneId, and the system will automatically select the zone of the specified vSwitch.

* `create_extra_param` - (Optional, Available since v1.252.0) Reserved parameters are not supported.

### `data_disk`

The data_disk supports the following:
* `category` - (Optional, ForceNew) The type of data disk. Valid values:
  - `cloud_efficiency`: Ultra disk.
  - `cloud_ssd`: SSD cloud disk.
  - `cloud_essd` (default): ESSD cloud disk.
  - `cloud_auto`: High-performance cloud disk.
* `performance_level` - (Optional, ForceNew) The performance level for an ESSD cloud disk. For information about performance differences among ESSD cloud disks, see [ESSD cloud disks](https://help.aliyun.com/document_detail/2859916.html). Valid values:
  - `PL0`
  - `PL1` (default)
  - `PL2`
  - `PL3`.
* `size` - (Optional, ForceNew, Int) The size of the data disk, in GiB. Valid values:
  - cloud_efficiency: 20 to 32,768.
  - cloud_ssd: 20 to 32,768.
  - cloud_auto: 1 to 65,536.
  - cloud_essd: The valid range depends on the value of **DataDisk.PerformanceLevel**.
  - PL0: 1 to 65,536.
  - PL1: 20 to 65,536.
  - PL2: 461 to 65,536.
  - PL3: 1,261 to 65,536.

### `system_disk`

The system_disk supports the following:
* `category` - (Optional, ForceNew) The system disk category. Valid values:
  - `cloud_efficiency`: ultra disk.
  - `cloud_ssd`: standard SSD.
  - `cloud_essd` (default): ESSD.
  - `cloud_auto`: high-performance cloud disk.
* `size` - (Optional, ForceNew) The size of the system disk, in GiB. The value must be greater than or equal to the size of the image specified by the `ImageId` parameter.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `region_id` - The region ID.
* `system_disk_id` - The ID of the system disk attached to the Custom instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Custom.
* `delete` - (Defaults to 5 mins) Used when delete the Custom.
* `update` - (Defaults to 7 mins) Used when update the Custom.

## Import

RDS Custom can be imported using the id, e.g.

```shell
$ terraform import alicloud_rds_custom.example <instance_id>
```
