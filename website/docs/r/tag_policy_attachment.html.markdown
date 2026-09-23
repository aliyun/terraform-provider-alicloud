---
subcategory: "TAG"
layout: "alicloud"
page_title: "Alicloud: alicloud_tag_policy_attachment"
sidebar_current: "docs-alicloud-resource-tag-policy-attachment"
description: |-
  Provides a Alicloud Tag Policy Attachment resource.
---

# alicloud_tag_policy_attachment

Provides a Tag Policy Attachment resource to attaches a policy to an object. After you attach a tag policy to an object.
For information about Tag Policy Attachment and how to use it,
see [What is Policy Attachment](https://www.alibabacloud.com/help/en/resource-management/latest/attach-policy).

-> **NOTE:** Available since v1.204.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_tag_policy_attachment&exampleId=11c764cf-ddd6-30ae-3c67-f4dacc8f4e0ddd8db27f&activeTab=example&spm=docs.r.tag_policy_attachment.0.11c764cfdd&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
provider "alicloud" {
  region = "cn-shanghai"
}
data "alicloud_account" "default" {}
resource "alicloud_tag_policy" "example" {
  policy_name    = var.name
  policy_desc    = var.name
  user_type      = "USER"
  policy_content = <<EOF
		{"tags":{"CostCenter":{"tag_value":{"@@assign":["Beijing","Shanghai"]},"tag_key":{"@@assign":"CostCenter"}}}}
    EOF
}

resource "alicloud_tag_policy_attachment" "example" {
  policy_id   = alicloud_tag_policy.example.id
  target_id   = data.alicloud_account.default.id
  target_type = "USER"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_tag_policy_attachment&spm=docs.r.tag_policy_attachment.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, ForceNew) The ID of the tag policy.
* `target_id` - (Required, ForceNew) The ID of the object.
* `target_type` - (Required, ForceNew) The type of the object. Valid values: `USER`, `ROOT`, `FOLDER`, `ACCOUNT`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Policy Attachment. It formats as `<policy_id>`:`<target_id>`:`<target_type>`.

## Import

Tag Policy Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_tag_policy_attachment.example <policy_id>:<target_id>:<target_type>
```