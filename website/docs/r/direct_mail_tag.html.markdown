---
subcategory: "Direct Mail"
layout: "alicloud"
page_title: "Alicloud: alicloud_direct_mail_tag"
sidebar_current: "docs-alicloud-resource-direct-mail-tag"
description: |-
  Provides a Alicloud Direct Mail Tag resource.
---

# alicloud_direct_mail_tag

Provides a Direct Mail Tag resource.

For information about Direct Mail Tag and how to use it, see [What is Tag](https://www.alibabacloud.com/help/en/directmail/latest/createtag).

-> **NOTE:** Available since v1.144.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_direct_mail_tag&exampleId=f391bc81-2652-ed15-f0a9-dcf710b14c129e5cd237&activeTab=example&spm=docs.r.direct_mail_tag.0.f391bc8126&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "example"
}
provider "alicloud" {
  region = "cn-hangzhou"
}
resource "alicloud_direct_mail_tag" "example" {
  tag_name = var.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_direct_mail_tag&spm=docs.r.direct_mail_tag.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `tag_name` - (Required) The name of the tag. The name must be `1` to `50` characters in length, and can contain letters and digits.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Tag.

## Import

Direct Mail Tag can be imported using the id, e.g.

```shell
$ terraform import alicloud_direct_mail_tag.example <id>
```