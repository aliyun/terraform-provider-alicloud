---
subcategory: "EAIS"
layout: "alicloud"
page_title: "Alicloud: alicloud_eais_client_instance_attachment"
description: |-
  Provides a Alicloud EAIS Client Instance Attachment resource.
---

# alicloud_eais_client_instance_attachment

Provides a EAIS Client Instance Attachment resource. Bind an ECS or ECI instance.

For information about EAIS Client Instance Attachment and how to use it, see [What is Client Instance Attachment](https://www.alibabacloud.com/help/en/resource-orchestration-service/latest/aliyun-eais-clientinstanceattachment).

-> **NOTE:** Available since v1.210.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/21"
  availability_zone = "cn-hangzhou-j"
  name              = var.name
}

resource "alicloud_security_group" "default" {
  name        = var.name
  description = "tf test"
  vpc_id      = alicloud_vpc.default.id
}

resource "alicloud_eais_instance" "default" {
  instance_type     = "eais.ei-a6.2xlarge"
  instance_name     = var.name
  security_group_id = alicloud_security_group.default.id
  vswitch_id        = alicloud_vswitch.default.id
}

data "alicloud_instance_types" "default" {
  availability_zone                 = "cn-hangzhou-j"
  system_disk_category              = "cloud_efficiency"
  cpu_core_count                    = 4
  minimum_eni_ipv6_address_quantity = 1
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

resource "alicloud_instance" "default" {
  image_id                   = "${data.alicloud_images.default.images.0.id}"
  instance_type              = "${data.alicloud_instance_types.default.instance_types.0.id}"
  instance_name              = "${var.name}"
  security_groups            = "${alicloud_security_group.default.*.id}"
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = "${alicloud_vswitch.default.id}"
}

resource "alicloud_eais_client_instance_attachment" "default" {
  instance_id        = alicloud_eais_instance.default.id
  client_instance_id = alicloud_instance.default.id
}
```

## Argument Reference

The following arguments are supported:
* `client_instance_id` - (Required, ForceNew) The ID of the ECS or ECI instance bound to the EAIS instance.
* `instance_id` - (Required, ForceNew) The EAIS instance ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<instance_id>:<client_instance_id>`.
* `create_time` - The creation time of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Client Instance Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Client Instance Attachment.

## Import

EAIS Client Instance Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_eais_client_instance_attachment.example <instance_id>:<client_instance_id>
```