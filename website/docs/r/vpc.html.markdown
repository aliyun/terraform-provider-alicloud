---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc"
sidebar_current: "docs-alicloud-resource-vpc"
description: |-
  Provides a Alicloud VPC resource.
---

# alicloud\_vpc

Provides a VPC resource.

-> **NOTE:** Terraform will auto build a router and a route table while it uses `alicloud_vpc` to build a vpc resource.

## Example Usage

Basic Usage

```terraform
resource "alicloud_vpc" "vpc" {
  vpc_name   = "tf_test_foo"
  cidr_block = "172.16.0.0/12"
}
```

## Module Support

You can use the existing [vpc module](https://registry.terraform.io/modules/alibaba/vpc/alicloud) 
to create a VPC and several VSwitches one-click.

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Optional) The CIDR block for the VPC. The `cidr_block` is Optional and default value is `172.16.0.0/12` after v1.119.0+.
* `vpc_name` - (Optional, Available in v1.119.0+) The name of the VPC. Defaults to null.
* `name` - (Optional, Deprecated in v1.119.0+) Field `name` has been deprecated from provider version 1.119.0. New field `vpc_name` instead.
* `description` - (Optional) The VPC description. Defaults to null.
* `resource_group_id` - (Optional, Available in 1.40.0+, Modifiable in 1.115.0+) The Id of resource group which the VPC belongs.
* `tags` - (Optional, Available in v1.55.3+) A mapping of tags to assign to the resource.
* `secondary_cidr_blocks` - (Optional,Available in v1.112.0+) The secondary CIDR blocks for the VPC.
* `dry_run` - (Optional, ForceNew, Available in v1.119.0+) Specifies whether to precheck this request only. Valid values: `true` and `false`.
* `user_cidrs` - (Optional, ForceNew, Available in v1.119.0+) The user cidrs of the VPC.
* `enable_ipv6` - (Optional, Available in v1.119.0+) Specifies whether to enable the IPv6 CIDR block. Valid values: `false` (Default): disables IPv6 CIDR blocks. `true`: enables IPv6 CIDR blocks. If the `enable_ipv6` is `true`, the system will automatically create a free version of an IPv6 gateway for your private network and assign an IPv6 network segment assigned as /56.

-> **NOTE:** Currently, the IPv4 / IPv6 dual-stack VPC function is under public testing. Only the following regions support IPv4 / IPv6 dual-stack VPC: `cn-hangzhou`, `cn-shanghai`, `cn-shenzhen`, `cn-beijing`, `cn-huhehaote`, `cn-hongkong` and `ap-southeast-1`, and need to apply for public beta qualification. To use, please [submit an application](https://help.aliyun.com/document_detail/100334.html).

### Timeouts

-> **NOTE:** Available in 1.79.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the vpc (until it reaches the initial `Available` status). 
* `delete` - (Defaults to 10 mins) Used when terminating the vpc. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the VPC.
* `cidr_block` - The CIDR block for the VPC.
* `name` - The name of the VPC.
* `description` - The description of the VPC.
* `router_id` - The ID of the router created by default on VPC creation.
* `route_table_id` - The route table ID of the router created by default on VPC creation.
* `router_table_id` - (Deprecated) It has been deprecated and replaced with `route_table_id`.
* `ipv6_cidr_block` - (Available in v1.119.0+) ) The ipv6 cidr block of VPC.
* `status` - The status of the VPC.

## Import

VPC can be imported using the id, e.g.

```
$ terraform import alicloud_vpc.example vpc-abc123456
```

