---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_host_account_share_key_attachment"
sidebar_current: "docs-alicloud-resource-bastionhost-host-account-share-key-attachment"
description: |-
  Provides a Alicloud Bastion Host Host Account Share Key Attachment resource.
---

# alicloud\_bastionhost\_host\_account\_share\_key\_attachment

Provides a Bastion Host Account Share Key Attachment resource.

For information about Bastion Host Host Account Share Key Attachment and how to use it, see [What is Host Account Share Key Attachment](https://www.alibabacloud.com/help/en/bastion-host/latest/attachhostaccountstohostsharekey).

-> **NOTE:** Available in v1.165.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tfacc_host_account_share_key_attachment"
}

data "alicloud_bastionhost_instances" "default" {
}

resource "alicloud_bastionhost_host_share_key" "default" {
  host_share_key_name = "example_name"
  instance_id         = data.alicloud_bastionhost_instances.default.instances.0.id
  pass_phrase         = "example_value"
  private_key         = "example_value"
}

resource "alicloud_bastionhost_host" "default" {
  instance_id          = data.alicloud_bastionhost_instances.default.ids.0
  host_name            = var.name
  active_address_type  = "Private"
  host_private_address = "172.16.0.10"
  os_type              = "Linux"
  source               = "Local"
}

resource "alicloud_bastionhost_host_account" "default" {
  instance_id       = data.alicloud_bastionhost_instances.default.ids.0
  host_account_name = var.name
  host_id           = alicloud_bastionhost_host.default.host_id
  protocol_name     = "SSH"
  password          = "YourPassword12345"
}

resource "alicloud_bastionhost_host_account_share_key_attachment" "default" {
  instance_id       = data.alicloud_bastionhost_instances.default.instances.0.id
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

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Bastion Host Account Share Key Attachment.
* `delete` - (Defaults to 1 mins) Used when delete the Bastion Host Account Share Key Attachment.

## Import

Bastion Host Account Share Key Attachment can be imported using the id, e.g.

```
$ terraform import alicloud_bastionhost_host_account_share_key_attachment.example <instance_id>:<host_share_key_id>:<host_account_id>
```