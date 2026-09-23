---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_artifact_lifecycle_rules"
sidebar_current: "docs-alicloud-datasource-cr-artifact-lifecycle-rules"
description: |-
  Provides a list of Cr Artifact Lifecycle Rule owned by an Alibaba Cloud account.
---

# alicloud_cr_artifact_lifecycle_rules

This data source provides Cr Artifact Lifecycle Rule available to the user.[What is Artifact Lifecycle Rule](https://next.api.alibabacloud.com/document/cr/2018-12-01/CreateArtifactLifecycleRule)

-> **NOTE:** Available since v1.285.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_cr_ee_instance" "default" {
  default_oss_bucket = "true"
  instance_name      = var.name
  renewal_status     = "ManualRenewal"
  image_scanner      = "DISABLE"
  period             = "1"
  payment_type       = "Subscription"
  instance_type      = "Economy"
}

resource "alicloud_cr_ee_namespace" "default" {
  instance_id        = alicloud_cr_ee_instance.default.id
  name               = var.name
  auto_create        = false
  default_visibility = "PRIVATE"
}

resource "alicloud_cr_ee_repo" "default" {
  instance_id = alicloud_cr_ee_instance.default.id
  namespace   = alicloud_cr_ee_namespace.default.name
  name        = var.name
  repo_type   = "PRIVATE"
  summary     = "example repository for lifecycle rule"
}

resource "alicloud_cr_artifact_lifecycle_rule" "default" {
  auto                = true
  namespace_name      = alicloud_cr_ee_namespace.default.name
  retention_tag_count = "30"
  schedule_time       = "WEEK"
  scope               = "REPO"
  instance_id         = alicloud_cr_ee_instance.default.id
  tag_regexp          = ".*"
  repo_name           = alicloud_cr_ee_repo.default.name
}

data "alicloud_cr_artifact_lifecycle_rules" "default" {
  ids         = ["${alicloud_cr_artifact_lifecycle_rule.default.id}"]
  instance_id = alicloud_cr_ee_instance.default.id
}

output "alicloud_cr_artifact_lifecycle_rule_example_id" {
  value = data.alicloud_cr_artifact_lifecycle_rules.default.rules.0.id
}
```

## Argument Reference

The following arguments are supported:
* `instance_id` - (Required) Instance ID
* `ids` - (Optional, Computed) A list of Artifact Lifecycle Rule IDs. The value is formulated as `<instance_id>:<artifact_lifecycle_rule_id>`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Artifact Lifecycle Rule IDs.
* `rules` - A list of Artifact Lifecycle Rule Entries. Each element contains the following attributes:
    * `artifact_lifecycle_rule_id` - The first ID of the resource.
    * `auto` - Whether to execute automatically.
    * `create_time` - Creation time.
    * `enable_delete_tag` - Activate the delete tag function.
    * `enable_delete_untagged_manifest` - Open garbage collection.
    * `instance_id` - Instance ID.
    * `modified_time` - Change time.
    * `namespace_name` - Namespace name.
    * `repo_name` - Repository Name.
    * `retention_tag_count` - Number of Retention Tags.
    * `schedule_time` - Execution cycle.
    * `scope` - Scope of cleaning.
    * `tag_regexp` - Retain regular expressions for mirrored versions.
    * `id` - The ID of the resource supplied above.
