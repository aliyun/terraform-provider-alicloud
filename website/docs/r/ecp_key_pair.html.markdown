---
subcategory: "Elastic Cloud Phone (ECP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecp_key_pair"
sidebar_current: "docs-alicloud-resource-ecp-key-pair"
description: |-
  Provides a Alicloud Elastic Cloud Phone (ECP) Key Pair resource.
---

# alicloud\_ecp\_key\_pair

Provides a Elastic Cloud Phone (ECP) Key Pair resource.

For information about Elastic Cloud Phone (ECP) Key Pair and how to use it, see [What is Key Pair](https://help.aliyun.com/document_detail/257197.html).

-> **NOTE:** Available in v1.130.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecp_key_pair" "example" {
  key_pair_name   = "my-KeyPair"
  public_key_body = "ssh-rsa AAAAxxxxxxxxxxtyuudsfsg"
}

```

## Argument Reference

The following arguments are supported:

* `key_pair_name` - (Required, ForceNew) The Key Name.
* `public_key_body` - (Required) The public key body.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Key Pair. Its value is same as `key_pair_name`.

## Import

Elastic Cloud Phone (ECP) Key Pair can be imported using the id, e.g.

```
$ terraform import alicloud_ecp_key_pair.example <key_pair_name>
```
