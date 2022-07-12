---
subcategory: "Elastic Desktop Service(EDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_custom_properties"
sidebar_current: "docs-alicloud-datasource-ecd-custom-properties"
description: |-
  Provides a list of Ecd Custom Properties to the user.
---

# alicloud\_ecd\_custom\_properties

This data source provides the Ecd Custom Properties of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.176.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecd_custom_properties" "ids" {
  ids = ["example_id"]
}
output "ecd_custom_property_id_1" {
  value = data.alicloud_ecd_custom_properties.ids.properties.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Custom Property IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `properties` - A list of Ecd Custom Properties. Each element contains the following attributes:
	* `custom_property_id` - The first ID of the resource.
	* `id` - The ID of the Custom Property.
	* `property_key` - The Custom attribute key.
	* `property_values` - Custom attribute sets the value of.
		* `property_value` - The value of an attribute.
		* `property_value_id` - The value of an attribute id.