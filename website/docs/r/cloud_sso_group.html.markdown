---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_group"
sidebar_current: "docs-alicloud-resource-cloud-sso-group"
description: |-
  Provides a Alicloud Cloud Sso Group resource.
---

# alicloud\_cloud\_sso\_group

Provides a Cloud SSO Group resource.

For information about Cloud SSO Group and how to use it, see [What is Group](https://www.alibabacloud.com/help/doc-detail/264683.html).

-> **NOTE:** Available in v1.138.0+.

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
resource "alicloud_cloud_sso_group" "default" {
  directory_id = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id, [""])[0]
  group_name   = var.name
  description  = var.name
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The Description of the group. The description can be up to `1024` characters long.
* `directory_id` - (Required, ForceNew) The ID of the Directory.
* `group_name` - (Required) The Name of the group. The name must be `1` to `128` characters in length and can contain letters, digits, periods (.), underscores (_), and hyphens (-).

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Group. The value formats as `<directory_id>:<group_id>`.
* `group_id` - The GroupId of the group.

## Import

Cloud SSO Group can be imported using the id, e.g.

```
$ terraform import alicloud_cloud_sso_group.example <directory_id>:<group_id>
```
