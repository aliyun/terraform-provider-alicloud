---
subcategory: "GWLB"
layout: "alicloud"
page_title: "Alicloud: alicloud_gwlb_load_balancer"
description: |-
  Provides a Alicloud GWLB Load Balancer resource.
---

# alicloud_gwlb_load_balancer

Provides a GWLB Load Balancer resource.



For information about GWLB Load Balancer and how to use it, see [What is Load Balancer](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.234.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_gwlb_load_balancer&exampleId=048e6dcc-c045-2c92-6171-abb5ccd48e62d426865e&activeTab=example&spm=docs.r.gwlb_load_balancer.0.048e6dccc0&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-wulanchabu"
}

variable "region_id" {
  default = "cn-wulanchabu"
}

variable "zone_id2" {
  default = "cn-wulanchabu-c"
}

variable "zone_id1" {
  default = "cn-wulanchabu-b"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaulti9Axhl" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default9NaKmL" {
  vpc_id       = alicloud_vpc.defaulti9Axhl.id
  zone_id      = var.zone_id1
  cidr_block   = "10.0.0.0/24"
  vswitch_name = format("%s1", var.name)
}

resource "alicloud_vswitch" "defaultH4pKT4" {
  vpc_id       = alicloud_vpc.defaulti9Axhl.id
  zone_id      = var.zone_id2
  cidr_block   = "10.0.1.0/24"
  vswitch_name = format("%s2", var.name)
}


resource "alicloud_gwlb_load_balancer" "default" {
  vpc_id             = alicloud_vpc.defaulti9Axhl.id
  load_balancer_name = var.name
  zone_mappings {
    vswitch_id = alicloud_vswitch.default9NaKmL.id
    zone_id    = var.zone_id1
  }
  address_ip_version = "Ipv4"
}
```

## Argument Reference

The following arguments are supported:
* `address_ip_version` - (Optional, ForceNew, Computed) The IP version. Valid values:

  - `Ipv4`: IPv4 (default)
* `dry_run` - (Optional) Specifies whether to perform a dry run, without performing the actual request. Valid values:

  - `true`: performs only a dry run. The system checks the request for potential issues, including missing parameter values, incorrect request syntax, and service limits. If the request fails the dry run, an error code is returned. If the request passes the dry run, the `DryRunOperation` error code is returned.
  - `false` (default): performs a dry run and performs the actual request. If the request passes the dry run, a 2xx HTTP status code is returned and the operation is performed.
* `load_balancer_name` - (Optional) The GWLB instance name.

  The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (\_), and hyphens (-). The name must start with a letter.
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `tags` - (Optional, Map) The tag keys. You can specify at most 20 tags in each call.
* `vpc_id` - (Required, ForceNew) The virtual private cloud (VPC) ID.
* `zone_mappings` - (Required, Set) The mappings between zones and vSwitches. You must specify at least one zone. You can specify at most 20 zones. If the region supports two or more zones, we recommend that you select two or more zones. See [`zone_mappings`](#zone_mappings) below.

### `zone_mappings`

The zone_mappings supports the following:
* `vswitch_id` - (Required) The ID of the vSwitch in the zone. You can specify only one vSwitch (subnet) in each zone of a GWLB instance.
* `zone_id` - (Required) The zone ID. You can call the DescribeZones operation to query the most recent zone list.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the resource was created. The time follows the ISO 8601 standard in the **yyyy-MM-ddTHH:mm:ssZ** format. The time is displayed in UTC.
* `status` - The GWLB instance status. 

* `zone_mappings` - The mappings between zones and vSwitches. You must specify at least one zone. You can specify at most 20 zones. If the region supports two or more zones, we recommend that you select two or more zones.
  * `load_balancer_addresses` - The information about the IP addresses used by the GWLB instance.
    * `eni_id` - The ID of the elastic network interface (ENI) used by the GWLB instance.
    * `private_ipv4_address` - The private IPv4 address.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Load Balancer.
* `delete` - (Defaults to 5 mins) Used when delete the Load Balancer.
* `update` - (Defaults to 5 mins) Used when update the Load Balancer.

## Import

GWLB Load Balancer can be imported using the id, e.g.

```shell
$ terraform import alicloud_gwlb_load_balancer.example <id>
```