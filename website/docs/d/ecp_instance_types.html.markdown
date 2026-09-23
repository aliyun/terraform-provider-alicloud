---
subcategory: "Elastic Cloud Phone (ECP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecp_instance_types"
sidebar_current: "docs-alicloud-datasource-ecp-instance-types"
description: |-
  Provides a list of Ecp available instance types to the user.
---

# alicloud\_ecp\_instance\_types

This data source provides the available instance types with the Cloud Phone (ECP) Instance of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.158.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecp_instance_types" "default" {}

output "first_ecp_instance_types_instance_type" {
  value = "${data.alicloud_ecp_instance_types.default.instance_types.0.instance_type}"
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `instance_types` - A list of ecp Instance types. Each element contains the following attributes:
    * `instance_type` - The list of available instance type.
    * `default_resolution` - The default resolution of the current instance type.
    * `cpu_core_count` - The cpu core count of the current instance type.
    * `name` - The name of the current instance type.
    * `name_en` - The English name of the current instance type.