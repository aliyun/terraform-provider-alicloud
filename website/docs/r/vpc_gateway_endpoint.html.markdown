---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_gateway_endpoint"
description: |-
  Provides a Alicloud VPC Gateway Endpoint resource.
---

# alicloud_vpc_gateway_endpoint

Provides a VPC Gateway Endpoint resource.

VPC gateway endpoint.

For information about VPC Gateway Endpoint and how to use it, see [What is Gateway Endpoint](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/gateway-endpoint).

-> **NOTE:** Available since v1.208.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_gateway_endpoint&exampleId=b42f095c-5d47-dac8-1f50-091d9074e35eaa2b33cf&activeTab=example&spm=docs.r.vpc_gateway_endpoint.0.b42f095c5d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_vpc_gateway_endpoint&spm=docs.r.vpc_gateway_endpoint.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `gateway_endpoint_descrption` - (Optional) The description of the VPC gateway endpoint.
The length of the description information is between 1 and 255 characters.
* `gateway_endpoint_name` - (Optional) The name of the VPC gateway endpoint.
* `policy_document` - (Optional, JsonString) Access control policies for cloud services. This parameter is required when the cloud service is oss. For details about the syntax and structure of access policies, see [syntax and structure of permission Policies](https://help.aliyun.com/document_detail/93739.html).
* `resource_group_id` - (Optional, Computed) The ID of the resource group to which the instance belongs.
* `route_tables` - (Optional, Computed, Set, Available since v1.244.0) The ID list of the route table associated with the VPC gateway endpoint. **NOTE:** this argument cannot be set at the same time as `alicloud_vpc_gateway_endpoint_route_table_attachment`.
* `service_name` - (Required, ForceNew) The endpoint service name.
* `tags` - (Optional, Map) The tags of the resource.
* `vpc_id` - (Required, ForceNew) The ID of the VPC.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the VPC gateway endpoint.
* `status` - The status of VPC gateway endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 8 mins) Used when create the Gateway Endpoint.
* `delete` - (Defaults to 10 mins) Used when delete the Gateway Endpoint.
* `update` - (Defaults to 5 mins) Used when update the Gateway Endpoint.

## Import

VPC Gateway Endpoint can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_gateway_endpoint.example <id>
```