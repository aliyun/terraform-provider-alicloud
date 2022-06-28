---
subcategory: "Quotas"
layout: "alicloud"
page_title: "Alicloud: alicloud_quotas_quota_application"
sidebar_current: "docs-alicloud-resource-quotas-quota-application"
description: |-
  Provides a Alicloud Quotas Quota Application resource.
---

# alicloud\_quotas\_quota\_application

Provides a Quotas Quota Application resource.

For information about Quotas Quota Application and how to use it, see [What is Quota Application](https://help.aliyun.com/document_detail/171289.html).

-> **NOTE:** Available in v1.117.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_quotas_quota_application" "example" {
  notice_type       = "0"
  desire_value      = "100"
  product_code      = "ess"
  quota_action_code = "q_db_instance"
  reason            = "For Terraform Test"
  dimensions {
    key   = "regionId"
    value = "cn-hangzhou"
  }
}

```

## Argument Reference

The following arguments are supported:

* `desire_value` - (Required, ForceNew) The desire value of the quota application.
* `dimensions` - (Optional, ForceNew) The quota dimensions.
* `notice_type` - (Optional, ForceNew) The notice type. Valid values: `0`, `1`, `2`, `3`.
* `product_code` - (Required, ForceNew) The product code.
* `quota_action_code` - (Required, ForceNew) The ID of quota action.
* `quota_category` - (Optional, ForceNew) The quota category. Valid values: `CommonQuota`, `FlowControl`.
* `reason` - (Required, ForceNew) The reason of the quota application.
* `audit_mode` - (Required, ForceNew) The audit mode. Valid values: `Async`, `Sync`. Default to: `Async`.

#### Block dimensions

The dimensions supports the following: 

* `key` - (Optional, ForceNew) The key of dimensions.
* `value` - (Optional, ForceNew) The value of dimensions.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Quota Application.
* `status` - The status of the quota application.
* `approve_value` - The approve value of the quota application.
* `audit_reason` - The audit reason.
* `effective_time` - The effective time of the quota application.
* `expire_time` - The expire time of the quota application.
* `quota_description` - The description of the quota application.
* `quota_name` - The name of the quota application.
* `quota_unit` - The unit of the quota application.

## Import

Quotas Application Info can be imported using the id, e.g.

```
$ terraform import alicloud_quotas_quota_application.example <id>
```
