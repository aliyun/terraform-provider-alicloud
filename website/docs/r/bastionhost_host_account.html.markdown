---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_host_account"
sidebar_current: "docs-alicloud-resource-bastionhost-host-account"
description: |-
  Provides a Alicloud Bastion Host Host Account resource.
---

# alicloud_bastionhost_host_account

Provides a Bastion Host Host Account resource.

For information about Bastion Host Host Account and how to use it, see [What is Host Account](https://www.alibabacloud.com/help/en/doc-detail/204377.htm).

-> **NOTE:** Available since v1.135.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
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
resource "alicloud_bastionhost_instance" "default" {
  description        = var.name
  license_code       = "bhah_ent_50_asset"
  plan_code          = "cloudbastion"
  storage            = "5"
  bandwidth          = "5"
  period             = "1"
  vswitch_id         = alicloud_vswitch.default.id
  security_group_ids = [alicloud_security_group.default.id]
}

resource "alicloud_bastionhost_host" "default" {
  instance_id          = alicloud_bastionhost_instance.default.id
  host_name            = var.name
  active_address_type  = "Private"
  host_private_address = "172.16.0.10"
  os_type              = "Linux"
  source               = "Local"
}

resource "alicloud_bastionhost_host_account" "default" {
  host_account_name = var.name
  host_id           = alicloud_bastionhost_host.default.host_id
  instance_id       = alicloud_bastionhost_host.default.instance_id
  protocol_name     = "SSH"
  password          = "YourPassword12345"
}
```

## Argument Reference

The following arguments are supported:

* `host_account_name` - (Required) The name of the host account. The name can be up to 128 characters in length.
* `host_id` - (Required, ForceNew) The ID of the host for which you want to create an account.
* `instance_id` - (Required, ForceNew) The ID of the Bastionhost instance where you want to create an account for the host.
* `pass_phrase` - (Optional, Sensitive) The passphrase of the private key for the host account. **NOTE:** It is valid when the attribute `protocol_name` is `SSH`.
* `password` - (Optional, Sensitive) The password of the host account.
* `private_key` - (Optional, Sensitive) The private key of the host account. The value is a Base64-encoded string. **NOTE:** It is valid when the attribute `protocol_name` is `SSH`
* `protocol_name` - (Required, ForceNew) The protocol used by the host account. Valid values: SSH,RDP

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Host Account. The value formats as `<instance_id>:<host_account_id>`.
* `host_account_id` - Hosting account ID.

## Import

Bastion Host Host Account can be imported using the id, e.g.

```shell
$ terraform import alicloud_bastionhost_host_account.example <instance_id>:<host_account_id>
```
