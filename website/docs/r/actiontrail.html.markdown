---
subcategory: "Actiontrail"
layout: "alicloud"
page_title: "Alicloud: alicloud_actiontrail"
sidebar_current: "docs-alicloud-resource-actiontrail"
description: |-
  Provides Alibaba Cloud ActionTrail Resource
---

# alicloud\_actiontrail

-> **DEPRECATED:**  This resource has been renamed to [alicloud_actiontrail_trail](https://www.terraform.io/docs/providers/alicloud/r/actiontrail_trail) from version 1.95.0.

Provides a new resource to manage [Action Trail](https://www.alibabacloud.com/help/doc-detail/28804.htm).

-> **NOTE:** Available in 1.35.0+

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_actiontrail&exampleId=2005a93b-e268-b761-88bd-aa5750531b4be95fd3f1&activeTab=example&spm=docs.r.actiontrail.0.2005a93be2&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
# Create a new action trail.
resource "alicloud_actiontrail" "foo" {
  name            = "action-trail"
  event_rw        = "Write-test"
  oss_bucket_name = alicloud_oss_bucket.bucket.id
  role_name       = alicloud_ram_role_policy_attachment.attach.role_name
  oss_key_prefix  = "at-product-account-audit-B"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the trail to be created, which must be unique for an account.
* `event_rw` - (Optional) Indicates whether the event is a read or a write event. Valid values: Read, Write, and All. Default value: Write.
* `oss_bucket_name` - (Required) The OSS bucket to which the trail delivers logs. Ensure that this is an existing OSS bucket.
* `role_name` - (Required) The RAM role in ActionTrail permitted by the user.
* `oss_key_prefix` - (Optional) The prefix of the specified OSS bucket name. This parameter can be left empty.
* `sls_project_arn` - (Optional) The unique ARN of the Log Service project.
* `sls_write_role_arn` - (Optional) The unique ARN of the Log Service role.

-> **NOTE:** `sls_project_arn` and `sls_write_role_arn` should be set or not set at the same time when actiontrail delivers logs.

## Attributes Reference

The following attributes are exported:

* `id` - The action trail id. The value is same as its name.

## Import

Action trail can be imported using the id, e.g.

```shell
$ terraform import alicloud_actiontrail.foo abc12345678
```
