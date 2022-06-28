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

```
$ terraform import alicloud_bastionhost_host_account.example <instance_id>:<host_account_id>
```
