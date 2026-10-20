---
subcategory: "Max Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_maxcompute_packages"
description: |-
  Provides a datasource of Max Compute Package owned by an Alibaba Cloud account.
---

# alicloud_maxcompute_packages

This data source provides Max Compute Package available to the user.

-> **NOTE:** Available since v1.294.0.

## Example Usage

```terraform
variable "name" {
  default = "tf_example_acc"
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

data "alicloud_maxcompute_packages" "default" {
  project_name = alicloud_maxcompute_project.default.project_name
  name_regex   = var.name
}

output "alicloud_maxcompute_package_example_id" {
  value = data.alicloud_maxcompute_packages.default.packages.0.package_name
}
```

## Argument Reference

The following arguments are supported:

* `project_name` - (Required) The name of the MaxCompute project whose packages are listed.
* `ids` - (Optional) A list of package IDs, formatted as `<project_name>:<package_name>`.
* `name_regex` - (Optional) A regex string to filter results by the package name.
* `output_file` - (Optional) File name where to write data source results after running `terraform plan`.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of package IDs, formatted as `<project_name>:<package_name>`.
* `names` - A list of package names.
* `packages` - A list of packages. Each element contains the following sub-attributes:
  * `project_name` - The name of the MaxCompute project that owns the package.
  * `package_name` - The name of the package.
  * `source_project` - The source project from which the package is installed.
  * `status` - The installation status of the package.
  * `install_time` - The time when the package was installed.
