---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_alias"
sidebar_current: "docs-alicloud-resource-kms-alias"
description: |-
  Provides a Alicloud KMS Alias resource.
---

# alicloud\_kms\_alias

Create an alias for the master key (CMK).

-> **NOTE:** Available in v1.77.0+.

## Example Usage

Basic Usage

```
resource "alicloud_kms_key" "this" {}

resource "alicloud_kms_alias" "this" {
  alias_name = "alias/test_kms_alias"
  key_id     = alicloud_kms_key.this.id
}
```

## Argument Reference

The following arguments are supported:

* `alias_name` - (Required, ForceNew) The alias of CMK. `Encrypt`、`GenerateDataKey`、`DescribeKey` can be called using aliases. Length of characters other than prefixes: minimum length of 1 character and maximum length of 255 characters. Must contain prefix `alias/`.
* `key_id` - (Required) The id of the key.

-> **NOTE:** Each alias represents only one master key(CMK).

-> **NOTE:** Within an area of the same user, alias is not reproducible.

-> **NOTE:** UpdateAlias can be used to update the mapping relationship between alias and master key(CMK).


## Attributes Reference

* `id` - The ID of the alias.

## Import

KMS alias can be imported using the id, e.g.

```
$ terraform import alicloud_kms_alias.example alias/test_kms_alias
```
