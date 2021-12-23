---
subcategory: "Open Search"
layout: "alicloud"
page_title: "Alicloud: alicloud_open_search_app_groups"
sidebar_current: "docs-alicloud-datasource-open-search-app-groups"
description: |-
  Provides a list of Open Search App Groups to the user.
---

# alicloud\_open\_search\_app\_groups

This data source provides the Open Search App Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.136.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_testacc"
}
resource "alicloud_open_search_app_group" "default" {
  app_group_name = var.name
  payment_type   = "PayAsYouGo"
  type           = "standard"
  quota {
    doc_size         = 1
    compute_resource = 20
    spec             = "opensearch.share.common"
  }
}
data "alicloud_open_search_app_groups" "default" {
  ids = [alicloud_open_search_app_group.default.id]
}
output "app_groups" {
  value = data.alicloud_open_search_app_groups.default.groups
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of App Group IDs. Its element value is same as App Group Name.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by App Group name.
* `instance_id` - (Optional, ForceNew) Instance ID.
* `resource_group_id` - (Optional, ForceNew) Resource Group Id.
* `type` - (Optional, ForceNew) The Type of AppGroup. Valid values: `standard`,`enhanced`.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `app_group_id` - The ID of the App Group.
* `id` - The resource ID in terraform of App Group. It is the same as `app_group_name`.
* `app_group_name` - (Required,ForceNew) Application Group Name.
* `payment_type` - (Required,ForceNew) The billing method of the resource. Valid values: `Subscription` and `PayAsYouGo`.
* `charge_way` - (Optional) Billing model. Valid values:`compute_resource` and `qps`.
* `commodity_code` - The commodity code.
* `create_time` - The time of creation.
* `current_version` - The version of Application Group Name.
* `description` - The description of the resource.
* `domain` - Domain name.
* `expire_on` - Expiration Time.
* `first_rank_algo_deployment_id` - Coarse deployment ID.
* `has_pending_quota_review_task` - Whether the quota status is under approval. Valid status:
    * `0`: normal
    * `1`: Approving.
* `instance_id` - The Instance ID.
* `lock_mode` - Locked state. Valid status: `Unlock`,`LockByExpiration`,`ManualLock`.
* `locked_by_expiration` - Instance is automatically locked after expiration.
* `pending_second_rank_algo_deployment_id` - Refine deployment ID in deployment.
* `processing_order_id` - Unfinished order number.
* `produced` - Whether the production is completed. Valid values:
    * `0`: producing.
    * `1`: completed.
* `project_id` - The Project ID.
* `quota` - Quota information.
* `resource_group_id` - The Resource Group ID.
* `second_rank_algo_deployment_id` - Refine deployment ID.
* `status` - The status of the resource. Valid values: `producing`,`review_pending`,`config_pending`,`normal`,`frozen`.
* `switched_time` - The Switched time.
* `type` - Application type. Valid Values: `standard`, `enhanced`.


#### quota
The quota supports the following:

* `doc_size` - (Required) Storage Size. Unit: GB.
* `compute_resource` - (Required) Computing resources. Unit: LCU.
* `qps` - (Required) Search request. Unit: times/second.
* `spec` - (Required) Specification. Valid values:
	* `opensearch.share.junior`: Entry-level.
	* `opensearch.share.common`: Shared universal.
	* `opensearch.share.compute`: Shared computing.
	* `opensearch.share.storage`: Shared storage type.
	* `opensearch.private.common`: Exclusive universal type.
	* `opensearch.private.compute`: Exclusive computing type.
	* `opensearch.private.storage`: Exclusive storage type
