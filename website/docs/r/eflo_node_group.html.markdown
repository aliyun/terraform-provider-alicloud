---
subcategory: "Eflo"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_node_group"
description: |-
  Provides a Alicloud Eflo Node Group resource.
---

# alicloud_eflo_node_group

Provides a Eflo Node Group resource.

Node group. Divide a cluster into multiple node groups, each containing multiple nodes.

For information about Eflo Node Group and how to use it, see [What is Node Group](https://next.api.alibabacloud.com/document/eflo-controller/2022-12-15/CreateNodeGroup).

-> **NOTE:** Available since v1.246.0.

## Example Usage

Basic Usage

```terraform
# Before executing this example, you need to confirm with the product team whether the resources are sufficient or you will get an error message with "Failure to check order before create instance"
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "create_vpc" {
  cidr_block = "192.168.0.0/16"
  vpc_name   = "cluster-resoure-example"
}

resource "alicloud_vswitch" "create_vswitch" {
  vpc_id       = alicloud_vpc.create_vpc.id
  zone_id      = "cn-hangzhou-b"
  cidr_block   = "192.168.0.0/24"
  vswitch_name = "cluster-resoure-example"
}

resource "alicloud_security_group" "create_security_group" {
  description         = "sg"
  security_group_name = "cluster-resoure-example"
  security_group_type = "normal"
  vpc_id              = alicloud_vpc.create_vpc.id
}

resource "alicloud_eflo_cluster" "default" {
  cluster_description  = "cluster-resource-example"
  open_eni_jumbo_frame = "false"
  hpn_zone             = "B1"
  nimiz_vswitches = [
    "1111"
  ]
  ignore_failed_node_tasks = "true"
  resource_group_id        = data.alicloud_resource_manager_resource_groups.default.ids.1
  node_groups {
    image_id               = "i198448731735114628708"
    zone_id                = "cn-hangzhou-b"
    node_group_name        = "cluster-resource-example"
    node_group_description = "cluster-resource-example"
    machine_type           = "efg2.C48cA3sen"
  }

  networks {
    tail_ip_version = "ipv4"
    new_vpd_info {
      monitor_vpc_id     = alicloud_vpc.create_vpc.id
      monitor_vswitch_id = alicloud_vswitch.create_vswitch.id
      cen_id             = "11111"
      cloud_link_id      = "1111"
      vpd_cidr           = "111"
      vpd_subnets {
        zone_id     = "1111"
        subnet_cidr = "111"
        subnet_type = "111"
      }
      cloud_link_cidr = "169.254.128.0/23"
    }

    security_group_id = alicloud_security_group.create_security_group.id
    vswitch_zone_id   = "cn-hangzhou-b"
    vpc_id            = alicloud_vpc.create_vpc.id
    vswitch_id        = alicloud_vswitch.create_vswitch.id
    vpd_info {
      vpd_id = "111"
      vpd_subnets = [
        "111"
      ]
    }
    ip_allocation_policy {
      bond_policy {
        bond_default_subnet = "111"
        bonds {
          name   = "111"
          subnet = "111"
        }
      }
      machine_type_policy {
        bonds {
          name   = "111"
          subnet = "111"
        }
        machine_type = "111"
      }
      node_policy {
        bonds {
          name   = "111"
          subnet = "111"
        }
        node_id = "111"
      }
    }
  }

  cluster_name = "tfacceflo7165"
  cluster_type = "Lite"
}

resource "alicloud_eflo_node" "default" {
  period           = "36"
  discount_level   = "36"
  billing_cycle    = "1month"
  classify         = "gpuserver"
  zone             = "cn-hangzhou-b"
  product_form     = "instance"
  payment_ratio    = "0"
  hpn_zone         = "B1"
  server_arch      = "bmserver"
  computing_server = "efg1.nvga1n"
  stage_num        = "36"
  renewal_status   = "AutoRenewal"
  renew_period     = "36"
  status           = "Unused"
}

resource "alicloud_eflo_node_group" "default" {
  nodes {
    node_id        = alicloud_eflo_node.default.id
    vpc_id         = alicloud_vpc.create_vpc.id
    vswitch_id     = alicloud_vswitch.create_vswitch.id
    hostname       = "jxyhostname"
    login_password = "Alibaba@2025"
  }

  ignore_failed_node_tasks = "true"
  cluster_id               = alicloud_eflo_cluster.default.id
  image_id                 = "i195048661660874208657"
  zone_id                  = "cn-hangzhou-b"
  vpd_subnets = [
    "example"
  ]
  user_data       = "YWxpLGFsaSxhbGliYWJh"
  vswitch_zone_id = "cn-hangzhou-b"
  ip_allocation_policy {
    bond_policy {
      bond_default_subnet = "example"
      bonds {
        name   = "example"
        subnet = "example"
      }
    }
    machine_type_policy {
      bonds {
        name   = "example"
        subnet = "example"
      }
      machine_type = "example"
    }
    node_policy {
      node_id = alicloud_eflo_node.default.id
      bonds {
        name   = "example"
        subnet = "example"
      }
    }
  }
  machine_type           = "efg1.nvga1"
  az                     = "cn-hangzhou-b"
  node_group_description = "resource-example1"
  node_group_name        = "tfacceflo63657_update"
}
```

## Argument Reference

The following arguments are supported:
* `az` - (Required, ForceNew) Az
* `cluster_id` - (Required, ForceNew) Cluster ID
* `ignore_failed_node_tasks` - (Optional) Whether to allow skipping failed nodes. Default value: False
* `image_id` - (Required, ForceNew) Image ID
* `ip_allocation_policy` - (Optional, List) IP address combination policy: only one policy type can be selected for each policy, and multiple policies can be combined. See [`ip_allocation_policy`](#ip_allocation_policy) below.
* `machine_type` - (Required, ForceNew) Machine type
* `node_group_description` - (Optional, ForceNew) NodeGroupDescription
* `node_group_name` - (Required) The name of the resource
* `nodes` - (Optional, Set) Node List See [`nodes`](#nodes) below.
* `user_data` - (Optional) Custom Data
* `vswitch_zone_id` - (Optional) Zone ID of the switch
* `vpd_subnets` - (Optional, List) Cluster subnet list
* `zone_id` - (Optional) Zone ID

### `ip_allocation_policy`

The ip_allocation_policy supports the following:
* `bond_policy` - (Optional, List) Specify the cluster subnet ID based on the bond name See [`bond_policy`](#ip_allocation_policy-bond_policy) below.
* `machine_type_policy` - (Optional, List) Model Assignment Policy See [`machine_type_policy`](#ip_allocation_policy-machine_type_policy) below.
* `node_policy` - (Optional, List) Node allocation policy See [`node_policy`](#ip_allocation_policy-node_policy) below.

### `ip_allocation_policy-bond_policy`

The ip_allocation_policy-bond_policy supports the following:
* `bond_default_subnet` - (Optional) Default bond cluster subnet
* `bonds` - (Optional, List) Bond information See [`bonds`](#ip_allocation_policy-bond_policy-bonds) below.

### `ip_allocation_policy-machine_type_policy`

The ip_allocation_policy-machine_type_policy supports the following:
* `bonds` - (Optional, List) Bond information See [`bonds`](#ip_allocation_policy-machine_type_policy-bonds) below.
* `machine_type` - (Optional) Model

### `ip_allocation_policy-node_policy`

The ip_allocation_policy-node_policy supports the following:
* `bonds` - (Optional, List) Bond information See [`bonds`](#ip_allocation_policy-node_policy-bonds) below.
* `node_id` - (Optional) Node ID

### `ip_allocation_policy-node_policy-bonds`

The ip_allocation_policy-node_policy-bonds supports the following:
* `name` - (Optional) The bond name
* `subnet` - (Optional) IP source cluster subnet

### `ip_allocation_policy-machine_type_policy-bonds`

The ip_allocation_policy-machine_type_policy-bonds supports the following:
* `name` - (Optional) The bond name
* `subnet` - (Optional) IP source cluster subnet

### `ip_allocation_policy-bond_policy-bonds`

The ip_allocation_policy-bond_policy-bonds supports the following:
* `name` - (Optional) The bond name
* `subnet` - (Optional) IP source cluster subnet

### `nodes`

The nodes supports the following:
* `hostname` - (Optional) Host name
* `login_password` - (Optional) Login Password
* `node_id` - (Optional) Node ID
* `vswitch_id` - (Optional) Switch ID
* `vpc_id` - (Optional) VPC ID

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<cluster_id>:<node_group_id>`.
* `create_time` - Create time
* `node_group_id` - The first ID of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Node Group.
* `delete` - (Defaults to 5 mins) Used when delete the Node Group.
* `update` - (Defaults to 120 mins) Used when update the Node Group.

## Import

Eflo Node Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_eflo_node_group.example <cluster_id>:<node_group_id>
```