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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecs_network_interface_permission&exampleId=2a7df141-e78e-51e0-5330-16a0934f7158ca1040f0&activeTab=example&spm=docs.r.ecs_network_interface_permission.0.2a7df141e7&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_account" "default" {}
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "default" {
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = "terraform-example"
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name   = "terraform-example"
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_ecs_network_interface" "default" {
  network_interface_name = "terraform-example"
  vswitch_id             = alicloud_vswitch.default.id
  security_group_ids     = [alicloud_security_group.default.id]
  description            = "terraform-example"
  primary_ip_address     = cidrhost(alicloud_vswitch.default.cidr_block, 100)
  tags = {
    Created = "TF",
    For     = "example",
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}

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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Network Interface Permission (until it reaches the initial `Granted` status).
* `delete` - (Defaults to 1 mins) Used when deleting the Network Interface Permission.

## Import

ECS Network Interface Permission can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_network_interface_permission.example <id>
```