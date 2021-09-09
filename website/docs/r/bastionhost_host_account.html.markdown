---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_host_account"
sidebar_current: "docs-alicloud-resource-bastionhost-host-account"
description: |-
  Provides a Alicloud Bastion Host Host Account resource.
---

# alicloud\_bastionhost\_host\_account

Provides a Bastion Host Host Account resource.

For information about Bastion Host Host Account and how to use it, see [What is Host Account](https://www.alibabacloud.com/help/en/doc-detail/204377.htm).

-> **NOTE:** Available in v1.135.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_bastionhost_host_account" "example" {
  host_account_name = "example_value"
  host_id           = "15"
  instance_id       = "bastionhost-cn-tl32bh0no30"
  protocol_name     = "SSH"
  password          = "YourPassword12345"
}

```

## Argument Reference

The following arguments are supported:

* `host_account_name` - (Required) Specify the new hosting account's name, support the longest 128 characters.
* `host_id` - (Required, ForceNew) Specifies the database where you want to create your hosting account's host ID.
* `instance_id` - (Required, ForceNew) Specifies the database where you want to create your hosting account's host bastion host ID of.
* `pass_phrase` - (Optional, Sensitive) Specify the new hosting account's private key password. **NOTE:** It is valid when the attribute `protocol_name` is `SSH`.
* `password` - (Optional, Sensitive) Specify the new hosting account's password.
* `private_key` - (Optional, Sensitive) Specify the new hosting account's private key using Base64 encoded string. **NOTE:** It is valid when the attribute `protocol_name` is `SSH`
* `protocol_name` - (Required, ForceNew) Specify the new hosting account of the agreement name. Valid values: SSH and RDP.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Host Account. The value formats as `<instance_id>:<host_account_id>`.
* `host_account_id` - Hosting account ID.

## Import

Bastion Host Host Account can be imported using the id, e.g.

```
$ terraform import alicloud_bastionhost_host_account.example <instance_id>:<host_account_id>
```