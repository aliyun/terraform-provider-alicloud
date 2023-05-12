---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_public_ip_address_pool"
sidebar_current: "docs-alicloud-resource-vpc-public-ip-address-pool"
description: |-
  Provides a Alicloud Vpc Public Ip Address Pool resource.
---

# alicloud_vpc_public_ip_address_pool

Provides a Vpc Public Ip Address Pool resource.

For information about Vpc Public Ip Address Pool and how to use it, see [What is Public Ip Address Pool](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/createpublicipaddresspool).

-> **NOTE:** Available in v1.186.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_resource_manager_resource_group" "defaultRg" {
  display_name        = "tf-test-acc-publicaddresspool-383"
  resource_group_name = "tf-test-acc-publicaddresspool-855"
}

resource "alicloud_resource_manager_resource_group" "changeRg" {
  display_name        = "tf-testacc-publicaddresspool-change-368"
  resource_group_name = "tf-testacc-publicaddresspool-change-499"
}


resource "alicloud_vpc_public_ip_address_pool" "default" {
  description                 = "rdk-test"
  public_ip_address_pool_name = "rdk-test"
  isp                         = "BGP"
  resource_group_id           = alicloud_resource_manager_resource_group.defaultRg.id
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) Description.
* `isp` - (ForceNew, Computed, Optional) The Internet service provider. Valid values: `BGP`, `BGP_PRO`, `ChinaTelecom`, `ChinaUnicom`, `ChinaMobile`, `ChinaTelecom_L2`, `ChinaUnicom_L2`, `ChinaMobile_L2`, `BGP_FinanceCloud`. Default Value: `BGP`.
* `public_ip_address_pool_name` - (Optional) The name of the VPC Public IP address pool.
* `resource_group_id` - (Computed, Optional) The resource group ID of the VPC Public IP address pool.
* `tags` - (Optional, Map) The tags of PrefixList.



## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `create_time` - The creation time of the resource.
* `ip_address_remaining` - Whether there is a free IP address.
* `public_ip_address_pool_id` - The resource ID in terraform of VPC Public Ip Address Pool.
* `status` - The status of the VPC Public IP address pool.
* `total_ip_num` - The total number of public IP address pools.
* `used_ip_num` - The number of used IP addresses in the public IP address pool.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Public Ip Address Pool.
* `delete` - (Defaults to 5 mins) Used when delete the Public Ip Address Pool.
* `update` - (Defaults to 5 mins) Used when update the Public Ip Address Pool.

## Import

Vpc Public Ip Address Pool can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_public_ip_address_pool.example <id>
```