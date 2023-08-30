---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_load_balancer_security_group_attachment"
sidebar_current: "docs-alicloud-resource-nlb-load-balancer-security-group-attachment"
description: |-
  Provides a Alicloud Nlb Load Balancer Security Group Attachment resource.
---

# alicloud_nlb_load_balancer_security_group_attachment

Provides a Nlb Load Balancer Security Group Attachment resource.

For information about Nlb Load Balancer Security Group Attachment and how to use it, see [What is Load Balancer Security Group Attachment](https://www.alibabacloud.com/help/en/server-load-balancer/latest/loadbalancerjoinsecuritygroup).

-> **NOTE:** Available since v1.198.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_resource_manager_resource_groups" "default" {}
data "alicloud_nlb_zones" "default" {}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_nlb_zones.default.zones.0.id
}
resource "alicloud_vswitch" "default1" {
  vswitch_name = var.name
  cidr_block   = "10.4.1.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_nlb_zones.default.zones.1.id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_nlb_load_balancer" "default" {
  load_balancer_name = var.name
  resource_group_id  = data.alicloud_resource_manager_resource_groups.default.ids.0
  load_balancer_type = "Network"
  address_type       = "Internet"
  address_ip_version = "Ipv4"
  vpc_id             = alicloud_vpc.default.id
  tags = {
    Created = "TF",
    For     = "example",
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.default.id
    zone_id    = data.alicloud_nlb_zones.default.zones.0.id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.default1.id
    zone_id    = data.alicloud_nlb_zones.default.zones.1.id
  }
}

resource "alicloud_nlb_load_balancer_security_group_attachment" "default" {
  security_group_id = alicloud_security_group.default.id
  load_balancer_id  = alicloud_nlb_load_balancer.default.id
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, ForceNew) The ID of the network-based server load balancer instance to be bound to the security group.
* `security_group_id` - (Required, ForceNew) The ID of security groups.
* `dry_run` - (Optional) Whether to PreCheck this request only. Value:-**true**: sends a check request and does not bind a security group to the instance. Check items include whether required parameters, request format, and business restrictions have been filled in. If the check fails, the corresponding error is returned. If the check passes, the error code 'DryRunOperation' is returned '.-**false** (default): Sends a normal request, returns the HTTP 2xx status code after the check, and directly performs the operation.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Load Balancer Security Group Attachment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Load Balancer Security Group Attachment.
* `update` - (Defaults to 1 mins) Used when delete the Load Balancer Security Group Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Load Balancer Security Group Attachment.

## Import

Nlb Load Balancer Security Group Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_nlb_load_balancer_security_group_attachment.example <LoadBalancerId>:<SecurityGroupId>
```