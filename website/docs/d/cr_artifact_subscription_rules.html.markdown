---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_artifact_subscription_rules"
sidebar_current: "docs-alicloud-datasource-cr-artifact-subscription-rules"
description: |-
  Provides a list of Cr Artifact Subscription Rule owned by an Alibaba Cloud account.
---

# alicloud_cr_artifact_subscription_rules

This data source provides Cr Artifact Subscription Rule available to the user.[What is Artifact Subscription Rule](https://next.api.alibabacloud.com/document/cr/2018-12-01/CreateArtifactSubscriptionRule)

-> **NOTE:** Available since v1.287.0.

## Example Usage

```terraform
data "alicloud_cr_artifact_subscription_rules" "example" {
  instance_id    = "cri-xxx"
  enable_details = true
}
```

## Argument Reference

The following arguments are supported:
* `instance_id` - (Required, ForceNew) Instance ID
* `namespace_name` - (ForceNew, Optional) Namespace name
* `repo_name` - (ForceNew, Optional) Repository name
* `ids` - (Optional, Computed) A list of Artifact Subscription Rule IDs. The value is formulated as `<instance_id>:<artifact_subscription_rule_id>`.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Artifact Subscription Rule IDs.
* `rules` - A list of Artifact Subscription Rule Entries. Each element contains the following attributes:
  * `accelerate` - Whether to enable acceleration.
  * `artifact_subscription_rule_id` - The first ID of the resource.
  * `create_time` - Creation time.
  * `instance_id` - Instance ID.
  * `modified_time` - Modification time.
  * `namespace_name` - Namespace name.
  * `override` - Whether to override existing tags.
  * `platform` - Subscription platform list.
  * `region_id` - **NOTE:** This field is only available when `enable_details` is `true`. The region ID of the resource.
  * `repo_name` - Repository name.
  * `source_domain` - Source domain.
  * `source_namespace_name` - Source namespace name.
  * `source_provider` - Source image registry provider, e.
  * `source_repo_name` - Source repository name.
  * `tag_count` - Number of tags to subscribe.
  * `tag_regexp` - Regular expression for subscribing tags.
  * `id` - The ID of the resource supplied above.
