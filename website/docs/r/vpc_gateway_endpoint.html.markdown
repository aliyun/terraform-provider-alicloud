---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_gateway_endpoint"
description: |-
  Provides a Alicloud VPC Gateway Endpoint resource.
---

# alicloud_vpc_gateway_endpoint

Provides a VPC Gateway Endpoint resource. VPC gateway endpoint.

For information about VPC Gateway Endpoint and how to use it, see [What is Gateway Endpoint](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/gateway-endpoint).

-> **NOTE:** This argument `route_tables` should not be used with the resource type `alicloud_vpc_gateway_endpoint_route_table_attachment` at the same time.

-> **NOTE:** Available since v1.208.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

variable "demoin" {
  default = "com.aliyun.cn-hangzhou.oss"
}

resource "alicloud_vpc" "defaultvpc" {
}

resource "alicloud_resource_manager_resource_group" "changeRg" {
  display_name        = "tf-exampleacc-848"
  resource_group_name = var.name
}

resource "alicloud_route_table" "defaultRt" {
  vpc_id      = alicloud_vpc.defaultvpc.id
  description = "tf-exampleacc"
}

resource "alicloud_route_table" "defaultrt1" {
  vpc_id      = alicloud_vpc.defaultvpc.id
  description = "tf-exampleacc1"
}

resource "alicloud_route_table" "defaultrt2" {
  vpc_id      = alicloud_vpc.defaultvpc.id
  description = "tf-exampleacc2"
}

resource "alicloud_vpc_gateway_endpoint" "default" {
  gateway_endpoint_descrption = "example-gateway-endpoint"
  gateway_endpoint_name       = var.name
  service_name                = var.demoin
  vpc_id                      = alicloud_vpc.defaultvpc.id
  policy_document             = "{ \"Version\" : \"1\", \"Statement\" : [ { \"Effect\" : \"Allow\", \"Resource\" : [ \"*\" ], \"Action\" : [ \"*\" ], \"Principal\" : [ \"*\" ] } ] }"
  resource_group_id           = alicloud_vpc.defaultvpc.resource_group_id
  route_tables                = ["${alicloud_route_table.defaultRt.id}", "${alicloud_route_table.defaultrt1.id}", "${alicloud_route_table.defaultrt2.id}"]
}
```

## Argument Reference

The following arguments are supported:
* `gateway_endpoint_descrption` - (Optional) The description of the gateway endpoint.
* `gateway_endpoint_name` - (Optional) The name of the gateway endpoint.
* `policy_document` - (Optional) Access control policies for cloud services. This parameter is required when the cloud service is oss. For details about the syntax and structure of access policies, see [syntax and structure of permission Policies](https://help.aliyun.com/document_detail/93739.html).
* `resource_group_id` - (Optional, Computed) The ID of the resource group to which the instance belongs.
* `route_tables` - (Optional, Computed) The route table id.
* `service_name` - (Required, ForceNew) The name of endpoint service.
* `tags` - (Optional, Map) The tags of the resource.
* `vpc_id` - (Required, ForceNew) The ID of the VPC.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the gateway endpoint.
* `status` - The status of VPC gateway endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Gateway Endpoint.
* `delete` - (Defaults to 5 mins) Used when delete the Gateway Endpoint.
* `update` - (Defaults to 5 mins) Used when update the Gateway Endpoint.

## Import

VPC Gateway Endpoint can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_gateway_endpoint.example <id>
```