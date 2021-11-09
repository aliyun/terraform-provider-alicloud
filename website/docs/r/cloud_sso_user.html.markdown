---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_user"
sidebar_current: "docs-alicloud-resource-cloud-sso-user"
description: |-
  Provides a Alicloud Cloud SSO User resource.
---

# alicloud\_cloud\_sso\_user

Provides a Cloud SSO User resource.

For information about Cloud SSO User and how to use it, see [What is User](https://www.alibabacloud.com/help/zh/doc-detail/264683.htm).

-> **NOTE:** Available in v1.140.0+.

-> **NOTE:** Cloud SSO Only Support `cn-shanghai` And `us-west-1` Region

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "example-value"
}
data "alicloud_cloud_sso_directories" "default" {}
resource "alicloud_cloud_sso_directory" "default" {
  count          = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? 0 : 1
  directory_name = var.name
}
resource "alicloud_cloud_sso_user" "default" {
  directory_id = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id, [""])[0]
  user_name    = var.name
}

```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of user.
* `directory_id` - (Required, ForceNew) The ID of the Directory.
* `display_name` - (Optional, Sensitive) The display name of user.
* `email` - (Optional, Sensitive) The User's Contact Email Address.
* `first_name` - (Optional, Sensitive) The first name of user.
* `last_name` - (Optional, Sensitive) The last name of user.
* `status` - (Optional) The status of user. Valid values: `Disabled`, `Enabled`.
* `user_name` - (Required from 1.141.0, ForceNew, Sensitive) The name of user.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of User. The value formats as `<directory_id>:<user_id>`.
* `user_id` - The User ID of the group.

## Import

Cloud SSO User can be imported using the id, e.g.

```
$ terraform import alicloud_cloud_sso_user.example <directory_id>:<user_id>
```
