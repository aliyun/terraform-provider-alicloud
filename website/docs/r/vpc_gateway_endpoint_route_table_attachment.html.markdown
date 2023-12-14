---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_gateway_endpoint_route_table_attachment"
description: |-
  Provides a Alicloud VPC Gateway Endpoint Route Table Attachment resource.
---

# alicloud_vpc_gateway_endpoint_route_table_attachment

Provides a VPC Gateway Endpoint Route Table Attachment resource. VPC gateway node association route.

For information about VPC Gateway Endpoint Route Table Attachment and how to use it, see [What is Gateway Endpoint Route Table Attachment](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/311148).

-> **NOTE:** Available since v1.208.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "terraform-example"
}

resource "alicloud_vpc" "defaulteVpc" {
  description = "test"
}

resource "alicloud_vpc_gateway_endpoint" "defaultGE" {
  service_name                = "com.aliyun.cn-hangzhou.oss"
  policy_document             = <<EOF
        {
          "Version": "1",
          "Statement": [{
            "Effect": "Allow",
            "Resource": ["*"],
            "Action": ["*"],
            "Principal": ["*"]
          }]
        }
        EOF
  vpc_id                      = alicloud_vpc.defaulteVpc.id
  gateway_endpoint_descrption = "test-gateway-endpoint"
  gateway_endpoint_name       = "${var.name}1"
}

resource "alicloud_route_table" "defaultRT" {
  vpc_id           = alicloud_vpc.defaulteVpc.id
  route_table_name = "${var.name}2"
}


resource "alicloud_vpc_gateway_endpoint_route_table_attachment" "default" {
  gateway_endpoint_id = alicloud_vpc_gateway_endpoint.defaultGE.id
  route_table_id      = alicloud_route_table.defaultRT.id
}
```

## Argument Reference

The following arguments are supported:
* `gateway_endpoint_id` - (Required, ForceNew) The ID of the gateway endpoint instance to which you want to associate the route table.
* `route_table_id` - (Required, ForceNew) Routing table ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<gateway_endpoint_id>:<route_table_id>`.
* `status` - Status of the gateway endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Gateway Endpoint Route Table Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Gateway Endpoint Route Table Attachment.

## Import

VPC Gateway Endpoint Route Table Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_gateway_endpoint_route_table_attachment.example <gateway_endpoint_id>:<route_table_id>
```