---
subcategory: "Max Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_maxcompute_package"
description: |-
  Provides a Alicloud Max Compute Package resource.
---

# alicloud_maxcompute_package

Provides a Max Compute Package resource.

A MaxCompute package groups together resources (tables, resources, and functions) from a source project so that they can be shared with other projects.

For information about Max Compute Package and how to use it, see [What is Package](https://www.alibabacloud.com/help/en/maxcompute/user-guide/).

-> **NOTE:** Available since v1.294.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example_package"
}

resource "alicloud_maxcompute_project" "default" {
  project_name = "${var.name}_project"
  product_type = "PayAsYouGo"
}

resource "alicloud_maxcompute_package" "default" {
  project_name = alicloud_maxcompute_project.default.project_name
  package_name = "${var.name}_pkg"
  is_install   = true
}
```

## Argument Reference

The following arguments are supported:

* `project_name` - (Required, ForceNew) The name of the MaxCompute project that owns the package.
* `package_name` - (Required, ForceNew) The name of the package. It is sent as the raw request body when the package is created.
* `is_install` - (Optional, ForceNew) Whether to install the package when creating it. Default value: `false`. Valid values:
  - `true`: install.
  - `false`: do not install.
* `body` - (Optional) The package content update request body. Only takes effect on resource update; a change to this field triggers `UpdatePackage` and is sent as the raw request body.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of the package, formatted as `<project_name>:<package_name>`.
* `project_name` - The name of the MaxCompute project that owns the package.
* `package_name` - The name of the package.

## Import

Max Compute Package can be imported using the id, e.g.

```shell
terraform import alicloud_maxcompute_package.example <project_name>:<package_name>
```
