---
subcategory: "Eflo"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_node_group_attachment"
description: |-
  Provides a Alicloud Eflo Node Group Attachment resource.
---

# alicloud_eflo_node_group_attachment

Provides a Eflo Node Group Attachment resource.

Node Association Node Group Resources.

For information about Eflo Node Group Attachment and how to use it, see [What is Node Group Attachment](https://next.api.alibabacloud.com/document/eflo-controller/2022-12-15/ExtendCluster).

-> **NOTE:** Available since v1.255.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_eflo_node_group_attachment&exampleId=a1e5cf03-c0c8-71a4-27bb-f267258771d4b1637688&activeTab=example&spm=docs.r.eflo_node_group_attachment.0.a1e5cf03c0&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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


resource "alicloud_eflo_node_group_attachment" "default" {
  vswitch_id     = "vsw-uf63gbmvwgreao66opmie"
  hostname       = "attachment-example-e01-cn-smw4d1bzd0a"
  login_password = "G7f$2kL9@vQx3Zp5*"
  cluster_id     = "i118976621753269898628"
  node_group_id  = "i127582271753269898630"
  node_id        = "e01-cn-smw4d1bzd0a"
  vpc_id         = "vpc-uf6t73bb01dfprb2qvpqa"
}
```

VPD Cluster Usage

For a VPD cluster, omit `vpc_id` and `vswitch_id` and configure `vpd_subnets` together with `ip_allocation_policy`, so that the node receives its IP addresses from the cluster subnets of the VPD. `vpd_subnets` lists the cluster subnet IDs that can be used for IP allocation, and `ip_allocation_policy` assigns those subnets to bonds, machine types or specific nodes. The cluster, node group and node must already exist, and the subnet IDs must belong to the VPD of that cluster. For a VPC cluster, keep `vpc_id` and `vswitch_id` configured and omit both VPD arguments.

The following example attaches an existing node of a VPD cluster. Replace the variable defaults with the IDs of your own VPD cluster, node group, node and cluster subnet; `machine_type` must match the machine type of that node, and bond names such as `Bond0` must match the bond names reported for the node.

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "cluster_id" {
  description = "ID of an existing VPD cluster."
  default     = "i-xxxxxxxxxxxxxxxx"
}

variable "node_group_id" {
  description = "ID of the node group that owns the node."
  default     = "i-xxxxxxxxxxxxxxxx"
}

variable "node_id" {
  description = "ID of an existing node that is not attached yet."
  default     = "i-xxxxxxxxxxxxxxxx"
}

variable "vpd_subnet_id" {
  description = "ID of a cluster subnet of the VPD that the cluster belongs to."
  default     = "subnet-xxxxxxxx"
}

variable "machine_type" {
  description = "Machine type of the node, for example efg1.nvga1."
  default     = "efg1.nvga1"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_eflo_node_group_attachment" "default" {
  cluster_id    = var.cluster_id
  node_group_id = var.node_group_id
  node_id       = var.node_id
  hostname      = "attachment-vpd-example"

  vpd_subnets = [var.vpd_subnet_id]

  ip_allocation_policy {
    bond_policy {
      bond_default_subnet = var.vpd_subnet_id
      bonds {
        name   = "Bond0"
        subnet = var.vpd_subnet_id
      }
    }
  }
  ip_allocation_policy {
    machine_type_policy {
      machine_type = var.machine_type
      bonds {
        name   = "Bond0"
        subnet = var.vpd_subnet_id
      }
    }
  }
  ip_allocation_policy {
    node_policy {
      node_id = var.node_id
      bonds {
        name   = "Bond0"
        subnet = var.vpd_subnet_id
      }
    }
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_eflo_node_group_attachment&spm=docs.r.eflo_node_group_attachment.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `cluster_id` - (Optional, ForceNew, Computed) Cluster ID
* `data_disk` - (Optional, List) The data disk of the cloud disk to be attached to the node. See [`data_disk`](#data_disk) below.
* `hostname` - (Required, ForceNew) Node hostname
* `ip_allocation_policy` - (Optional, ForceNew, Set, Available since v2.3.0) The IP address allocation policies that assign the cluster subnets of a VPD to bonds, machine types or nodes. Each policy can select only one policy type, and multiple policies can be combined. Changing this argument replaces the resource. This argument is only sent when the resource is created and is not read back, so an imported resource has no value for it. See [`ip_allocation_policy`](#ip_allocation_policy) below.
* `login_password` - (Optional) Node login password
* `node_group_id` - (Optional, ForceNew, Computed) Node group ID
* `node_id` - (Optional, ForceNew, Computed) Node ID
* `user_data` - (Optional) User-defined data
* `vpd_subnets` - (Optional, ForceNew, Set, Available since v2.3.0) The cluster subnet IDs used to allocate IP addresses for the node in a VPD cluster. Configure it together with `ip_allocation_policy` and omit `vpc_id` and `vswitch_id` for a VPD cluster. Changing this argument replaces the resource. This argument is only sent when the resource is created and is not read back, so an imported resource has no value for it.
* `vswitch_id` - (Optional, ForceNew, Computed, String) vswitch id. Required for a VPC cluster; can be omitted for a VPD cluster.
* `vpc_id` - (Optional, ForceNew, Computed, String) Vpc id. Required for a VPC cluster; can be omitted for a VPD cluster.

### `data_disk`

The data_disk supports the following:
* `category` - (Optional) Type
* `delete_with_node` - (Optional) Indicate whether the data disk is released with the node. true indicates that the data disk will be released together when the node unsubscribes.
* `performance_level` - (Optional) Performance level
* `size` - (Optional, Int) Data disk size

### `ip_allocation_policy`

The ip_allocation_policy supports the following:
* `bond_policy` - (Optional, List, Available since v2.3.0) Assigns a cluster subnet to each bond. See [`bond_policy`](#ip_allocation_policy-bond_policy) below.
* `machine_type_policy` - (Optional, Set, Available since v2.3.0) Assigns cluster subnets to bonds of a machine type. See [`machine_type_policy`](#ip_allocation_policy-machine_type_policy) below.
* `node_policy` - (Optional, Set, Available since v2.3.0) Assigns cluster subnets to bonds of a specific node. See [`node_policy`](#ip_allocation_policy-node_policy) below.

### `ip_allocation_policy-bond_policy`

The ip_allocation_policy-bond_policy supports the following:
* `bond_default_subnet` - (Optional, String, Available since v2.3.0) The ID of the default cluster subnet of the bond.
* `bonds` - (Optional, Set, Available since v2.3.0) The bond information. See [`bonds`](#ip_allocation_policy-bond_policy-bonds) below.

### `ip_allocation_policy-machine_type_policy`

The ip_allocation_policy-machine_type_policy supports the following:
* `bonds` - (Optional, Set, Available since v2.3.0) The bond information. See [`bonds`](#ip_allocation_policy-machine_type_policy-bonds) below.
* `machine_type` - (Optional, String, Available since v2.3.0) The machine type of the node, for example `efg1.nvga1`.

### `ip_allocation_policy-node_policy`

The ip_allocation_policy-node_policy supports the following:
* `bonds` - (Optional, Set, Available since v2.3.0) The bond information. See [`bonds`](#ip_allocation_policy-node_policy-bonds) below.
* `hostname` - (Optional, String, Available since v2.3.0) The hostname of the node.
* `node_id` - (Optional, String, Available since v2.3.0) The ID of the node, for example `i-3fdodw2`.

### `ip_allocation_policy-bond_policy-bonds`

The ip_allocation_policy-bond_policy-bonds supports the following:
* `name` - (Optional, String, Available since v2.3.0) The bond name, for example `Bond0`.
* `subnet` - (Optional, String, Available since v2.3.0) The ID of the cluster subnet that provides IP addresses for the bond.

### `ip_allocation_policy-machine_type_policy-bonds`

The ip_allocation_policy-machine_type_policy-bonds supports the following:
* `name` - (Optional, String, Available since v2.3.0) The bond name, for example `Bond0`.
* `subnet` - (Optional, String, Available since v2.3.0) The ID of the cluster subnet that provides IP addresses for the bond.

### `ip_allocation_policy-node_policy-bonds`

The ip_allocation_policy-node_policy-bonds supports the following:
* `name` - (Optional, String, Available since v2.3.0) The bond name, for example `Bond0`.
* `subnet` - (Optional, String, Available since v2.3.0) The ID of the cluster subnet that provides IP addresses for the bond.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<cluster_id>:<node_group_id>:<node_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 3605 mins) Used when create the Node Group Attachment.
* `delete` - (Defaults to 3605 mins) Used when delete the Node Group Attachment.

## Import

Eflo Node Group Attachment can be imported using the id, which consists of cluster_id, node_group_id and node_id, e.g.

```shell
$ terraform import alicloud_eflo_node_group_attachment.example <cluster_id>:<node_group_id>:<node_id>
```