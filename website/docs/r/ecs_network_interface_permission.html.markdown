---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_network_interface_permission"
sidebar_current: "docs-alicloud-resource-ecs-network-interface-permission"
description: |-
  Provides a Alicloud ECS Network Interface Permission resource.
---

# alicloud\_ecs\_network\_interface\_permission

Provides a ECS Network Interface Permission resource.

For information about ECS Network Interface Permission and how to use it, see [What is Network Interface Permission](https://www.alibabacloud.com/help/en/elastic-compute-service/latest/elastic-network-interfaces-overview).

-> **NOTE:** Available in v1.166.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_ecs_network_interface" "default" {
  network_interface_name = var.name
  vswitch_id             = data.alicloud_vswitches.default.ids.0
  security_group_ids     = [alicloud_security_group.default.id]
  description            = "Basic test"
  primary_ip_address     = cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 100)
  tags = {
    Created = "TF",
    For     = "Test",
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}
data "alicloud_account" "default" {}
resource "alicloud_ecs_network_interface_permission" "example" {
  account_id           = data.alicloud_account.default.id
  network_interface_id = alicloud_ecs_network_interface.default.id
  permission           = "InstanceAttach"
  force                = true
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Required, ForceNew) Alibaba Cloud Partner (Certified ISV) account ID or individual user ID.
* `network_interface_id` - (Required, ForceNew) The ID of the network interface.
* `permission` - (Required, ForceNew) The permissions of the Network Interface. Valid values: `InstanceAttach`. `InstanceAttach`: Allows authorized users to mount your ENI to the other ECS instance. The ECS instance must be in the same zone as the ENI.
* `force` - (Optional, ForceNew) Whether to force deletion of Network Interface Permission. Default value: `true`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Network Interface Permission.
* `status` - The Status of the Network Interface Permissions. Valid values: `Pending`, `Granted`, `Revoking`, `Revoked`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Network Interface Permission (until it reaches the initial `Granted` status).
* `delete` - (Defaults to 1 mins) Used when deleting the Network Interface Permission.

## Import

ECS Network Interface Permission can be imported using the id, e.g.

```
$ terraform import alicloud_ecs_network_interface_permission.example <id>
```