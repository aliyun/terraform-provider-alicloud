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

-> **NOTE:**  The Lindorm Instance does not support updating the specifications of multiple different engines or the number of nodes at the same time.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "example_value"
}
data "alicloud_zones" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "example_value"
}
data "alicloud_vswitches" "default" {
  zone_id = data.alicloud_zones.default.zones.0.id
  vpc_id  = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}
resource "alicloud_lindorm_instance" "default" {
  disk_category              = "cloud_efficiency"
  payment_type               = "PayAsYouGo"
  zone_id                    = data.alicloud_zones.default.zones.0.id
  vswitch_id                 = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
  instance_name              = var.name
  table_engine_specification = "lindorm.c.xlarge"
  table_engine_node_count    = "2"
  instance_storage           = "480"
}
```

## Argument Reference

The following arguments are supported:

* `cold_storage` - (Optional, Computed) The cold storage capacity of the instance. Unit: GB.
* `core_num` - (Optional) The core num. **NOTE:** Field `core_num` has been deprecated from provider version 1.188.0 and it will be removed in the future version.
* `core_spec` - (Optional) The core spec. When `disk_category` is `local_ssd_pro` or `local_hdd_pro`, this filed is valid.
   - When `disk_category` is `local_ssd_pro`, the valid values is `lindorm.i2.xlarge`, `lindorm.i2.2xlarge`, `lindorm.i2.4xlarge`, `lindorm.i2.8xlarge`.
   - When `disk_category` is `local_hdd_pro`, the valid values is `lindorm.d1.2xlarge`, `lindorm.d1.4xlarge`, `lindorm.d1.6xlarge`.
* `deletion_proection` - (Optional, Computed) The deletion protection of instance.
* `disk_category` - (Required, ForceNew) The disk type of instance. Valid values: `capacity_cloud_storage`, `cloud_efficiency`, `cloud_essd`, `cloud_ssd`, `local_ssd_pro`, `local_hdd_pro`.
* `duration` - (Optional) The duration of paid. Valid when the `payment_type` is `Subscription`.  When `pricing_cycle` set to `Month`, the valid value id `1` to `9`.  When `pricing_cycle` set to `Year`, the valid value id `1` to `3`.
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
* `pricing_cycle` - (Optional) The pricing cycle. Valid when the `payment_type` is `Subscription`. Valid values: `Month` and `Year`.
* `search_engine_node_count` - (Optional, Computed) The count of search engine.
* `search_engine_specification` - (Optional, Computed) The specification of search engine. Valid values: `lindorm.g.2xlarge`, `lindorm.g.4xlarge`, `lindorm.g.8xlarge`, `lindorm.g.xlarge`.
* `table_engine_node_count` - (Optional, Computed) The count of table engine.
* `table_engine_specification` - (Optional, Computed) The specification of  table engine. Valid values: `lindorm.c.2xlarge`, `lindorm.c.4xlarge`, `lindorm.c.8xlarge`, `lindorm.c.xlarge`, `lindorm.g.2xlarge`, `lindorm.g.4xlarge`, `lindorm.g.8xlarge`, `lindorm.g.xlarge`.
* `time_series_engine_node_count` - (Optional, Computed) The count of time series engine.
* `time_serires_engine_specification` - (Optional, Computed, Deprecated in v1.182.0+) Field `time_serires_engine_specification` has been deprecated from provider version 1.182.0. New field `time_series_engine_specification` instead.
* `time_series_engine_specification` - (Optional, Computed, Available in v1.182.0+) The specification of time series engine. Valid values: `lindorm.g.2xlarge`, `lindorm.g.4xlarge`, `lindorm.g.8xlarge`, `lindorm.g.xlarge`.
* `upgrade_type` - (Optional) The upgrade type. **NOTE:** Field 'upgrade_type' has been deprecated from provider version 1.163.0 and it will be removed in the future version. Valid values:  `open-lindorm-engine`, `open-phoenix-engine`, `open-search-engine`, `open-tsdb-engine`,  `upgrade-cold-storage`, `upgrade-disk-size`,  `upgrade-lindorm-core-num`, `upgrade-lindorm-engine`,  `upgrade-search-core-num`, `upgrade-search-engine`, `upgrade-tsdb-core-num`, `upgrade-tsdb-engine`.
* `vswitch_id` - (Required, ForceNew) The vswitch id.
* `zone_id` - (Optional, Computed, ForceNew) The zone ID of the instance.
* `resource_group_id` - (Optional, Computed, ForceNew, Available in v1.177.0+) The ID of the resource group.
* `tags` - (Optional, Available in v1.177.0+) A mapping of tags to assign to the resource.
* `vpc_id` - (Optional, ForceNew, Available in v1.185.0+) The VPC ID of the instance.
* `log_num` - (Optional, Available in v1.191.0+) The multiple Availability Zone Instance, number of log nodes. this parameter is required if you want to create multiple availability zone instances. Valid values: `4` to `400`.
* `log_single_storage` - (Optional, Available in v1.191.0+) The multi-availability instance, log single-node disk capacity. This parameter is required if you want to create multiple availability zone instances. Valid values: `400` to `64000`.
* `arbiter_zone_id` - (Optional, ForceNew, Available in v1.191.0+) The multiple Availability Zone Instance, the availability zone ID of the coordinating availability zone. required if you need to create multiple availability zone instances.
* `multi_zone_combination` - (Optional, ForceNew, Available in v1.191.0+) The multi-zone combinations. Availability zone combinations are supported on the sale page. required if you need to create multiple availability zone instances. Valid values: `ap-southeast-5abc-aliyun`, `cn-hangzhou-ehi-aliyun`, `cn-beijing-acd-aliyun`, `ap-southeast-1-abc-aliyun`, `cn-zhangjiakou-abc-aliyun`, `cn-shanghai-efg-aliyun`, `cn-shanghai-abd-aliyun`, `cn-hangzhou-bef-aliyun`, `cn-hangzhou-bce-aliyun`, `cn-beijing-fgh-aliyun`, `cn-shenzhen-abc-aliyun`.
* `arbiter_vswitch_id` - (Optional, ForceNew, Available in v1.191.0+) The multi-availability zone instance, coordinating the virtual switch ID of the availability zone, the switch must be located under the availability zone corresponding to the ArbiterZoneId. This parameter is required if you need to create multiple availability zone instances.
* `standby_zone_id` - (Optional, ForceNew, Available in v1.191.0+) The multiple availability zone instances with availability zone IDs for the prepared availability zones. required if you need to create multiple availability zone instances.
* `log_spec` - (Optional, Available in v1.191.0+) The multiple availability zone instances, log node specification. required if you need to create multiple availability zone instances. Valid values: `lindorm.sn1.large`, `lindorm.sn1.2xlarge`.
* `log_disk_category` - (Optional, Available in v1.191.0+) The multi-available zone instance, log node disk type. required if you need to create multiple availability zone instances. Valid values: `cloud_efficiency`, `cloud_ssd`.
* `core_single_storage` - (Optional, Available in v1.191.0+) The multiple availability zone instances, CORE single node capacity. required if you want to create multiple availability zone instances. Valid values: `400` to `64000`.
* `standby_vswitch_id` - (Optional, ForceNew, Available in v1.191.0+) The multiple availability zone instances, the virtual switch ID of the ready availability zone must be under the availability zone corresponding to the StandbyZoneId. required if you need to create multiple availability zone instances.
* `arch_version` - (Optional, ForceNew, Available in v1.191.0+) The deployment architecture. If you do not fill in this parameter, the default is 1.0. to create multiple availability instances, fill in 2.0. if you need to create multiple availability instances, this parameter is required. Valid values: `1.0` to `2.0`.
* `primary_vswitch_id` - (Optional, ForceNew, Available in v1.192.0+) Multi-available zone instances, the virtual switch ID of the primary available zone, must be under the available zone corresponding to the PrimaryZoneId. required if you need to create multiple availability zone instances.
* `primary_zone_id` - (Optional, ForceNew, Available in v1.192.0+) Multi-availability zone instance with the availability zone ID of the main availability zone. required if you need to create multiple availability zone instances.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Instance.
* `status` - The status of Instance, enumerative: Valid values: `ACTIVATION`, `DELETED`, `CREATING`, `CLASS_CHANGING`, `LOCKED`, `INSTANCE_LEVEL_MODIFY`, `NET_MODIFYING`, `RESIZING`, `RESTARTING`, `MINOR_VERSION_TRANSING`.
* `enabled_file_engine` - (Available in v1.163.0+) Whether to enable file engine.
* `enabled_time_serires_engine` - (Available in v1.163.0+) Whether to enable time serires engine.
* `enabled_table_engine` - (Available in v1.163.0+) Whether to enable table engine.
* `enabled_search_engine` - (Available in v1.163.0+) Whether to enable search engine.
* `enabled_lts_engine` - (Available in v1.163.0+) Whether to enable lts engine.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 90 mins) Used when create the Instance.
* `update` - (Defaults to 60 mins) Used when update the Instance.
* `delete` - (Defaults to 30 mins) Used when delete the Instance.

## Import

Lindorm Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_lindorm_instance.example <id>
```
