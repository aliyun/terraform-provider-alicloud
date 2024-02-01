---
subcategory: "DCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_dcdn_kv_namespace"
sidebar_current: "docs-alicloud-resource-dcdn-kv-namespace"
description: |-
  Provides a Alicloud Dcdn Kv Namespace resource.
---

# alicloud_dcdn_kv_namespace

Provides a Dcdn Kv Namespace resource.

For information about Dcdn Kv Namespace and how to use it, see [What is Kv Namespace](https://www.alibabacloud.com/help/en/dcdn/developer-reference/api-dcdn-2018-01-15-putdcdnkvnamespace).

-> **NOTE:** Available since v1.198.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_dcdn_kv_namespace" "default" {
  description = var.name
  namespace   = "${var.name}-${random_integer.default.result}"
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Required, ForceNew) Namespace description information
* `namespace` - (Required, ForceNew) Namespace name. The name can contain letters, digits, hyphens (-), and underscores (_).


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Kv Namespace.
* `delete` - (Defaults to 5 mins) Used when delete the Kv Namespace.

## Import

Dcdn Kv Namespace can be imported using the id, e.g.

```shell
$ terraform import alicloud_dcdn_kv_namespace.example 
```