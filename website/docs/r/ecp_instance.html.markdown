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

For information about Elastic Cloud Phone (ECP) Instance and how to use it, see [What is Instance](https://next.api.aliyun.com/document/cloudphone/2020-12-30/RunInstances).

-> **NOTE:** Available since v1.158.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecp_instance&exampleId=c9a35006-8eb1-424d-2d4c-9bc1b65acae9da3b203f&activeTab=example&spm=docs.r.ecp_instance.0.c9a350068e&intl_lang=EN_US" target="_blank">
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

data "alicloud_ecp_zones" "default" {
}

data "alicloud_ecp_instance_types" "default" {
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_vpc" "default" {
  vpc_name   = "${var.name}-${random_integer.default.result}"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = "${var.name}-${random_integer.default.result}"
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.192.0/24"
  zone_id      = data.alicloud_ecp_zones.default.zones.0.zone_id
}

resource "alicloud_security_group" "default" {
  name   = "${var.name}-${random_integer.default.result}"
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_ecp_key_pair" "default" {
  key_pair_name   = "${var.name}-${random_integer.default.result}"
  public_key_body = "ssh-rsa AAAAB3Nza12345678qwertyuudsfsg"
}

resource "alicloud_ecp_instance" "default" {
  instance_type     = data.alicloud_ecp_instance_types.default.instance_types.0.instance_type
  image_id          = "android-image-release5501072_a11_20240530.raw"
  vswitch_id        = alicloud_vswitch.default.id
  security_group_id = alicloud_security_group.default.id
  key_pair_name     = alicloud_ecp_key_pair.default.key_pair_name
  vnc_password      = "Ecp123"
  payment_type      = "PayAsYouGo"
  instance_name     = var.name
  description       = var.name
  force             = "true"
}
```

## Argument Reference

The following arguments are supported:

* `instance_type` - (Required, ForceNew) The specifications of the ECP instance.
* `image_id` - (Required, ForceNew) The ID of the image.
* `vswitch_id` - (Required, ForceNew) The ID of the vSwitch.
* `security_group_id` - (Required, ForceNew) The ID of the security group.
* `eip_bandwidth` - (Optional, ForceNew, Int) The bandwidth of the elastic IP address (EIP). **NOTE:** From version 1.232.0, `eip_bandwidth` cannot be modified.
* `resolution` - (Optional) The resolution that you want to select for the ECP instance. **NOTE:** From version 1.232.0, `resolution` can be modified.
* `key_pair_name` - (Optional) The name of the key pair that you want to use to connect to the instance.
* `vnc_password` - (Optional) The VNC password of the instance. The password must be `6` characters in length and can contain only uppercase letters, lowercase letters, and digits.
* `payment_type` - (Optional, ForceNew) The billing method of the ECP instance. Default value: `PayAsYouGo`. Valid values: `PayAsYouGo`,`Subscription`. **NOTE:** From version 1.232.0, `payment_type` cannot be modified.
* `auto_pay` - (Optional, Bool) Specifies whether to enable the auto-payment feature. Valid values:
  - `true`: Enable.
  - `false`: Disable.
* `period` - (Optional) The subscription duration. Default value: `1`. Valid values:
  - If `period_unit` is set to `Month`. Valid values: `1`, `2`, `3`, and `6`.
  - If `period_unit` is set to `Year`. Valid values: `1` to `5`.
* `period_unit` - (Optional) The unit of the subscription duration. Default value: `Month`. Valid values: `Month`, `Year`.
* `auto_renew` - (Optional, Bool) Specifies whether to enable the auto-renewal feature. Valid values:
  - `true`: Enable.
  - `false`: Disable.
* `instance_name` - (Optional) The name of the ECP instance. The name must be `2` to `128` characters in length. It must start with a letter but cannot start with `http://` or `https://`. It can contain letters, digits, colons (:), underscores (_), periods (.), and hyphens (-).
* `description` - (Optional) The description of the ECP instance. The description must be `2` to `256` characters in length and cannot start with `http://` or `https://`.
* `force` - (Optional, Bool) Specifies whether to forcefully stop and release the instance. Default value: `false`. Valid values:
  - `true`: Enable.
  - `false`: Disable.
-> **NOTE:** If you want to destroy `alicloud_ecp_instance`, `force` must be set to `true`, and if `force` set to `true`, when `status` is set to `Stopped`， cache data that is not written to storage in the instance will be lost, which is similar to the effect of a power-off action.
* `status` - (Optional) The status of the Instance. Valid values: `Running`, `Stopped`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Instance.
* `update` - (Defaults to 3 mins) Used when update the Instance.
* `delete` - (Defaults to 3 mins) Used when delete the Instance.

## Import

Elastic Cloud Phone (ECP) Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecp_instance.example <id>
```
