---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_key_pair"
sidebar_current: "docs-alicloud-resource-ens-key-pair"
description: |-
  Provides a Alicloud ENS Key Pair resource.
---

# alicloud_ens_key_pair

Provides a ENS Key Pair resource.

For information about ENS Key Pair and how to use it, see [What is Key Pair](https://www.alibabacloud.com/help/en/ens/latest/createkeypair).

-> **NOTE:** Available since v1.133.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}
resource "alicloud_ens_key_pair" "example" {
  key_pair_name = var.name
  version       = "2017-11-10"
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

```shell
$ terraform import alicloud_ens_key_pair.example <key_pair_name>:<version>
```
