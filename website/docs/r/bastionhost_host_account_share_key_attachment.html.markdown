---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_host_account_share_key_attachment"
sidebar_current: "docs-alicloud-resource-bastionhost-host-account-share-key-attachment"
description: |-
  Provides a Alicloud Bastion Host Host Account Share Key Attachment resource.
---

# alicloud_bastionhost_host_account_share_key_attachment

Provides a Bastion Host Account Share Key Attachment resource.

For information about Bastion Host Host Account Share Key Attachment and how to use it, see [What is Host Account Share Key Attachment](https://www.alibabacloud.com/help/en/bastion-host/latest/attachhostaccountstohostsharekey).

-> **NOTE:** Available since v1.165.0.

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

variable "private_key" {
  default = "LS0tLS1CR*******"
}
resource "alicloud_bastionhost_host_share_key" "default" {
  host_share_key_name = var.name
  instance_id         = alicloud_bastionhost_instance.default.id
  pass_phrase         = "NTIxeGlubXU="
  private_key         = var.private_key
}

resource "alicloud_bastionhost_host_account_share_key_attachment" "default" {
  instance_id       = alicloud_bastionhost_instance.default.id
  host_share_key_id = alicloud_bastionhost_host_share_key.default.host_share_key_id
  host_account_id   = alicloud_bastionhost_host_account.default.host_account_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the Bastion machine instance.
* `host_account_id` - (Required, ForceNew) The ID list of the host account.
* `host_share_key_id` - (Required, ForceNew) The ID of the host shared key.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Host Account Share Key Attachment. The value formats as `<instance_id>:<host_share_key_id>:<host_account_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Bastion Host Account Share Key Attachment.
* `delete` - (Defaults to 1 mins) Used when delete the Bastion Host Account Share Key Attachment.

## Import

Bastion Host Account Share Key Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_bastionhost_host_account_share_key_attachment.example <instance_id>:<host_share_key_id>:<host_account_id>
```