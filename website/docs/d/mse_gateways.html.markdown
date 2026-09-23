---
subcategory: "Microservice Engine (MSE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_mse_gateways"
sidebar_current: "docs-alicloud-datasource-mse-gateways"
description: |-
  Provides a list of Mse Gateways to the user.
---

# alicloud\_mse\_gateways

This data source provides the Mse Gateways of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.157.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_mse_gateways" "ids" {
  ids = ["example_id"]
}
output "mse_gateway_id_1" {
  value = data.alicloud_mse_gateways.ids.gateways.0.id
}

data "alicloud_mse_gateways" "nameRegex" {
  name_regex = "^my-Gateway"
}
output "mse_gateway_id_2" {
  value = data.alicloud_mse_gateways.nameRegex.gateways.0.id
}

data "alicloud_mse_gateways" "status" {
  status = "2"
}
output "mse_gateway_id_3" {
  value = data.alicloud_mse_gateways.status.gateways.0.id
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `gateway_name` - (Optional, ForceNew) The name of the Gateway.
* `ids` - (Optional, ForceNew, Computed)  A list of Gateway IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Gateway name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the gateway. Valid values: `0`, `1`, `2`, `3`, `4`, `6`, `8`, `9`, `10`, `11`, `12`, `13`.
* `vpc_id` - (Optional, ForceNew) The ID of the VPC.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Gateway names.
* `gateways` - A list of Mse Gateways. Each element contains the following attributes:
  * `backup_vswitch_id` - The backup vswitch id.
  * `gateway_name` - The name of the Gateway.
  * `gateway_unique_id` - Gateway unique identification.
  * `id` - The ID of the Gateway.
  * `payment_type` - The payment type of the resource.
  * `replica` - Number of Gateway Nodes.
  * `spec` - Gateway Node Specifications.
  * `status` - The status of the gateway.
  * `vpc_id` - The ID of the vpc.
  * `vswitch_id` - The ID of the vswitch.
  * `slb_list` - A list of gateway Slb.
    * `associate_id` - The associate id.
    * `slb_id` - The ID of the gateway slb.
    * `slb_ip` - The ip of the gateway slb.
    * `slb_port` - The port of the gateway slb.
    * `type` - The type of the gateway slb.
    * `gmt_create` - The creation time of the gateway slb.
    * `gateway_slb_mode` - The Mode of the gateway slb.
    * `gateway_slb_status` - The Status of the gateway slb.