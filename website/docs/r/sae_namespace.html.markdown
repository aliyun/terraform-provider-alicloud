---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_namespace"
sidebar_current: "docs-alicloud-resource-sae-namespace"
description: |-
  Provides a Alicloud Serverless App Engine (SAE) Namespace resource.
---

# alicloud_sae_namespace

Provides a Serverless App Engine (SAE) Namespace resource.

For information about SAE Namespace and how to use it, see [What is Namespace](https://www.alibabacloud.com/help/en/sae/latest/createnamespace).

-> **NOTE:** Available since v1.129.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
}
data "alicloud_regions" "default" {
  current = true
}
resource "random_integer" "default" {
  max = 99999
  min = 10000
}
resource "alicloud_sae_namespace" "example" {
  namespace_id              = "${data.alicloud_regions.default.regions.0.id}:example${random_integer.default.result}"
  namespace_name            = var.name
  namespace_description     = var.name
  enable_micro_registration = false
}
```

## Argument Reference

The following arguments are supported:

* `namespace_name` - (Required) The Name of Namespace.
* `namespace_id` - (Optional, ForceNew) The ID of the Namespace. It can contain 2 to 32 lowercase characters. The value is in format `{RegionId}:{namespace}`.
* `namespace_short_id` - (Optional, ForceNew, Available since v1.206.0) The short ID of the Namespace. You do not need to specify a region ID. The value of `namespace_short_id` can be up to 20 characters in length and can contain only lowercase letters and digits.
* `namespace_description` - (Optional) The Description of Namespace.
* `enable_micro_registration` - (Optional, Available since v1.206.0) Specifies whether to enable the SAE built-in registry. If you do not use the built-in registry, you can set `enable_micro_registration` to `false` to accelerate the creation of the namespace. Default value: `true`. Valid values:
  - `true`: Enable.
  - `false`: Disable.
-> **NOTE:** From version 1.206.0, You should specify one of the `namespace_id` and `namespace_short_id`, and `namespace_short_id` is recommended.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Namespace. Its value is same as `namespace_id`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `delete` - (Defaults to 1 mins) Used when delete the Namespace.

## Import

Serverless App Engine (SAE) Namespace can be imported using the id, e.g.

```shell
$ terraform import alicloud_sae_namespace.example <namespace_id>
```
