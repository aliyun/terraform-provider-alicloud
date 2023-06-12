---
subcategory: "Quotas"
layout: "alicloud"
page_title: "Alicloud: alicloud_quotas_template_quota"
sidebar_current: "docs-alicloud-resource-quotas-template-quota"
description: |-
  Provides a Alicloud Quotas Template Quota resource.
---

# alicloud_quotas_template_quota

Provides a Quotas Template Quota resource. 

For information about Quotas Template Quota and how to use it, see [What is Template Quota](https://www.alibabacloud.com/help/en/quota-center/latest/api-doc-quotas-2020-05-10-api-doc-createtemplatequotaitem).

-> **NOTE:** Available since v1.206.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}


resource "alicloud_quotas_template_quota" "default" {
  quota_action_code = "q_desktop-count"
  product_code      = "gws"
  notice_type       = 3
  dimensions {
    key   = "regionId"
    value = "cn-hangzhou"
  }
  desire_value   = 1001
  env_language   = "zh"
  quota_category = "CommonQuota"
}
```


## Argument Reference

The following arguments are supported:
* `desire_value` - (Required) Quota application value.
* `dimensions` - (Optional) The Quota Dimensions. See [`dimensions`](#dimensions) below.
* `effective_time` - (Optional) The UTC time when the quota takes effect.
* `env_language` - (Optional, Computed) The language of the quota alert notification. Value:
  - zh: Chinese.
  - en: English.
* `expire_time` - (Optional) The UTC time when the quota expires.
* `notice_type` - (Optional, Computed) Whether to notify the result of quota promotion application. Value:
  - 0: No.
  - 3: Yes.
* `product_code` - (Required, ForceNew) The abbreviation of the cloud service name.
* `quota_action_code` - (Required, ForceNew) The quota ID.
* `quota_category` - (Optional) Type of quota. Value:
  - CommonQuota : Generic quota.
  - WhiteListLabel: Equity quota.
  - FlowControl:API rate quota.

### `dimensions`

The Dimensions supports the following:
* `key` - (Optional) The Key of quota_dimensions.
* `value` - (Optional) The Value of quota_dimensions.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Template Quota.
* `delete` - (Defaults to 5 mins) Used when delete the Template Quota.
* `update` - (Defaults to 5 mins) Used when update the Template Quota.

## Import

Quotas Template Quota can be imported using the id, e.g.

```shell
$ terraform import alicloud_quotas_template_quota.example <id>
```