---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_key_pair"
sidebar_current: "docs-alicloud-resource-ens-key-pair"
description: |-
  Provides a Alicloud ENS Key Pair resource.
---

# alicloud\_ens\_key\_pair

Provides a ENS Key Pair resource.

For information about ENS Key Pair and how to use it, see [What is Key Pair](https://help.aliyun.com/product/62684.html).

-> **NOTE:** Available in v1.133.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ens_key_pair" "example" {
  key_pair_name = "example_value"
  version       = "example_value"
}
```

## Argument Reference

The following arguments are supported:

* `key_pair_name` - (Required, ForceNew) The name of the key pair.
* `version` - (Required, ForceNew) The version number.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Key Pair. The value formats as `<key_pair_name>:<version>`.

## Import

ENS Key Pair can be imported using the id, e.g.

```
$ terraform import alicloud_ens_key_pair.example <key_pair_name>:<version>
```
