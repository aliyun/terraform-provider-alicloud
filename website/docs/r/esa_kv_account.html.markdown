---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_kv_account"
description: |-
  Provides a Alicloud ESA Kv Account resource.
---

# alicloud_esa_kv_account

Provides a ESA Kv Account resource.



For information about ESA Kv Account and how to use it, see [What is Kv Account](https://next.api.alibabacloud.com/document/ESA/2024-09-10/OpenErService).

-> **NOTE:** Available since v1.259.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_esa_kv_account" "open" {
}
```

### Deleting `alicloud_esa_kv_account` or removing it from your configuration

Terraform cannot destroy resource `alicloud_esa_kv_account`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as account id.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Kv Account.

## Import

ESA Kv Account can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_kv_account.example 
```