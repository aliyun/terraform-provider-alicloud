---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_network"
description: |-
  Provides a Alicloud ENS Network resource.
---

# alicloud_ens_network

Provides a ENS Network resource. 

For information about ENS Network and how to use it, see [What is Network](https://www.alibabacloud.com/help/en/ens/developer-reference/api-createnetwork-1).

-> **NOTE:** Available since v1.213.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_ens_network" "default" {
  network_name = var.name

  description   = var.name
  cidr_block    = "192.168.2.0/24"
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
}
```

## Argument Reference

The following arguments are supported:
* `cidr_block` - (Required, ForceNew) The network segment of the network. You can use the following network segments or a subset of them as the network segment: `10.0.0.0/8` (default), `172.16.0.0/12`, `192.168.0.0/16`.
* `description` - (Optional) Description information.Rules:It must be 2 to 256 characters in length and must start with a letter or Chinese, but cannot start with `http://` or `https://`. Example value: this is my first network.
* `ens_region_id` - (Required, ForceNew) Ens node IDExample value: cn-beijing-telecom.
* `network_name` - (Optional) Name of the network instanceThe naming rules are as follows: 1. Length is 2~128 English or Chinese characters; 2. It must start with a large or small letter or Chinese, not with `http://` and `https://`; 3. Can contain numbers, colons (:), underscores (_), or dashes (-).

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Creation time, timestamp (MS).
* `status` - The status of the network instance. Pending: Configuring, Available: Available.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Network.
* `delete` - (Defaults to 5 mins) Used when delete the Network.
* `update` - (Defaults to 5 mins) Used when update the Network.

## Import

ENS Network can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_network.example <id>
```