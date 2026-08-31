---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_vpc_firewall_cens"
sidebar_current: "docs-alicloud-datasource-cloud-firewall-vpc-firewall-cens"
description: |-
  Provides a list of Cloud Firewall Vpc Firewall Cens to the user.
---

# alicloud_cloud_firewall_vpc_firewall_cens

This data source provides Cloud Firewall Vpc Firewall Cen available to the user.[What is Vpc Firewall Cen](https://www.alibabacloud.com/help/en/cloud-firewall/latest/describevpcfirewallcenlist)

-> **NOTE:** Available since v1.194.0.

## Example Usage

```terraform
data "alicloud_cloud_firewall_vpc_firewall_cens" "default" {
  ids               = ["${alicloud_cloud_firewall_vpc_firewall_cen.default.id}"]
  cen_id            = "cen-cjok7uyb5w2b27573v"
  member_uid        = "1415189284827022"
  status            = "closed"
  vpc_firewall_name = "tf-test"
}

output "alicloud_cloud_firewall_vpc_firewall_cen_example_id" {
  value = data.alicloud_cloud_firewall_vpc_firewall_cens.default.cens.0.id
}
```

## Argument Reference

The following arguments are supported:
* `cen_id` - (ForceNew,Optional) The ID of the CEN instance.
* `lang` - (ForceNew,Optional) The language type of the requested and received messages. Value:-**zh** (default): Chinese.-**en**: English.
* `member_uid` - (ForceNew,Optional) The UID of the member account (other Alibaba Cloud account) of the current Alibaba cloud account.
* `network_instance_id` - (ForceNew,Optional) The ID of the VPC instance that created the VPC firewall.
* `status` - (ForceNew,Optional) Firewall switch status
  - `opened`: Enabled.
  - `closed`: closed.
  - `notconfigured`: indicates that the VPC boundary firewall is not configured. 
  - `configured`: indicates that the VPC boundary firewall is configured but not enabled.
* `vpc_firewall_id` - (ForceNew,Optional) VPC firewall ID
* `vpc_firewall_name` - (ForceNew,Optional) The name of the VPC firewall instance.
* `ids` - (Optional, ForceNew, Computed) A list of Vpc Firewall Cen IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Vpc Firewall Cen IDs.
* `cens` - A list of Vpc Firewall Cen Entries. Each element contains the following attributes:
  * `id` - The ID of the CEN instance.
  * `cen_id` - The ID of the CEN instance.
  * `connect_type` - Intercommunication type, value: `expressconnect`: Express Channel `cen`: Cloud Enterprise Network
  * `local_vpc` - The details of the VPC.
    * `defend_cidr_list` - The list of network segments protected by the VPC firewall.
    * `manual_v_switch_id` - The ID of the vSwitch specified when the routing mode is manual mode.
    * `network_instance_id` - The ID of the VPC instance that created the VPC firewall.
    * `network_instance_name` - The name of the network instance.
    * `network_instance_type` - The type of the network instance. Value: **VPC * *.
    * `owner_id` - The UID of the Alibaba Cloud account to which the VPC belongs.
    * `region_no` - The region ID of the VPC.
    * `route_mode` - Routing mode,. Value:-auto: indicates automatic mode.-manual: indicates manual mode.
    * `support_manual_mode` - Whether routing mode supports manual mode. Value:-**1**: Supported.-**0**: Not supported.
    * `transit_router_type` - The version of the cloud enterprise network forwarding router (CEN-TR). Value:-**Basic**: Basic Edition.-**Enterprise**: Enterprise Edition.
    * `vpc_cidr_table_list` - The VPC network segment list.
      * `route_entry_list` - The list of route entries in the VPC.
        * `destination_cidr` - The target network segment of the VPC.
        * `next_hop_instance_id` - The ID of the next hop instance in the VPC.
      * `route_table_id` - The ID of the route table of the VPC.
    * `vpc_id` - The ID of the VPC instance.
    * `vpc_name` - The instance name of the VPC.
  * `member_uid` - The UID of the member account (other Alibaba Cloud account) of the current Alibaba cloud account.
  * `status` - Firewall switch status
  * `vpc_firewall_id` - VPC firewall ID
  * `vpc_firewall_name` - The name of the VPC firewall instance.
