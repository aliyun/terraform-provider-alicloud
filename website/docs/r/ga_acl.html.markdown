---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_acl"
sidebar_current: "docs-alicloud-resource-ga-acl"
description: |-
  Provides a Alicloud Global Accelerator (GA) Acl resource.
---

# alicloud_ga_acl

Provides a Global Accelerator (GA) Acl resource.

For information about Global Accelerator (GA) Acl and how to use it, see [What is Acl](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-createacl).

-> **NOTE:** Available since v1.150.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ga_acl" "default" {
  acl_name           = "terraform-example"
  address_ip_version = "IPv4"
}

resource "alicloud_ga_acl_entry_attachment" "default" {
  acl_id            = alicloud_ga_acl.default.id
  entry             = "192.168.1.1/32"
  entry_description = "terraform-example"
}
```

## Argument Reference

The following arguments are supported:

* `address_ip_version` - (Required, ForceNew) The IP version. Valid values: `IPv4` and `IPv6`.
* `acl_name` - (Optional) The name of the ACL. The name must be `2` to `128` characters in length, and can contain letters, digits, periods (.), hyphens (-) and underscores (_). It must start with a letter.
* `acl_entries` - (Optional, Computed, Set, Deprecated since v1.190.0) The entries of the Acl. See [`acl_entries`](#acl_entries) below. **NOTE:** "Field `acl_entries` has been deprecated from provider version 1.190.0 and it will be removed in the future version. Please use the new resource `alicloud_ga_acl_entry_attachment`."
* `tags` - (Optional, Available since v1.207.1) A mapping of tags to assign to the resource.
* `dry_run` - (Optional, Bool) The dry run.

### `acl_entries`

The acl_entries supports the following: 

* `entry` - (Optional, Computed) The IP address(192.168.XX.XX) or CIDR(10.0.XX.XX/24) block that you want to add to the network ACL.
* `entry_description` - (Optional, Computed) The description of the IP entry. The description must be `1` to `256` characters in length, and can contain letters, digits, hyphens (-), forward slashes (/), periods (.),and underscores (_).

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Acl. Its value is same as `acl_id`.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Acl.
* `update` - (Defaults to 5 mins) Used when update the Acl.
* `delete` - (Defaults to 5 mins) Used when delete the Acl.

## Import

Global Accelerator (GA) Acl can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_acl.example <id>
```
