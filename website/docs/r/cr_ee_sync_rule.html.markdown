---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_ee_sync_rule"
sidebar_current: "docs-alicloud-resource-cr-ee-sync-rule"
description: |-
  Provides a Alicloud Container Registry Enterprise Edition Sync Rule resource.
---

# alicloud_cr_ee_sync_rule

Provides a Container Registry Enterprise Edition Sync Rule resource.

For information about Container Registry Enterprise Edition Sync Rule and how to use it, see [What is Sync Rule](https://www.alibabacloud.com/help/en/acr/developer-reference/api-cr-2018-12-01-createreposyncrule)

-> **NOTE:** Available since v1.90.0.

-> **NOTE:** You need to set your registry password in Container Registry Enterprise Edition console before use this resource.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cr_ee_sync_rule&exampleId=3a398d4e-fdcf-c865-e1fc-5fbcabdf27b5a4eba53d&activeTab=example&spm=docs.r.cr_ee_sync_rule.0.3a398d4efd&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_regions" "default" {
  current = true
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_cr_ee_instance" "source" {
  payment_type   = "Subscription"
  period         = 1
  renew_period   = 0
  renewal_status = "ManualRenewal"
  instance_type  = "Advanced"
  instance_name  = "${var.name}-source-${random_integer.default.result}"
}

resource "alicloud_cr_ee_instance" "target" {
  payment_type   = "Subscription"
  period         = 1
  renew_period   = 0
  renewal_status = "ManualRenewal"
  instance_type  = "Advanced"
  instance_name  = "${var.name}-target-${random_integer.default.result}"
}

resource "alicloud_cr_ee_namespace" "source" {
  instance_id        = alicloud_cr_ee_instance.source.id
  name               = "${var.name}-${random_integer.default.result}"
  auto_create        = false
  default_visibility = "PUBLIC"
}

resource "alicloud_cr_ee_namespace" "target" {
  instance_id        = alicloud_cr_ee_instance.target.id
  name               = "${var.name}-${random_integer.default.result}"
  auto_create        = false
  default_visibility = "PUBLIC"
}

resource "alicloud_cr_ee_repo" "source" {
  instance_id = alicloud_cr_ee_instance.source.id
  namespace   = alicloud_cr_ee_namespace.source.name
  name        = "${var.name}-${random_integer.default.result}"
  summary     = "this is summary of my new repo"
  repo_type   = "PUBLIC"
}

resource "alicloud_cr_ee_repo" "target" {
  instance_id = alicloud_cr_ee_instance.target.id
  namespace   = alicloud_cr_ee_namespace.target.name
  name        = "${var.name}-${random_integer.default.result}"
  summary     = "this is summary of my new repo"
  repo_type   = "PUBLIC"
}

resource "alicloud_cr_ee_sync_rule" "default" {
  instance_id           = alicloud_cr_ee_instance.source.id
  namespace_name        = alicloud_cr_ee_namespace.source.name
  name                  = "${var.name}-${random_integer.default.result}"
  target_instance_id    = alicloud_cr_ee_instance.target.id
  target_namespace_name = alicloud_cr_ee_namespace.target.name
  target_region_id      = data.alicloud_regions.default.regions.0.id
  tag_filter            = ".*"
  repo_name             = alicloud_cr_ee_repo.source.name
  target_repo_name      = alicloud_cr_ee_repo.target.name
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the Container Registry Enterprise Edition source instance.
* `namespace_name` - (Required, ForceNew) The namespace name of the source instance.
* `name` - (Required, ForceNew) The name of the sync rule.
* `target_instance_id` - (Required, ForceNew) The ID of the destination instance.
* `target_namespace_name` - (Required, ForceNew) The namespace name of the destination instance.
* `target_region_id` - (Required, ForceNew) The region ID of the destination instance.
* `tag_filter` - (Required, ForceNew) The regular expression used to filter image tags.
* `repo_name` - (Optional, ForceNew) The image repository name of the source instance.
* `target_repo_name` - (Optional, ForceNew) The image repository name of the destination instance.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Sync Rule. It formats as `<instance_id>:<namespace_name>:<rule_id>`.
* `rule_id` - The ID of the sync rule.
* `sync_direction` - The synchronization direction.
* `sync_scope` - The synchronization scope.

## Import

Container Registry Enterprise Edition Sync Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_cr_ee_sync_rule.example <instance_id>:<namespace_name>:<rule_id>
```
