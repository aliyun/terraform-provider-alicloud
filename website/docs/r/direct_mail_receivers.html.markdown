---
subcategory: "Direct Mail"
layout: "alicloud"
page_title: "Alicloud: alicloud_direct_mail_receivers"
sidebar_current: "docs-alicloud-resource-direct-mail-receivers"
description: |-
  Provides a Alicloud Direct Mail Receivers resource.
---

# alicloud\_direct\_mail\_receivers

Provides a Direct Mail Receivers resource.

For information about Direct Mail Receivers and how to use it, see [What is Direct Mail Receivers](https://www.alibabacloud.com/help/en/doc-detail/29414.htm).

-> **NOTE:** Available in v1.125.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_direct_mail_receivers" "example" {
  receivers_alias = "tf-vme8@onaliyun.com"
  receivers_name  = "vme8"
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

```
$ terraform import alicloud_direct_mail_receivers.example <id>
```
