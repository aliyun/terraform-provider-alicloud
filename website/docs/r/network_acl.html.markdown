---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_network_acl"
sidebar_current: "docs-alicloud-resource-network-acl"
description: |-
  Provides a Alicloud Network Acl resource.
---

# alicloud\_network_acl

Provides a network acl resource to add network acls.

-> **NOTE:** Available in 1.43.0+. Currently, the resource are only available in Hongkong(cn-hongkong), India(ap-south-1), and Indonesia(ap-southeast-1) regions.

## Example Usage

Basic Usage

```
resource "alicloud_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  name       = "VpcConfig"
}

resource "alicloud_network_acl" "default" {
  vpc_id      = alicloud_vpc.default.id
  name        = "network_acl"
  description = "network_acl"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) The vpc_id of the network acl, the field can't be changed.
* `name` - (Optional) The name of the network acl.
* `description` - (Optional) The description of the network acl instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the network acl instance id.

## Import

The network acl can be imported using the id, e.g.

```
$ terraform import alicloud_network_acl.default nacl-abc123456
```


