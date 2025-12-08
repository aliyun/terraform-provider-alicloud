---
subcategory: "Eflo"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_node"
description: |-
  Provides a Alicloud Eflo Node resource.
---

# alicloud_eflo_node

Provides a Eflo Node resource.

Large computing node.

For information about Eflo Node and how to use it, see [What is Node](https://next.api.alibabacloud.com/document/BssOpenApi/2017-12-14/CreateInstance).

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

resource "alicloud_eflo_node" "default" {
  period         = "36"
  discount_level = "36"
  billing_cycle  = "1month"
  classify       = "gpuserver"
  zone           = "cn-hangzhou-b"
  product_form   = "instance"
  payment_ratio  = "0"
  hpn_zone       = "B1"
  server_arch    = "bmserver"
  machine_type   = "efg1.nvga1n"
  stage_num      = "36"
  renewal_status = "AutoRenewal"
  renew_period   = "36"
  status         = "Unused"
}
```
Creating a PayAsYouGo eflo node
```terraform
resource "alicloud_eflo_node" "payasyougo" {
  machine_type   = "efg1.nvga8n"
  payment_type   = "PayAsYouGo"
  hpn_zone       = "A1"
  product_form   = "instance"
  renewal_status = "ManualRenewal"
  zone           = "cn-wulanchabu-a"
  tags = {
    From = "Terraform"
  }
  # status = "Unused"
  cluster_id     = "i11922307xxxxxxx"
  node_group_id  = "i1254705xxxxxxxx"
  hostname       = "terraform-example"
  login_password = "xxxxxxxx"
  data_disk {
    size              = 120
    category          = "cloud_essd"
    performance_level = "PL0"
  }
  data_disk {
    size              = 120
    category          = "cloud_essd"
    performance_level = "PL1"
  }

  ip_allocation_policy {
    machine_type_policy {
      machine_type = "efg1.nvga8n"

      bonds {
        subnet = "subnet-x1xxx"
        name   = "example01"
      }
      bonds {
        subnet = "subnet-xxxx"
        name   = "example02"
      }
      bonds {
        subnet = "subnet-xxxx"
        name   = "example03"
      }
      bonds {
        subnet = "subnet-xxxx"
        name   = "example04"
      }
      bonds {
        subnet = "subnet-xxxx"
        name   = "example05"
      }
    }
  }
}
```

### Deleting `alicloud_eflo_node` or removing it from your configuration

The `alicloud_eflo_node` resource allows you to manage  `payment_type = "Subscription"`  instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the Instance.
You can resume managing the subscription instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:

* `install_pai` - (Optional) Whether to buy PAI. default value `false`.
* `billing_cycle` - (Optional) Billing cycle

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `classify` - (Optional) Classification

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `cluster_id` - (Optional, Available since v1.265.0) Cluster id
* `computing_server` - (Optional, ForceNew, Computed, Deprecated since v1.265.0) Node Model
* `discount_level` - (Optional) Offer Information

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `hostname` - (Optional, ForceNew, Available since v1.265.0) Host name
* `hpn_zone` - (Optional, ForceNew) Cluster Number
* `ip_allocation_policy` - (Optional, ForceNew, List, Available since v1.265.0) IP address combination policy: only one policy type can be selected for each policy, and multiple policies can be combined. See [`ip_allocation_policy`](#ip_allocation_policy) below.
* `login_password` - (Optional, ForceNew, Available since v1.265.0) Login Password
* `machine_type` - (Optional, ForceNew, Computed, Available since v1.261.0) Model
* `node_group_id` - (Optional, Available since v1.265.0) node group id
* `node_type` - (Optional, Computed, Available since v1.265.0) node type
* `payment_ratio` - (Optional) Down payment ratio

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `payment_type` - (Optional, ForceNew, Computed, Available since v1.261.0) The payment method of the node. Value range: Subscription: fixed fee installment; PayAsYouGo: pay by volume.
The default is Subscription.
* `period` - (Optional, Int) Prepaid cycle. The unit is Month, please enter an integer multiple of 12 for the annual payment product.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `product_form` - (Optional) Form

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `renew_period` - (Optional, Int) Automatic renewal period, in months.

-> **NOTE:**  When setting `RenewalStatus` to `AutoRenewal`, it must be set.

* `renewal_status` - (Optional) Automatic renewal status, value:
  - AutoRenewal: automatic renewal.
  - ManualRenewal: manual renewal.

The default ManualRenewal.
* `resource_group_id` - (Optional, Computed) The ID of the resource group
* `server_arch` - (Optional) Architecture

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `stage_num` - (Optional) Number of stages

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `status` - (Optional, Computed) The status of the resource
* `tags` - (Optional, Map) The tag of the resource
* `user_data` - (Optional, ForceNew, Available since v1.265.0) Custom Data
* `vswitch_id` - (Optional, ForceNew, Available since v1.265.0) Switch ID
* `vpc_id` - (Optional, ForceNew, Available since v1.265.0) VPC ID
* `zone` - (Optional, ForceNew) Availability Zone
* `data_disk` - (Optional, List, Available since v1.265.0) The data disk of the cloud disk to be attached to the node. See [`data_disk`](#data_disk) below.

### `data_disk`

The data_disk supports the following:
* `category` - (Optional) Data disk type
* `performance_level` - (Optional) Performance level
* `size` - (Optional, Int) Data disk size

### `ip_allocation_policy`

The ip_allocation_policy supports the following:
* `bond_policy` - (Optional, ForceNew, List, Available since v1.265.0) Specify the cluster subnet ID based on the bond name See [`bond_policy`](#ip_allocation_policy-bond_policy) below.
* `machine_type_policy` - (Optional, ForceNew, List, Available since v1.265.0) Model Assignment Policy See [`machine_type_policy`](#ip_allocation_policy-machine_type_policy) below.
* `node_policy` - (Optional, ForceNew, List, Available since v1.265.0) Node allocation policy See [`node_policy`](#ip_allocation_policy-node_policy) below.

### `ip_allocation_policy-bond_policy`

The ip_allocation_policy-bond_policy supports the following:
* `bond_default_subnet` - (Optional, ForceNew, Available since v1.265.0) Default bond cluster subnet
* `bonds` - (Optional, ForceNew, List, Available since v1.265.0) Bond information See [`bonds`](#ip_allocation_policy-bond_policy-bonds) below.

### `ip_allocation_policy-machine_type_policy`

The ip_allocation_policy-machine_type_policy supports the following:
* `bonds` - (Optional, ForceNew, List, Available since v1.265.0) Bond information See [`bonds`](#ip_allocation_policy-machine_type_policy-bonds) below.
* `machine_type` - (Optional, ForceNew, Available since v1.265.0) Model

### `ip_allocation_policy-node_policy`

The ip_allocation_policy-node_policy supports the following:
* `bonds` - (Optional, ForceNew, List, Available since v1.265.0) Bond information See [`bonds`](#ip_allocation_policy-node_policy-bonds) below.
* `hostname` - (Optional, ForceNew, Available since v1.265.0) Host name
* `node_id` - (Optional, ForceNew, Available since v1.265.0) Node ID

### `ip_allocation_policy-node_policy-bonds`

The ip_allocation_policy-node_policy-bonds supports the following:
* `name` - (Optional, ForceNew, Available since v1.265.0) Bond Name
* `subnet` - (Optional, ForceNew, Available since v1.265.0) IP source cluster subnet

### `ip_allocation_policy-machine_type_policy-bonds`

The ip_allocation_policy-machine_type_policy-bonds supports the following:
* `name` - (Optional, ForceNew, Available since v1.265.0) Bond Name
* `subnet` - (Optional, ForceNew, Available since v1.265.0) IP source cluster subnet

### `ip_allocation_policy-bond_policy-bonds`

The ip_allocation_policy-bond_policy-bonds supports the following:
* `name` - (Optional, ForceNew, Available since v1.265.0) Bond Name
* `subnet` - (Optional, ForceNew, Available since v1.265.0) IP source cluster subnet

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource
* `region_id` - The region ID of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 10 mins) Used when create the Node.
* `delete` - (Defaults to 5 mins) Used when delete the Node.
* `update` - (Defaults to 6 mins) Used when update the Node.

## Import

Eflo Node can be imported using the id, e.g.

```shell
$ terraform import alicloud_eflo_node.example <id>
```