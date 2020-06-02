---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_ee_namespace"
sidebar_current: "docs-alicloud-resource-container-registry"
description: |-
  Provides a Alicloud resource to manage Container Registry Enterprise namespaces.
---

# alicloud\_cr_ee\_namespace

This resource will help you to manager Container Registry Enterprise namespaces.

-> **NOTE:** Available in v1.85.0+.

-> **NOTE:** You need to set your registry password in Container Registry Enterprise console before use this resource.

## Example Usage

Basic Usage

```
resource "alicloud_cr_ee_namespace" "my-namespace" {
  instance_id        = "cri-xxx"
  name               = "my-namespace"
  auto_create        = false
  default_visibility = "PUBLIC"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) ID of CR EE instance.
* `name` - (Required, ForceNew) Name of CR EE namespace. It can contain 2 to 30 characters.
* `auto_create` - (Required) Boolean, when it set to true, repositories are automatically created when pushing new images. If it set to false, you create repository for images before pushing.
* `default_visibility` - (Required) `PUBLIC` or `PRIVATE`, default repository visibility in this namespace.

## Attributes Reference

The following attributes are exported:

* `id` - ID of CR EE namespace. The value is in format `{instance_id}/{namespace}` .

## Import

CR EE namespace can be imported using the `{instance_id}/{namespace}`, e.g.

```
$ terraform import alicloud_cr_ee_namespace.default cri-xxx/my-namespace
```
