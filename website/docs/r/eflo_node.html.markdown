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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_eflo_node&exampleId=5da29aac-f596-2411-6e2e-d5071f67231f9e9f9657&activeTab=example&spm=docs.r.eflo_node.0.5da29aacf5&intl_lang=EN_US" target="_blank">
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
```

## Argument Reference

The following arguments are supported:
* `billing_cycle` - (Optional) Billing cycle
* `classify` - (Optional) Classification
* `computing_server` - (Optional) Node Model
* `discount_level` - (Optional) Offer Information
* `hpn_zone` - (Optional) Cluster Number
* `payment_ratio` - (Optional) Down payment ratio
* `period` - (Optional, Int) Prepaid cycle. The unit is Month, please enter an integer multiple of 12 for the annual payment product.
* `product_form` - (Optional) Form
* `renew_period` - (Optional, Int) Automatic renewal period, in months.

-> **NOTE:**  When setting `RenewalStatus` to `AutoRenewal`, it must be set.

* `renewal_status` - (Optional) Automatic renewal status, value:
  - AutoRenewal: automatic renewal.
  - ManualRenewal: manual renewal.

The default ManualRenewal.
* `resource_group_id` - (Optional, Computed) The ID of the resource group
* `server_arch` - (Optional) Architecture
* `stage_num` - (Optional) Number of stages
* `status` - (Optional, Computed) The status of the resource
* `tags` - (Optional, Map) The tag of the resource
* `zone` - (Optional) Availability Zone

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Node.
* `delete` - (Defaults to 5 mins) Used when delete the Node.
* `update` - (Defaults to 5 mins) Used when update the Node.

## Import

Eflo Node can be imported using the id, e.g.

```shell
$ terraform import alicloud_eflo_node.example <id>
```