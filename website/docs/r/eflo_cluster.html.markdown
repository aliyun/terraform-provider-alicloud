---
subcategory: "Eflo"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_cluster"
description: |-
  Provides a Alicloud Eflo Cluster resource.
---

# alicloud_eflo_cluster

Provides a Eflo Cluster resource.

Large computing cluster.

For information about Eflo Cluster and how to use it, see [What is Cluster](https://next.api.alibabacloud.com/document/eflo-controller/2022-12-15/CreateCluster).

-> **NOTE:** Available since v1.246.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_eflo_cluster&exampleId=9eb8c140-bb26-5016-39e0-391b050fbc9deb8009f4&activeTab=example&spm=docs.r.eflo_cluster.0.9eb8c140bb&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_eflo_cluster&spm=docs.r.eflo_cluster.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `cluster_description` - (Optional, ForceNew) cluster description
* `cluster_name` - (Optional, ForceNew) ClusterName
* `cluster_type` - (Optional, ForceNew) cluster type
* `components` - (Optional, List) Component (software instance) See [`components`](#components) below.
* `hpn_zone` - (Optional) Cluster Number
* `ignore_failed_node_tasks` - (Optional) Whether to allow skipping failed nodes. Default value: False
* `networks` - (Optional, List) Network Information See [`networks`](#networks) below.
* `nimiz_vswitches` - (Optional, List) Node virtual switch
* `node_groups` - (Optional, List) Node Group List See [`node_groups`](#node_groups) below.
* `open_eni_jumbo_frame` - (Optional) Whether the network interface supports jumbo frames
* `resource_group_id` - (Optional, Computed) The ID of the resource group
* `tags` - (Optional, Map) tag

### `components`

The components supports the following:
* `component_config` - (Optional, List) Component Configuration See [`component_config`](#components-component_config) below.
* `component_type` - (Optional) Component Type

### `components-component_config`

The components-component_config supports the following:
* `basic_args` - (Optional) Component Basic Parameters
* `node_units` - (Optional, List) Node pool configuration, and is used to establish the corresponding relationship between node groups and node pools. When
ComponentType = "ACKEdge" is required. Other values are empty.

### `networks`

The networks supports the following:
* `ip_allocation_policy` - (Optional, List) IP allocation policy See [`ip_allocation_policy`](#networks-ip_allocation_policy) below.
* `new_vpd_info` - (Optional, List) Vpd configuration information See [`new_vpd_info`](#networks-new_vpd_info) below.
* `security_group_id` - (Optional) Security group ID
* `tail_ip_version` - (Optional) IP version
* `vswitch_id` - (Optional) Switch ID
* `vswitch_zone_id` - (Optional) Switch ZoneID
* `vpc_id` - (Optional) VPC ID
* `vpd_info` - (Optional, List) Multiplexing VPD information See [`vpd_info`](#networks-vpd_info) below.

### `networks-ip_allocation_policy`

The networks-ip_allocation_policy supports the following:
* `bond_policy` - (Optional, List) Bond policy See [`bond_policy`](#networks-ip_allocation_policy-bond_policy) below.
* `machine_type_policy` - (Optional, List) Model Assignment Policy See [`machine_type_policy`](#networks-ip_allocation_policy-machine_type_policy) below.
* `node_policy` - (Optional, List) Node allocation policy See [`node_policy`](#networks-ip_allocation_policy-node_policy) below.

### `networks-new_vpd_info`

The networks-new_vpd_info supports the following:
* `cen_id` - (Optional) Cloud Enterprise Network ID
* `cloud_link_cidr` - (Optional) Cloud chain cidr
* `cloud_link_id` - (Optional) Cloud chain ID
* `monitor_vpc_id` - (Optional) Proprietary Network
* `monitor_vswitch_id` - (Optional) Proprietary network switch
* `vpd_cidr` - (Optional) Cluster network segment
* `vpd_subnets` - (Optional, List) Cluster Subnet See [`vpd_subnets`](#networks-new_vpd_info-vpd_subnets) below.

### `networks-vpd_info`

The networks-vpd_info supports the following:
* `vpd_id` - (Optional) VPC ID
* `vpd_subnets` - (Optional, List) List of cluster subnet ID

### `networks-new_vpd_info-vpd_subnets`

The networks-new_vpd_info-vpd_subnets supports the following:
* `subnet_cidr` - (Optional) Subnet cidr
* `subnet_type` - (Optional) Subnet Type
* `zone_id` - (Optional) Zone ID

### `networks-ip_allocation_policy-bond_policy`

The networks-ip_allocation_policy-bond_policy supports the following:
* `bond_default_subnet` - (Optional) Default bond cluster subnet
* `bonds` - (Optional, List) Bond information See [`bonds`](#networks-ip_allocation_policy-bond_policy-bonds) below.

### `networks-ip_allocation_policy-machine_type_policy`

The networks-ip_allocation_policy-machine_type_policy supports the following:
* `bonds` - (Optional, List) Bond information See [`bonds`](#networks-ip_allocation_policy-machine_type_policy-bonds) below.
* `machine_type` - (Optional) Model

### `networks-ip_allocation_policy-node_policy`

The networks-ip_allocation_policy-node_policy supports the following:
* `bonds` - (Optional, List) Bond information See [`bonds`](#networks-ip_allocation_policy-node_policy-bonds) below.
* `node_id` - (Optional) Node ID

### `networks-ip_allocation_policy-node_policy-bonds`

The networks-ip_allocation_policy-node_policy-bonds supports the following:
* `name` - (Optional) The bond name
* `subnet` - (Optional) IP source cluster subnet

### `networks-ip_allocation_policy-machine_type_policy-bonds`

The networks-ip_allocation_policy-machine_type_policy-bonds supports the following:
* `name` - (Optional) The bond name
* `subnet` - (Optional) IP source cluster subnet

### `networks-ip_allocation_policy-bond_policy-bonds`

The networks-ip_allocation_policy-bond_policy-bonds supports the following:
* `name` - (Optional) The bond name
* `subnet` - (Optional) IP source cluster subnet

### `node_groups`

The node_groups supports the following:
* `image_id` - (Optional) System Image ID
* `machine_type` - (Optional) Model
* `node_group_description` - (Optional) Node Group Description
* `node_group_name` - (Optional) Node Group Name
* `nodes` - (Optional, List) Node List See [`nodes`](#node_groups-nodes) below.
* `user_data` - (Optional) Instance custom data. It needs to be encoded in Base64 mode, and the original data is at most 16KB.
* `zone_id` - (Optional) Zone ID

### `node_groups-nodes`

The node_groups-nodes supports the following:
* `hostname` - (Optional) Host name
* `login_password` - (Optional) Login Password
* `node_id` - (Optional) Node ID
* `vswitch_id` - (Optional) Virtual Switch ID
* `vpc_id` - (Optional) VPC ID

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Cluster.
* `delete` - (Defaults to 5 mins) Used when delete the Cluster.
* `update` - (Defaults to 5 mins) Used when update the Cluster.

## Import

Eflo Cluster can be imported using the id, e.g.

```shell
$ terraform import alicloud_eflo_cluster.example <id>
```