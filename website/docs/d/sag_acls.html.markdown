---
subcategory: "Smart Access Gateway (Smartag)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sag_acls"
sidebar_current: "docs-alicloud-resource-sag-acls"
description: |-
  Provides the access control list (ACL) function in the form of whitelists and blacklists for different SAG instances.
---

# alicloud_sag_acls

This data source provides Sag Acls available to the user.

-> **NOTE:** Available since v1.60.0+

-> **NOTE:** Only the following regions support create Cloud Connect Network. [`cn-shanghai`, `cn-shanghai-finance-1`, `cn-hongkong`, `ap-southeast-1`, `ap-southeast-3`, `ap-southeast-5`, `ap-northeast-1`, `eu-central-1`]

## Example Usage

Basic Usage

```terraform
data "alicloud_sag_acls" "default" {
  ids        = ["${alicloud_sag_acl.default.id}"]
  name_regex = "^terraform-example"
}
resource "alicloud_sag_acl" "default" {
  name = "terraform-example"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Sag Acl IDs.
* `name_regex` - (Optional) A regex string to filter Sag Acl instances by name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Sag Acl IDs.
* `names` - A list of Sag Acls names. 
* `acls` - A list of Sag Acls. Each element contains the following attributes:
  * `id` - The ID of the ACL. For example "acl-xxx".
  * `name` - The name of the Acl.
