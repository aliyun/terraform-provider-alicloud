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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_cr_ee_sync_rule&exampleId=bbb51f23-85cb-11bd-570e-c8ee82b7b1f979f7fa82&activeTab=example&spm=docs.r.cr_ee_sync_rule.0.bbb51f2385&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
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

data "alicloud_regions" "default" {
  current = true
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

* `instance_id` - (Required, ForceNew) ID of Container Registry Enterprise Edition source instance.
* `namespace_name` - (Required, ForceNew) Name of Container Registry Enterprise Edition source namespace. It can contain 2 to 30 characters.
* `name` - (Required, ForceNew) Name of Container Registry Enterprise Edition sync rule.
* `target_region_id` - (Required, ForceNew) The target region to be synchronized.
* `target_instance_id` - (Required, ForceNew) ID of Container Registry Enterprise Edition target instance to be synchronized.
* `target_namespace_name` - (Required, ForceNew) Name of Container Registry Enterprise Edition target namespace to be synchronized. It can contain 2 to 30 characters.
* `tag_filter` - (Required, ForceNew) The regular expression used to filter image tags for synchronization in the source repository.
* `repo_name` - (Optional, ForceNew) Name of the source repository which should be set together with `target_repo_name`, if empty means that the synchronization scope is the entire namespace level.
* `target_repo_name` - (Optional, ForceNew) Name of the target repository.
* `rule_id` - (Optional, ForceNew) The uuid of Container Registry Enterprise Edition sync rule.

## Attributes Reference

The following attributes are exported:

* `id` - The resource id of Container Registry Enterprise Edition sync rule. The value is in format `{instance_id}:{namespace_name}:{rule_id}`.
* `sync_direction` - `FROM` or `TO`, the direction of synchronization. `FROM` means source instance, `TO` means target instance.
* `sync_scope` - `REPO` or `NAMESPACE`,the scope that the synchronization rule applies.

## Import

Container Registry Enterprise Edition sync rule can be imported using the id. Format to `{instance_id}:{namespace_name}:{rule_id}`, e.g.

```shell
$ terraform import alicloud_cr_ee_sync_rule.default `cri-xxx:my-namespace:crsr-yyy`
```
