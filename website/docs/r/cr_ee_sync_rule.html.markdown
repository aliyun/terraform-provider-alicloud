---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_ee_sync_rule"
sidebar_current: "docs-alicloud-resource-cr-ee-sync-rule"
description: |-
  Provides a Alicloud resource to manage Container Registry Enterprise Edition sync rules.
---

# alicloud_cr_ee_sync_rule

This resource will help you to manager Container Registry Enterprise Edition sync rules.

For information about Container Registry Enterprise Edition sync rules and how to use it, see [Create a Sync Rule](https://www.alibabacloud.com/help/en/acr/developer-reference/api-cr-2018-12-01-createreposynctaskbyrule)

-> **NOTE:** Available since v1.90.0.

-> **NOTE:** You need to set your registry password in Container Registry Enterprise Edition console before use this resource.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_regions" "default" {
  current = true
}

resource "alicloud_cr_ee_instance" "source" {
  payment_type   = "Subscription"
  period         = 1
  renew_period   = 0
  renewal_status = "ManualRenewal"
  instance_type  = "Advanced"
  instance_name  = "${var.name}-source"
}

resource "alicloud_cr_ee_instance" "target" {
  payment_type   = "Subscription"
  period         = 1
  renew_period   = 0
  renewal_status = "ManualRenewal"
  instance_type  = "Advanced"
  instance_name  = "${var.name}-target"
}

resource "alicloud_cr_ee_namespace" "source" {
  instance_id        = alicloud_cr_ee_instance.source.id
  name               = var.name
  auto_create        = false
  default_visibility = "PUBLIC"
}

resource "alicloud_cr_ee_namespace" "target" {
  instance_id        = alicloud_cr_ee_instance.target.id
  name               = var.name
  auto_create        = false
  default_visibility = "PUBLIC"
}

resource "alicloud_cr_ee_repo" "source" {
  instance_id = alicloud_cr_ee_instance.source.id
  namespace   = alicloud_cr_ee_namespace.source.name
  name        = var.name
  summary     = "this is summary of my new repo"
  repo_type   = "PUBLIC"
  detail      = "this is a public repo"
}

resource "alicloud_cr_ee_repo" "target" {
  instance_id = alicloud_cr_ee_instance.target.id
  namespace   = alicloud_cr_ee_namespace.target.name
  name        = var.name
  summary     = "this is summary of my new repo"
  repo_type   = "PUBLIC"
  detail      = "this is a public repo"
}

resource "alicloud_cr_ee_sync_rule" "default" {
  instance_id           = alicloud_cr_ee_instance.source.id
  namespace_name        = alicloud_cr_ee_namespace.source.name
  name                  = var.name
  target_region_id      = data.alicloud_regions.default.regions.0.id
  target_instance_id    = alicloud_cr_ee_instance.target.id
  target_namespace_name = alicloud_cr_ee_namespace.target.name
  tag_filter            = ".*"
  repo_name             = alicloud_cr_ee_repo.source.name
  target_repo_name      = alicloud_cr_ee_repo.target.name
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the Container Registry Enterprise Edition source instance.
* `namespace_name` - (Required, ForceNew) The name of the Container Registry Enterprise Edition source namespace. It can contain 2 to 30 characters.
* `target_instance_id` - (Required, ForceNew) The ID of the Container Registry Enterprise Edition target instance to be synchronized.
* `target_namespace_name` - (Required, ForceNew) The name of the Container Registry Enterprise Edition target namespace to be synchronized. It can contain 2 to 30 characters.
* `target_region_id` - (Required, ForceNew) The target region to be synchronized.
* `name` - (Required, ForceNew) The name of the Container Registry Enterprise Edition sync rule.
* `tag_filter` - (Required, ForceNew) The regular expression used to filter image tags for synchronization in the source repository.
* `repo_name` - (Optional, ForceNew) The name of the source repository which should be set together with `target_repo_name`, if empty means that the synchronization scope is the entire namespace level.
* `target_repo_name` - (Optional, ForceNew) The name of the target repository.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Container Registry Enterprise Edition sync rule. It formats as `<instance_id>:<namespace_name>:<rule_id>`.
* `rule_id` - (Optional, ForceNew) The ID of the synchronization rule.
* `sync_direction` - The direction of the synchronization rule.
* `sync_scope` - The scope of the synchronization rule.

## Timeouts

-> **NOTE:** Available since v1.214.1.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Enterprise Edition sync rule.
* `delete` - (Defaults to 5 mins) Used when delete the Enterprise Edition sync rule.

## Import

Container Registry Enterprise Edition sync rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_cr_ee_sync_rule.example <instance_id>:<namespace_name>:<rule_id>
```
