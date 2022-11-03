---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_public_ip_address_pool"
sidebar_current: "docs-alicloud-resource-vpc-public-ip-address-pool"
description: |-
  Provides a Alicloud VPC Public Ip Address Pool resource.
---

# alicloud\_vpc\_public\_ip\_address\_pool

Provides a VPC Public Ip Address Pool resource.

For information about VPC Public Ip Address Pool and how to use it, see [What is Public Ip Address Pool](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/createpublicipaddresspool).

-> **NOTE:** Available in v1.186.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_vpc_public_ip_address_pool" "default" {
  public_ip_address_pool_name = "example_value"
  isp                         = "BGP_PRO"
  description                 = "example_value"
}
```

## Argument Reference

The following arguments are supported:

* `public_ip_address_pool_name` - (Optional) The name of the VPC Public IP address pool.
* `isp` - (Optional, ForceNew, Computed) The Internet service provider. Valid values: `BGP`, `BGP_PRO`, `ChinaTelecom`, `ChinaUnicom`, `ChinaMobile`, `ChinaTelecom_L2`, `ChinaUnicom_L2`, `ChinaMobile_L2`, `BGP_FinanceCloud`. Default Value: `BGP`.
* `description` - (Optional) The description of the VPC Public IP address pool.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of VPC Public Ip Address Pool.
* `status` - The status of the VPC Public IP address pool.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the VPC Public Ip Address Pool.
* `update` - (Defaults to 3 mins) Used when update the VPC Public Ip Address Pool.
* `delete` - (Defaults to 3 mins) Used when delete the VPC Public Ip Address Pool.

## Import

VPC Public Ip Address Pool can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_public_ip_address_pool.example <id>
```