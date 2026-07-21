---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_router_interface_connection"
sidebar_current: "docs-alicloud-resource-router-interface-connection"
description: |-
  Provides a VPC router interface connection resource to connect two VPCs.
---

# alicloud_router_interface_connection

-> **DEPRECATED:**  This resource has been deprecated from version `1.199.0`. Please use new resource [alicloud_express_connect_router_interface](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/express_connect_router_interface).

Provides a VPC router interface connection resource to connect two router interfaces which are in two different VPCs.
After that, all of the two router interfaces will be active.

-> **NOTE:** At present, Router interface does not support changing opposite router interface, the connection delete action is only deactivating it to inactive, not modifying the connection to empty.

-> **NOTE:** If you want to changing opposite router interface, you can delete router interface and re-build them.

-> **NOTE:** A integrated router interface connection tunnel requires both InitiatingSide and AcceptingSide configuring opposite router interface.

-> **NOTE:** Please remember to add a `depends_on` clause in the router interface connection from the InitiatingSide to the AcceptingSide, because the connection from the AcceptingSide to the InitiatingSide must be done first.

## Example Usage

```terraform
provider "alicloud" {
  region = var.region
}

variable "region" {
  default = "cn-hangzhou"
}

variable "name" {
  default = "alicloudRouterInterfaceConnectionBasic"
}

resource "alicloud_vpc" "foo" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vpc" "bar" {
  provider   = alicloud
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_router_interface" "initiate" {
  opposite_region      = var.region
  router_type          = "VRouter"
  router_id            = alicloud_vpc.foo.router_id
  role                 = "InitiatingSide"
  specification        = "Large.2"
  name                 = var.name
  description          = var.name
  instance_charge_type = "PostPaid"
}

resource "alicloud_router_interface" "opposite" {
  provider        = alicloud
  opposite_region = var.region
  router_type     = "VRouter"
  router_id       = alicloud_vpc.bar.router_id
  role            = "AcceptingSide"
  specification   = "Large.1"
  name            = "${var.name}-opposite"
  description     = "${var.name}-opposite"
}

// A integrated router interface connection tunnel requires both InitiatingSide and AcceptingSide configuring opposite router interface.
resource "alicloud_router_interface_connection" "foo" {
  interface_id          = alicloud_router_interface.initiate.id
  opposite_interface_id = alicloud_router_interface.opposite.id
  depends_on            = [alicloud_router_interface_connection.bar] // The connection must start from the accepting side.
}

resource "alicloud_router_interface_connection" "bar" {
  provider              = alicloud
  interface_id          = alicloud_router_interface.opposite.id
  opposite_interface_id = alicloud_router_interface.initiate.id
}

```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_router_interface_connection&spm=docs.r.router_interface_connection.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `interface_id` - (Required, ForceNew) One side router interface ID.
* `opposite_interface_id` - (Required, ForceNew) Another side router interface ID. It must belong the specified "opposite_interface_owner_id" account.
* `opposite_interface_owner_id` - (Optional, ForceNew) Another side router interface account ID. Log on to the Alibaba Cloud console, select User Info > Account Management to check the account ID. Default to [Provider account_id](https://www.terraform.io/docs/providers/alicloud/index.html#account_id).
* `opposite_router_id` - (Optional, ForceNew) Another side router ID. It must belong the specified "opposite_interface_owner_id" account. It is valid when field "opposite_interface_owner_id" is specified.
* `opposite_router_type` - (Optional, ForceNew) Another side router Type. Optional value: VRouter, VBR. It is valid when field "opposite_interface_owner_id" is specified.

-> **NOTE:** The value of "opposite_interface_owner_id" or "account_id" must be main account and not be sub account.

## Attributes Reference

The following attributes are exported:

* `id` - Router interface ID. The value is equal to "interface_id".

## Import

The router interface connection can be imported using the id, e.g.

```shell
$ terraform import alicloud_router_interface_connection.foo ri-abc123456
```
