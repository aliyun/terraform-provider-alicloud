---
subcategory: "Direct Mail"
layout: "alicloud"
page_title: "Alicloud: alicloud_direct_mail_receivers"
sidebar_current: "docs-alicloud-resource-direct-mail-receivers"
description: |-
  Provides a Alicloud Direct Mail Receivers resource.
---

# alicloud_direct_mail_receivers

Provides a Direct Mail Receivers resource.

For information about Direct Mail Receivers and how to use it, see [What is Direct Mail Receivers](https://www.alibabacloud.com/help/en/doc-detail/29414.htm).

-> **NOTE:** Available since v1.125.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_direct_mail_receivers&exampleId=7b6aaaa3-d2ef-eada-2e66-054b9f84ffc64eb29b5e&activeTab=example&spm=docs.r.direct_mail_receivers.0.7b6aaaa3d2&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tfexample"
}
provider "alicloud" {
  region = "cn-hangzhou"
}
resource "alicloud_direct_mail_receivers" "example" {
  receivers_alias = format("%s@onaliyun.com", var.name)
  receivers_name  = var.name
  description     = var.name
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional, ForceNew) The description of receivers and 1-50 characters in length.
* `receivers_alias` - (Required, ForceNew) The alias of receivers. Must email address and less than 30 characters in length.
* `receivers_name` - (Required, ForceNew) The name of the resource. The length that cannot be repeated is 1-30 characters.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Receivers.
* `status` - The status of the resource. `0` means uploading, `1` means upload completed. 

## Import

Direct Mail Receivers can be imported using the id, e.g.

```shell
$ terraform import alicloud_direct_mail_receivers.example <id>
```
