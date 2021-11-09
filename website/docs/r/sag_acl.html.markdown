---
subcategory: "Smart Access Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_sag_acl"
sidebar_current: "docs-alicloud-resource-sag-acl"
description: |-
  Provides a Sag Acl resource.
---

# alicloud\_sag\_acl

Provides a Sag Acl resource. Smart Access Gateway (SAG) provides the access control list (ACL) function in the form of whitelists and blacklists for different SAG instances.

For information about Sag Acl and how to use it, see [What is access control list (ACL)](https://www.alibabacloud.com/help/doc-detail/111518.htm).

-> **NOTE:** Available in 1.60.0+

-> **NOTE:** Only the following regions support create Cloud Connect Network. [`cn-shanghai`, `cn-shanghai-finance-1`, `cn-hongkong`, `ap-southeast-1`, `ap-southeast-2`, `ap-southeast-3`, `ap-southeast-5`, `ap-northeast-1`, `eu-central-1`]

## Example Usage

Basic Usage

```
resource "alicloud_sag_acl" "default" {
  name        = "tf-testAccSagAclName"
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

```
$ terraform import alicloud_sag_acl.example acl-abc123456
```

