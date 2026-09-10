---
subcategory: "Schedulerx"
layout: "alicloud"
page_title: "Alicloud: alicloud_schedulerx_namespaces"
sidebar_current: "docs-alicloud-datasource-schedulerx-namespaces"
description: |-
  Provides a list of Schedulerx Namespaces to the user.
---

# alicloud\_schedulerx\_namespaces

This data source provides the Schedulerx Namespaces of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.173.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_schedulerx_namespaces" "ids" {}
output "schedulerx_namespace_id_1" {
  value = data.alicloud_schedulerx_namespaces.ids.namespaces.0.id
}

data "alicloud_schedulerx_namespaces" "nameRegex" {
  name_regex = "^my-Namespace"
}
output "schedulerx_namespace_id_2" {
  value = data.alicloud_schedulerx_namespaces.nameRegex.namespaces.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Namespace IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Namespace name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Namespace names.
* `namespaces` - A list of Schedulerx Namespaces. Each element contains the following attributes:
	* `description` - The description of the resource.
	* `id` - The ID of the resource.
	* `namespace_id` - The ID of the Namespace.
	* `namespace_name` - The name of the resource.