---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc"
description: |-
  Provides a Alicloud Vpc Vpc resource.
---

# alicloud_vpc

Provides a Vpc Vpc resource. A VPC instance creates a VPC. You can fully control your own VPC, such as selecting IP address ranges, configuring routing tables, and gateways. You can use Alibaba cloud resources such as cloud servers, apsaradb for RDS, and load balancer in your own VPC. 

-> **NOTE:** Available since v1.0.0.

-> **NOTE:** This resource will auto build a router and a route table while it uses `alicloud_vpc` to build a vpc resource. 

## Module Support

You can use the existing [vpc module](https://registry.terraform.io/modules/alibaba/vpc/alicloud) 
to create a VPC and several VSwitches one-click.

For information about Vpc Vpc and how to use it, see [What is Vpc](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/what-is-a-vpc).

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}


resource "alicloud_vpc" "default" {
  ipv6_isp    = "BGP"
  description = "test"
  cidr_block  = "10.0.0.0/8"
  vpc_name    = var.name
  enable_ipv6 = true
}
```


## Argument Reference

The following arguments are supported:
* `cidr_block` - (Optional, Computed) The CIDR block for the VPC. The `cidr_block` is Optional and default value is `172.16.0.0/12` after v1.119.0+.
* `classic_link_enabled` - (Optional) The status of ClassicLink function.
* `description` - (Optional) The VPC description. Defaults to null.
* `dry_run` - (Optional, Available since v1.119.0) Whether to PreCheck only this request. Value:
  - **true**: The check request is sent without creating a VPC. Check items include whether required parameters, request format, and business restrictions are filled in. If the check does not pass, the corresponding error is returned. If the check passes, the error code 'DryRunOperation' is returned '.
  - **false** (default): Sends a normal request, returns an HTTP 2xx status code and directly creates a VPC.
* `enable_ipv6` - (Optional, Available since v1.119.0) Whether to enable the IPv6 network segment. Value:
  - **false** (default): not enabled.
  - **true**: on.
* `ipv6_isp` - (Optional) The IPv6 address segment type of the VPC. Value:
  - **BGP** (default): Alibaba Cloud BGP IPv6.
  - **ChinaMobile**: China Mobile (single line).
  - **ChinaUnicom**: China Unicom (single line).
  - **ChinaTelecom**: China Telecom (single line).
-> **NOTE:**  If a single-line bandwidth whitelist is enabled, this field can be set to **ChinaTelecom** (China Telecom), **ChinaUnicom** (China Unicom), or **ChinaMobile** (China Mobile).
* `resource_group_id` - (Optional, Computed, Available since v1.115) The ID of the resource group to which the VPC belongs.
* `secondary_cidr_blocks` - (Optional, Computed, Deprecated since v1.185.0) Field 'secondary_cidr_blocks' has been deprecated from provider version 1.185.0 and it will be removed in the future version. Please use the new resource 'alicloud_vpc_ipv4_cidr_block'. `secondary_cidr_blocks` attributes and `alicloud_vpc_ipv4_cidr_block` resource cannot be used at the same time.
* `tags` - (Optional, Map, Available since v1.55.3) The tags of Vpc.
* `user_cidrs` - (Optional, ForceNew, Computed, Available since v1.119.0) A list of user CIDRs.
* `vpc_name` - (Optional, Available since v1.119.0) The name of the VPC. Defaults to null.

The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Deprecated since v1.119.0). Field 'name' has been deprecated from provider version 1.119.0. New field 'vpc_name' instead.
* `router_table_id` - (Deprecated since v1.206.0+) Field 'router_table_id' has been deprecated from provider version 1.206.0. New field 'route_table_id' instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the VPC.
* `ipv6_cidr_block` - The ipv6 cidr block of vpc.
* `ipv6_cidr_blocks` - The IPv6 CIDR block information of the VPC.
  * `ipv6_cidr_block` - The IPv6 CIDR block of the VPC.
  * `ipv6_isp` - Valid values: **BGP** (default): Alibaba Cloud BGP IPv6.
    - **ChinaMobile**: China Mobile (single line).
    - **ChinaUnicom**: China Unicom (single line).
    - **ChinaTelecom**: China Telecom (single line).
* `route_table_id` - The route table ID of the router created by default on VPC creation.
* `router_id` - The ID of the router created by default on VPC creation.
* `status` - The status of the VPC.   **Pending**: The VPC is being configured. **Available**: The VPC is available.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vpc.
* `delete` - (Defaults to 5 mins) Used when delete the Vpc.
* `update` - (Defaults to 5 mins) Used when update the Vpc.

## Import

Vpc Vpc can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc.example <id>
```