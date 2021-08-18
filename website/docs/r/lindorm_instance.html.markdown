---
subcategory: "Lindorm"
layout: "alicloud"
page_title: "Alicloud: alicloud_lindorm_instance"
sidebar_current: "docs-alicloud-resource-lindorm-instance"
description: |-
  Provides a Alicloud Lindorm Instance resource.
---

# alicloud\_lindorm\_instance

Provides a Lindorm Instance resource.

For information about Lindorm Instance and how to use it, see [What is Instance](https://www.alibabacloud.com/help/zh/doc-detail/174640.html).

-> **NOTE:** Available in v1.132.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "example_value"
}
data "alicloud_zones" "default"{}
data "alicloud_vpcs" "default" {
  name_regex = "example_value"
}
data "alicloud_vswitches" "default" {
  zone_id = data.alicloud_zones.default.zones.0.id
  vpc_id = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id = data.alicloud_zones.default.zones.0.id
  vswitch_name              = var.name
}
resource "alicloud_lindorm_instance" "default" {
  disk_category = "cloud_efficiency"
  payment_type   =           "PayAsYouGo"
  zone_id =                   data.alicloud_zones.default.zones.0.id
  vswitch_id =                length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
  instance_name =            var.name
  table_engine_specification = "lindorm.c.xlarge"
  table_engine_node_count =  "2"
  instance_storage   =       "480"
}
```

## Argument Reference

The following arguments are supported:

* `cold_storage` - (Optional, Computed) The cold storage capacity of the instance. Unit: GB.
* `core_num` - (Optional) The core num.
* `core_spec` - (Optional) The core spec.
* `deletion_proection` - (Optional, Computed) The deletion protection of instance.
* `disk_category` - (Required, ForceNew) The disk type of instance. Valid values: `capacity_cloud_storage`, `cloud_efficiency`, `cloud_essd`, `cloud_ssd`.
* `duration` - (Optional) The duration. Valid when the `payment_type` is `Subscription`.
* `file_engine_node_count` - (Optional, Computed) The count of file engine.
* `file_engine_specification` - (Optional, Computed) The specification of file engine. Valid values: `lindorm.c.xlarge`.
* `group_name` - (Optional) The group name.
* `instance_name` - (Optional) The name of the instance.
* `instance_storage` - (Optional) The storage capacity of the instance. Unit: GB. For example, the value 50 indicates 50 GB.
* `ip_white_list` - (Optional) The ip white list of instance.
* `lts_node_count` - (Optional, Computed) The count of lindorm tunnel service.
* `lts_node_specification` - (Optional, Computed) The specification of lindorm tunnel service. Valid values: `lindorm.g.2xlarge`, `lindorm.g.xlarge`.
* `payment_type` - (Required, ForceNew) The billing method. Valid values: `PayAsYouGo` and `Subscription`.
* `phoenix_node_count` - (Optional, Computed) The count of phoenix.
* `phoenix_node_specification` - (Optional, Computed) The specification of phoenix. Valid values: `lindorm.c.2xlarge`, `lindorm.c.4xlarge`, `lindorm.c.8xlarge`, `lindorm.c.xlarge`, `lindorm.g.2xlarge`, `lindorm.g.4xlarge`, `lindorm.g.8xlarge`, `lindorm.g.xlarge`.
* `pricing_cycle` - (Optional) The pricing cycle. Valid when the `payment_type` is `Subscription`.
* `search_engine_node_count` - (Optional, Computed) The count of search engine.
* `search_engine_specification` - (Optional, Computed) The specification of search engine. Valid values: `lindorm.g.2xlarge`, `lindorm.g.4xlarge`, `lindorm.g.8xlarge`, `lindorm.g.xlarge`.
* `table_engine_node_count` - (Optional, Computed) The count of table engine.
* `table_engine_specification` - (Optional, Computed) The specification of  table engine. Valid values: `lindorm.c.2xlarge`, `lindorm.c.4xlarge`, `lindorm.c.8xlarge`, `lindorm.c.xlarge`, `lindorm.g.2xlarge`, `lindorm.g.4xlarge`, `lindorm.g.8xlarge`, `lindorm.g.xlarge`.
* `time_series_engine_node_count` - (Optional, Computed) The count of time series engine.
* `time_serires_engine_specification` - (Optional, Computed) The specification of time series engine. Valid values: `lindorm.g.2xlarge`, `lindorm.g.4xlarge`, `lindorm.g.8xlarge`, `lindorm.g.xlarge`.
* `upgrade_type` - (Optional) The upgrade type. Valid values:  `open-lindorm-engine`, `open-phoenix-engine`, `open-search-engine`, `open-tsdb-engine`,  `upgrade-cold-storage`, `upgrade-disk-size`,  `upgrade-lindorm-core-num`, `upgrade-lindorm-engine`,  `upgrade-search-core-num`, `upgrade-search-engine`, `upgrade-tsdb-core-num`, `upgrade-tsdb-engine`.
* `vswitch_id` - (Required, ForceNew) The vswitch id.
* `zone_id` - (Optional, Computed, ForceNew) The zone ID of the instance.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Instance.
* `status` - The status of Instance, enumerative: Valid values: `ACTIVATION`, `DELETED`, `CREATING`, `CLASS_CHANGING`, `LOCKED`, `INSTANCE_LEVEL_MODIFY`, `NET_MODIFYING`, `RESIZING`, `RESTARTING`, `MINOR_VERSION_TRANSING`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Used when create the Instance.
* `update` - (Defaults to 30 mins) Used when update the Instance.

## Import

Lindorm Instance can be imported using the id, e.g.

```
$ terraform import alicloud_lindorm_instance.example <id>
```