---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_public_ip_address_pool"
description: |-
  Provides a Alicloud VPC Public Ip Address Pool resource.
---

# alicloud_vpc_public_ip_address_pool

Provides a VPC Public Ip Address Pool resource.



For information about VPC Public Ip Address Pool and how to use it, see [What is Public Ip Address Pool](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/createpublicipaddresspool).

-> **NOTE:** Available since v1.186.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_public_ip_address_pool&exampleId=c2f8af64-2f4b-a68c-3b9f-bf200e9dc7007f8123fd&activeTab=example&spm=docs.r.vpc_public_ip_address_pool.0.c2f8af642f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}
resource "alicloud_vpc_public_ip_address_pool" "default" {
  description                 = var.name
  public_ip_address_pool_name = var.name
  isp                         = "BGP"
  resource_group_id           = data.alicloud_resource_manager_resource_groups.default.ids.0
}
```

## Argument Reference

The following arguments are supported:
* `biz_type` - (Optional, ForceNew, Computed) The name of the VPC Public IP address pool.
* `description` - (Optional) Description.
* `isp` - (Optional, ForceNew, Computed) The Internet service provider. Valid values: `BGP`, `BGP_PRO`, `ChinaTelecom`, `ChinaUnicom`, `ChinaMobile`, `ChinaTelecom_L2`, `ChinaUnicom_L2`, `ChinaMobile_L2`, `BGP_FinanceCloud`. Default Value: `BGP`.
* `public_ip_address_pool_name` - (Optional) The name of the VPC Public IP address pool.
* `resource_group_id` - (Optional, Computed) The resource group ID of the VPC Public IP address pool.
* `security_protection_types` - (Optional, ForceNew, Available since v1.227.1) Security protection level.
  - If the configuration is empty, the default value is DDoS protection (Basic edition).
  - `AntiDDoS_Enhanced` indicates DDoS protection (enhanced version).
* `tags` - (Optional, Map) The tags of PrefixList.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource
* `ip_address_remaining` - Whether there is a free IP address.
* `public_ip_address_pool_id` - The resource ID in terraform of VPC Public Ip Address Pool.
* `status` - The status of the VPC Public IP address pool.
* `total_ip_num` - The total number of public IP address pools.
* `used_ip_num` - The number of used IP addresses in the public IP address pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Public Ip Address Pool.
* `delete` - (Defaults to 5 mins) Used when delete the Public Ip Address Pool.
* `update` - (Defaults to 5 mins) Used when update the Public Ip Address Pool.

## Import

VPC Public Ip Address Pool can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_public_ip_address_pool.example <id>
```