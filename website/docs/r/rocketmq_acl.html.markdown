---
subcategory: "RocketMQ"
layout: "alicloud"
page_title: "Alicloud: alicloud_rocketmq_acl"
description: |-
  Provides a Alicloud RocketMQ Acl resource.
---

# alicloud_rocketmq_acl

Provides a RocketMQ Acl resource.



For information about RocketMQ Acl and how to use it, see [What is Acl](https://www.alibabacloud.com/help/en/apsaramq-for-rocketmq/cloud-message-queue-rocketmq-5-x-series/developer-reference/api-rocketmq-2022-08-01-createinstanceacl).

-> **NOTE:** Available since v1.245.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_rocketmq_acl&exampleId=1482156b-0ca9-9982-5f30-22479a2a824188f77221&activeTab=example&spm=docs.r.rocketmq_acl.0.1482156b0c&intl_lang=EN_US" target="_blank">
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

resource "alicloud_vpc" "defaultrqDtGm" {
  description = "1111"
  cidr_block  = "192.168.0.0/16"
  vpc_name    = "pop-example-vpc"
}

resource "alicloud_vswitch" "defaultjUrTYm" {
  vpc_id       = alicloud_vpc.defaultrqDtGm.id
  zone_id      = "cn-hangzhou-j"
  cidr_block   = "192.168.0.0/24"
  vswitch_name = "pop-example-vswitch"
}

resource "alicloud_rocketmq_instance" "defaultKJZNVM" {
  product_info {
    msg_process_spec       = "rmq.p2.4xlarge"
    send_receive_ratio     = "0.3"
    message_retention_time = "70"
  }
  service_code    = "rmq"
  series_code     = "professional"
  payment_type    = "PayAsYouGo"
  instance_name   = var.name
  sub_series_code = "cluster_ha"
  remark          = "example"
  network_info {
    vpc_info {
      vpc_id = alicloud_vpc.defaultrqDtGm.id
      vswitches {
        vswitch_id = alicloud_vswitch.defaultjUrTYm.id
      }
    }
    internet_info {
      internet_spec      = "enable"
      flow_out_type      = "payByBandwidth"
      flow_out_bandwidth = "5"
    }
  }
  acl_info {
    default_vpc_auth_free = false
    acl_types             = ["default", "apache_acl"]
  }
}

resource "alicloud_rocketmq_account" "defaultMeNlxe" {
  account_status = "ENABLE"
  instance_id    = alicloud_rocketmq_instance.defaultKJZNVM.id
  username       = "tfexample"
  password       = "123456"
}

resource "alicloud_rocketmq_topic" "defaultVA0zog" {
  instance_id  = alicloud_rocketmq_instance.defaultKJZNVM.id
  message_type = "NORMAL"
  topic_name   = "tfexample"
}

resource "alicloud_rocketmq_acl" "default" {
  actions       = ["Pub", "Sub"]
  instance_id   = alicloud_rocketmq_instance.defaultKJZNVM.id
  username      = alicloud_rocketmq_account.defaultMeNlxe.username
  resource_name = alicloud_rocketmq_topic.defaultVA0zog.topic_name
  resource_type = "Topic"
  decision      = "Deny"
  ip_whitelists = ["192.168.5.5"]
}
```

## Argument Reference

The following arguments are supported:
* `actions` - (Required, List) The type of operations that can be performed on the resource. Valid values:
  - If `resource_type` is set to `Topic`. Valid values: `Pub`, `Sub`.
  - If `resource_type` is set to `Group`. Valid values: `Sub`.
* `decision` - (Required) The decision result of the authorization. Valid values: `Deny`, `Allow`.
* `instance_id` - (Required, ForceNew) The instance ID.
* `ip_whitelists` - (Optional, List) The IP address whitelists.
* `resource_name` - (Required, ForceNew) The name of the resource on which you want to grant permissions.
* `resource_type` - (Required, ForceNew) The type of the resource on which you want to grant permissions. Valid values: `Group`, `Topic`.
* `username` - (Required, ForceNew) The username of the account.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Acl. It formats as `<instance_id>:<username>:<resource_type>:<resource_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Acl.
* `delete` - (Defaults to 5 mins) Used when delete the Acl.
* `update` - (Defaults to 5 mins) Used when update the Acl.

## Import

RocketMQ Acl can be imported using the id, e.g.

```shell
$ terraform import alicloud_rocketmq_acl.example <instance_id>:<username>:<resource_type>:<resource_name>
```
