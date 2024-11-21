---
subcategory: "Smart Access Gateway (Smartag)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sag_acl"
sidebar_current: "docs-alicloud-resource-sag-acl"
description: |-
  Provides a Sag Acl resource.
---

# alicloud_sag_acl

Provides a Sag Acl resource. Smart Access Gateway (SAG) provides the access control list (ACL) function in the form of whitelists and blacklists for different SAG instances.

For information about Sag Acl and how to use it, see [What is access control list (ACL)](https://www.alibabacloud.com/help/en/smart-access-gateway/latest/createacl).

-> **NOTE:** Available since v1.60.0.

-> **NOTE:** Only the following regions support create Cloud Connect Network. [`cn-shanghai`, `cn-shanghai-finance-1`, `cn-hongkong`, `ap-southeast-1`, `ap-southeast-3`, `ap-southeast-5`, `ap-northeast-1`, `eu-central-1`]

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_sag_acl&exampleId=a0603b71-6545-ff68-9819-0ed03ba7fcafd25895a3&activeTab=example&spm=docs.r.sag_acl.0.a0603b7165&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-shanghai"
}

resource "alicloud_sag_acl" "default" {
  name = "terraform-example"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the ACL instance. The name can contain 2 to 128 characters including a-z, A-Z, 0-9, periods, underlines, and hyphens. The name must start with an English letter, but cannot start with http:// or https://.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the ACL. For example "acl-xxx".

## Import

The Sag Acl can be imported using the id, e.g.

```shell
$ terraform import alicloud_sag_acl.example acl-abc123456
```

