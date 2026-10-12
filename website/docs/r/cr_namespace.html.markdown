---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_namespace"
sidebar_current: "docs-alicloud-resource-container-registry"
description: |-
  Provides a Alicloud resource to manage Container Registry namespaces.
---

# alicloud_cr_namespace

Provides a CR Namespace resource for an Enterprise Container Registry instance.

For information about CR Namespace and how to use it, see [CreateNamespace](https://www.alibabacloud.com/help/en/acr/developer-reference/api-cr-2018-12-01-createnamespace).

-> **NOTE:** Available since v1.287.0. The resource is rebuilt against the `cr` 2018-12-01 POP API for instance-scoped (Enterprise Edition) namespaces. The legacy personal-edition fields `name` / `auto_create` / `default_visibility` are no longer supported — use `namespace_name` / `auto_create_repo` / `default_repo_type` together with `instance_id` instead. The resource ID format changed from the namespace name to `<instance_id>:<namespace_name>`, and import now requires both segments.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_cr_ee_instance" "example" {
  instance_name  = "terraform-example"
  instance_type  = "Basic"
  payment_type   = "Subscription"
  period         = 1
}

resource "alicloud_cr_namespace" "example" {
  instance_id       = alicloud_cr_ee_instance.example.id
  namespace_name    = "${var.name}-ns"
  auto_create_repo  = false
  default_repo_type = "PRIVATE"

  default_repo_configuration {
    repo_type        = "PRIVATE"
    tag_immutability = false
    artifact_build_rule_parameters {
      image_index_only = false
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the Enterprise Container Registry instance.
* `namespace_name` - (Required, ForceNew) The name of the namespace.
* `auto_create_repo` - (Optional) Whether to automatically create a repository when a new image is pushed.
* `default_repo_type` - (Optional) Default repository type in this namespace. Valid values: `PUBLIC`, `PRIVATE`.
* `default_repo_configuration` - (Optional) Default repository configuration of the namespace. See [`default_repo_configuration`](#default_repo_configuration) below.

### default_repo_configuration

The `default_repo_configuration` block supports:

* `repo_type` - (Required) Default repository type. Valid values: `PRIVATE`, `PUBLIC`.
* `tag_immutability` - (Required) Whether the tag is immutable.
* `artifact_build_rule_parameters` - (Required) Artifact build rule parameters. See [`artifact_build_rule_parameters`](#artifact_build_rule_parameters) below.

-> **NOTE:** `default_repo_configuration` is supplemented from the CloudSpec pre model (`alicloud.CR.cr.v20181201/Namespace.cspec`) because the code generator does not emit struct-reference attributes. It is serialized as a single JSON object on create/update and flattened from the `GetNamespace` response on read.

#### artifact_build_rule_parameters

The `artifact_build_rule_parameters` block supports:

* `image_index_only` - (Required) Whether only the image index is built.
* `priority_file` - (Optional) Priority file.来源于 CloudSpec pre 模型(线上 prod RepoConfiguration schema 未含此字段,留空时不发送)。

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. The value is formulated as `<instance_id>:<namespace_name>`.
* `namespace_id` - The ID of the namespace.
* `resource_group_id` - The resource group ID.
* `status` - The status of the namespace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when creating the namespace.
* `delete` - (Defaults to 5 mins) Used when deleting the namespace.
* `update` - (Defaults to 5 mins) Used when updating the namespace.

## Import

CR Namespace can be imported using the id, e.g.

```shell
$ terraform import alicloud_cr_namespace.example <instance_id>:<namespace_name>
```
