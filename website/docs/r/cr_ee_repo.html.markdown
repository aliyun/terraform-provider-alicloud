---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_ee_repo"
sidebar_current: "docs-alicloud-resource-cr-ee-repo"
description: |-
  Provides a Alicloud resource to manage Container Registry Enterprise Edition repositories.
---

# alicloud_cr_ee_repo

This resource will help you to manager Container Registry Enterprise Edition repositories.

For information about Container Registry Enterprise Edition repository and how to use it, see [Create a Repository](https://www.alibabacloud.com/help/en/acr/developer-reference/api-cr-2018-12-01-createrepository)

-> **NOTE:** Available since v1.86.0.

-> **NOTE:** You need to set your registry password in Container Registry Enterprise Edition console before use this resource.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}
resource "alicloud_cr_ee_instance" "example" {
  payment_type   = "Subscription"
  period         = 1
  renew_period   = 0
  renewal_status = "ManualRenewal"
  instance_type  = "Advanced"
  instance_name  = var.name
}

resource "alicloud_cr_ee_namespace" "example" {
  instance_id        = alicloud_cr_ee_instance.example.id
  name               = var.name
  auto_create        = false
  default_visibility = "PUBLIC"
}

resource "alicloud_cr_ee_repo" "example" {
  instance_id = alicloud_cr_ee_instance.example.id
  namespace   = alicloud_cr_ee_namespace.example.name
  name        = var.name
  summary     = "this is summary of my new repo"
  repo_type   = "PUBLIC"
  detail      = "this is a public repo"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) ID of Container Registry Enterprise Edition instance.
* `namespace` - (Required, ForceNew) Name of Container Registry Enterprise Edition namespace where repository is located. It can contain 2 to 30 characters.
* `name` - (Required, ForceNew) Name of Container Registry Enterprise Edition repository. It can contain 2 to 64 characters.
* `summary` - (Required) The repository general information. It can contain 1 to 100 characters.
* `repo_type` - (Required) `PUBLIC` or `PRIVATE`, repo's visibility.
* `detail` - (Optional) The repository specific information. MarkDown format is supported, and the length limit is 2000.

## Attributes Reference

The following attributes are exported:

* `id` - The resource id of Container Registry Enterprise Edition repository. The value is in format `{instance_id}:{namespace}:{repository}`.
* `repo_id` - The uuid of Container Registry Enterprise Edition repository.

## Import

Container Registry Enterprise Edition repository can be imported using the `{instance_id}:{namespace}:{repository}`, e.g.

```shell
$ terraform import alicloud_cr_ee_repo.default `cri-xxx:my-namespace:my-repo`
```
