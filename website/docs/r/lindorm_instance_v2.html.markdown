---
subcategory: "Lindorm"
layout: "alicloud"
page_title: "Alicloud: alicloud_lindorm_instance_v2"
description: |-
  Provides a Alicloud Lindorm Instance V2 resource.
---

# alicloud_lindorm_instance_v2

Provides a Lindorm Instance V2 resource.

Cloud-native multi-model database.

For information about Lindorm Instance V2 and how to use it, see [What is Instance V2](https://next.api.alibabacloud.com/document/hitsdb/2020-06-15/CreateLindormV2Instance).

-> **NOTE:** Available since v1.260.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

resource "alicloud_vpc" "defaultR8vXlP" {
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default9umuzwH" {
  vpc_id     = alicloud_vpc.defaultR8vXlP.id
  zone_id    = "cn-beijing-h"
  cidr_block = "172.16.0.0/24"
}

resource "alicloud_vswitch" "defaultgOFAo3L" {
  vpc_id     = alicloud_vpc.defaultR8vXlP.id
  zone_id    = "cn-beijing-l"
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vswitch" "defaultTAbr2pJ" {
  vpc_id     = alicloud_vpc.defaultR8vXlP.id
  zone_id    = "cn-beijing-j"
  cidr_block = "172.16.2.0/24"
}


resource "alicloud_lindorm_instance_v2" "default" {
  standby_zone_id = "cn-beijing-l"
  engine_list {
    engine_type = "TABLE"
    node_group {
      node_count          = "4"
      node_spec           = "lindorm.g.2xlarge"
      resource_group_name = "cx-mz-rg"
    }
  }
  cloud_storage_size = "400"
  primary_zone_id    = "cn-beijing-h"
  zone_id            = "cn-beijing-h"
  cloud_storage_type = "PerformanceStorage"
  arch_version       = "2.0"
  vswitch_id         = alicloud_vswitch.default9umuzwH.id
  standby_vswitch_id = alicloud_vswitch.defaultgOFAo3L.id
  primary_vswitch_id = alicloud_vswitch.default9umuzwH.id
  arbiter_vswitch_id = alicloud_vswitch.defaultTAbr2pJ.id
  vpc_id             = alicloud_vpc.defaultR8vXlP.id
  instance_alias     = "preTest-MZ"
  payment_type       = "POSTPAY"
  arbiter_zone_id    = "cn-beijing-j"
  auto_renewal       = false
}
```

## Argument Reference

The following arguments are supported:
* `arbiter_vswitch_id` - (Optional, ForceNew) Coordination Zone VswitchId
* `arbiter_zone_id` - (Optional, ForceNew) Coordination Zone ZoneId
* `arch_version` - (Required) Deployment Scenario

-> **NOTE:**  Enumeration value

-> **NOTE:**  - 1.0 Single AZ

-> **NOTE:**  - 2.0 Multi-AZ Basic

-> **NOTE:**  - 3.0 Multi-AZ High Availability Edition


-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `auto_renew_duration` - (Optional, Available since v1.262.0) Automatic renewal duration. Unit: Month.

Value range: `1` to **12 * *.

-> **NOTE:**  This item takes effect only when `AutoRenewal` is **true.


-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `auto_renewal` - (Optional) Whether the instance is automatically renewed. Enumerated values:
  - `true`: Automatic renewal.
  - `false`: does not renew automatically.

The default value is false

-> **NOTE:**  This parameter is valid only when the `PayType` value is `PREPAY` (Subscription).


-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `cloud_storage_size` - (Optional, Int) 

-> **NOTE:**  Cloud storage capacity in GB

* `cloud_storage_type` - (Optional, ForceNew) 

-> **NOTE:** >

-> **NOTE:**  - StandardStorage: Standard cloud storage

-> **NOTE:**  - PerformanceStorage: performance-based cloud storage

-> **NOTE:** - capacity storage: Capacity-based cloud storage

* `deletion_protection` - (Optional, Computed) Whether to enable deletion protection
* `duration` - (Optional, Int, Available since v1.262.0) The specified duration when the resource is purchased. Only the subscription instances are valid.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `engine_list` - (Required, Set) Engine List See [`engine_list`](#engine_list) below.
* `instance_alias` - (Required) Instance name
* `payment_type` - (Required) The payment type of the resource
* `pricing_cycle` - (Optional, Available since v1.262.0) Purchase duration unit: Month, Year

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `primary_vswitch_id` - (Optional, ForceNew) Primary zone VswitchId
* `primary_zone_id` - (Optional, ForceNew) Primary zone ZoneID
* `standby_vswitch_id` - (Optional, ForceNew) Standby zone VswitchId
* `standby_zone_id` - (Optional, ForceNew) Standby zone ZoneID
* `vpc_id` - (Required, ForceNew) VpcId
* `vswitch_id` - (Required, ForceNew) VswitchId
* `zone_id` - (Required, ForceNew) The zone ID  of the resource

### `engine_list`

The engine_list supports the following:
* `engine_type` - (Required) Engine
* `node_group` - (Optional, Set) Node Group List See [`node_group`](#engine_list-node_group) below.

### `engine_list-node_group`

The engine_list-node_group supports the following:
* `node_count` - (Required, Int) Number of nodes
* `node_disk_size` - (Optional, Int) Local cloud disk storage capacity
* `node_disk_type` - (Optional) Node Disk Type
* `node_spec` - (Required) Node Specifications
* `resource_group_name` - (Required) Resource group name

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `engine_list` - Engine List
  * `connect_address_list` - Connect Address List
    * `address` - Connect Address
    * `port` - Connect Port
    * `type` - Connect Type:
  * `is_last_version` - Whether it is the latest version
  * `latest_version` - Latest Version
  * `node_group` - Node Group List
    * `category` - Node Type
    * `cpu_core_count` - Number of CPU cores
    * `enable_attach_local_disk` - Whether to mount  local cloud disks
    * `memory_size_gi_b` - Node memory size
    * `spec_id` - Spec Id
    * `status` - Node Status
  * `version` - Engine Version
* `region_id` - The region ID of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 101 mins) Used when create the Instance V2.
* `delete` - (Defaults to 20 mins) Used when delete the Instance V2.
* `update` - (Defaults to 1001 mins) Used when update the Instance V2.

## Import

Lindorm Instance V2 can be imported using the id, e.g.

```shell
$ terraform import alicloud_lindorm_instance_v2.example <id>
```