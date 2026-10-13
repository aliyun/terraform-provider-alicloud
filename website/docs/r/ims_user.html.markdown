---
subcategory: "IMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ims_user"
sidebar_current: "docs-alicloud-resource-ims-user"
description: |-
  Provides a IMS User resource.
---

# alicloud_ims_user

Provides a IMS User resource.

For information about IMS User and how to use it, see [What is User](https://www.alibabacloud.com/help/en/ram/developer-reference/api-ims-2019-08-15-createuser).

-> **NOTE:** Available since v2.0.0-beta5.

## Example Usage

Basic Usage

```terraform
data "alicloud_ims_default_domain" "example" {}

resource "alicloud_ims_user" "example" {
  user_principal_name = "terraform-example@${data.alicloud_ims_default_domain.example.default_domain}"
  display_name        = "Terraform Example"
  comments            = "Managed by Terraform"
}
```

## Argument Reference

The following arguments are supported:

* `user_principal_name` - (Required) The logon name of the RAM user, in the format `<username>@<AccountAlias>.onaliyun.com`. The total length is 1 to 128 characters, and `<username>` is 1 to 64 characters in length. Only letters, digits, periods (.), hyphens (-) and underscores (_) are allowed. The default domain can be read from the `alicloud_ims_default_domain` data source.
* `display_name` - (Required) The display name of the RAM user. The length is 1 to 24 characters.
* `email` - (Optional) The email address of the RAM user. This parameter is available only on the China site.
* `mobile_phone` - (Optional) The mobile phone number of the RAM user, in the format `<international area code>-<number>`, for example `86-1444444****`. This parameter is available only on the China site.
* `comments` - (Optional) The remarks of the RAM user. The length is 1 to 128 characters.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. Same as `user_id`.
* `user_id` - The ID of the RAM user, for example `123456789000****`.
* `user_name` - The name of the RAM user, the part of the logon name before the at sign (@).
* `provision_type` - How the RAM user was created. Valid values: `Manual`, `SCIM` and `CloudSSO`.
* `create_date` - The time when the RAM user was created, in RFC 3339 format (UTC).
* `update_date` - The time when the RAM user was last updated, in RFC 3339 format (UTC).
* `last_login_date` - The time when the RAM user last logged on to the console, in RFC 3339 format (UTC). Empty if the user has never logged on.

## Import

IMS User can be imported using the id, e.g.

```shell
$ terraform import alicloud_ims_user.example <user_id>
```
