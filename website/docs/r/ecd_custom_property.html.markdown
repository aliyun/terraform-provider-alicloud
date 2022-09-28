---
subcategory: "Elastic Desktop Service(EDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_custom_property"
sidebar_current: "docs-alicloud-resource-ecd-custom-property"
description: |-
  Provides a Alicloud ECD Custom Property resource.
---

# alicloud\_ecd\_custom\_property

Provides a ECD Custom Property resource.

For information about ECD Custom Property and how to use it, see [What is Custom Property](https://help.aliyun.com/document_detail/436381.html).

-> **NOTE:** Available in v1.176.0+.

-> **NOTE:** Up to 10 different attributes can be created under an alibaba cloud account. Up to 50 different attribute values can be added under an attribute.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecd_custom_property" "example" {
  property_key = "example_key"
  property_values {
    property_value = "example_value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `property_key` - (Required) The Custom attribute key.
* `property_values` - (Optional) Custom attribute sets the value of. See the following `Block property_values`.

#### Block property_values

The property_values supports the following: 

* `property_value` - (Optional) The value of an attribute.
* `property_value_id` - (Computed) The value of an attribute id.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Custom Property.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Custom Property.
* `delete` - (Defaults to 1 mins) Used when delete the Custom Property.
* `update` - (Defaults to 1 mins) Used when update the Custom Property.



## Import

ECD Custom Property can be imported using the id, e.g.

```
$ terraform import alicloud_ecd_custom_property.example <id>
```