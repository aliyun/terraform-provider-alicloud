---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_namespaces"
sidebar_current: "docs-alicloud-datasource-cms-namespaces"
description: |-
  Provides a list of Cms Namespaces to the user.
---

# alicloud\_cms\_namespaces

This data source provides the Cms Namespaces of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.171.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_namespaces" "ids" {
  ids = ["example_id"]
}
output "cms_namespace_id_1" {
  value = data.alicloud_cms_namespaces.ids.namespaces.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Namespace IDs.
* `keyword` - (Optional, ForceNew) The keywords of the `namespace` or `description` of the namespace.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `namespaces` - A list of Cms Namespaces. Each element contains the following attributes:
 * `create_time` - Create the timestamp of the indicator warehouse.
 * `description` - Description of indicator warehouse.
 * `id` - The ID of the Namespace.
 * `namespace_id` - The ID of the Namespace.
 * `modify_time` - The timestamp of the last modification indicator warehouse.
 * `namespace` - Indicator warehouse name.
 * `specification` - Data storage duration.