---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_namespace"
sidebar_current: "docs-alicloud-resource-container-registry"
description: |-
  Provides a Alicloud resource to manage Container Registry namespaces.
---

# alicloud_cr_namespace

This resource will help you to manager Container Registry namespaces, see [What is Namespace](https://www.alibabacloud.com/help/en/acr/developer-reference/api-cr-2018-12-01-createnamespace).

-> **NOTE:** Available since v1.34.0.

-> **NOTE:** You need to set your registry password in Container Registry console before use this resource.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}
resource "alicloud_cr_namespace" "example" {
  name               = var.name
  auto_create        = false
  default_visibility = "PUBLIC"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) Name of Container Registry namespace.
* `auto_create` - (Required) Boolean, when it set to true, repositories are automatically created when pushing new images. If it set to false, you create repository for images before pushing.
* `default_visibility` - (Required) `PUBLIC` or `PRIVATE`, default repository visibility in this namespace.

## Attributes Reference

The following attributes are exported:

* `id` - The id of Container Registry namespace. The value is same as its name.

## Import

Container Registry namespace can be imported using the namespace, e.g.

```shell
$ terraform import alicloud_cr_namespace.default my-namespace
```
