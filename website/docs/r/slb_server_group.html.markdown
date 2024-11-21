---
subcategory: "Classic Load Balancer (SLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_server_group"
sidebar_current: "docs-alicloud-resource-slb-server-group"
description: |-
  Provides a Alicloud Load Balancer Virtual Backend Server Group resource.
---

# alicloud_slb_server_group

Provides a Load Balancer Virtual Backend Server Group resource.

For information about Load Balancer Virtual Backend Server Group and how to use it, see [What is Virtual Backend Server Group](https://www.alibabacloud.com/help/en/doc-detail/35215.html).

-> **NOTE:** Available since v1.6.0.

-> **NOTE:** One ECS instance can be added into multiple virtual server groups.

-> **NOTE:** One virtual server group can be attached with multiple listeners in one load balancer.

-> **NOTE:** One Classic and Internet load balancer, its virtual server group can add Classic and VPC ECS instances.

-> **NOTE:** One Classic and Intranet load balancer, its virtual server group can only add Classic ECS instances.

-> **NOTE:** One VPC load balancer, its virtual server group can only add the same VPC ECS instances.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_slb_server_group&exampleId=4373d1f9-7b66-0d11-0be5-e55c439b3c19f7f21ef3&activeTab=example&spm=docs.r.slb_server_group.0.4373d1f97b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/16"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = var.name
  vswitch_id         = alicloud_vswitch.default.id
  load_balancer_spec = "slb.s2.small"
}

resource "alicloud_slb_server_group" "default" {
  load_balancer_id = alicloud_slb_load_balancer.default.id
  name             = var.name
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, ForceNew) The ID of the Server Load Balancer (SLB) instance.
* `name` - (Optional) The name of the vServer group. Default value: `tf-server-group`.
* `delete_protection_validation` - (Optional, Bool, Available since v1.63.0) Checking DeleteProtection of SLB instance before deleting. Default value: `false`. If `delete_protection_validation` is set to `true`, this resource will not be deleted when its SLB instance enabled DeleteProtection.
* `tags` - (Optional, Available since v1.227.1) A mapping of tags to assign to the resource.
* `servers` - (Deprecated since v1.163.0) The list of backend servers to be added. See [`servers`](#servers) below.
-> **NOTE:** Field `servers` has been deprecated from provider version 1.163.0, and it will be removed in the future version. Please use the new resource `alicloud_slb_server_group_server_attachment`.

### `servers`

The servers mapping supports the following:

* `type` - (Optional, Available since v1.51.0) Specify the type of the backend server. Default value: `ecs`. Valid values: `ecs`, `eni`.
* `port` - (Required, Int) The port used by the backend server. Valid values: `1` to `65535`.
* `weight` - (Optional, Int) Weight of the backend server. Default value: `100`. Valid values: `0` to `100`.
* `server_ids` - (Required, List) The list of Elastic Compute Service (ECS) Ids or Elastic Network Interface (ENI) Ids.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Virtual Backend Server Group.

## Import

Load Balancer Virtual Backend Server Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_slb_server_group.example <id>
```
