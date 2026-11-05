---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_vpd"
description: |-
  Provides a Alicloud ENS Vpd resource.
---

# alicloud_ens_vpd

Provides a ENS Vpd resource.

For information about ENS Vpd and how to use it, see [What is Vpd](https://www.alibabacloud.com/help/en/ens/).

-> **NOTE:** Available since v2.3.0.

-> **NOTE:** The ENS VPD related APIs are private. Please confirm that the ENS VPD feature has been opened for your account before using this resource.

## Example Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_ens_network" "default" {
  network_name  = var.name
  description   = var.name
  cidr_block    = "192.168.2.0/24"
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_vpd" "default" {
  ens_region_id = alicloud_ens_network.default.ens_region_id
  network_id    = alicloud_ens_network.default.id
  name          = var.name
  description   = var.name
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) The description of the VPD. It must be 1 to 256 characters in length and cannot start with `http://` or `https://`.
* `ens_region_id` - (Required, ForceNew) The ENS node ID, e.g. `cn-chenzhou-telecom_unicom_cmcc`.
* `name` - (Optional) The name of the VPD. It must be 1 to 128 characters in length and cannot start with `http://` or `https://`.
* `network_id` - (Required, ForceNew) The ID of the ENS Network to which the VPD belongs.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `cidr_block` - The CIDR block of the VPD.
* `creation_time` - The creation time of the VPD.
* `instance_ids` - The list of instance IDs attached to the VPD.
* `status` - The status of the VPD. Valid values: `Creating`, `Available`, `Deleting`.
* `vpd_vswitch_ids` - The list of VSwitch IDs under the VPD.

## Import

ENS Vpd can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_vpd.example <id>
```
