---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_kv_namespace"
description: |-
  Provides a Alicloud ESA Kv Namespace resource.
---

# alicloud_esa_kv_namespace

Provides a ESA Kv Namespace resource.



For information about ESA Kv Namespace and how to use it, see [What is Kv Namespace](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateKvNamespace).

-> **NOTE:** Available since v1.244.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_esa_kv_namespace" "default" {
  description  = "this is a example namespace."
  kv_namespace = "example_namespace"
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional, ForceNew) The description of the namespace.
* `kv_namespace` - (Required, ForceNew) KV storage space name

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - KV storage space State

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Kv Namespace.
* `delete` - (Defaults to 5 mins) Used when delete the Kv Namespace.

## Import

ESA Kv Namespace can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_kv_namespace.example <id>
```