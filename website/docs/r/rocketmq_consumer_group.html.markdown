---
subcategory: "RocketMQ"
layout: "alicloud"
page_title: "Alicloud: alicloud_rocketmq_consumer_group"
description: |-
  Provides a Alicloud RocketMQ Consumer Group resource.
---

# alicloud_rocketmq_consumer_group

Provides a RocketMQ Consumer Group resource. 

For information about RocketMQ Consumer Group and how to use it, see [What is Consumer Group](https://www.alibabacloud.com/help/en/apsaramq-for-rocketmq/cloud-message-queue-rocketmq-5-x-series/developer-reference/api-rocketmq-2022-08-01-createconsumergroup).

-> **NOTE:** Available since v1.212.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_rocketmq_consumer_group&exampleId=ba15d9e6-acb1-e257-a427-5ec60b95b914a8f69563&activeTab=example&spm=docs.r.rocketmq_consumer_group.0.ba15d9e6ac&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "createVpc" {
  description = "example"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name
}

resource "alicloud_vswitch" "createVswitch" {
  description  = "example"
  vpc_id       = alicloud_vpc.createVpc.id
  zone_id      = data.alicloud_zones.default.zones.0.id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
}

resource "alicloud_rocketmq_instance" "createInstance" {
  product_info {
    msg_process_spec       = "rmq.u2.10xlarge"
    send_receive_ratio     = "0.3"
    message_retention_time = "70"
  }
  service_code    = "rmq"
  payment_type    = "PayAsYouGo"
  instance_name   = var.name
  sub_series_code = "cluster_ha"
  remark          = "example"
  ip_whitelists   = ["192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"]
  software {
    maintain_time = "02:00-06:00"
  }
  tags = {
    Created = "TF"
    For     = "example"
  }
  series_code = "ultimate"
  network_info {
    vpc_info {
      vpc_id = alicloud_vpc.createVpc.id
      vswitches {
        vswitch_id = alicloud_vswitch.createVswitch.id
      }
    }
    internet_info {
      internet_spec      = "enable"
      flow_out_type      = "payByBandwidth"
      flow_out_bandwidth = "30"
    }
  }
}

resource "alicloud_rocketmq_consumer_group" "default" {
  consumer_group_id = var.name
  instance_id       = alicloud_rocketmq_instance.createInstance.id
  consume_retry_policy {
    retry_policy    = "DefaultRetryPolicy"
    max_retry_times = "10"
  }

  delivery_order_type = "Concurrently"
  remark              = "example"
}
```

## Argument Reference

The following arguments are supported:
* `consume_retry_policy` - (Required, Set) Consumption retry strategy. See [`consume_retry_policy`](#consume_retry_policy) below.
* `consumer_group_id` - (Required, ForceNew) The first ID of the resource.
* `delivery_order_type` - (Optional) Delivery order.
* `instance_id` - (Required, ForceNew) Instance ID.
* `max_receive_tps` - (Optional, Int, Available since v1.247.0) Maximum received message tps.
* `remark` - (Optional) Custom remarks.

### `consume_retry_policy`

The consume_retry_policy supports the following:
* `dead_letter_target_topic` - (Optional, Available since v1.247.0) The dead-letter topic. If the consumer fails to consume a message in an abnormal situation and the message is still unsuccessful after retrying, the message will be delivered to the dead letter Topic for subsequent business recovery or backtracking.
* `max_retry_times` - (Optional) Maximum number of retries.
* `retry_policy` - (Optional) Consume retry policy.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<instance_id>:<consumer_group_id>`.
* `create_time` - The creation time of the resource.
* `region_id` - (Available since v1.247.0) The ID of the region in which the instance resides.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Consumer Group.
* `delete` - (Defaults to 5 mins) Used when delete the Consumer Group.
* `update` - (Defaults to 5 mins) Used when update the Consumer Group.

## Import

RocketMQ Consumer Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_rocketmq_consumer_group.example <instance_id>:<consumer_group_id>
```