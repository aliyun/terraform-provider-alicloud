---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_load_balancer_security_group_attachment"
description: |-
  Provides a Alicloud NLB Load Balancer Security Group Attachment resource.
---

# alicloud_nlb_load_balancer_security_group_attachment

Provides a NLB Load Balancer Security Group Attachment resource. Security Group mount.

For information about NLB Load Balancer Security Group Attachment and how to use it, see [What is Load Balancer Security Group Attachment](https://www.alibabacloud.com/help/en/server-load-balancer/latest/loadbalancerjoinsecuritygroup).

-> **NOTE:** Available since v1.198.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "vpc" {
  cidr_block = "192.168.0.0/16"
  vpc_name   = var.name

}

resource "alicloud_vswitch" "vswtich" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name

  cidr_block = "192.168.10.0/24"
}

resource "alicloud_vswitch" "vswtich2" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = data.alicloud_zones.default.zones.1.id
  vswitch_name = var.name

  cidr_block = "192.168.30.0/24"
}

resource "alicloud_nlb_load_balancer" "nlb" {
  zone_mappings {
    vswitch_id = alicloud_vswitch.vswtich2.id
    zone_id    = alicloud_vswitch.vswtich2.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vswtich.id
    zone_id    = alicloud_vswitch.vswtich.zone_id
  }
  ipv6_address_type  = "Intranet"
  load_balancer_type = "Network"
  vpc_id             = alicloud_vpc.vpc.id
  address_type       = "Internet"
  address_ip_version = "Ipv4"
}

resource "alicloud_security_group" "securityGroup" {
  security_group_name = var.name

  vpc_id = alicloud_vpc.vpc.id
}


resource "alicloud_nlb_load_balancer_security_group_attachment" "default" {
  load_balancer_id  = alicloud_nlb_load_balancer.nlb.id
  security_group_id = alicloud_security_group.securityGroup.id
}
```

## Argument Reference

The following arguments are supported:
* `dry_run` - (Optional) Whether to PreCheck this request only. Value:
  - **true**: sends a check request and does not bind a security group to the instance. Check items include whether required parameters, request format, and business restrictions have been filled in. If the check fails, the corresponding error is returned. If the check passes, the error code 'DryRunOperation' is returned '.
  - **false** (default): Sends a normal request, returns the HTTP 2xx status code after the check, and directly performs the operation.
* `load_balancer_id` - (Required, ForceNew) The ID of the network-based server load balancer instance to be bound to the security group.
* `security_group_id` - (Optional, ForceNew, Computed) The ID of the security group.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<load_balancer_id>:<security_group_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Load Balancer Security Group Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Load Balancer Security Group Attachment.

## Import

NLB Load Balancer Security Group Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_nlb_load_balancer_security_group_attachment.example <load_balancer_id>:<security_group_id>
```