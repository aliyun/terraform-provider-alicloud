---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_umodel"
description: |-
  Provides a Alicloud Cms Umodel resource.
---

# alicloud_cms_umodel

Provides a Cms Umodel resource.

The Umodel is a singleton configuration under a CMS workspace that manages common schema references used by Cloud Monitoring 2.0.

For information about Cms Umodel and how to use it, see [CreateUmodel](https://next.api.alibabacloud.com/document/Cms/2024-03-30/CreateUmodel).

-> **NOTE:** Available since v1.292.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = var.name
  sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_umodel" "default" {
  workspace   = alicloud_cms_workspace.default.workspace_name
  description = "terraform-example"
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required, ForceNew) The name of the workspace to which the umodel belongs. The umodel is a singleton per workspace, so this value also identifies the resource.
* `description` - (Optional) The description of the umodel. This is a write-only field: it is accepted by Create/Update but not returned by Get, so it is not refreshed from the server after apply.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource supplied above. It is the workspace name, since the umodel is a singleton per workspace.
* `region_id` - The region ID of the resource.
* `common_schema_ref` - The common schema references managed by the system. See [`common_schema_ref`](#common_schema_ref) below.

<a name="common_schema_ref"></a>

### `common_schema_ref`

The common_schema_ref exports the following:

* `group` - The group of the common schema reference.
* `items` - The item list of the common schema reference.
* `version` - The version of the common schema reference.

## Import

Cms Umodel can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_umodel.example <workspace>
```
