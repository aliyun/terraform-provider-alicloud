---
subcategory: "Star Rocks"
layout: "alicloud"
page_title: "Alicloud: alicloud_star_rocks_node_group"
description: |-
  Provides a Alicloud Star Rocks Node Group resource.
---

# alicloud_star_rocks_node_group

Provides a Star Rocks Node Group resource.



For information about Star Rocks Node Group and how to use it, see [What is Node Group](https://next.api.alibabacloud.com/document/starrocks/2022-10-19/CreateNodeGroup).

-> **NOTE:** Available since v1.262.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_vpc" "defaultq6pcFe" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "example-vpc-487"
}

resource "alicloud_vswitch" "defaultujlpyG" {
  vpc_id       = alicloud_vpc.defaultq6pcFe.id
  zone_id      = "cn-hangzhou-i"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "sr-example-ng"
}

resource "alicloud_star_rocks_instance" "defaultvjnpM0" {
  cluster_zone_id = "cn-hangzhou-i"
  encrypted       = false
  auto_renew      = false
  pay_type        = "postPaid"
  frontend_node_groups {
    cu                        = "8"
    storage_size              = "100"
    storage_performance_level = "pl1"
    disk_number               = "1"
    zone_id                   = "cn-hangzhou-i"
    spec_type                 = "standard"
    resident_node_number      = "1"
  }
  instance_name = "t1"
  vswitches {
    zone_id    = "cn-hangzhou-i"
    vswitch_id = alicloud_vswitch.defaultujlpyG.id
  }
  vpc_id                  = alicloud_vpc.defaultq6pcFe.id
  version                 = "3.3"
  run_mode                = "shared_data"
  package_type            = "official"
  oss_accessing_role_name = "AliyunEMRStarRocksAccessingOSSRolecn"
  admin_password          = "1qaz@QAZ"
  backend_node_groups {
    cu                        = "8"
    storage_size              = "200"
    zone_id                   = "cn-hangzhou-i"
    spec_type                 = "standard"
    resident_node_number      = "3"
    disk_number               = "1"
    storage_performance_level = "pl1"
  }
}


resource "alicloud_star_rocks_node_group" "default" {
  description                 = "example_desc"
  node_group_name             = "ng_676"
  instance_id                 = alicloud_star_rocks_instance.defaultvjnpM0.id
  spec_type                   = "standard"
  storage_performance_level   = "pl1"
  pricing_cycle               = "1"
  auto_renew                  = false
  storage_size                = "200"
  duration                    = "1"
  pay_type                    = "postPaid"
  cu                          = "8"
  disk_number                 = "1"
  resident_node_number        = "1"
  local_storage_instance_type = "non_local_storage"
  promotion_option_no         = "blank"
}
```

## Argument Reference

The following arguments are supported:
* `auto_renew` - (Optional) Whether auto-renewal is enabled.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `cu` - (Optional, Int) Number of CUs. CU (Compute Unit) is the basic unit of service measurement, where 1 CU = 1 vCPU + 4 GiB memory. When SpecType is memory-optimized, 1 CU = 1 vCPU + 8 GiB memory.
* `description` - (Optional, ForceNew) Description of node group.
* `disk_number` - (Optional, Int) Number of disks.
* `duration` - (Optional, Int) Duration of node group.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `fast_mode` - (Optional) Whether to restart in fast restart mode. The default is false.
  - true: Reboots the compute node in fast restart mode. Restart computing nodes in multiple batches, restart in parallel within a batch, and execute serially between batches;
  - false: Restarts the compute node in rolling restart mode.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `instance_id` - (Optional, ForceNew, Computed) The ID of the instance.
* `local_storage_instance_type` - (Optional, ForceNew) Node group local SSD instance specification. This value is only relevant when based on ECS instances and SpecType is set to local SSD/large-scale storage.
* `node_group_name` - (Optional, ForceNew) The name of the node group.
* `pay_type` - (Optional, ForceNew) Payment type:
  - PrePaid: Subscription (prepaid).
  - PostPaid: Pay-as-you-go (postpaid).
* `pricing_cycle` - (Optional) Unit of purchase duration:
  - Month
  - Year

This is only applicable when payType is set to PrePaid.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `promotion_option_no` - (Optional) ID of promotion option.

-> **NOTE:** This parameter only applies during resource creation, update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `resident_node_number` - (Optional, Int) Number of nodes.
* `spec_type` - (Optional) Node group spec types include the following:
  - standard: Standard edition.
  - localSSD: Local SSD.
  - bigData: Large-scale storage.
  - ramEnhanced: Memory-enhanced type.
  - networkEnhanced: Network-enhanced type.
* `storage_performance_level` - (Optional) Performance levels of the cloud disk. Includes the following values:
  - pl0: Maximum random read/write IOPS of a single disk is 10,000.
  - pl1: Maximum random read/write IOPS of a single disk is 50,000.
  - pl2: Maximum random read/write IOPS of a single disk is 100,000.
  - pl3: Maximum random read/write IOPS of a single disk is 1,000,000.
* `storage_size` - (Optional, Int) Storage size, measured in GiB.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<instance_id>:<node_group_id>`.
* `create_time` - The creation time of the node group.
* `node_group_id` - The ID of the node group.
* `region_id` - The region ID of the node group.
* `status` - Node group status, including the following values:_FAILED: Creation failed._CONFIG: Modifying configuration._TIMEZONE: Modifying timezone._SCALING_OUT: Elastic scaling out._SCALING_IN: Elastic scaling in._OUT: Scaling out._IN: Scaling in._UP: Scaling up (upgrading configuration)._DOWN: Scaling down (downgrading configuration)._PUBLIC_NETWORK: Enabling public network._PUBLIC_NETWORK: Disabling public network._AZ: Switching availability zones.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 20 mins) Used when create the Node Group.
* `delete` - (Defaults to 60 mins) Used when delete the Node Group.
* `update` - (Defaults to 20 mins) Used when update the Node Group.

## Import

Star Rocks Node Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_star_rocks_node_group.example <instance_id>:<node_group_id>
```