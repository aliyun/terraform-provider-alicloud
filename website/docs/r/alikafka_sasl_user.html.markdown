---
subcategory: "Alikafka"
layout: "alicloud"
page_title: "Alicloud: alicloud_alikafka_sasl_user"
sidebar_current: "docs-alicloud-resource-alikafka-sasl_user"
description: |-
  Provides a Alicloud Alikafka Sasl User resource.
---

# alicloud\_alikafka\_sasl\_user

Provides an ALIKAFKA sasl user resource.

-> **NOTE:** Available in 1.66.0+

-> **NOTE:**  Only the following regions support create alikafka sasl user.
[`cn-hangzhou`,`cn-beijing`,`cn-shenzhen`,`cn-shanghai`,`cn-qingdao`,`cn-hongkong`,`cn-huhehaote`,`cn-zhangjiakou`,`ap-southeast-1`,`ap-south-1`,`ap-southeast-5`]

## Example Usage

Basic Usage

```
variable "username" {
  default = "testusername"
}

variable "password" {
  default = "testpassword"
}

data "alicloud_zones" "default" {
    available_resource_creation= "VSwitch"
}
resource "alicloud_vpc" "default" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_alikafka_instance" "default" {
  name = "tf-testacc-alikafkainstance"
  topic_quota = "50"
  disk_type = "1"
  disk_size = "500"
  deploy_type = "5"
  io_max = "20"
  vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_alikafka_sasl_user" "default" {
  instance_id = "${alicloud_alikafka_instance.default.id}"
  username = "${var.username}"
  password = "${var.password}"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) ID of the ALIKAFKA Instance that owns the groups.
* `username` - (Required, ForceNew) Username for the sasl user. The length should between 1 to 64 characters. The characters can only contain 'a'-'z', 'A'-'Z', '0'-'9', '_' and '-'.
* `password` - (Required, ForceNew) Password for the sasl user. The length should between 1 to 64 characters.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `<instance_id>:<username>`.

## Import

ALIKAFKA GROUP can be imported using the id, e.g.

```
$ terraform import alicloud_alikafka_sasl_user.user alikafka_post-cn-123455abc:username
```
