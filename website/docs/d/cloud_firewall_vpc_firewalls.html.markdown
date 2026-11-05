---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_vpc_firewalls"
sidebar_current: "docs-alicloud-datasource-cloud-firewall-vpc-firewalls"
description: |-
  Provides a list of Cloud Firewall Vpc Firewalls to the user.
---

# alicloud_cloud_firewall_vpc_firewalls

This data source provides Cloud Firewall Vpc Firewall available to the user.[What is Vpc Firewall](https://www.alibabacloud.com/help/doc-detail/342932.htm)

-> **NOTE:** Available since v1.194.0.

## Example Usage

```terraform
data "alicloud_cloud_firewall_vpc_firewalls" "default" {
  ids               = ["id-example"]
  member_uid        = "1415189284827022"
  status            = "closed"
  vpc_firewall_name = "tf-example"
}

output "alicloud_cfw_vpc_firewall_example_id" {
  value = data.alicloud_cloud_firewall_vpc_firewalls.default.firewalls.0.id
}
```

## Argument Reference

The following arguments are supported:
* `lang` - (ForceNew, Optional) The language type of the requested and received messages. Value:-**zh** (default): Chinese.-**en**: English.
* `member_uid` - (ForceNew, Optional) The UID of the Alibaba Cloud member account.
* `region_no` - (Optional) The region ID.
* `status` - (ForceNew, Optional) The status of the resource
* `vpc_firewall_id` - (ForceNew, Optional) The ID of the VPC firewall instance.
* `vpc_firewall_name` - (ForceNew, Optional) The name of the VPC firewall instance.
* `vpc_id` - (Optional) The ID of the VPC instance.
* `ids` - (Optional, ForceNew, Computed) A list of Vpc Firewall IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Vpc Firewall IDs.
* `names` - A list of name of Vpc Firewalls.
* `firewalls` - A list of Vpc Firewall Entries. Each element contains the following attributes:
  * `bandwidth` - Bandwidth specifications for high-speed channels. Unit: Mbps.
  * `connect_type` - The communication type of the VPC firewall. Valid value: **expressconnect**, which indicates Express Connect.
  * `local_vpc` - The details of the local VPC.
    * `eni_id` - The ID of the ENI.
    * `eni_private_ip_address` - The private IP address of the ENI.
    * `local_vpc_cidr_table_list` - The network segment list of the local VPC.
      * `local_route_entry_list` - The list of route entries of the local VPC.
        * `local_destination_cidr` - The target network segment of the local VPC.
        * `local_next_hop_instance_id` - The ID of the next-hop instance in the local VPC.
      * `local_route_table_id` - The ID of the route table of the local VPC.
    * `region_no` - The region ID of the local VPC.
    * `router_interface_id` - The ID of the router interface.
    * `vpc_id` - The ID of the local VPC instance.
    * `vpc_name` - The instance name of the local VPC.
  * `peer_vpc` - The details of the peer VPC.
    * `eni_id` - The ID of the ENI.
    * `eni_private_ip_address` - The private IP address of the ENI.
    * `peer_vpc_cidr_table_list` - The network segment list of the peer VPC.
      * `peer_route_entry_list` - Peer VPC route entry list information.
        * `peer_destination_cidr` - The target network segment of the peer VPC.
        * `peer_next_hop_instance_id` - The ID of the next-hop instance in the peer VPC.
      * `peer_route_table_id` - The ID of the route table of the peer VPC.
    * `region_no` - The region ID of the peer VPC.
    * `router_interface_id` - The ID of the router interface.
    * `vpc_id` - The ID of the peer VPC instance.
    * `vpc_name` - The instance name of the peer VPC.
  * `region_status` - The region is open. Value:-**enable**: is enabled, indicating that VPC firewall can be configured in this region.-**disable**: indicates that VPC firewall cannot be configured in this region.
  * `status` - The status of the resource
  * `vpc_firewall_id` - The ID of the VPC firewall instance.
  * `vpc_firewall_name` - The name of the VPC firewall instance.
  * `id` - The name of the VPC firewall instance and the value same as `vpc_firewall_id`.
