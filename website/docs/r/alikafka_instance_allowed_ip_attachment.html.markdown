---
subcategory: "Alikafka"
layout: "alicloud"
page_title: "Alicloud: alicloud_ali_kafka_instance_allowed_ip_attachment"
sidebar_current: "docs-alicloud-resource-alikafka-instance-allowed-ip-attachment"
description: |-
  Provides a Alicloud AliKafka Instance Allowed Ip Attachment resource.
---

# alicloud\_alikafka\_instance\_allowed\_ip\_attachment

Provides a AliKafka Instance Allowed Ip Attachment resource.

For information about Ali Kafka Instance Allowed Ip Attachment and how to use it, see [What is Instance Allowed Ip Attachment](https://www.alibabacloud.com/help/en/doc-detail/68151.html).

-> **NOTE:** Available in v1.163.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tftest"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_alikafka_instance" "default" {
  name           = var.name
  topic_quota    = 50
  disk_type      = 1
  disk_size      = 500
  deploy_type    = 5
  io_max         = 20
  vswitch_id     = data.alicloud_vswitches.default.ids.0
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
* `allowed_type` - (Required, ForceNew) The type of whitelist. Valid Value: `vpc`.
* `instance_id` - (Required, ForceNew) The ID of the instance.
* `port_range` - (Required, ForceNew) The Port range.  Valid Value: `9092/9092`.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Instance Allowed Ip Attachment. The value formats as `<instance_id>:<allowed_type>:<port_range>:<allowed_ip>`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the resource.
* `delete` - (Defaults to 1 mins) Used when delete the resource.


## Import

AliKafka Instance Allowed Ip Attachment can be imported using the id, e.g.

```
$ terraform import alicloud_ali_kafka_instance_allowed_ip_attachment.example <instance_id>:<allowed_type>:<port_range>:<allowed_ip>
```