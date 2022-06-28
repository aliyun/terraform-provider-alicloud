---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_config_map"
sidebar_current: "docs-alicloud-resource-sae-config-map"
description: |-
  Provides a Alicloud Serverless App Engine (SAE) Config Map resource.
---

# alicloud\_sae\_config\_map

Provides a Serverless App Engine (SAE) Config Map resource.

For information about Serverless App Engine (SAE) Config Map and how to use it, see [What is Config Map](https://help.aliyun.com/document_detail/97792.html).

-> **NOTE:** Available in v1.130.0+.

## Example Usage

Basic Usage

```terraform
variable "ConfigMapName" {
  default = "examplename"
}
resource "alicloud_sae_config_map" "example" {
  data         = jsonencode({ "env.home" : "/root", "env.shell" : "/bin/sh" })
  name         = var.ConfigMapName
  namespace_id = alicloud_sae_namespace.example.namespace_id
}
resource "alicloud_sae_namespace" "example" {
  namespace_id          = "cn-hangzhou:yourname"
  namespace_name        = "example_value"
  namespace_description = "your_description"
}

```

## Argument Reference

The following arguments are supported:

* `data` - (Required) ConfigMap instance data.
* `description` - (Optional) The Description of ConfigMap.
* `name` - (Required, ForceNew) ConfigMap instance name.
* `namespace_id` - (Required, ForceNew) The NamespaceId of ConfigMap.It can contain 2 to 32 lowercase characters.The value is in format `{RegionId}:{namespace}`

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Config Map.

## Import

Serverless App Engine (SAE) Config Map can be imported using the id, e.g.

```
$ terraform import alicloud_sae_config_map.example <id>
```
