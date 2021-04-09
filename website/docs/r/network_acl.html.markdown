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

```terraform
resource "alicloud_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "VpcConfig"
}

resource "alicloud_network_acl" "default" {
  vpc_id           = alicloud_vpc.default.id
  network_acl_name = "network_acl"
  description      = "network_acl"
  ingress_acl_entries {
    description            = "tf-testacc"
    network_acl_entry_name = "tcp23"
    source_cidr_ip         = "196.168.2.0/21"
    policy                 = "accept"
    port                   = "22/80"
    protocol               = "tcp"
  }
  egress_acl_entries {
    description            = "tf-testacc"
    network_acl_entry_name = "tcp23"
    destination_cidr_ip    = "0.0.0.0/0"
    policy                 = "accept"
    port                   = "-1/-1"
    protocol               = "all"
  }
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) The vpc_id of the network acl, the field can't be changed.
* `name` - (Optional, Deprecated from v1.122.0) Field `name` has been deprecated from provider version 1.122.0. New field `network_acl_name` instead.
* `network_acl_name` - (Optional, Available in 1.122.0+) The name of the network acl.
* `description` - (Optional) The description of the network acl instance.
* `ingress_acl_entries` - (Optional, Computed, Available in 1.122.0+) List of the ingress entries of the network acl. The order of the ingress entries determines the priority. The details see Block `ingress_acl_entries`.
* `egress_acl_entries` - (Optional, Computed, Available in 1.122.0+) List of the egress entries of the network acl. The order of the egress entries determines the priority. The details see Block `egress_acl_entries`.

### Block ingress_acl_entries

* `description` - (Optional) The description of ingress entries.
* `network_acl_entry_name` - (Optional) The entry name of ingress entries. 
* `policy` - (Optional) The policy of ingress entries. Valid values `accept` and `drop`.
* `port` - (Optional) The port of ingress entries.
* `protocol` - (Optional) The protocol of ingress entries. Valid values `icmp`,`gre`,`tcp`,`udp`, and `all`.
* `source_cidr_ip` - (Optional) The source cidr ip of ingress entries.

### Block egress_acl_entries

* `description` - (Optional) The description of egress entries.
* `network_acl_entry_name` - (Optional) The entry name of egress entries. 
* `policy` - (Optional) The policy of egress entries. Valid values `accept` and `drop`.
* `port` - (Optional) The port of egress entries.
* `protocol` - (Optional) The protocol of egress entries. Valid values `icmp`,`gre`,`tcp`,`udp`, and `all`.
* `destination_cidr_ip` - (Optional) The destination cidr ip of egress entries.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the network acl instance id.
* `status` - (Available in 1.122.0+) The status of the network acl.

### Timeouts

-> **NOTE:** Available in 1.122.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the Network ACL. (until it reaches the initial `Available` status). 
* `update` - (Defaults to 10 mins) Used when updating the Network ACL. (until it reaches the initial `Available` status). 
* `delete` - (Defaults to 10 mins) Used when terminating the Network ACL.

## Import

The network acl can be imported using the id, e.g.

```
$ terraform import alicloud_network_acl.default nacl-abc123456
```


