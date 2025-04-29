---
subcategory: "RocketMQ"
layout: "alicloud"
page_title: "Alicloud: alicloud_rocketmq_topic"
description: |-
  Provides a Alicloud RocketMQ Topic resource.
---

# alicloud_rocketmq_topic

Provides a RocketMQ Topic resource. 

For information about RocketMQ Topic and how to use it, see [What is Topic](https://www.alibabacloud.com/help/en/apsaramq-for-rocketmq/cloud-message-queue-rocketmq-5-x-series/developer-reference/api-rocketmq-2022-08-01-createtopic).

-> **NOTE:** Available since v1.211.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_rocketmq_topic&exampleId=cc041f76-f391-834e-7ab8-73ce42034f311cd55635&activeTab=example&spm=docs.r.rocketmq_topic.0.cc041f76f3&intl_lang=EN_US" target="_blank">
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

resource "alicloud_rocketmq_topic" "default" {
  remark       = "example"
  instance_id  = alicloud_rocketmq_instance.createInstance.id
  message_type = "NORMAL"
  topic_name   = var.name
}
```

## Argument Reference

The following arguments are supported:
* `instance_id` - (Required, ForceNew) Instance ID.
* `max_send_tps` - (Optional, Int, Available since v1.247.0) The maximum TPS for message sending.
* `message_type` - (Optional, ForceNew) Message type.
* `remark` - (Optional) Custom remarks.
* `topic_name` - (Required, ForceNew) Topic name and identification.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<instance_id>:<topic_name>`.
* `create_time` - The creation time of the resource.
* `region_id` - (Available since v1.247.0) The region ID to which the instance belongs.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Topic.
* `delete` - (Defaults to 5 mins) Used when delete the Topic.
* `update` - (Defaults to 5 mins) Used when update the Topic.

## Import

RocketMQ Topic can be imported using the id, e.g.

```shell
$ terraform import alicloud_rocketmq_topic.example <instance_id>:<topic_name>
```