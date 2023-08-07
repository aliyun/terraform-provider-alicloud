---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_user_attachment"
sidebar_current: "docs-alicloud-resource-cloud-sso-user-attachment"
description: |-
  Provides a Alicloud Cloud SSO User Attachment resource.
---

# alicloud_cloud_sso_user_attachment

Provides a Cloud SSO User Attachment resource.

For information about Cloud SSO User Attachment and how to use it, see [What is User Attachment](https://www.alibabacloud.com/help/en/cloudsso/latest/api-cloudsso-2021-05-15-addusertogroup).

-> **NOTE:** Available since v1.141.0.

-> **NOTE:** Cloud SSO Only Support `cn-shanghai` And `us-west-1` Region

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}
provider "alicloud" {
  region = "cn-shanghai"
}
data "alicloud_cloud_sso_directories" "default" {}

resource "alicloud_cloud_sso_directory" "default" {
  count          = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? 0 : 1
  directory_name = var.name
}

locals {
  directory_id = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id, [""])[0]
}
resource "alicloud_cloud_sso_user" "default" {
  directory_id = local.directory_id
  user_name    = var.name
}

resource "alicloud_cloud_sso_group" "default" {
  directory_id = local.directory_id
  group_name   = var.name
  description  = var.name
}

resource "alicloud_cloud_sso_user_attachment" "default" {
  directory_id = local.directory_id
  user_id      = alicloud_cloud_sso_user.default.user_id
  group_id     = alicloud_cloud_sso_group.default.group_id
}
```

## Argument Reference

The following arguments are supported:

* `directory_id` - (Required, ForceNew) The ID of the Directory.
* `group_id` - (Required, ForceNew) The Group ID.
* `user_id` - (Required, ForceNew) The User ID.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of User Attachment. The value formats as `<directory_id>:<group_id>:<user_id>`.

## Import

Cloud SSO User Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_sso_user_attachment.example <directory_id>:<group_id>:<user_id>
```
