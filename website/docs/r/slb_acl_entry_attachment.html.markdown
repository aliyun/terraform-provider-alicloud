---
subcategory: "Classic Load Balancer (SLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_acl_entry_attachment"
sidebar_current: "docs-alicloud-resource-slb-acl-entry-attachment"
description: |-
  Provides a Acl entry attachment resource.
---

# alicloud_slb_acl_entry_attachment

-> **NOTE:** Available since v1.162.0.

-> **NOTE:** The `entries` attribute is available since v2.1.0. In batch mode, the attachment takes ownership of all entries of the acl: entries added out of band or by other `alicloud_slb_acl_entry_attachment` resources attached to the same acl are removed on the next apply. Do not manage the entries of the same acl from multiple resources.

-> **NOTE:** The maximum number of entries per acl is 300.

-> **NOTE:** Exactly one of `entry` and `entries` must be specified. Switching between them replaces the resource. At least one entry block is required; to remove all the entries, remove the resource.

For information about acl entry attachment and how to use it, see [Configure an acl entry](https://www.alibabacloud.com/help/en/doc-detail/70023.html).


## Example Usage

### Batch mode

The `entries` attribute manages all entries of the acl in one resource. The entries are added and removed in batches of at most `50` entries per API call.

```terraform
resource "alicloud_slb_acl" "attachment" {
  name       = "forSlbAclEntryAttachment"
  ip_version = "ipv4"
}

resource "alicloud_slb_acl_entry_attachment" "attachment" {
  acl_id = alicloud_slb_acl.attachment.id

  entries {
    entry   = "168.10.10.0/24"
    comment = "first"
  }

  entries {
    entry   = "168.10.20.0/24"
    comment = "second"
  }
}
```

### Single entry

`entry` is deprecated, please use `entries` instead.

```terraform
resource "alicloud_slb_acl" "attachment" {
  name       = "forSlbAclEntryAttachment"
  ip_version = "ipv4"
}

resource "alicloud_slb_acl_entry_attachment" "attachment" {
  acl_id  = alicloud_slb_acl.attachment.id
  entry   = "168.10.10.0/24"
  comment = "second"
}
```

## Argument Reference

The following arguments are supported:

* `acl_id` - (Required, ForceNew) The ID of the Acl.
* `entry` - (Optional, ForceNew, Deprecated) The CIDR block of the acl entry. Field `entry` has been deprecated from provider version 2.1.0 and it will be removed in the future version. Please use the new field `entries`. Exactly one of `entry` and `entries` must be specified.
* `comment` - (Optional, ForceNew) The comment of the entry. Only valid when `entry` is set.
* `entries` - (Optional, Available since v2.1.0) One or more entry blocks. Exactly one of `entry` and `entries` must be specified. The order of the blocks is not significant. See [`entries`](#entries) below for details.

### `entries`

The entries supports the following:

* `entry` - (Required) The CIDR block of the acl entry.
* `comment` - (Optional) The comment of the acl entry.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. The value formats as `<acl_id>:<entry>` when `entry` is set, or `<acl_id>` when `entries` is set.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the resource.
* `update` - (Defaults to 5 mins, Available since v2.1.0) Used when update the resource.
* `delete` - (Defaults to 5 mins) Used when delete the resource.

## Import

Acl entry attachment can be imported using the id. For a single entry the id consists of acl_id and entry, e.g.

```shell
$ terraform import alicloud_slb_acl_entry_attachment.example <acl_id>:<entry>
```

When `entries` is used the id is the acl id, e.g.

```shell
$ terraform import alicloud_slb_acl_entry_attachment.example <acl_id>
```
