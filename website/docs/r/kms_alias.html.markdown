---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_alias"
sidebar_current: "docs-alicloud-resource-kms-alias"
description: |-
  Provides a Alicloud KMS Alias resource.
---

# alicloud_kms_alias

Create an alias for the master key (CMK).

-> **NOTE:** Available since v1.77.0+.

## Example Usage
<div class="oics-button" style="float: right;margin: 0 0 -40px 0;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_kms_alias&exampleId=bb7c42c5-6bc8-0711-f7fe-492295da3043f2d20abe&activeTab=example&spm=docs.r.kms_alias.0.bb7c42c56b" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>

Basic Usage

```terraform
resource "alicloud_kms_key" "this" {
  pending_window_in_days = 7
}

resource "alicloud_kms_alias" "this" {
  alias_name = "alias/example_kms_alias"
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

```shell
$ terraform import alicloud_kms_alias.example alias/test_kms_alias
```
