---
subcategory: "Security Center"
layout: "alicloud"
page_title: "Alicloud: alicloud_security_center_group"
sidebar_current: "docs-alicloud-resource-security-center-group"
description: |-
  Provides a Alicloud Security Center Group resource.
---

# alicloud_security_center_group

Provides a Security Center Group resource.

For information about Security Center Group and how to use it, see [What is Group](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-sas-2018-12-03-createorupdateassetgroup).

-> **NOTE:** Available since v1.133.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_security_center_group&exampleId=16e39a32-3945-dbc5-1207-e45939b63e2d4609fb6b&activeTab=example&spm=docs.r.security_center_group.0.16e39a3239&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}
resource "alicloud_security_center_group" "example" {
  group_name = var.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_security_center_group&spm=docs.r.security_center_group.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `group_id` - (Optional, ForceNew) GroupId.
* `group_name` - (Optional) GroupName.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Group. Its value is same as `group_id`.

## Timeouts

-> **NOTE:** Available since v1.163.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Security Center Group.
* `update` - (Defaults to 1 mins) Used when update the Security Center Group.
* `delete` - (Defaults to 1 mins) Used when delete the Security Center Group.

## Import

Security Center Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_security_center_group.example <group_id>
```
