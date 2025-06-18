---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_custom"
description: |-
  Provides a Alicloud RDS Custom resource.
---

# alicloud_rds_custom

Provides a RDS Custom resource.

Dedicated RDS User host.

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
* `amount` - (Optional, Int) Represents the number of instances created
* `auto_pay` - (Optional) Whether to pay automatically. Value range:
  - `true` (default): automatic payment. You need to ensure that your account balance is sufficient.
  - `false`: only orders are generated without deduction.

-> **NOTE:**  If the balance of your payment method is insufficient, you can set the parameter AutoPay to false, and an unpaid order will be generated. You can log on to the RDS management console to pay by yourself.

-> **NOTE:** >

* `auto_renew` - (Optional) Whether the instance is automatically renewed. Valid values: true/false. The default is false.
* `create_extra_param` - (Optional, Available since v1.252.0) Reserved parameters are not supported.
* `create_mode` - (Optional) Whether to allow joining the ACK cluster. When this parameter is set to `1`, the created instance can be added to the ACK cluster through The `AttachRCInstances` API to efficiently manage container applications.
  - `1`: Yes.
  - `0` (default): No.
* `data_disk` - (Optional, ForceNew, List) Data disk See [`data_disk`](#data_disk) below.
* `deployment_set_id` - (Optional, ForceNew) The ID of the deployment set.
* `description` - (Optional, ForceNew) Instance description. It must be 2 to 256 characters in length and cannot start with http:// or https.
* `direction` - (Optional) Instance configuration type, value range:

-> **NOTE:**  This parameter does not need to be uploaded, and the system can automatically determine whether to upgrade or downgrade. If you want to upload, please follow the following logic rules.
  - `Up` (default): upgrade the instance specification. Please ensure that your account balance is sufficient.
  - `Down`: Downgrade instance specifications. When the instance type set to InstanceType is lower than the current instance type, set Direction = down.
* `dry_run` - (Optional) Whether to pre-check the operation of creating an instance. Valid values:
  - `true`: The PreCheck operation is performed without creating an instance. Check items include request parameters, request formats, business restrictions, and inventory.
  - `false` (default): Sends a normal request and directly creates an instance after the check is passed.
* `force` - (Optional) Whether to forcibly release the running instance. Value: true/false
* `force_stop` - (Optional) Whether to force shutdown. Value range:
  - `true`: Force shutdown.
  - `false` (default): Normal shutdown.
* `host_name` - (Optional) The instance host name.
* `image_id` - (Optional) The ID of the image used by the instance.
* `instance_charge_type` - (Optional) The Payment type. Currently, only `Prepaid` (package year and month) types are supported.
* `instance_type` - (Required) The type of the created RDS Custom dedicated host instance.
* `internet_charge_type` - (Optional) Reserved parameters are not supported.
* `internet_max_bandwidth_out` - (Optional, Int) Reserved parameters are not supported.
* `io_optimized` - (Optional) Reserved parameters are not supported.
* `key_pair_name` - (Optional) The key pair name. Only flyer names are supported.
* `password` - (Optional) The account and password of the instance.
* `period` - (Optional, Int) Prepaid renewal duration, unit: Month/Year. 
* `period_unit` - (Optional) The unit of duration of the year-to-month billing method. Value range:
  - `Year`: Year
  - `Month` (default): Month
* `resource_group_id` - (Optional, Computed) The ID of the resource group
* `security_enhancement_strategy` - (Optional) Reserved parameters are not supported.
* `security_group_ids` - (Optional, ForceNew, List) Security group list
* `spot_strategy` - (Optional, Available since v1.252.0) The bidding strategy for pay-as-you-go instances. This parameter takes effect when the value of `InstanceChargeType` is set to **PostPaid. Value range:
  - `NoSpot`: normal pay-as-you-go instances.
  - `SpotAsPriceGo`: The system automatically bids and follows the actual price in the current market.

Default value: **NoSpot * *.
* `status` - (Optional, Computed) The status of the resource
* `support_case` - (Optional, Available since v1.252.0) Supported scenarios: createMode:supportCase, for example: NATIVE("0", "eni"),RCK("1", "rck"),ACK_EDGE("1", "edge");
* `system_disk` - (Optional, List) System disk specifications. See [`system_disk`](#system_disk) below.
* `tags` - (Optional, Map) The tag of the resource
* `vswitch_id` - (Required, ForceNew) The ID of the virtual switch. The zone in which the vSwitch is located must correspond to the zone ID entered in ZoneId.
The network type InstanceNetworkType must be VPC.
* `zone_id` - (Optional, ForceNew) The zone ID  of the resource

### `data_disk`

The data_disk supports the following:
* `category` - (Optional, ForceNew) Instance storage type
local_ssd: local SSD disk
cloud_essd:ESSD PL1 cloud disk
* `performance_level` - (Optional, ForceNew) Cloud Disk Performance
Currently only supports PL1
* `size` - (Optional, ForceNew, Int) Instance storage space. Unit: GB.

### `system_disk`

The system_disk supports the following:
* `category` - (Optional) The cloud disk type of the system disk. Currently, only `cloud_essd`(ESSD cloud disk) is supported.
* `size` - (Optional) System disk size, unit: GiB. Only ESSD PL1 is supported. Valid values range from 20 to 2048.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `region_id` - The region ID. Callable DescribeRegions to get.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Custom.
* `delete` - (Defaults to 5 mins) Used when delete the Custom.
* `update` - (Defaults to 7 mins) Used when update the Custom.

## Import

RDS Custom can be imported using the id, e.g.

```shell
$ terraform import alicloud_rds_custom.example <id>
```