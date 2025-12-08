---
subcategory: "Eflo"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_hyper_node"
description: |-
  Provides a Alicloud Eflo Hyper Node resource.
---

# alicloud_eflo_hyper_node

Provides a Eflo Hyper Node resource.

Hyper computing node.

For information about Eflo Hyper Node and how to use it, see [What is Hyper Node](https://next.api.alibabacloud.com/document/BssOpenApi/2017-12-14/CreateInstance).

-> **NOTE:** Available since v1.265.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_eflo_hyper_node" "default" {
  zone_id          = "ap-southeast-7a"
  product_code     = "bccluster"
  machine_type     = "efg3.GN9A.ch72"
  product_type     = "bccluster_eflocomputing_public_cn"
  hpn_zone         = "A1"
  server_arch      = "bmserver"
  payment_duration = "1"
  payment_type     = "Subscription"
  stage_num        = "1"
}
```

## Argument Reference

The following arguments are supported:
* `cluster_id` - (Optional) Cluster ID
* `data_disk` - (Optional, List) List of disk information See [`data_disk`](#data_disk) below.

-> **NOTE:** This parameter only applies during resource update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `hostname` - (Optional) Host name
* `hpn_zone` - (Optional, ForceNew) Number of the cluster to which the supercompute node belongs
* `login_password` - (Optional) Login Password
* `machine_type` - (Optional, ForceNew) The model used by the super computing node
* `node_group_id` - (Optional) Node group ID
* `payment_duration` - (Optional, Int) The duration of the instance purchase, in units.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `payment_type` - (Required, ForceNew) The payment type of the resource
* `renewal_duration` - (Optional, Int) Number of auto-renewal cycles
* `renewal_status` - (Optional, Computed) Automatic renewal status. Value: AutoRenewal: automatic renewal. ManualRenewal: manual renewal. The default ManualRenewal.
* `resource_group_id` - (Optional, Computed) The ID of the resource group
* `server_arch` - (Optional, ForceNew) Super Node Architecture
* `stage_num` - (Optional) The number of installments of the supercomputing node of the fixed fee installment.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `tags` - (Optional, Map) The tag of the resource
* `user_data` - (Optional) Custom Data
* `vswitch_id` - (Optional) Switch ID
* `vpc_id` - (Optional) VPC ID
* `zone_id` - (Optional, ForceNew) The zone where the super compute node is located

### `data_disk`

The data_disk supports the following:
* `bursting_enabled` - (Optional) Whether to enable Burst (performance Burst).
* `category` - (Optional) The disk type. Value range:
  - cloud_essd:ESSD cloud disk.
* `delete_with_node` - (Optional) Whether the data disk is unsubscribed and deleted with the node.
* `performance_level` - (Optional) When creating an ESSD cloud disk to use as a system disk, set the performance level of the cloud disk. Value range:
  - PL0: maximum random read/write IOPS 10000 for a single disk.
  - PL1: maximum random read/write IOPS 50000 for a single disk.
* `provisioned_iops` - (Optional, Int) ESSD AutoPL cloud disk (single disk) pre-configuration performance of IOPS.
* `size` - (Optional, Int) The size of the disk. The unit is GiB.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource
* `region_id` - The region ID of the resource
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Hyper Node.
* `delete` - (Defaults to 5 mins) Used when delete the Hyper Node.
* `update` - (Defaults to 38 mins) Used when update the Hyper Node.

## Import

Eflo Hyper Node can be imported using the id, e.g.

```shell
$ terraform import alicloud_eflo_hyper_node.example <id>
```