---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_repository"
description: |-
  Provides an Alicloud CR Repository resource.
---

# alicloud_cr_repository

Provides a Container Registry (CR) Repository resource.

For information about CR Repository and how to use it, see [What is Repository](https://next.api.alibabacloud.com/document/cr/2018-12-01/CreateRepository).

-> **NOTE:** Available since v2.2.0.

## Example Usage

Basic Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_cr_ee_instance" "default" {
  instance_name  = var.name
  payment_type   = "Subscription"
  renewal_status = "ManualRenewal"
  image_scanner  = "ACR"
  period         = "1"
  instance_type  = "Basic"
}

resource "alicloud_cr_ee_namespace" "default" {
  instance_id        = alicloud_cr_ee_instance.default.id
  name               = var.name
  auto_create        = false
  default_visibility = "PRIVATE"
}

resource "alicloud_cr_repository" "default" {
  instance_id         = alicloud_cr_ee_instance.default.id
  repo_namespace_name = alicloud_cr_ee_namespace.default.name
  repository_name     = var.name
  repo_type           = "PUBLIC"
  summary             = "terraform-example"
  detail              = "terraform-example"
  tag_immutability    = false
}
```

## Argument Reference

The following arguments are supported:

* `detail` - (Optional, String) The description of the repository.
* `instance_id` - (Required, ForceNew, String) The ID of the Container Registry instance.
* `repo_namespace_name` - (Required, ForceNew, String) The name of the namespace the repository belongs to.
* `repo_type` - (Required, ForceNew, String) The repository type. Valid values: `PUBLIC` and `PRIVATE`.
* `repository_name` - (Required, ForceNew, String) The name of the repository.
* `summary` - (Required, String) The summary of the repository.
* `tag_immutability` - (Optional, Bool) Specifies whether image tags in the repository are immutable.

## Attributes Reference

The following attributes are exported:

* `id` - (String) The ID of the resource. The value is formulated as `<instance_id>:<repository_id>`.
* `create_time` - (String) The time when the repository was created.
* `modified_time` - (String) The time when the repository was last modified.
* `repo_build_type` - (String) The build type of the repository.
* `repository_id` - (String) The ID of the repository.
* `resource_group_id` - (String) The ID of the resource group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (String) Defaults to `5m`. Used when creating the Repository.
* `update` - (String) Defaults to `5m`. Used when updating the Repository.
* `delete` - (String) Defaults to `5m`. Used when deleting the Repository.

## Import

CR Repository can be imported using the id, e.g.

```shell
$ terraform import alicloud_cr_repository.example <instance_id>:<repository_id>
```
