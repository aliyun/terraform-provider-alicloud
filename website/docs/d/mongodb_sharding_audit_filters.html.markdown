---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_sharding_audit_filters"
description: |-
  Provides a list of Mongodb Sharding Audit Filter to user.
---

# alicloud_mongodb_sharding_audit_filters

This data source provides the Mongodb Sharding Audit Filter of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.290.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_mongodb_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_mongodb_zones.default.zones.0.id
}

resource "alicloud_mongodb_sharding_instance" "default" {
  zone_id        = data.alicloud_mongodb_zones.default.zones.0.id
  vswitch_id     = data.alicloud_vswitches.default.ids[0]
  engine_version = "4.2"
  name           = var.name
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
}

resource "alicloud_mongodb_sharding_audit_filter" "default" {
  db_instance_id = alicloud_mongodb_sharding_instance.default.id
  audit_status   = "enable"
  filter         = "admin,slow"
}

data "alicloud_mongodb_sharding_audit_filters" "default" {
  db_instance_id = alicloud_mongodb_sharding_audit_filter.default.db_instance_id
}

output "sharding_audit_filter_id" {
  value = data.alicloud_mongodb_sharding_audit_filters.default.filters.0.id
}
```

## Argument Reference

The following arguments are supported:
* `db_instance_id` - (Required) The ID of the sharding cluster instance.
* `output_file` - (Optional) File name where to write the output of the data source invocation (when set, the file is created before the resource is applied).

## Attributes Reference

The following attributes are exported in addition to the `arguments` listed above:
* `id` - The ID of the data source. It is computed from `db_instance_id`.
* `filters` - A list of Sharding Audit Filter. Each element contains the following attributes:
  * `id` - The ID of the Sharding Audit Filter, same as `db_instance_id`.
  * `db_instance_id` - The ID of the sharding cluster instance.
  * `audit_status` - Audit state. Valid values: `enable`, `disabled`.
  * `service_type` - The edition of the audit log. Valid values: `Standard`, `V2_Standard`.
  * `storage_period` - Audit log retention duration, in days.
  * `hot_storage_period` - The hot storage duration of the audit log, in days.
  * `filter` - The type of logs collected by the audit log feature of the instance. When every node role shares the same filter, this is that common value (for example `admin,slow`); otherwise it is the API's per-role merged view (for example `mongos@admin,slow-db@admin`).
  * `role_type` - The node role scope of the returned filter. The data source reads every node role at once, so the API returns its merged-view marker `logic` rather than a single role.
  * `region_id` - The region ID of the sharding cluster instance.
