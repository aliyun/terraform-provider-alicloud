---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_acl"
sidebar_current: "docs-alicloud-resource-ga-acl"
description: |-
  Provides a Alicloud Global Accelerator (GA) Acl resource.
---

# alicloud\_ga\_acl

Provides a Global Accelerator (GA) Acl resource.

For information about Global Accelerator (GA) Acl and how to use it, see [What is Acl](https://www.alibabacloud.com/help/en/doc-detail/258289.html).

-> **NOTE:** Available in v1.150.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ga_acl" "default" {
  acl_name           = "tf-testAccAcl"
  address_ip_version = "IPv4"
  acl_entries {
    entry             = "192.168.1.0/24"
    entry_description = "tf-test1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `acl_entries` - (Optional) The entries of the Acl. See the following `Block acl_entries`.
* `acl_name` - (Optional) The name of the ACL. The name must be `2` to `128` characters in length, and can contain letters, digits, periods (.), hyphens (-) and underscores (_). It must start with a letter.
* `address_ip_version` - (Required, ForceNew) The IP version. Valid values: `IPv4` and `IPv6`.
* `dry_run` - (Optional) The dry run.

#### Block acl_entries

The acl_entries supports the following: 

* `entry` - (Optional) The IP entry that you want to add to the ACL.
* `entry_description` - (Optional) The description of the IP entry. The description must be `1` to `256` characters in length, and can contain letters, digits, hyphens (-), forward slashes (/), periods (.),and underscores (_).

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Acl. Its value is same as `acl_id`.
* `status` - The status of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Acl.
* `delete` - (Defaults to 5 mins) Used when delete the Acl.
* `update` - (Defaults to 5 mins) Used when update the Acl.

## Import

Global Accelerator (GA) Acl can be imported using the id, e.g.

```
$ terraform import alicloud_ga_acl.example <id>
```