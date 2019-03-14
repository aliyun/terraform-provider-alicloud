---
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_repo"
sidebar_current: "docs-alicloud-resource-container-registry"
description: |-
  Provides a Alicloud resource to manage container registry repos.
---

# alicloud\_cr\_repo

This resource will help you to manager container registry repos.

-> **NOTE:** Available in v1.35.0+.

## Example Usage

Basic Usage

```
resource "alicloud_cr_namespace" "my-namespace" {
    name = "my-namespace"
    auto_create = false
    default_visibility = "PUBLIC"
}

resource "alicloud_cr_repo" "my-repo" {
    namespace = "${alicloud_cr_namespace.my-namespace.name}"
    name = "my-repo"
    summary = "this is summary of my new repo"
    repo_type = "PUBLIC"
    detail  = "this is a public repo"
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Required, ForceNew) Name of container registry namespace where repo is located.
* `name` - (Required, ForceNew) Name of container registry repo.
* `summary` - (Required) The repository general information. It can contain 1 to 80 characters.
* `repo_type` - (Required) `PUBLIC` or `PRIVATE`, repo's visibility.
* `detail` - (Optional) The repository spesific information. MarkDown format is supported, and the length limit is 2000.

## Attributes Reference

The following attributes are exported:

* `id` - The id of container registry repo. The value is in format `namespace/name`.
* `domain_list` - The repository domain list.
  * `public` - Domain of public endpoint.
  * `internal` - Domain of internal endpoint, only in some regions.
  * `vpc` - Domain of vpc endpoint.

## Import

Container Registry Repo can be imported using the `namespace/name`, e.g.

```
$ terraform import alicloud_cr_repo.default `my-namespace/my-repo`
```