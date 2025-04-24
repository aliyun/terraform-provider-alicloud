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

For information about AliKafka Instance Allowed Ip Attachment and how to use it, see [What is Instance Allowed Ip Attachment](https://www.alibabacloud.com/help/en/message-queue-for-apache-kafka/latest/api-alikafka-2019-09-16-updateallowedip).

-> **NOTE:** Available since v1.163.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alikafka_instance_allowed_ip_attachment&exampleId=c1d0d26f-b82d-d0a3-3030-db4a26ec5b5e891ada72&activeTab=example&spm=docs.r.alikafka_instance_allowed_ip_attachment.0.c1d0d26fb8&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
  instance_id  = alicloud_alikafka_instance.default.id
  allowed_type = "vpc"
  port_range   = "9092/9092"
  allowed_ip   = "114.237.9.78/32"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the instance.
* `allowed_type` - (Required, ForceNew) The type of the whitelist. Valid Value: `vpc`, `internet`. **NOTE:** From version 1.179.0, `allowed_type` can be set to `internet`.
  - `vpc`: A whitelist for access from a VPC.
  - `internet`: A whitelist for access from the Internet.
* `port_range` - (Required, ForceNew) The Port range. Valid Value: `9092/9092`, `9093/9093`, `9094/9094`, `9095/9095`. **NOTE:** From version 1.179.0, `port_range` can be set to `9093/9093`. From version 1.219.0, `port_range` can be set to `9094/9094`, `9095/9095`.
  - `9092/9092`: The port range for access from virtual private clouds (VPCs) by using the default endpoint.
  - `9093/9093`: The port range for access from the Internet.
  - `9094/9094`: The port range for access from VPCs by using the Simple Authentication and Security Layer (SASL) endpoint.
  - `9095/9095`: The port range for access from VPCs by using the Secure Sockets Layer (SSL) endpoint.
* `allowed_ip` - (Required, ForceNew) The IP address whitelist. It can be a CIDR block.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Instance Allowed Ip Attachment. It formats as `<instance_id>:<allowed_type>:<port_range>:<allowed_ip>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Instance Allowed Ip Attachment.
* `delete` - (Defaults to 1 mins) Used when delete the Instance Allowed Ip Attachment.


## Import

AliKafka Instance Allowed Ip Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_alikafka_instance_allowed_ip_attachment.example <instance_id>:<allowed_type>:<port_range>:<allowed_ip>
```
