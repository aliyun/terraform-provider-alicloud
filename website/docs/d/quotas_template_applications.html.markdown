---
subcategory: "Quotas"
layout: "alicloud"
page_title: "Alicloud: alicloud_quotas_template_applications"
sidebar_current: "docs-alicloud-datasource-quotas-template-applications"
description: |-
  Provides a list of Quotas Template Applications owned by an Alibaba Cloud account.
---

# alicloud_quotas_template_applications

This data source provides Quotas Template Applications available to the user.[What is Template Applications](https://www.alibabacloud.com/help/en/quota-center/developer-reference/api-quotas-2020-05-10-createquotaapplicationsfortemplate)

-> **NOTE:** Available since v1.214.0.

## Example Usage

```terraform
data "alicloud_resource_manager_accounts" "default" {
  status = "CreateSuccess"
}

resource "alicloud_quotas_template_applications" "default" {
  quota_action_code = "vpc_whitelist/ha_vip_whitelist"
  product_code      = "vpc"
  quota_category    = "FlowControl"
  aliyun_uids       = ["${data.alicloud_resource_manager_accounts.default.ids.0}"]
  desire_value      = 6
  notice_type       = "0"
  env_language      = "zh"
  reason            = "example"
  dimensions {
    key   = "apiName"
    value = "GetProductQuotaDimension"
  }
  dimensions {
    key   = "apiVersion"
    value = "2020-05-10"
  }
  dimensions {
    key   = "regionId"
    value = "cn-hangzhou"
  }
}

data "alicloud_quotas_template_applications" "default" {
  ids               = ["${alicloud_quotas_template_applications.default.id}"]
  product_code      = "vpc"
  quota_action_code = "vpc_whitelist/ha_vip_whitelist"
  quota_category    = "FlowControl"
}

output "alicloud_quotas_template_applications_example_id" {
  value = data.alicloud_quotas_template_applications.default.applications.0.id
}
```

## Argument Reference

The following arguments are supported:
* `batch_quota_application_id` - (ForceNew, Optional) The ID of the quota application batch.
* `product_code` - (ForceNew, Optional) Cloud service name abbreviation.> For more information about cloud services that support quota centers, see [Cloud services that support quota centers](~~ 182368 ~~).
* `quota_action_code` - (ForceNew, Optional) The quota ID.
* `quota_category` - (ForceNew, Optional) The quota type. Value:-CommonQuota (default): Generic quota.-FlowControl:API rate quota.-WhiteListLabel: Equity quota.
* `ids` - (Optional, ForceNew, Computed) A list of Template Applications IDs.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Template Applications IDs.
* `applications` - A list of Template Applications Entries. Each element contains the following attributes:
  * `id` - The ID of the quota application batch.
  * `aliyun_uids` - The list of Alibaba Cloud accounts (primary accounts) of the resource directory members to which the quota is applied.> Only 50 members can apply for quota increase in batch at a time. For more information about the members of the resource directory, see [Query the list of all members in the resource directory](~~ 604207 ~~).
  * `apply_time` - The UTC time of the quota increase application.
  * `audit_status_vos` - Quantity of requisitions in different approval statuses.
    * `count` - Approval document quantity.
    * `status` - The approval status of the quota promotion application. Value:-Disagree: reject.-Approve: approved.-Process: under review.-Cancel: Closed.
  * `batch_quota_application_id` - The ID of the quota application batch.
  * `desire_value` - The value of the quota request.> The quota request is approved by the technical support of each cloud service. If you want to increase the chance of passing, please fill in a reasonable application value and detailed application reasons when applying for quota.
  * `dimensions` - Quota dimension.
    * `key` - Quota dimension Key.
    * `value` - Quota dimension Value.
  * `effective_time` - The UTC time when the quota takes effect. This parameter applies only to the equity quota (WhiteListLabel).> If the current account does not select the effective time, the default is the submission time.
  * `expire_time` - The UTC time when the quota expires. This parameter applies only to the equity quota (WhiteListLabel).> If No Expiration Time is selected for the current account, the expiration time is 99 years from the effective time of the current quota.
  * `product_code` - Cloud service name abbreviation.> For more information about cloud services that support quota centers, see [Cloud services that support quota centers](~~ 182368 ~~).
  * `quota_action_code` - The quota ID.
  * `quota_category` - The quota type. Value:-CommonQuota (default): Generic quota.-FlowControl:API rate quota.-WhiteListLabel: Equity quota.
  * `reason` - Reason for quota application.> The quota request is approved by the technical support of each cloud service. If you want to increase the chance of passing, please fill in a reasonable application value and detailed application reasons when applying for quota.
