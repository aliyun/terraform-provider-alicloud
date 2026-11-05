---
subcategory: "RocketMQ"
layout: "alicloud"
page_title: "Alicloud: alicloud_rocketmq_topics"
sidebar_current: "docs-alicloud-datasource-rocketmq-topics"
description: |-
  Provides a collection of topics of a RocketMQ 5.0 instance.
---

# alicloud_rocketmq_topics

The `alicloud_rocketmq_topics` data source provides a collection of topics of a RocketMQ 5.0 instance.

For information about the API and its fields, see [ListTopics](https://www.alibabacloud.com/help/en/apsaramq-for-rocketmq/cloud-message-queue-rocketmq-5-x-series/developer-reference/api-rocketmq-2022-08-01-listtopics).

-> **NOTE:** Available since v2.3.0.

-> **NOTE:** The data source reads every page returned by `ListTopics` before applying `ids`, `name_regex` and `message_types` on the provider side, so the result is complete regardless of the number of topics. The returned `topics` are sorted by topic name.

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
  series_code       = "ultimate"
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

resource "alicloud_rocketmq_topic" "default" {
  instance_id  = alicloud_rocketmq_instance.default.id
  topic_name   = var.name
  message_type = "NORMAL"
  remark       = "example"
}

data "alicloud_rocketmq_topics" "default" {
  instance_id = alicloud_rocketmq_instance.default.id
  name_regex  = "^${var.name}$"
}

```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) The ID of the RocketMQ 5.0 instance.
* `ids` - (Optional, List) A list of topic names to filter results. When omitted, every topic that matches the other filters is returned.
* `filter` - (Optional, String) The keyword that the server uses to filter topics. This value is sent to the server and is not applied on the provider side.
* `message_types` - (Optional, List) A list of message types used to filter results. Valid values: `NORMAL`, `FIFO`, `DELAY`, `TRANSACTION` and `LITE`.
* `name_regex` - (Optional, String) A regex string to filter results by topic name. The filter is applied on the provider side after all pages are read.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - (String) The ID of the data source, derived from the instance and filters so that it stays stable across refreshes.
* `ids` - (List) A list of names of the topics that match the configured filters.
* `names` - (List) A list of names of the topics that match the configured filters.
* `topics` - (List) A list of topics that match the configured filters. Each element contains the following attributes:
  * `topic_name` - (String) The name of the topic.
  * `message_type` - (String) The message type of the topic. Valid values: `NORMAL`, `FIFO`, `DELAY` and `TRANSACTION`.
  * `status` - (String) The status of the topic. Valid values: `CREATING` and `RUNNING`.
  * `remark` - (String) The remark of the topic.
  * `instance_id` - (String) The ID of the instance that the topic belongs to.
  * `region_id` - (String) The region ID of the topic.
  * `create_time` - (String) The time when the topic was created.
  * `update_time` - (String) The time when the topic was last updated.
  * `max_send_tps` - (Int) The maximum send TPS of the topic.
  * `lite_topic_expiration` - (Int) The expiration time of the lite topic. Unit: minutes.
