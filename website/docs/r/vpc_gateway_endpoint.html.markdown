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

variable "domain" {
  default = "com.aliyun.cn-hangzhou.oss"
}

resource "alicloud_vpc" "defaultVpc" {
  description = "tf-example"
}

resource "alicloud_resource_manager_resource_group" "defaultRg" {
  display_name        = "tf-example-497"
  resource_group_name = var.name
}

resource "alicloud_vpc_gateway_endpoint" "default" {
  gateway_endpoint_descrption = "test-gateway-endpoint"
  gateway_endpoint_name       = var.name
  vpc_id                      = alicloud_vpc.defaultVpc.id
  resource_group_id           = alicloud_resource_manager_resource_group.defaultRg.id
  service_name                = var.domain
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
}
```

## Argument Reference

The following arguments are supported:
* `gateway_endpoint_descrption` - (Optional) The description of the gateway endpoint.
* `gateway_endpoint_name` - (Optional) The name of the gateway endpoint.
* `policy_document` - (Optional) Access control policies for cloud services. This parameter is required when the cloud service is oss. For details about the syntax and structure of access policies, see [syntax and structure of permission Policies](https://help.aliyun.com/document_detail/93739.html).
* `resource_group_id` - (Optional, Computed) The ID of the resource group to which the instance belongs.
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