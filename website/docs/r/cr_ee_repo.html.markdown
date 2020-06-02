---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_ee_repo"
sidebar_current: "docs-alicloud-resource-container-registry"
description: |-
  Provides a Alicloud resource to manage Container Registry Enterprise repositories.
---

# alicloud\_cr_ee\_repo

This resource will help you to manager Container Registry Enterprise repositories.

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

resource "alicloud_cr_ee_repo" "my-repo" {
  instance_id   = "${alicloud_cr_ee_namespace.my-namespace.instance_id}"
  namespace     = "${alicloud_cr_ee_namespace.my-namespace.name}"
  name          = "my-repo"
  summary       = "this is summary of my new repo"
  repo_type     = "PUBLIC"
  detail        = "this is a public repo"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) ID of CR EE instance.
* `namespace` - (Required, ForceNew) Name of CR EE namespace where repository is located. It can contain 2 to 30 characters.
* `name` - (Required, ForceNew) Name of CR EE repository. It can contain 2 to 64 characters.
* `summary` - (Required) The repository general information. It can contain 1 to 100 characters.
* `repo_type` - (Required) `PUBLIC` or `PRIVATE`, repo's visibility.
* `detail` - (Optional) The repository specific information. MarkDown format is supported, and the length limit is 2000.

## Attributes Reference

The following attributes are exported:

* `id` - The resource id of CR EE repository. The value is in format `{instance_id}/{namespace}/{repository}`.
* `repo_id` - The uuid of CR EE repository.

## Import

CR EE repository can be imported using the `{instance_id}/{namespace}/{repository}`, e.g.

```
$ terraform import alicloud_cr_ee_repo.default `cri-xxx/my-namespace/my-repo`
```
