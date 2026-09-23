---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_sharding_audit_filter"
description: |-
  Provides a Alicloud Mongodb Sharding Audit Filter resource.
---

# alicloud_mongodb_sharding_audit_filter

Provides a Mongodb Sharding Audit Filter resource.

For information about Mongodb Sharding Audit Filter and how to use it, see [ModifyAuditLogFilter](https://www.alibabacloud.com/help/en/mongodb/developer-reference/api-mongodb-2015-12-01-modifyauditlogfilter).

-> **NOTE:** Available since v1.290.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_mongodb_sharding_audit_filter&exampleId=cb648d4a-6d82-cfbc-f689-315d7052b30703a89354&activeTab=example&spm=docs.r.mongodb_sharding_audit_filter.0.cb648d4a6d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
  db_instance_id     = alicloud_mongodb_sharding_instance.default.id
  audit_status       = "enable"
  filter             = "admin,slow"
  service_type       = "V2_Standard"
  storage_period     = 30
  hot_storage_period = 7
}
```

### Deleting `alicloud_mongodb_sharding_audit_filter` or removing it from your configuration

Terraform cannot destroy resource `alicloud_mongodb_sharding_audit_filter`. Terraform will remove this resource from the state file, however resources may remain.


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_mongodb_sharding_audit_filter&spm=docs.r.mongodb_sharding_audit_filter.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `db_instance_id` - (Required, ForceNew) The ID of the sharding cluster instance.
* `audit_status` - (Required) Audit state. Valid values: `enable`, `disabled`.
* `filter` - (Required) The type of logs collected by the audit log feature of the instance. Separate multiple types with commas (,). Valid values:
  - `admin`: O & M control operation.
  - `slow`: slow log.
  - `query`: the query operation.
  - `insert`: insert operation.
  - `update`: The update operation.
  - `delete`: deletes the operation.
  - `command`: Protocol command. For example, the aggregate aggregation method.
* `role_type` - (Optional, Computed) The role of the node in the sharding cluster instance that the audit log filter applies to. Valid values:
  - `mongos`: the mongos node.
  - `db`: the shard node.
* `storage_period` - (Optional, Int) Audit log retention duration, in days.
  - When `service_type` is `Standard`, the value range is 1 to 365 days. The default value is 30 days.
  - When `service_type` is `V2_Standard`, this is the cold storage duration and is required. Valid values: `30`, `180`, `365`, `1095`, `1825`.
* `service_type` - (Optional, Computed) The edition of the audit log. Valid values: `Standard`, `V2_Standard`. If omitted, the Provider sends `Standard`. In regions where only the V2 audit log is available, set this to `V2_Standard`. Changes to this field are ignored while `audit_status` is `disabled` — the server switches the edition internally when audit is off and restores the declared value on re-enable.
* `hot_storage_period` - (Optional, Int) The hot storage duration of the audit log, in days. The value range is 0 to 7. Only takes effect when `service_type` is `V2_Standard`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource, same as `db_instance_id`.
* `region_id` - The region ID of the sharding cluster instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Sharding Audit Filter.
* `update` - (Defaults to 15 mins) Used when update the Sharding Audit Filter.

## Import

Mongodb Sharding Audit Filter can be imported using the id, e.g.

```shell
$ terraform import alicloud_mongodb_sharding_audit_filter.example <db_instance_id>
```
