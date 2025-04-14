---
subcategory: "Eflo"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_resource"
description: |-
  Provides a Alicloud Eflo Resource resource.
---

# alicloud_eflo_resource

Provides a Eflo Resource resource.



For information about Eflo Resource and how to use it, see [What is Resource](https://www.alibabacloud.com/help/en/pai/developer-reference/api-eflo-cnp-2023-08-28-createresource).

-> **NOTE:** Available since v1.248.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-wulanchabu"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_eflo_resource" "default" {
  user_access_param {
    access_id    = "your_access_id"
    access_key   = "your_access_key"
    workspace_id = "your_workspace_id"
    endpoint     = "your_endpoint"
  }
  cluster_id = "terraform-${random_integer.default.result}"
  machine_types {
    memory_info  = "32x 64GB DDR4 4800 Memory"
    type         = "Private"
    bond_num     = "5"
    node_count   = "1"
    cpu_info     = "2x Intel Saphhire Rapid 8469C 48C CPU"
    network_info = "1x 200Gbps Dual Port BF3 DPU for VPC 4x 200Gbps Dual Port EIC"
    gpu_info     = "8x OAM 810 GPU"
    disk_info    = "2x 480GB SATA SSD 4x 3.84TB NVMe SSD"
    network_mode = "net"
    name         = "lingjun"
  }
  cluster_name = var.name
  cluster_desc = var.name
}
```

### Deleting `alicloud_eflo_resource` or removing it from your configuration

Terraform cannot destroy resource `alicloud_eflo_resource`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `cluster_desc` - (Optional, ForceNew) Used to provide a description or comment on the compute cluster.
* `cluster_id` - (Required, ForceNew) Used to uniquely identify a computing cluster.
* `cluster_name` - (Required, ForceNew) Represents the name of the compute cluster, usually including the model number.
* `machine_types` - (Required, ForceNew, Set) Generally refers to the type or instance type of a computing resource. See [`machine_types`](#machine_types) below.
* `user_access_param` - (Required, Set) Used to define the access parameters for the user. See [`user_access_param`](#user_access_param) below.

### `machine_types`

The machine_types supports the following:
* `bond_num` - (Optional, ForceNew, Int) This property specifies the number of network bindings, which relates to the number of physical or virtual network cards connected to the network through the network interface card (NIC). Multiple network bindings can increase bandwidth and redundancy and improve network reliability.
* `cpu_info` - (Required, ForceNew) Provides CPU details, including the number of cores, number of threads, clock frequency, and architecture type. This information helps to evaluate the processing power and identify whether it can meet the performance requirements of a particular application.
* `disk_info` - (Optional, ForceNew) Displays information about the storage device, including the disk type (such as SSD or HDD), capacity, and I/O performance. Storage performance is critical in data-intensive applications such as big data processing and databases.
* `gpu_info` - (Required, ForceNew) Provides detailed information about the GPU, including the number, model, memory size, and computing capability. This information is particularly important for tasks such as deep learning, scientific computing, and graph processing, helping users understand the graph processing capabilities of nodes.
* `memory_info` - (Optional, ForceNew) This property provides memory details, including total memory, available memory, and usage. This helps users understand the memory processing capabilities of compute nodes, especially when running heavy-duty applications.
* `name` - (Optional, ForceNew) Specification Name.
* `network_info` - (Optional, ForceNew) Contains detailed information about the network interface, such as network bandwidth, latency, protocol types supported by the network, IP addresses, and network topology. Optimizing network information is essential to ensure efficient data transmission and low latency.
* `network_mode` - (Optional, ForceNew) Specifies the network mode, such as bridge mode, NAT mode, or direct connection mode. Different network modes affect the network configuration and data transmission performance of nodes, and affect the network access methods of computing instances.
* `node_count` - (Optional, ForceNew, Int) Specifies the total number of compute nodes. This property is particularly important in distributed computing and cluster environments, because the number of nodes often directly affects the computing power and the ability to parallel processing.
* `type` - (Optional, ForceNew) Usually refers to a specific resource type (such as virtual machine, physical server, container, etc.), which is used to distinguish different computing units or resource categories.

### `user_access_param`

The user_access_param supports the following:
* `access_id` - (Required) Access keys are important credentials for authentication.
* `access_key` - (Required) A Secret Key is a Secret credential paired with an access Key to verify a user's identity and protect the security of an interface.
* `endpoint` - (Required) An Endpoint is a network address for accessing a service or API, usually a URL to a specific service instance.
* `workspace_id` - (Required) A Workspace generally refers to a separate space created by a user on a particular computing environment or platform.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `resource_id` - The ID of the Resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Resource.
* `update` - (Defaults to 5 mins) Used when update the Resource.

## Import

Eflo Resource can be imported using the id, e.g.

```shell
$ terraform import alicloud_eflo_resource.example <id>
```

