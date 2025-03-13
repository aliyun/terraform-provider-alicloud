---
subcategory: "RocketMQ"
layout: "alicloud"
page_title: "Alicloud: alicloud_rocketmq_account"
description: |-
  Provides a Alicloud RocketMQ Account resource.
---

# alicloud_rocketmq_account

Provides a RocketMQ Account resource.



For information about RocketMQ Account and how to use it, see [What is Account](https://www.alibabacloud.com/help/en/apsaramq-for-rocketmq/cloud-message-queue-rocketmq-5-x-series/developer-reference/api-rocketmq-2022-08-01-createinstanceaccount).

-> **NOTE:** Available since v1.245.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_vpc" "defaultg6ZXs2" {
  description = "111"
  cidr_block  = "192.168.0.0/16"
  vpc_name    = "pop-example-vpc"
}

resource "alicloud_vswitch" "defaultvMQbCy" {
  vpc_id       = alicloud_vpc.defaultg6ZXs2.id
  zone_id      = "cn-hangzhou-j"
  cidr_block   = "192.168.0.0/24"
  vswitch_name = "pop-example-vswitch"
}

resource "alicloud_rocketmq_instance" "default9hAb83" {
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
  software {
    maintain_time = "02:00-06:00"
  }
  tags = {
    Created = "TF"
    For     = "example"
  }
  network_info {
    vpc_info {
      vpc_id = alicloud_vpc.defaultg6ZXs2.id
      vswitches {
        vswitch_id = alicloud_vswitch.defaultvMQbCy.id
      }
    }
    internet_info {
      internet_spec      = "enable"
      flow_out_type      = "payByBandwidth"
      flow_out_bandwidth = "30"
    }
  }
  acl_info {
    default_vpc_auth_free = false
    acl_types             = ["default", "apache_acl"]
  }
}

resource "alicloud_rocketmq_account" "default" {
  account_status = "ENABLE"
  instance_id    = alicloud_rocketmq_instance.default9hAb83.id
  username       = "tfexample"
  password       = "1741835136"
}
```

## Argument Reference

The following arguments are supported:
* `account_status` - (Optional) The status of the account. Valid values: `DISABLE`, `ENABLE`.
* `instance_id` - (Required, ForceNew) The instance ID.
* `password` - (Required) The password of the account.
* `username` - (Required, ForceNew) The username of the account.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Account. It formats as `<instance_id>:<username>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Account.
* `delete` - (Defaults to 5 mins) Used when delete the Account.
* `update` - (Defaults to 5 mins) Used when update the Account.

## Import

RocketMQ Account can be imported using the id, e.g.

```shell
$ terraform import alicloud_rocketmq_account.example <instance_id>:<username>
```