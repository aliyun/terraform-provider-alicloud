---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_host_group_account_user_attachment"
sidebar_current: "docs-alicloud-resource-bastionhost-host-group-account-user-attachment"
description: |-
  Provides a Alicloud Bastion Host Host Group Account attachment resource.
---

# alicloud_bastionhost_host_group_account_user_attachment

Provides a Bastion Host Host Account Attachment resource to add list host accounts into one user and one host group.

-> **NOTE:** Available in v1.135.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_bastionhost_host" "default" {
  instance_id          = "bastionhost-cn-tl3xxxxxxx"
  host_name            = var.name
  active_address_type  = "Private"
  host_private_address = "172.16.0.10"
  os_type              = "Linux"
  source               = "Local"
}
resource "alicloud_bastionhost_host_account" "default" {
  count             = 3
  instance_id       = alicloud_bastionhost_host.default.instance_id
  host_account_name = "example_value-${count.index}"
  host_id           = alicloud_bastionhost_host.default.host_id
  protocol_name     = "SSH"
  password          = "YourPassword12345"
}
resource "alicloud_bastionhost_user" "default" {
  instance_id         = alicloud_bastionhost_host.default.instance_id
  mobile_country_code = "CN"
  mobile              = "13312345678"
  password            = "YourPassword-123"
  source              = "Local"
  user_name           = "my-local-user"
}
resource "alicloud_bastionhost_host_group" "default" {
  host_group_name = "example_value"
  instance_id     = "bastionhost-cn-tl3xxxxxxx"
}
resource "alicloud_bastionhost_host_group_account_user_attachment" "default" {
  instance_id        = alicloud_bastionhost_host.default.instance_id
  user_id            = alicloud_bastionhost_user.default.user_id
  host_group_id      = alicloud_bastionhost_host_group.default.host_group_id
  host_account_names = alicloud_bastionhost_host_account.default.*.host_account_name
}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Required, ForceNew) The ID of the user that you want to authorize to manage the specified hosts and host accounts.
* `host_group_id` - (Required, ForceNew) The ID of the host group.
* `instance_id` - (Required, ForceNew) The ID of the Bastionhost instance where you want to authorize the user to manage the specified hosts and host accounts.
* `host_account_names` - (Required, List) A list names of the host account.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Host Account. The value formats as `<instance_id>:<user_id>:<host_group_id>`.

## Import

Bastion Host Host Account can be imported using the id, e.g.

```
$ terraform import alicloud_bastionhost_host_account.example <instance_id>:<user_id>:<host_group_id>
```
