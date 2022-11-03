---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_namespace"
sidebar_current: "docs-alicloud-resource-container-registry"
description: |-
  Provides a Alicloud resource to manage Container Registry namespaces.
---

# alicloud\_cr\_namespace

This resource will help you to manager Container Registry namespaces.

-> **NOTE:** Available in v1.34.0+.

-> **NOTE:** You need to set your registry password in Container Registry console before use this resource.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cr_namespace" "my-namespace" {
  name               = "my-namespace"
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
