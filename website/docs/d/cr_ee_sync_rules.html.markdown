---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_ee_sync_rules"
sidebar_current: "docs-alicloud-datasource-cr-ee-sync-rules"
description: |-
  Provides a list of Container Registry Enterprise Edition sync rules.
---

# alicloud\_cr\_ee\_sync\_rules

This data source provides a list Container Registry Enterprise Edition sync rules on Alibaba Cloud.

-> **NOTE:** Available in v1.90.0+

## Example Usage

```
# Declare the data source
data "alicloud_cr_ee_sync_rules" "my_sync_rules" {
  instance_id = "cri-xxx"
  namespace_name = "test-namespace"
  repo_name = "test-repo"
  target_instance_id = "cri-yyy"
  name_regex = "test-rule"
}

output "output" {
  value = data.alicloud_cr_ee_sync_rules.my_sync_rules.rules.*.id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of Container Registry Enterprise Edition local instance.
* `namespace_name` - (Optional) Name of Container Registry Enterprise Edition local namespace.
* `repo_name` - (Optional) Name of Container Registry Enterprise Edition local repo.
* `target_instance_id` - (Optional) ID of Container Registry Enterprise Edition target instance.
* `name_regex` - (Optional) A regex string to filter results by sync rule name.
* `ids` - (Optional) A list of ids to filter results by sync rule id.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of matched Container Registry Enterprise Edition sync rules. Its element is a sync rule uuid.
* `names` - A list of sync rule names.
* `rules` - A list of matched Container Registry Enterprise Edition sync rules. Each element contains the following attributes:
  * `id` - ID of Container Registry Enterprise Edition sync rule.
  * `name` - Name of Container Registry Enterprise Edition sync rule.
  * `region_id` - Region of Container Registry Enterprise Edition local instance.
  * `instance_id` - ID of Container Registry Enterprise Edition local instance.
  * `namespace_name` - Name of Container Registry Enterprise Edition local namespace.
  * `repo_name` - Name of Container Registry Enterprise Edition local repo.
  * `target_region_id` - Region of Container Registry Enterprise Edition target instance.
  * `target_instance_id` - ID of Container Registry Enterprise Edition target instance.
  * `target_namespace_name` - Name of Container Registry Enterprise Edition target namespace.
  * `target_repo_name` - Name of Container Registry Enterprise Edition target repo.
  * `tag_filter` - The regular expression used to filter image tags for synchronization in the source repository.
  * `sync_direction` - `FROM` or `TO`, the direction of synchronization. `FROM` indicates that the local instance is the source instance. `TO` indicates that the local instance is the target instance to be synchronized.
  * `sync_scope` - `REPO` or `NAMESPACE`,the scope that the synchronization rule applies.
  * `sync_trigger` - `PASSIVE` or `INITIATIVE`, the policy configured to trigger the synchronization rule.
