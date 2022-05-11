---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_host_share_key"
sidebar_current: "docs-alicloud-resource-bastionhost-host-share-key"
description: |-
  Provides a Alicloud Bastion Host Host Share Key resource.
---

# alicloud\_bastionhost\_host\_share\_key

Provides a Bastion Host Share Key resource.

For information about Bastion Host Host Share Key and how to use it, see [What is Host Share Key](https://www.alibabacloud.com/help/en/bastion-host/latest/createhostsharekey).

-> **NOTE:** Available in v1.165.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_bastionhost_instances" "default" {
}

resource "alicloud_bastionhost_host_share_key" "default" {
  host_share_key_name = "example_name"
  instance_id         = data.alicloud_bastionhost_instances.default.instances.0.id
  pass_phrase         = "example_value"
  private_key         = "example_value"
}
```

## Argument Reference

The following arguments are supported:

* `host_share_key_name` - (Required) The name of the host shared key to be added. The name can be a maximum of 128 characters in length.
* `instance_id` - (Required, ForceNew) The ID of the Bastion instance.
* `pass_phrase` - (Optional, Sensitive) The password of the private key. The value is a Base64-encoded string.
* `private_key` - (Required, Sensitive) The private key. The value is a Base64-encoded string.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Host Share Key. The value formats as `<instance_id>:<host_share_key_id>`.
* `host_share_key_id` - The first ID of the resource.
* `private_key_finger_print` - The fingerprint of the private key.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Bastion Host Share Key.
* `update` - (Defaults to 1 mins) Used when update the Bastion Host Share Key.
* `delete` - (Defaults to 1 mins) Used when delete the Bastion Host Share Key.


## Import

Bastion Host Share Key can be imported using the id, e.g.

```
$ terraform import alicloud_bastionhost_host_share_key.example <instance_id>:<host_share_key_id>
```