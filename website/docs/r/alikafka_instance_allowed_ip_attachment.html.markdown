---
subcategory: "AliKafka"
layout: "alicloud"
page_title: "Alicloud: alicloud_alikafka_instance_allowed_ip_attachment"
sidebar_current: "docs-alicloud-resource-alikafka-instance-allowed-ip-attachment"
description: |-
  Provides a Alicloud AliKafka Instance Allowed Ip Attachment resource.
---

# alicloud_alikafka_instance_allowed_ip_attachment

Provides a AliKafka Instance Allowed Ip Attachment resource.

For information about Ali Kafka Instance Allowed Ip Attachment and how to use it, see [What is Instance Allowed Ip Attachment](https://www.alibabacloud.com/help/en/message-queue-for-apache-kafka/latest/api-alikafka-2019-09-16-updateallowedip).

-> **NOTE:** Available since v1.163.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_alikafka_instance" "default" {
  name           = "${var.name}-${random_integer.default.result}"
  partition_num  = 50
  disk_type      = 1
  disk_size      = 500
  deploy_type    = 5
  io_max         = 20
  vswitch_id     = alicloud_vswitch.default.id
  security_group = alicloud_security_group.default.id
}

resource "alicloud_alikafka_instance_allowed_ip_attachment" "default" {
  allowed_ip   = "114.237.9.78/32"
  allowed_type = "vpc"
  instance_id  = alicloud_alikafka_instance.default.id
  port_range   = "9092/9092"
}
```

## Argument Reference

The following arguments are supported:

* `allowed_ip` - (Required, ForceNew) The allowed ip. It can be a CIDR block.
* `allowed_type` - (Required, ForceNew) The type of whitelist. Valid Value: `vpc`, `internet`. **NOTE:** From version 1.179.0, `allowed_type` can be set to `internet`.
  - `vpc`: IP address whitelist for VPC access.
  - `internet`: IP address whitelist for Internet access.
* `instance_id` - (Required, ForceNew) The ID of the instance.
* `port_range` - (Required, ForceNew) The Port range.  Valid Value: `9092/9092`, `9093/9093`. **NOTE:** From version 1.179.0, `port_range` can be set to `9093/9093`.
  - `9092/9092`: port range for a VPC whitelist.
  - `9093/9093`: port range for an Internet whitelist.
  
## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Instance Allowed Ip Attachment. The value formats as `<instance_id>:<allowed_type>:<port_range>:<allowed_ip>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the resource.
* `delete` - (Defaults to 1 mins) Used when delete the resource.


## Import

AliKafka Instance Allowed Ip Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_ali_kafka_instance_allowed_ip_attachment.example <instance_id>:<allowed_type>:<port_range>:<allowed_ip>
```