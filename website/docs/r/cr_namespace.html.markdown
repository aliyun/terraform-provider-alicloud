---
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_namespace"
sidebar_current: "docs-alicloud-resource-container-registry"
description: |-
  Provides a Alicloud resource to manage container registry namespaces.
---

# alicloud\_cr\_namespace

This resource will help you to manager container registry namespaces.

-> **NOTE:** Available in v1.34.0+.

## Example Usage

Basic Usage

```
resource "alicloud_cr_namespace" "my-namespace" {
    name = "my-namespace"
    auto_create = false
    default_visibility = "PUBLIC"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) Name of container registry namespace.
* `auto_create` - (Required, ForceNew) Boolean, when it set to true, repositories are automatically created when pushing new images. If it set to false, you create repository for images before pushing.
* `default_visibility` - (Required, ForceNew) `PUBLIC` or `PRIVATE`, default repository visibility in this namespace.

## Attributes Reference

The following attributes are exported:

* `id` - The id of container registry namespace. The value is same as its name.

## Import

Container Registry Namespace can be imported using the namespace, e.g.

```
$ terraform import alicloud_cr_namespace.default my-namespace
```
