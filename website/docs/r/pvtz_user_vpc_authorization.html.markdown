---
subcategory: "Private Zone"
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_user_vpc_authorization"
sidebar_current: "docs-alicloud-resource-pvtz-user-vpc-authorization"
description: |-
  Provides a Alicloud Private Zone User Vpc Authorization resource.
---

# alicloud_pvtz_user_vpc_authorization

Provides a Private Zone User Vpc Authorization resource.

-> **NOTE:** Available since v1.138.0.

## Example Usage

Basic Usage

```terraform
variable "authorized_user_id" {
  default = 123456789
}
resource "alicloud_pvtz_user_vpc_authorization" "example" {
  authorized_user_id = var.authorized_user_id
  auth_channel       = "RESOURCE_DIRECTORY"
}
```

## Argument Reference

The following arguments are supported:

* `auth_channel` - (Optional) The auth channel. Valid values: `RESOURCE_DIRECTORY`.
* `authorized_user_id` - (Required, ForceNew) The primary account ID of the user who authorizes the resource.
* `auth_type` - (Optional, ForceNew) The type of Authorization. Valid values: `NORMAL` and `CLOUD_PRODUCT`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of User Vpc Authorization. The value formats as `<authorized_user_id>:<auth_type>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the User Vpc Authorization.
* `delete` - (Defaults to 2 mins) Used when delete the User Vpc Authorization.

## Import

Private Zone User Vpc Authorization can be imported using the id, e.g.

```shell
$ terraform import alicloud_pvtz_user_vpc_authorization.example <authorized_user_id>:<auth_type>
```
