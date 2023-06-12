---
subcategory: "Quotas"
layout: "alicloud"
page_title: "Alicloud: alicloud_quotas_quota_application"
description: |-
  Provides a Alicloud Quotas Quota Application resource.
---

# alicloud_quotas_quota_application

Provides a Quotas Quota Application resource. Details of Quota Application.

For information about Quotas Quota Application and how to use it, see [What is Quota Application](https://www.alibabacloud.com/help/en/quota-center/latest/api-doc-quotas-2020-05-10-api-doc-createquotaapplication).

-> **NOTE:** Available since v1.117.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}


resource "alicloud_quotas_quota_application" "default" {
  quota_action_code = "q_desktop-count"
  product_code      = "gws"
  quota_category    = "CommonQuota"
  notice_type       = 3
  dimensions {
    key   = "regionId"
    value = "cn-hangzhou"
  }
  desire_value = 1001
  reason       = "测试"
  env_language = "zh"
}
```

### Deleting `alicloud_quotas_quota_application` or removing it from your configuration

Terraform cannot destroy resource `alicloud_quotas_quota_application`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `audit_mode` - (Optional, ForceNew, Computed) Quota audit mode. Value:
  - Sync: Synchronize auditing. The quota center automatically approves, and the approval result is returned immediately, but the probability of application passing is lower than that of asynchronous approval, and the validity period of the increase quota is 1 hour.
  - Async: Asynchronous auditing. Manual review, the probability of application passing is relatively high, and the validity period of the increase quota is 1 month.
-> **NOTE:**  This parameter takes effect only for the ECS specification quota of the cloud server.
* `desire_value` - (Required, ForceNew) The desire value of the quota application.
* `dimensions` - (Optional, ForceNew) QuotaDimensions. See [`dimensions`](#dimensions) below.
* `effective_time` - (Optional, ForceNew) The effective time of the quota application.
* `env_language` - (Optional, ForceNew, Available in v1.206.0+) The language of the quota alert notification. Value:
  - zh (default): Chinese.
  - en: English.
* `expire_time` - (Optional, ForceNew) The expired time of the quota application.
* `notice_type` - (Optional, ForceNew, Computed) Specifies whether to send a notification about the application result. Valid values:0: sends a notification about the application result.3: A notification about the application result is sent.
* `product_code` - (Required, ForceNew) The product code.
* `quota_action_code` - (Required, ForceNew) The ID of quota action.
* `quota_category` - (Optional, ForceNew) The quota type.
  - CommonQuota (default): Generic quota.
  - FlowControl:API rate quota.
  - WhiteListLabel: Equity quota.
* `reason` - (Required, ForceNew) The reason of the quota application.


### `dimensions`

The dimensions support the following:
* `key` - (Optional, ForceNew) Key.
* `value` - (Optional, ForceNew) Value.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `approve_value` - The approve value of the quota application.
* `audit_reason` - The audit reason.
* `create_time` - Resource attribute field representing creation time.
* `quota_description` - The description of the quota application.
* `quota_name` - The name of the quota application.
* `quota_unit` - The unit of the quota application.
* `status` - Application Status:
  - Disagree: reject.
  - Agree: Approved.
  - Process: under review.
  - Cancel: Closed.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Quota Application.

## Import

Quotas Quota Application can be imported using the id, e.g.

```shell
$ terraform import alicloud_quotas_quota_application.example <id>
```