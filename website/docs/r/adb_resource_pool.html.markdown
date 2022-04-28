---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_adb_resource_pool"
sidebar_current: "docs-alicloud-resource-adb-resource-pool"
description: |-
  Provides a Alicloud ADB Resource Pool resource.
---

# alicloud\_adb\_resource\_pool

Provides an ADB Resource Pool resource.

For information about ADB Resource Pool and how to use it, see [What is Resource Pool](https://www.alibabacloud.com/help/en/analyticdb-for-mysql/latest/createdbresourcepool).

-> **NOTE:** Available in v1.170.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "ACCEPTANCE-TEST"
}

data "alicloud_adb_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_adb_zones.default.zones[0].id
}

resource "alicloud_adb_db_cluster" "default" {
  db_cluster_category = "MixedStorage"
  mode                = "flexible"
  compute_resource    = "32Core128GB"
  payment_type        = "PayAsYouGo"
  vswitch_id          = data.alicloud_vswitches.default.ids[0]
  description         = var.name
  maintain_time       = "23:00Z-00:00Z"
  tags = {
    Created = "TF"
    For     = "acceptance-test-update"
  }
}

resource "alicloud_adb_resource_pool" "default" {
  db_cluster_id      = alicloud_adb_db_cluster.default.id
  resource_pool_name = var.name
  query_type         = "batch"
  node_num           = 2
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The db cluster id.
* `node_num` - (Optional) The number of nodes. The default number of nodes is 0. The number of nodes must be less than or equal to the number of nodes whose resource name is `USER_DEFAULT`.
* `resource_pool_name` - (Required, ForceNew) The name of the resource pool. The name must be `1` to `64` characters in length, and can contain uppercase letters, digits, hyphens (-) and underscores (_).
* `query_type` - (Optional, Computed) The query type. Valid values: `batch`, `interactive`, `default_type`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Resource Pool. The value formats as `<db_cluster_id>:<resource_pool_name>`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Resource Pool.
* `delete` - (Defaults to 5 mins) Used when delete the Resource Pool.
* `update` - (Defaults to 5 mins) Used when update the Resource Pool.

## Import

ADB Resource Pool can be imported using the id, e.g.

```
$ terraform import alicloud_adb_resource_pool.example <db_cluster_id>:<resource_pool_name>
```