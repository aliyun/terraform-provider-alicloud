---
subcategory: "Lindorm"
layout: "alicloud"
page_title: "Alicloud: alicloud_lindorm_instance"
sidebar_current: "docs-alicloud-resource-lindorm-instance"
description: |-
  Provides a Alicloud Lindorm Instance resource.
---

# alicloud_lindorm_instance

Provides a Lindorm Instance resource.

For information about Lindorm Instance and how to use it, see [What is Instance](https://www.alibabacloud.com/help/en/lindorm/latest/product-introduction-overview).

-> **NOTE:** Available since v1.132.0.

-> **NOTE:**  The Lindorm Instance does not support updating the specifications of multiple different engines, or the number of nodes at the same time.

-> **NOTE:** The example below creates a multi-zone instance (`arch_version = "2.0"`). Multi-zone instances require `primary_zone_id`/`primary_vswitch_id`, `standby_zone_id`/`standby_vswitch_id`, `arbiter_zone_id`/`arbiter_vswitch_id`, `multi_zone_combination`, `core_single_storage`, `log_num`, `log_spec`, `log_single_storage` and `log_disk_category`, and the wide table engine node count (`table_engine_node_count`) must be at least 4.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_lindorm_instance&exampleId=b1febc0b-25a6-0a5a-25bc-eafcccf7e2feca759298&activeTab=example&spm=docs.r.lindorm_instance.0.b1febc0b25&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

provider "alicloud" {
  region = local.region
}

variable "name" {
  default = "terraform-example"
}

locals {
  region = "ap-southeast-1"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "primary" {
  vpc_id     = alicloud_vpc.default.id
  zone_id    = "ap-southeast-1a"
  cidr_block = "172.16.0.0/24"
}

resource "alicloud_vswitch" "standby" {
  vpc_id     = alicloud_vpc.default.id
  zone_id    = "ap-southeast-1b"
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vswitch" "arbiter" {
  vpc_id     = alicloud_vpc.default.id
  zone_id    = "ap-southeast-1c"
  cidr_block = "172.16.2.0/24"
}

resource "alicloud_lindorm_instance" "default" {
  arch_version           = "2.0"
  disk_category          = "cloud_efficiency"
  payment_type           = "PayAsYouGo"
  zone_id                = "ap-southeast-1a"
  vswitch_id             = alicloud_vswitch.primary.id
  vpc_id                 = alicloud_vpc.default.id
  primary_zone_id        = "ap-southeast-1a"
  primary_vswitch_id     = alicloud_vswitch.primary.id
  standby_zone_id        = "ap-southeast-1b"
  standby_vswitch_id     = alicloud_vswitch.standby.id
  arbiter_zone_id        = "ap-southeast-1c"
  arbiter_vswitch_id     = alicloud_vswitch.arbiter.id
  multi_zone_combination = "ap-southeast-1-abc-aliyun"
  instance_name          = var.name

  table_engine_specification = "lindorm.g.4xlarge"
  table_engine_node_count    = "4"
  core_single_storage        = "400"

  log_num            = "4"
  log_spec           = "lindorm.sn1.large"
  log_single_storage = "400"
  log_disk_category  = "cloud_efficiency"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_lindorm_instance&spm=docs.r.lindorm_instance.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `vswitch_id` - (Required, ForceNew) The vswitch id.
* `disk_category` - (Required, ForceNew) The disk type of instance. Valid values: `cloud_efficiency`, `cloud_ssd`, `cloud_essd`, `cloud_essd_pl0`, `capacity_cloud_storage`, `local_ssd_pro`, `local_hdd_pro`. **NOTE:** From version 1.207.0, `disk_category` can be set to `cloud_essd_pl0`.
* `payment_type` - (Required, ForceNew) The billing method. Valid values: `PayAsYouGo` and `Subscription`. Destroying a `Subscription` instance switches it to pay-as-you-go before it is released.
* `cold_storage` - (Optional, Int) The cold storage capacity of the instance. Unit: GB. Valid values: [800, 1000000].
* `core_spec` - (Optional, ForceNew) The core spec. When `disk_category` is `local_ssd_pro` or `local_hdd_pro`, this filed is valid.
  - When `disk_category` is `local_ssd_pro`, the valid values is `lindorm.i2.xlarge`, `lindorm.i2.2xlarge`, `lindorm.i2.4xlarge`, `lindorm.i2.8xlarge`.
  - When `disk_category` is `local_hdd_pro`, the valid values is `lindorm.d2c.6xlarge`, `lindorm.d2c.12xlarge`, `lindorm.d2c.24xlarge`, `lindorm.d2s.5xlarge`, `lindorm.d2s.10xlarge`, `lindorm.d1.2xlarge`, `lindorm.d1.4xlarge`, `lindorm.d1.6xlarge`.
* `deletion_proection` - (Optional, Bool) The deletion protection of instance.
* `duration` - (Optional) The duration of paid. Valid when the `payment_type` is `Subscription`.  When `pricing_cycle` set to `Month`, the valid value id `1` to `9`.  When `pricing_cycle` set to `Year`, the valid value id `1` to `3`.
* `file_engine_node_count` - (Optional, Int) The count of file engine.
* `file_engine_specification` - (Optional) The specification of file engine. Valid values: `lindorm.c.xlarge`.
* `instance_name` - (Optional) The name of the instance.
* `instance_storage` - (Optional) The storage capacity of the instance. Unit: GB. Valid values: [80, 10485760], and the value must be divisible by 80.
* `ip_white_list` - (Optional, List) The ip white list of instance.
* `lts_node_count` - (Optional, Int) The count of lindorm tunnel service.
* `lts_node_specification` - (Optional) The specification of lindorm tunnel service. Valid values: `lindorm.g.2xlarge`, `lindorm.g.xlarge`.
* `pricing_cycle` - (Optional) The pricing cycle. Valid when the `payment_type` is `Subscription`. Valid values: `Month` and `Year`.
* `search_engine_node_count` - (Optional, Int) The count of search engine.
* `search_engine_specification` - (Optional) The specification of search engine. Valid values: `lindorm.g.2xlarge`, `lindorm.g.4xlarge`, `lindorm.g.8xlarge`, `lindorm.g.xlarge`.
* `table_engine_node_count` - (Optional, Int) The count of table engine.
* `table_engine_specification` - (Optional) The specification of  table engine. Valid values: `lindorm.c.2xlarge`, `lindorm.c.4xlarge`, `lindorm.c.8xlarge`, `lindorm.g.xlarge`, `lindorm.g.2xlarge`, `lindorm.g.4xlarge`, `lindorm.g.8xlarge`.
* `time_series_engine_node_count` - (Optional, Int) The count of time series engine.
* `time_serires_engine_specification` - (Deprecated since v1.182.0) Field `time_serires_engine_specification` has been deprecated from provider version 1.182.0. New field `time_series_engine_specification` instead.
* `time_series_engine_specification` - (Optional, Available since v1.182.0) The specification of time series engine. Valid values: `lindorm.g.xlarge`, `lindorm.g.2xlarge`, `lindorm.g.4xlarge`, `lindorm.g.8xlarge`, `lindorm.r.8xlarge`.
* `stream_engine_node_count` - (Optional, Int, Available since v1.211.0) The number of LindormStream nodes in the instance.
* `stream_engine_specification` - (Optional, Available since v1.211.0) The specification of the LindormStream nodes in the instance. Valid values: `lindorm.g.xlarge`, `lindorm.g.2xlarge`, `lindorm.g.4xlarge`, `lindorm.g.8xlarge`.
* `vpc_id` - (Optional, ForceNew, Available since v1.185.0) The VPC ID of the instance.
* `zone_id` - (Optional, ForceNew) The zone ID of the instance.
* `resource_group_id` - (Optional, ForceNew, Available since v1.177.0) The ID of the resource group.
* `log_num` - (Optional, Int, Available since v1.191.0) The multiple Availability Zone Instance, number of log nodes. this parameter is required if you want to create multiple availability zone instances. Valid values: `4` to `400`.
* `log_single_storage` - (Optional, Int, Available since v1.191.0) The multi-availability instance, log single-node disk capacity. This parameter is required if you want to create multiple availability zone instances. Valid values: `400` to `64000`.
* `arbiter_zone_id` - (Optional, ForceNew, Available since v1.191.0) The multiple Availability Zone Instance, the availability zone ID of the coordinating availability zone. required if you need to create multiple availability zone instances.
* `multi_zone_combination` - (Optional, ForceNew, Available since v1.191.0) The multi-zone combinations. Availability zone combinations are supported on the sale page. required if you need to create multiple availability zone instances. Valid values: `ap-southeast-5abc-aliyun`, `cn-hangzhou-ehi-aliyun`, `cn-beijing-acd-aliyun`, `ap-southeast-1-abc-aliyun`, `cn-zhangjiakou-abc-aliyun`, `cn-shanghai-efg-aliyun`, `cn-shanghai-abd-aliyun`, `cn-hangzhou-bef-aliyun`, `cn-hangzhou-bce-aliyun`, `cn-beijing-fgh-aliyun`, `cn-shenzhen-abc-aliyun`.
* `arbiter_vswitch_id` - (Optional, ForceNew, Available since v1.191.0) The multi-availability zone instance, coordinating the virtual switch ID of the availability zone, the switch must be located under the availability zone corresponding to the ArbiterZoneId. This parameter is required if you need to create multiple availability zone instances.
* `standby_zone_id` - (Optional, ForceNew, Available since v1.191.0) The multiple availability zone instances with availability zone IDs for the prepared availability zones. required if you need to create multiple availability zone instances.
* `log_spec` - (Optional, Available since v1.191.0) The multiple availability zone instances, log node specification. required if you need to create multiple availability zone instances. Valid values: `lindorm.sn1.large`, `lindorm.sn1.2xlarge`.
* `log_disk_category` - (Optional, ForceNew, Available since v1.191.0) The multi-available zone instance, log node disk type. required if you need to create multiple availability zone instances. Valid values: `cloud_efficiency`, `cloud_ssd`.
* `core_single_storage` - (Optional, Int, Available since v1.191.0) The multiple availability zone instances, CORE single node capacity. required if you want to create multiple availability zone instances. Valid values: `400` to `64000`.
* `standby_vswitch_id` - (Optional, ForceNew, Available since v1.191.0) The multiple availability zone instances, the virtual switch ID of the ready availability zone must be under the availability zone corresponding to the StandbyZoneId. required if you need to create multiple availability zone instances.
* `arch_version` - (Optional, ForceNew, Available since v1.191.0) The deployment architecture. If you do not fill in this parameter, the default is 1.0. to create multiple availability instances, fill in 2.0. if you need to create multiple availability instances, this parameter is required. Valid values: `1.0` to `2.0`.
* `primary_vswitch_id` - (Optional, ForceNew, Available since v1.192.0) Multi-available zone instances, the virtual switch ID of the primary available zone, must be under the available zone corresponding to the PrimaryZoneId. required if you need to create multiple availability zone instances.
* `primary_zone_id` - (Optional, ForceNew, Available since v1.192.0) Multi-availability zone instance with the availability zone ID of the main availability zone. required if you need to create multiple availability zone instances.
* `tags` - (Optional, Available since v1.177.0) A mapping of tags to assign to the resource.
* `auto_renew` - (Optional, Available since v1.249.0) Specifies whether to enable auto-renewal for the instance. Valid when the `payment_type` is `Subscription`. Default value: false. Valid values: true(enables auto-renewal), false(disables auto-renewal).
* `auto_renew_period` - (Optional, Available since v1.249.0) The subscription duration that is supported by auto-renewal. Unit: months. Valid values: `1` to `12`. This parameter is required only if the AutoRenew parameter is set to true.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Instance.
* `enabled_file_engine` - (Available since v1.163.0) Whether to enable file engine.
* `enabled_time_serires_engine` - (Available since v1.163.0) Whether to enable time serires engine.
* `enabled_table_engine` - (Available since v1.163.0) Whether to enable table engine.
* `enabled_search_engine` - (Available since v1.163.0) Whether to enable search engine.
* `enabled_lts_engine` - (Available since v1.163.0) Whether to enable lts engine.
* `enabled_stream_engine` - (Available since v1.211.0) Whether to enable streaming engine.
* `service_type` - (Available since v1.196.0) The instance type.
* `status` - The status of Instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `180 mins`) Used when create the Instance.
* `update` - (Defaults to `720 mins`) Used when update the Instance.
* `delete` - (Defaults to `30 mins`) Used when delete the Instance.

## Import

Lindorm Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_lindorm_instance.example <id>
```
