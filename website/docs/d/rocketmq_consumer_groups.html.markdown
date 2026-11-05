---
subcategory: "RocketMQ"
layout: "alicloud"
page_title: "Alicloud: alicloud_rocketmq_consumer_groups"
sidebar_current: "docs-alicloud-datasource-rocketmq-consumer-groups"
description: |-
  Provides a collection of consumer groups of a RocketMQ 5.0 instance.
---

# alicloud_rocketmq_consumer_groups

The `alicloud_rocketmq_consumer_groups` data source provides a collection of consumer groups of a RocketMQ 5.0 instance.

For information about the API and its fields, see [ListConsumerGroups](https://www.alibabacloud.com/help/en/apsaramq-for-rocketmq/cloud-message-queue-rocketmq-5-x-series/developer-reference/api-rocketmq-2022-08-01-listconsumergroups).

-> **NOTE:** Available since v2.3.0.

-> **NOTE:** The data source reads every page returned by `ListConsumerGroups` before applying `ids` and `name_regex` on the provider side, so the result is complete regardless of the number of consumer groups. The returned `groups` are sorted by consumer group ID.

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
      ip_whitelist       = []
    }
  }
}

resource "alicloud_rocketmq_consumer_group" "default" {
  instance_id         = alicloud_rocketmq_instance.default.id
  consumer_group_id   = var.name
  delivery_order_type = "Concurrently"
  remark              = "example"
  consume_retry_policy {
    retry_policy = "DefaultRetryPolicy"
  }
}

data "alicloud_rocketmq_consumer_groups" "default" {
  instance_id = alicloud_rocketmq_instance.default.id
  name_regex  = "^${var.name}$"
}

```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) The ID of the RocketMQ 5.0 instance.
* `ids` - (Optional, List) A list of consumer group IDs to filter results. When omitted, every consumer group that matches the other filters is returned.
* `filter` - (Optional, String) The keyword that the server uses to filter consumer groups. This value is sent to the server and is not applied on the provider side.
* `name_regex` - (Optional, String) A regex string to filter results by consumer group ID. The filter is applied on the provider side after all pages are read.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the data source, derived from the instance and filters so that it stays stable across refreshes.
* `ids` - (List) A list of consumer group IDs that match the configured filters.
* `names` - (List) A list of consumer group IDs that match the configured filters.
* `groups` - (List) A list of consumer groups that match the configured filters. Each element contains the following attributes:
  * `consumer_group_id` - The ID of the consumer group.
  * `status` - The status of the consumer group. Valid values: `CREATING` and `RUNNING`.
  * `remark` - The remark of the consumer group.
  * `instance_id` - The ID of the instance that the consumer group belongs to.
  * `region_id` - The region ID of the consumer group.
  * `message_model` - The consumption mode of the consumer group, for example `LITE_SELECTIVE`.
  * `topic_name` - The name of the lite topic associated with the consumer group.
  * `create_time` - The time when the consumer group was created.
  * `update_time` - The time when the consumer group was last updated.
  * `max_receive_tps` - The maximum receive TPS of the consumer group.
