---
subcategory: "Quotas"
layout: "alicloud"
page_title: "Alicloud: alicloud_quotas_quota_applications"
sidebar_current: "docs-alicloud-datasource-quotas-quota-applications"
description: |-
  Provides a list of Quotas Quota Applications to the user.
---

# alicloud_quotas_quota_applications

This data source provides the Quotas Quota Applications of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.117.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_quotas_quota_application" "default" {
  product_code      = "vpc"
  notice_type       = "3"
  effective_time    = "2023-05-22T16:00:00Z"
  expire_time       = "2024-09-15T00:08:32Z"
  desire_value      = "1"
  reason            = ""
  quota_action_code = "vpc_whitelist/ha_vip_whitelist"
  audit_mode        = "Sync"
  env_language      = "zh"
  quota_category    = "WhiteListLabel"
}

data "alicloud_quotas_quota_applications" "default" {
  product_code   = "vpc"
  enable_details = "true"
  quota_category = alicloud_quotas_quota_application.default.quota_category
  ids = [
    "${alicloud_quotas_quota_application.default.id}"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `dimensions` - (Optional) The quota dimensions.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Application Info IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `product_code` - (Required, ForceNew) The product code.
* `quota_action_code` - (Optional, ForceNew) The ID of quota action.
* `quota_category` - (Optional, ForceNew) The quota category. Valid values: `CommonQuota`, `FlowControl`, `WhiteListLabel`.
* `status` - (Optional, ForceNew) The status of the quota application. Valid Values: `Agree`, `Disagree` and `Process`.

#### Block dimensions

The dimensions supports the following: 

* `key` - (Optional, ForceNew) The key of dimensions.
* `value` - (Optional, ForceNew) The value of dimensions.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `applications` - A list of Quotas Quota Applications. Each element contains the following attributes:
  * `application_id` - The first ID of the resource.
  * `approve_value` - The approve value.
  * `audit_reason` - The audit reason.
  * `desire_value` - The desire value of the quota application.
  * `dimensions` - The quota dimensions.
    * `key` - The key of dimensions.
    * `value` - The value of dimensions.
  * `effective_time` - The effective time.
  * `expire_time` - The expire time.
  * `id` - The ID of the Application Info.
  * `notice_type` - The notice type.
  * `product_code` - The product code.
  * `quota_action_code` - The ID of quota action..
  * `quota_description` - The description of the quota.
  * `quota_name` - The name of the quota.
  * `quota_unit` - The quota unit.
  * `reason` - The reason of the quota application.
  * `status` - The status of the quota application.
