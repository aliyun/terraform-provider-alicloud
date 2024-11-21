---
subcategory: "Quotas"
layout: "alicloud"
page_title: "Alicloud: alicloud_quotas_quota_application"
description: |-
  Provides a Alicloud Quotas Quota Application resource.
---

# alicloud_quotas_quota_application

Provides a Quotas Quota Application resource. Details of Quota Application.

For information about Quotas Quota Application and how to use it, see [What is Quota Application](https://www.alibabacloud.com/help/en/quota-center/developer-reference/api-quotas-2020-05-10-createquotaapplication).

-> **NOTE:** Available since v1.117.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_quotas_quota_application&exampleId=3814f509-3245-3eca-4c43-9628e29a123cade386e9&activeTab=example&spm=docs.r.quotas_quota_application.0.3814f50932&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
* `audit_mode` - (Optional, ForceNew, Computed) This parameter is discontinued and is not recommended. The mode in which you want the application to be reviewed. Valid values:
  - Sync: The application is reviewed in a synchronous manner. Quota Center automatically reviews the application. The result is returned immediately after you submit the application. However, the chance of an approval for an application that is reviewed in Sync mode is lower than the chance of an approval for an application that is reviewed in Async mode. The validity period of the new quota value is 1 hour.
  - Async: The application is reviewed in an asynchronous manner. An Alibaba Cloud support engineer reviews the application. The chance of an approval for an application that is reviewed in Async mode is higher than the chance of an approval for an application that is reviewed in Sync mode. The validity period of the new quota value is one month.
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
* `key` - (Optional, ForceNew) The key of the dimension. You must configure `dimensions.N.key` and `dimensions.N.value` at the same time. The value range of N varies based on the number of dimensions that are supported by the related Alibaba Cloud service. You can call the [ListProductQuotaDimensions](https://next.api.aliyun.com/document/quotas/2020-05-10/ListProductQuotaDimensions) operation to query the dimensions that are supported by an Alibaba Cloud service. The number of elements in the returned array is N.
* `value` - (Optional, ForceNew) The value of the dimension. You must configure `dimensions.N.key` and `dimensions.N.value` at the same time. The value range of N varies based on the number of dimensions that are supported by the related Alibaba Cloud service. You can call the [ListProductQuotaDimensions](https://next.api.aliyun.com/document/quotas/2020-05-10/ListProductQuotaDimensions) operation to query the dimensions that are supported by an Alibaba Cloud service. The number of elements in the returned array is N.


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