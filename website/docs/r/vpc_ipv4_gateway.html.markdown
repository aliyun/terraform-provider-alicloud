---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipv4_gateway"
sidebar_current: "docs-alicloud-resource-vpc-ipv4-gateway"
description: |-
  Provides a Alicloud Vpc Ipv4 Gateway resource.
---

# alicloud_vpc_ipv4_gateway

Provides a Vpc Ipv4 Gateway resource. 

For information about Vpc Ipv4 Gateway and how to use it, see [What is Ipv4 Gateway](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/createipv4gateway).

-> **NOTE:** Available since v1.181.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_ipv4_gateway&exampleId=94047563-35ae-6bfc-0c14-1599a9e9b84fc431a5d3&activeTab=example&spm=docs.r.vpc_ipv4_gateway.0.9404756335&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-testacc-example"
}

resource "alicloud_resource_manager_resource_group" "default" {
  display_name        = "tf-testAcc-rg665"
  resource_group_name = var.name
}

resource "alicloud_resource_manager_resource_group" "modify" {
  display_name        = "tf-testAcc-rg298"
  resource_group_name = "${var.name}1"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "${var.name}2"
  cidr_block = "10.0.0.0/8"
}


resource "alicloud_vpc_ipv4_gateway" "default" {
  ipv4_gateway_name        = var.name
  ipv4_gateway_description = "tf-testAcc-Ipv4Gateway"
  resource_group_id        = alicloud_resource_manager_resource_group.default.id
  vpc_id                   = alicloud_vpc.default.id
}
```

## Argument Reference

The following arguments are supported:
* `dry_run` - (Optional) Whether to PreCheck only this request. Value:-**true**: The check request is sent without creating an IPv4 Gateway. Check items include whether required parameters, request format, and business restrictions are filled in. If the check does not pass, the corresponding error is returned. If the check passes, the error code 'DryRunOperation' is returned '.-**false** (default): Sends a normal request, returns an HTTP 2xx status code and directly creates an IPv4 Gateway.
* `enabled` - (Optional, Computed, Available since v1.193.1) Whether the IPv4 gateway is active or not. Valid values are **true** and **false**.
* `ipv4_gateway_description` - (Optional) The description of the IPv4 gateway. The description must be 2 to 256 characters in length. It must start with a letter but cannot start with http:// or https://.
* `ipv4_gateway_name` - (Optional) The name of the IPv4 gateway. The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (_), and hyphens (-). It must start with a letter.
* `resource_group_id` - (Optional, Computed, Available since v1.205.0) The ID of the resource group to which the instance belongs.
* `tags` - (Optional, Map, Available since v1.205.0) The tags of the current resource.
* `vpc_id` - (Required, ForceNew) The ID of the virtual private cloud (VPC) where you want to create the IPv4 gateway. You can create only one IPv4 gateway in a VPC.



## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.
* `ipv4_gateway_id` - Resource primary key field.
* `ipv4_gateway_route_table_id` - ID of the route table associated with IPv4 Gateway.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Ipv4 Gateway.
* `delete` - (Defaults to 5 mins) Used when delete the Ipv4 Gateway.
* `update` - (Defaults to 5 mins) Used when update the Ipv4 Gateway.

## Import

Vpc Ipv4 Gateway can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_ipv4_gateway.example <id>
```