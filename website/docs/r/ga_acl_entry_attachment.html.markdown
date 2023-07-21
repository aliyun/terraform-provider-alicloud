---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_acl_entry_attachment"
sidebar_current: "docs-alicloud-resource-ga-acl-entry-attachment"
description: |-
  Provides a Alicloud Global Accelerator (GA) Acl entry attachment resource.
---

# alicloud_ga_acl_entry_attachment

Provides a Global Accelerator (GA) Acl entry attachment resource.

For information about Global Accelerator (GA) Acl entry attachment and how to use it, see [What is Acl entry attachment](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-addentriestoacl).

-> **NOTE:** Available since v1.190.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ga_acl" "default" {
  acl_name           = "tf-example-value"
  address_ip_version = "IPv4"
}

resource "alicloud_ga_acl_entry_attachment" "default" {
  acl_id            = alicloud_ga_acl.default.id
  entry             = "192.168.1.1/32"
  entry_description = "tf-example-value"
}
```

## Argument Reference

The following arguments are supported:

* `acl_id` - (Required, ForceNew) The ID of the global acceleration instance.
* `entry` - (Required, ForceNew) The IP address(192.168.XX.XX) or CIDR(10.0.XX.XX/24) block that you want to add to the network ACL.
* `entry_description` - (Optional, ForceNew) The description of the entry. The description must be 1 to 256 characters in length, and can contain letters, digits, hyphens (-), forward slashes (/), periods (.), and underscores (_).

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Acl entry attachment. The value formats as `<acl_id>:<entry>`.
* `status` - The status of the network ACL.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Acl entry attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Acl entry attachment.

## Import

Global Accelerator (GA) Acl entry attachment can be imported using the id.Format to `<acl_id>:<entry>`, e.g.

```shell
$ terraform import alicloud_ga_acl_entry_attachment.example your_acl_id:your_entry
```
