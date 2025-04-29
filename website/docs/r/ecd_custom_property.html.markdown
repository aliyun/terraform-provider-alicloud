---
subcategory: "Elastic Desktop Service (ECD)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_custom_property"
sidebar_current: "docs-alicloud-resource-ecd-custom-property"
description: |-
  Provides a Alicloud ECD Custom Property resource.
---

# alicloud_ecd_custom_property

Provides a ECD Custom Property resource.

For information about ECD Custom Property and how to use it, see [What is Custom Property](https://www.alibabacloud.com/help/en/wuying-workspace/developer-reference/api-eds-user-2021-03-08-createproperty-desktop).

-> **NOTE:** Available since v1.176.0.

-> **NOTE:** Up to 10 different attributes can be created under an alibaba cloud account. Up to 50 different attribute values can be added under an attribute.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecd_custom_property&exampleId=3db4efc6-3592-2358-c623-7557341bd10e8a65d189&activeTab=example&spm=docs.r.ecd_custom_property.0.3db4efc635&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-shanghai"
}

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
* `property_values` - (Optional) Custom attribute sets the value of. See [`property_values`](#property_values) below.

### `property_values`

The property_values supports the following: 

* `property_value` - (Optional) The value of an attribute.
* `property_value_id` - (Computed) The value of an attribute id.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Custom Property.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Custom Property.
* `delete` - (Defaults to 1 mins) Used when delete the Custom Property.
* `update` - (Defaults to 1 mins) Used when update the Custom Property.



## Import

ECD Custom Property can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecd_custom_property.example <id>
```