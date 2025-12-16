---
subcategory: "Quotas"
layout: "alicloud"
page_title: "Alicloud: alicloud_quotas_template_applications"
description: |-
  Provides a Alicloud Quotas Template Applications resource.
---

# alicloud_quotas_template_applications

Provides a Quotas Template Applications resource. Template Batch Application.

For information about Quotas Template Applications and how to use it, see [What is Template Applications](https://www.alibabacloud.com/help/en/quota-center/developer-reference/api-quotas-2020-05-10-createquotaapplicationsfortemplate).

-> **NOTE:** Available since v1.214.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_quotas_template_applications&exampleId=41f55dda-90b5-2862-333e-a87bc933f8dc24ec4757&activeTab=example&spm=docs.r.quotas_template_applications.0.41f55dda90&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_resource_manager_account" "account" {
  display_name = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_quotas_template_applications" "default" {
  env_language      = "zh"
  notice_type       = "0"
  quota_category    = "WhiteListLabel"
  desire_value      = "1"
  reason            = "example"
  quota_action_code = "quotas.label_multi/A"
  aliyun_uids = [
    "${alicloud_resource_manager_account.account.id}"
  ]
  product_code = "quotas"
}
```

### Deleting `alicloud_quotas_template_applications` or removing it from your configuration

Terraform cannot destroy resource `alicloud_quotas_template_applications`. Terraform will remove this resource from the state file, however resources may remain.

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_quotas_template_applications&spm=docs.r.quotas_template_applications.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `aliyun_uids` - (Required, ForceNew) The list of Alibaba Cloud accounts (primary accounts) of the resource directory members to which the quota is applied.
-> **NOTE:**  Only 50 members can apply for quota increase in batch at a time. For more information about the members of the resource directory, see [Query the list of all members in the resource directory](~~ 604207 ~~).
* `desire_value` - (Required, ForceNew) The value of the quota request.
-> **NOTE:**  The quota request is approved by the technical support of each cloud service. If you want to increase the chance of passing, please fill in a reasonable application value and detailed application reasons when applying for quota.
* `dimensions` - (Optional, ForceNew) Quota dimension. See [`dimensions`](#dimensions) below.
* `effective_time` - (Optional, ForceNew) The UTC time when the quota takes effect. This parameter applies only to the equity quota (WhiteListLabel).
-> **NOTE:**  If the current account does not select the effective time, the default is the submission time.
* `env_language` - (Optional) The language of the quota application result notification. Value:
  - zh (default): Chinese.
  - en: English.
* `expire_time` - (Optional, ForceNew) The UTC time when the quota expires. This parameter applies only to the equity quota (WhiteListLabel).
-> **NOTE:**  If No Expiration Time is selected for the current account, the expiration time is 99 years from the effective time of the current quota.
* `notice_type` - (Optional) Whether to send notification of quota application result. Value:
  - 0 (default): No.
  - 3: Yes.
* `product_code` - (Required, ForceNew) Cloud service name abbreviation.
-> **NOTE:**  For more information about cloud services that support quota centers, see [Cloud services that support quota centers](~~ 182368 ~~).
* `quota_action_code` - (Required, ForceNew) The quota ID.
* `quota_category` - (Required, ForceNew) The quota type. Value:
  - CommonQuota (default): Generic quota.
  - FlowControl:API rate quota.
  - WhiteListLabel: Equity quota.
* `reason` - (Required, ForceNew) Reason for quota application.
-> **NOTE:**  The quota request is approved by the technical support of each cloud service. If you want to increase the chance of passing, please fill in a reasonable application value and detailed application reasons when applying for quota.

### `dimensions`

The dimensions supports the following:
* `key` - (Optional, ForceNew) Quota dimension Key.
* `value` - (Optional, ForceNew) Quota dimension Value.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `quota_application_details` - List of quota application details.
  * `aliyun_uid` - Alibaba Cloud account (primary account).
  * `application_id` - The ID of the quota promotion request.
  * `approve_value` - The approved quota value of the quota increase request.
  * `audit_reason` - Approval comments on quota increase applications.
  * `env_language` - The language of the quota alert notification. Value: zh: Chinese, en: English.
  * `notice_type` - Whether to notify the quota increase application result. Value: 0: No, 3: Yes.
  * `period` - Quota calculation period.
    * `period_unit` - Quota calculation cycle unit.
    * `period_value` - The quota calculation period value.
  * `dimensions` - Quota dimension.
  * `quota_arn` - Quota ARN.
  * `quota_description` - The quota description.
  * `quota_name` - The quota name.
  * `quota_unit` - Quota unit.
  * `reason` - The reason for the quota increase application.
  * `status` - The approval status of the quota promotion application. Value:
  - Disagree: reject.
  - Approve: approved.
  - Process: under review.
  - Cancel: Closed.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Template Applications.

## Import

Quotas Template Applications can be imported using the id, e.g.

```shell
$ terraform import alicloud_quotas_template_applications.example <id>
```