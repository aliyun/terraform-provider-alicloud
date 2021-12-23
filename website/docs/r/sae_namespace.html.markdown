---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_namespace"
sidebar_current: "docs-alicloud-resource-sae-namespace"
description: |-
  Provides a Alicloud Serverless App Engine (SAE) Namespace resource.
---

# alicloud\_sae\_namespace

Provides a Serverless App Engine (SAE) Namespace resource.

For information about SAE Namespace and how to use it, see [What is Namespace](https://help.aliyun.com/document_detail/97792.html).

-> **NOTE:** Available in v1.129.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_sae_namespace" "example" {
  namespace_id          = "cn-hangzhou:yourname"
  namespace_name        = "example_value"
  namespace_description = "your_description"
}

```

## Argument Reference

The following arguments are supported:

* `namespace_description` - (Optional) The Description of Namespace.
* `namespace_id` - (Required, ForceNew) The Id of Namespace.It can contain 2 to 32 lowercase characters.The value is in format `{RegionId}:{namespace}`
* `namespace_name` - (Required) The Name of Namespace.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Namespace. Its value is same as `namespace_id`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `delete` - (Defaults to 1 mins) Used when delete the Namespace.

## Import

Serverless App Engine (SAE) Namespace can be imported using the id, e.g.

```
$ terraform import alicloud_sae_namespace.example <namespace_id>
```
