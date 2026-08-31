---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_nat_firewalls"
sidebar_current: "docs-alicloud-datasource-cloud-firewall-nat-firewalls"
description: |-
  Provides a list of Cloud Firewall Nat Firewall owned by an Alibaba Cloud account.
---

# alicloud_cloud_firewall_nat_firewalls

This data source provides Cloud Firewall Nat Firewall available to the user.[What is Nat Firewall](https://next.api.alibabacloud.com/document/Cloudfw/2017-12-07/CreateSecurityProxy)

-> **NOTE:** Available since v1.243.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_vpc" "defaultikZ0gD" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "defaultp4O7qi" {
  vpc_id       = alicloud_vpc.defaultikZ0gD.id
  cidr_block   = "172.16.6.0/24"
  vswitch_name = var.name
  zone_id      = "cn-shenzhen-e"
}

resource "alicloud_nat_gateway" "default2iRZpC" {
  description      = var.name
  nat_gateway_name = var.name
  eip_bind_mode    = "MULTI_BINDED"
  nat_type         = "Enhanced"
  vpc_id           = alicloud_vpc.defaultikZ0gD.id
  payment_type     = "PayAsYouGo"
  network_type     = "internet"
}

resource "alicloud_eip_address" "defaultyiRwgs" {
}

resource "alicloud_eip_association" "defaults2MTuO" {
  instance_id   = alicloud_nat_gateway.default2iRZpC.id
  allocation_id = alicloud_eip_address.defaultyiRwgs.allocation_id
  mode          = "NAT"
  instance_type = "NAT"
  vpc_id        = alicloud_nat_gateway.default2iRZpC.vpc_id
}

resource "alicloud_snat_entry" "defaultAKE43g" {
  snat_ip           = alicloud_eip_address.defaultyiRwgs.ip_address
  snat_table_id     = alicloud_nat_gateway.default2iRZpC.snat_table_ids[0]
  eip_affinity      = "1"
  source_vswitch_id = alicloud_vswitch.defaultp4O7qi.id
}


resource "alicloud_cloud_firewall_nat_firewall" "default" {
  region_no      = "cn-shenzhen"
  vswitch_auto   = "true"
  strict_mode    = "0"
  vpc_id         = alicloud_vpc.defaultikZ0gD.id
  proxy_name     = var.name
  lang           = "zh"
  nat_gateway_id = alicloud_nat_gateway.default2iRZpC.id
  nat_route_entry_list {
    nexthop_id       = alicloud_nat_gateway.default2iRZpC.id
    destination_cidr = "0.0.0.0/0"
    nexthop_type     = "NatGateway"
    route_table_id   = alicloud_vswitch.defaultp4O7qi.route_table_id
  }
  firewall_switch = "close"
  vswitch_cidr    = "172.16.5.0/24"
  status          = "closed"
  vswitch_id      = alicloud_vswitch.defaultp4O7qi.id
}

data "alicloud_cloud_firewall_nat_firewalls" "default" {
  ids            = ["${alicloud_cloud_firewall_nat_firewall.default.id}"]
  lang           = "zh"
  nat_gateway_id = alicloud_nat_gateway.default2iRZpC.id
  proxy_name     = var.name
  region_no      = "cn-shenzhen"
  status         = "closed"
  vpc_id         = alicloud_vpc.defaultikZ0gD.id
}

output "alicloud_cloud_firewall_nat_firewall_example_id" {
  value = data.alicloud_cloud_firewall_nat_firewalls.default.firewalls.0.id
}
```

## Argument Reference

The following arguments are supported:
* `lang` - (ForceNew, Optional) Lang
* `member_uid` - (ForceNew, Optional) Member Account ID
* `nat_gateway_id` - (ForceNew, Optional) NAT gateway ID
* `page_number` - (ForceNew, Optional) Page No
* `page_size` - (ForceNew, Optional) Page Size
* `proxy_id` - (ForceNew, Optional) NAT firewall ID
* `proxy_name` - (ForceNew, Optional) NAT firewall name
* `region_no` - (ForceNew, Optional) Region
* `status` - (ForceNew, Optional) The status of the resource
* `vpc_id` - (ForceNew, Optional) The ID of the VPC instance.
* `ids` - (Optional, ForceNew, Computed) A list of Nat Firewall IDs.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Nat Firewall IDs.
* `firewalls` - A list of Nat Firewall Entries. Each element contains the following attributes:
  * `ali_uid` - Alibaba Cloud account ID
  * `member_uid` - Member Account ID
  * `nat_gateway_id` - NAT gateway ID
  * `nat_gateway_name` - NAT Gateway name
  * `nat_route_entry_list` - The list of routes to be switched by the NAT gateway.
    * `destination_cidr` - The destination network segment of the default route.
    * `nexthop_id` - The next hop address of the original NAT gateway.
    * `nexthop_type` - The network type of the next hop. Value: NatGateway : NAT Gateway.
    * `route_table_id` - The route table where the default route of the NAT gateway is located.
  * `proxy_id` - NAT firewall ID
  * `proxy_name` - NAT firewall name
  * `strict_mode` - Whether strict mode is enabled1-Enable strict mode0-Disable strict mode
  * `vpc_id` - The ID of the VPC instance.
  * `id` - The ID of the resource supplied above.
