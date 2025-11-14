---
subcategory: "Milvus"
layout: "alicloud"
page_title: "Alicloud: alicloud_milvus_instance"
description: |-
  Provides a Alicloud Milvus Instance resource.
---

# alicloud_milvus_instance

Provides a Milvus Instance resource.



For information about Milvus Instance and how to use it, see [What is Instance](https://next.api.alibabacloud.com/document/milvus/2023-10-12/CreateInstance).

-> **NOTE:** Available since v1.263.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-zhangjiakou"
}

variable "region_id" {
  default = "cn-zhangjiakou"
}

variable "zone_id" {
  default = "cn-zhangjiakou-b"
}

resource "alicloud_vpc" "defaultILXuit" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "defaultN80M7S" {
  vpc_id       = alicloud_vpc.defaultILXuit.id
  zone_id      = var.zone_id
  cidr_block   = "172.16.1.0/24"
  vswitch_name = "milvus-example"
}


resource "alicloud_milvus_instance" "default" {
  zone_id = var.zone_id
  vswitch_ids {
    vsw_id  = alicloud_vswitch.defaultN80M7S.id
    zone_id = alicloud_vswitch.defaultN80M7S.zone_id
  }
  db_admin_password = "Test123456@"
  components {
    type    = "standalone"
    cu_num  = "8"
    replica = "1"
    cu_type = "general"
  }
  instance_name         = "镇远测试包年包月"
  db_version            = "2.4"
  vpc_id                = alicloud_vpc.defaultILXuit.id
  ha                    = false
  payment_type          = "Subscription"
  multi_zone_mode       = "Single"
  payment_duration_unit = "year"
  payment_duration      = "1"
}
```

### Deleting `alicloud_milvus_instance` or removing it from your configuration

The `alicloud_milvus_instance` resource allows you to manage  `payment_type = "Subscription"`  instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the Instance.
You can resume managing the subscription instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:
* `auto_backup` - (Optional, Computed) Auto Backup
* `components` - (Optional, Set) Components See [`components`](#components) below.
* `configuration` - (Optional) User Configuration
* `db_admin_password` - (Optional) DB Admin Password

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `db_version` - (Required, ForceNew) Query node size
* `encrypted` - (Optional, ForceNew) Whether to encrypt
* `ha` - (Optional) High Availability
* `instance_name` - (Required) Instance Name
* `kms_key_id` - (Optional, ForceNew) Kms encryption id
* `multi_zone_mode` - (Optional, ForceNew) Availability Zone
* `payment_duration` - (Optional, Int) Payment Duration

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `payment_duration_unit` - (Optional) Payment Duration Unit

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `payment_type` - (Required, ForceNew) The payment type of the resource
* `resource_group_id` - (Optional, Computed) Resource Group ID
* `tags` - (Optional, Map) Tags
* `vswitch_ids` - (Optional, ForceNew, List) VSwitch ID See [`vswitch_ids`](#vswitch_ids) below.
* `vpc_id` - (Required, ForceNew) VPC ID
* `zone_id` - (Optional, ForceNew) Availability Zone

### `components`

The components supports the following:
* `cu_num` - (Required, Int) CU Count
* `cu_type` - (Optional, ForceNew, Computed) Compute Type
* `disk_size_type` - (Optional, ForceNew, Computed) Disk Size
* `replica` - (Required, Int) Replica Count
* `type` - (Required) Component Type

### `vswitch_ids`

The vswitch_ids supports the following:
* `vsw_id` - (Optional, ForceNew) VswId
* `zone_id` - (Optional, ForceNew) Availability Zone

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Create Time
* `region_id` - Region
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 14 mins) Used when create the Instance.
* `delete` - (Defaults to 14 mins) Used when delete the Instance.
* `update` - (Defaults to 21 mins) Used when update the Instance.

## Import

Milvus Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_milvus_instance.example <id>
```