---
subcategory: "unreleased"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_repository"
description: |-
  Provides a Alicloud CR Repository resource.
---

# alicloud_cr_repository

Provides a CR Repository resource.



For information about CR Repository and how to use it, see [What is Repository](https://next.api.alibabacloud.com/document/cr/2018-12-01/CreateRepository).

-> **NOTE:** Available since v2.1.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_cr_ee_instance" "default1TiZyc" {
  instance_name  = "auto-example-0228-2"
  renewal_status = "ManualRenewal"
  image_scanner  = "ACR"
  period         = "1"
  instance_type  = "Basic"
}

resource "alicloud_cr_namespace" "defaultRZzh9M" {
}


resource "alicloud_cr_repository" "default" {
  repo_namespace_name = alicloud_cr_namespace.defaultRZzh9M.namespace_name
  tag_immutability    = false
  repo_type           = "PUBLIC"
  instance_id         = alicloud_cr_ee_instance.default1TiZyc.id
  summary             = "sdad"
  repository_name     = "auto-example-repo2"
  detail              = "ds"
}
```

## Argument Reference

The following arguments are supported:
* `detail` - (Optional) Detail
* `instance_id` - (Required, ForceNew) InstanceId
* `repo_namespace_name` - (Required, ForceNew) RepoNamespaceName
* `repo_type` - (Required, ForceNew) RepoType
* `repository_name` - (Required, ForceNew) RepoName
* `summary` - (Required) Summary
* `tag_immutability` - (Optional) TagImmutability

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<instance_id>:<repository_id>`.
* `create_time` - The creation time of the resource.
* `modified_time` - Last modification time.
* `repo_build_type` - repository build type.
* `repository_id` - The first ID of the resource.
* `resource_group_id` - The ID of the resource group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Repository.
* `delete` - (Defaults to 5 mins) Used when delete the Repository.
* `update` - (Defaults to 5 mins) Used when update the Repository.

## Import

CR Repository can be imported using the id, e.g.

```shell
$ terraform import alicloud_cr_repository.example <instance_id>:<repository_id>
```