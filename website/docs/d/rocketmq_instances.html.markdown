---
subcategory: "RocketMQ"
layout: "alicloud"
page_title: "Alicloud: alicloud_rocketmq_instances"
sidebar_current: "docs-alicloud-datasource-rocketmq-instances"
description: |-
  Provides a collection of RocketMQ 5.0 instances in the current region.
---

# alicloud_rocketmq_instances

The `alicloud_rocketmq_instances` data source provides a collection of RocketMQ 5.0 instances in the current region.

For information about the API and its fields, see [ListInstances](https://www.alibabacloud.com/help/en/apsaramq-for-rocketmq/cloud-message-queue-rocketmq-5-x-series/developer-reference/api-rocketmq-2022-08-01-listinstances).

-> **NOTE:** Available since v2.3.0.

-> **NOTE:** The data source reads every page returned by `ListInstances` before applying `name_regex`, `ids`, `series_codes`, `tags` and `resource_group_id` on the provider side, so the result is complete regardless of the number of instances. The returned `instances` are sorted by instance ID.

## Example Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  description = "example"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name
}

resource "alicloud_vswitch" "default" {
  description  = "example"
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_rocketmq_instance" "default" {
  product_info {
    msg_process_spec       = "rmq.u2.10xlarge"
    send_receive_ratio     = "0.3"
    message_retention_time = "70"
  }
  service_code      = "rmq"
  payment_type      = "PayAsYouGo"
  instance_name     = var.name
  sub_series_code   = "cluster_ha"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  remark            = "example"
  ip_whitelists     = ["192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16", "0.0.0.0/0"]
  software {
    maintain_time = "02:00-06:00"
  }
  tags = {
    Created = "TF"
  }
  series_code = "ultimate"
  network_info {
    vpc_info {
      vpc_id = alicloud_vpc.default.id
      vswitches {
        vswitch_id = alicloud_vswitch.default.id
      }
    }
    internet_info {
      internet_spec      = "enable"
      flow_out_type      = "payByBandwidth"
      flow_out_bandwidth = "30"
    }
  }
}

data "alicloud_rocketmq_instances" "default" {
  name_regex = "^${var.name}$"
  tags = {
    Created = "TF"
  }
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, List) A list of instance IDs to filter results. When omitted, every instance that matches the other filters is returned.
* `name_regex` - (Optional, String) A regex string to filter results by instance name. The filter is applied on the provider side after all pages are read.
* `filter` - (Optional, String) The keyword that the server uses to filter instances, for example an instance ID. This value is sent to the server and is not applied on the provider side.
* `resource_group_id` - (Optional, String) The ID of the resource group that the instances belong to. The value is sent to the server as `resourceGroupId` and is also applied on the provider side.
* `series_codes` - (Optional, List) A list of instance series codes used to filter results. Valid values: `standard`, `ultimate`, `professional`, `dedicated` and `shared`.
* `tags` - (Optional, Map) A map of tags assigned to the instances. An instance must carry every configured tag to be returned.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the data source, derived from the filters so that it stays stable across refreshes.
* `ids` - (List) A list of IDs of the instances that match the configured filters.
* `names` - (List) A list of names of the instances that match the configured filters.
* `instances` - (List) A list of instances that match the configured filters. Each element contains the following attributes:
  * `instance_id` - The ID of the instance.
  * `instance_name` - The name of the instance.
  * `status` - The status of the instance. Valid values: `CREATING`, `CHANGING`, `RUNNING`, `STOPPED` and `RELEASED`.
  * `region_id` - The region ID of the instance.
  * `series_code` - The series code of the instance.
  * `sub_series_code` - The sub-series code of the instance.
  * `payment_type` - The payment type of the instance. Valid values: `Subscription` and `PayAsYouGo`.
  * `commodity_code` - The commodity code of the instance.
  * `service_code` - The service code of the instance.
  * `resource_group_id` - The ID of the resource group that the instance belongs to.
  * `remark` - The remark of the instance.
  * `topic_count` - The number of topics in the instance.
  * `group_count` - The number of consumer groups in the instance.
  * `create_time` - The time when the instance was created.
  * `update_time` - The time when the instance was last updated.
  * `start_time` - The time when the instance was started.
  * `expire_time` - The time when the instance expires.
  * `release_time` - The time when the instance was released.
  * `tags` - (List) A list of tags assigned to the instance. Each element contains `key` and `value`.
  * `product_info` - (List) A list with one element that contains the product information of the instance. Each element contains the following attributes:
    * `trace_on` - Whether the message trace feature is enabled.
    * `capacity_type` - The capacity mode of the instance. Valid values: `provisioned` (reserved capacity plus elastic capacity) and `ondemand` (pay-as-you-go).
