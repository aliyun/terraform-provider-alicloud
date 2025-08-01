---
subcategory: "Star Rocks"
layout: "alicloud"
page_title: "Alicloud: alicloud_star_rocks_instance"
description: |-
  Provides a Alicloud Star Rocks Instance resource.
---

# alicloud_star_rocks_instance

Provides a Star Rocks Instance resource.

StarRocks resource instance.

For information about Star Rocks Instance and how to use it, see [What is Instance](https://next.api.alibabacloud.com/document/starrocks/2022-10-19/CreateInstanceV1).

-> **NOTE:** Available since v1.256.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_star_rocks_instance&exampleId=478e545f-d6d6-8b4e-15eb-e2aba9eb0894c32daacb&activeTab=example&spm=docs.r.star_rocks_instance.0.478e545fd6&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultB21JUD" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default106DkE" {
  vpc_id       = alicloud_vpc.defaultB21JUD.id
  cidr_block   = "172.16.1.0/24"
  vswitch_name = "sr-example"
  zone_id      = "cn-hangzhou-i"
}


resource "alicloud_star_rocks_instance" "default" {
  instance_name = "create-instance-1"
  auto_renew    = false
  frontend_node_groups {
    cu                          = "8"
    storage_size                = "100"
    resident_node_number        = "3"
    storage_performance_level   = "pl1"
    spec_type                   = "standard"
    disk_number                 = "1"
    zone_id                     = "cn-hangzhou-i"
    local_storage_instance_type = "null"
  }
  vswitches {
    vswitch_id = alicloud_vswitch.default106DkE.id
    zone_id    = "cn-hangzhou-i"
  }
  backend_node_groups {
    cu                          = "8"
    storage_size                = "100"
    resident_node_number        = "3"
    disk_number                 = "1"
    storage_performance_level   = "pl1"
    spec_type                   = "standard"
    zone_id                     = "cn-hangzhou-i"
    local_storage_instance_type = "null"
  }
  cluster_zone_id         = "cn-hangzhou-i"
  duration                = "1"
  pay_type                = "postPaid"
  vpc_id                  = alicloud_vpc.defaultB21JUD.id
  version                 = "3.3"
  run_mode                = "shared_data"
  package_type            = "official"
  admin_password          = "1qaz@QAZ"
  oss_accessing_role_name = "AliyunEMRStarRocksAccessingOSSRole"
  pricing_cycle           = "Month"
  kms_key_id              = "123"
  promotion_option_no     = "123"
  encrypted               = false
  observer_node_groups {
    cu                          = "8"
    storage_size                = "100"
    storage_performance_level   = "pl1"
    disk_number                 = "1"
    resident_node_number        = "1"
    spec_type                   = "standard"
    local_storage_instance_type = "null"
    zone_id                     = "cn-hangzhou-h"
  }
}
```

## Argument Reference

The following arguments are supported:
* `admin_password` - (Required) Password of admin user.
* `auto_renew` - (Optional) Whether to enable automatic renewal. This is only meaningful when payType is set to PrePaid. Disabled by default.
* `backend_node_groups` - (Optional, ForceNew, List) BackendNodeGroups See [`backend_node_groups`](#backend_node_groups) below.
* `cluster_zone_id` - (Required) ZoneId of instance.
* `duration` - (Optional, Int) Duration of purchase. It is only meaningful when payType is set to PrePaid.
* `encrypted` - (Optional, ForceNew) Whether encrypted
* `frontend_node_groups` - (Optional, ForceNew, List) FrontendNodeGroups See [`frontend_node_groups`](#frontend_node_groups) below.
* `instance_name` - (Required) The name of the instance.
* `kms_key_id` - (Optional, ForceNew) KmsKeyId
* `observer_node_groups` - (Optional, ForceNew, List) ObserverNodeGroups See [`observer_node_groups`](#observer_node_groups) below.
* `oss_accessing_role_name` - (Optional) Role name used for password-free access to OSS.
* `package_type` - (Required, ForceNew) The package type of the instance:
  - trial
  - official
* `pay_type` - (Required, ForceNew) The pay type of the instance:
  - prePaid
  - postPaid
* `pricing_cycle` - (Optional) The duration unit for purchasing:
  - Month
  - Year
This is only meaningful when PayType is set to PrePaid.
* `promotion_option_no` - (Optional) Promotion
* `resource_group_id` - (Optional, Computed) ResourceGroupId
* `run_mode` - (Required, ForceNew) The run mode of the instance:
  - shared_nothing
  - shared_data
  - lakehouse
* `tags` - (Optional, Map) Tag list of the instance.
* `version` - (Required, ForceNew) The version of the instance.
* `vpc_id` - (Required, ForceNew) The VPC ID of the instance.
* `vswitches` - (Optional, ForceNew, List) The VSwitches info of the instance. See [`vswitches`](#vswitches) below.

### `backend_node_groups`

The backend_node_groups supports the following:
* `cu` - (Optional, ForceNew, Int) Number of CUs. CU (Compute Unit) is the basic measurement unit of the service, where 1 CU = 1 CPU core + 4 GiB memory.
* `disk_number` - (Optional, ForceNew, Int) The number of disks.
* `local_storage_instance_type` - (Optional, ForceNew) Local SSD instance specifications.
* `resident_node_number` - (Optional, ForceNew, Int) Resident node number of node group.
* `spec_type` - (Optional, ForceNew) Compute group specification types include the following:
  - standard
  - localSSD
  - bigData
  - ramEnhanced
  - networkEnhanced
* `storage_performance_level` - (Optional, ForceNew) Performance levels of cloud disks include the following values:
  - pl0: Maximum random read/write IOPS per disk is 10,000.
  - pl1: Maximum random read/write IOPS per disk is 50,000.
  - pl2: Maximum random read/write IOPS per disk is 100,000.
  - pl3: Maximum random read/write IOPS per disk is 1,000,000.
* `storage_size` - (Optional, ForceNew, Int) Storage size, measured in GiB.
* `zone_id` - (Optional, ForceNew) Zone ID.

### `frontend_node_groups`

The frontend_node_groups supports the following:
* `cu` - (Optional, ForceNew, Int) Number of CUs. CU (Compute Unit) is the basic measurement unit of the service, where 1 CU = 1 CPU core + 4 GiB memory.
* `disk_number` - (Optional, ForceNew, Int) DiskNumber
* `local_storage_instance_type` - (Optional, ForceNew) Local SSD instance specifications.
* `resident_node_number` - (Optional, ForceNew, Int) Resident node number of node group.
* `spec_type` - (Optional, ForceNew) Compute group specification types include the following:
  - standard
  - ramEnhanced
* `storage_performance_level` - (Optional, ForceNew) Performance levels of cloud disks include the following values:
  - pl0: Maximum random read/write IOPS per disk is 10,000.
  - pl1: Maximum random read/write IOPS per disk is 50,000.
  - pl2: Maximum random read/write IOPS per disk is 100,000.
  - pl3: Maximum random read/write IOPS per disk is 1,000,000.
* `storage_size` - (Optional, ForceNew, Int) Storage size, measured in GiB.
* `zone_id` - (Optional, ForceNew) Zone ID.

### `observer_node_groups`

The observer_node_groups supports the following:
* `cu` - (Optional, ForceNew, Int) Number of CUs. CU (Compute Unit) is the basic measurement unit of the service, where 1 CU = 1 CPU core + 4 GiB memory.
* `disk_number` - (Optional, ForceNew, Int) DiskNumber
* `local_storage_instance_type` - (Optional, ForceNew) Local SSD instance specifications.
* `resident_node_number` - (Optional, ForceNew, Int) Resident node number of node group.
* `spec_type` - (Optional, ForceNew) Compute group specification types include the following:
  - standard
* `storage_performance_level` - (Optional, ForceNew) Performance levels of cloud disks include the following values:
  - pl0: Maximum random read/write IOPS per disk is 10,000.
  - pl1: Maximum random read/write IOPS per disk is 50,000.
  - pl2: Maximum random read/write IOPS per disk is 100,000.
  - pl3: Maximum random read/write IOPS per disk is 1,000,000.
* `storage_size` - (Optional, ForceNew, Int) Storage size, measured in GiB.
* `zone_id` - (Optional, ForceNew) Zone ID.

### `vswitches`

The vswitches supports the following:
* `vswitch_id` - (Required, ForceNew) ID of VSwitch.
* `zone_id` - (Optional, ForceNew) Zone ID of VSwitch.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the instance.
* `region_id` - The region ID of the instance.
* `status` - The status of the instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 20 mins) Used when create the Instance.
* `delete` - (Defaults to 10 mins) Used when delete the Instance.
* `update` - (Defaults to 5 mins) Used when update the Instance.

## Import

Star Rocks Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_star_rocks_instance.example <id>
```