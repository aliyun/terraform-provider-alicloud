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

For information about Eflo Hyper Node and how to use it, see [What is Hyper Node](https://www.alibabacloud.com/help/en/pai/developer-reference/api-eflo-controller-2022-12-15-overview).

-> **NOTE:** Available since v1.264.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "ap-southeast-7"
}

resource "alicloud_eflo_hyper_node" "default" {
  zone_id          = "ap-southeast-7a"
  machine_type     = "efg3.GN9A.ch72"
  hpn_zone         = "A1"
  server_arch      = "bmserver"
  payment_duration = "1"
  payment_type     = "Subscription"
  stage_num        = "1"
  renewal_duration = 2
  renewal_status   = "ManualRenewal"
  tags = {
    From = "Terraform"
    Env  = "Product"
  }
}
```

## Argument Reference

The following arguments are supported:
* `hpn_zone` - (Optional, ForceNew) Number of the cluster to which the supercompute node belongs
* `machine_type` - (Optional, ForceNew) The model used by the super computing node
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
* `zone_id` - (Optional, ForceNew) The zone where the super compute node is located

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
* `update` - (Defaults to 5 mins) Used when update the Hyper Node.

## Import

Eflo Hyper Node can be imported using the id, e.g.

```shell
$ terraform import alicloud_eflo_hyper_node.example <id>
```