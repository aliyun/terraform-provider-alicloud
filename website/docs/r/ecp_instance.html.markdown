---
subcategory: "Elastic Cloud Phone (ECP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecp_instance"
sidebar_current: "docs-alicloud-resource-ecp-instance"
description: |-
  Provides a Alicloud Elastic Cloud Phone (ECP) Instance resource.
---

# alicloud_ecp_instance

Provides a Elastic Cloud Phone (ECP) Instance resource.

For information about Elastic Cloud Phone (ECP) Instance and how to use it,
see [What is Instance](https://www.alibabacloud.com/help/en/cloudphone/latest/api-cloudphone-2020-12-30-runinstances).

-> **NOTE:** Available since v1.158.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_ecp_zones" "default" {}
data "alicloud_ecp_instance_types" "default" {}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.1.0.0/16"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_ecp_zones.default.zones.0.zone_id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}
resource "alicloud_ecp_key_pair" "default" {
  key_pair_name   = var.name
  public_key_body = "ssh-rsa AAAAB3Nza12345678qwertyuudsfsg"
}

resource "alicloud_ecp_instance" "default" {
  instance_name     = var.name
  description       = var.name
  key_pair_name     = alicloud_ecp_key_pair.default.key_pair_name
  security_group_id = alicloud_security_group.default.id
  vswitch_id        = alicloud_vswitch.default.id
  image_id          = "android_9_0_0_release_2851157_20211201.vhd"
  instance_type     = data.alicloud_ecp_instance_types.default.instance_types.1.instance_type
  vnc_password      = "Ecp123"
  payment_type      = "PayAsYouGo"
}
```

## Argument Reference

The following arguments are supported:

* `auto_pay` - (Optional) The auto pay.
* `auto_renew` - (Optional) The auto renew.
* `payment_type` - (Optional) The payment type.Valid values: `PayAsYouGo`,`Subscription`
* `description` - (Optional) Description of the instance. 2 to 256 English or Chinese characters in length and cannot
  start with `http://` and `https`.
* `eip_bandwidth` - (Optional) The eip bandwidth.
* `force` - (Optional) The force.
* `image_id` - (Required, ForceNew) The ID Of The Image.
* `instance_name` - (Optional) The name of the instance. It must be 2 to 128 characters in length and must start with an
  uppercase letter or Chinese. It cannot start with http:// or https. It can contain Chinese, English, numbers,
  half-width colons (:), underscores (_), half-width periods (.), or dashes (-). The default value is the InstanceId of
  the instance.
* `instance_type` - (Required, ForceNew) Instance Type.
* `key_pair_name` - (Optional) The name of the key pair of the mobile phone instance.
* `period` - (Optional) The period. It is valid when `period_unit` is 'Year'. Valid value: `1`, `2`, `3`, `4`, `5`. It
  is valid when `period_unit` is 'Month'. Valid value: `1`, `2`, `3`, `5`
* `period_unit` - (Optional) The duration unit that you will buy the resource. Valid value: `Year`,`Month`. Default
  to `Month`.
* `resolution` - (Optional, ForceNew) The selected resolution for the cloud mobile phone instance.
* `security_group_id` - (Required, ForceNew) The ID of the security group. The security group is the same as that of the
  ECS instance.
* `status` - (Optional) Instance status. Valid values: `Running`, `Stopped`.
* `vnc_password` - (Optional) Cloud mobile phone VNC password. The password must be six characters in length and must
  contain only uppercase, lowercase English letters and Arabic numerals.
* `vswitch_id` - (Required, ForceNew) The vswitch id.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Instance.

## Timeouts

The `timeouts` block allows you to
specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Instance.
* `update` - (Defaults to 3 mins) Used when update the Instance.

## Import

Elastic Cloud Phone (ECP) Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecp_instance.example <id>
```