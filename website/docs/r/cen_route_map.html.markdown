---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_route_map"
sidebar_current: "docs-alicloud-resource-cen-route-map"
description: |-
  Provides a Alicloud CEN manage route map resource.
---

# alicloud_cen_route_map

This topic provides an overview of the route map function of Cloud Enterprise Networks (CENs).
You can use the route map function to filter routes and modify route attributes.
By doing so, you can manage the communication between networks attached to a CEN. 

For information about CEN Route Map and how to use it, see [Manage CEN Route Map](https://www.alibabacloud.com/help/en/cloud-enterprise-network/latest/api-cbn-2017-09-12-createcenroutemap).

-> **NOTE:** Available since v1.82.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_route_map&exampleId=20f13b8e-3e2b-a139-a615-7a9170d9ee436d31d7b8&activeTab=example&spm=docs.r.cen_route_map.0.20f13b8e3e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "source_region" {
  default = "cn-hangzhou"
}
variable "destination_region" {
  default = "cn-shanghai"
}

provider "alicloud" {
  alias  = "hz"
  region = var.source_region
}
provider "alicloud" {
  alias  = "sh"
  region = var.destination_region
}

resource "alicloud_vpc" "example_hz" {
  provider   = alicloud.hz
  vpc_name   = "tf_example"
  cidr_block = "192.168.0.0/16"
}
resource "alicloud_vpc" "example_sh" {
  provider   = alicloud.sh
  vpc_name   = "tf_example"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_cen_instance" "example" {
  cen_instance_name = "tf_example"
  description       = "an example for cen"
}
resource "alicloud_cen_instance_attachment" "example_hz" {
  instance_id              = alicloud_cen_instance.example.id
  child_instance_id        = alicloud_vpc.example_hz.id
  child_instance_type      = "VPC"
  child_instance_region_id = var.source_region
}
resource "alicloud_cen_instance_attachment" "example_sh" {
  instance_id              = alicloud_cen_instance.example.id
  child_instance_id        = alicloud_vpc.example_sh.id
  child_instance_type      = "VPC"
  child_instance_region_id = var.destination_region
}

resource "alicloud_cen_route_map" "default" {
  cen_region_id                          = var.source_region
  cen_id                                 = alicloud_cen_instance.example.id
  description                            = "tf_example"
  priority                               = "1"
  transmit_direction                     = "RegionIn"
  map_result                             = "Permit"
  next_priority                          = "1"
  source_region_ids                      = [var.source_region]
  source_instance_ids                    = [alicloud_cen_instance_attachment.example_hz.child_instance_id]
  source_instance_ids_reverse_match      = "false"
  destination_instance_ids               = [alicloud_cen_instance_attachment.example_sh.child_instance_id]
  destination_instance_ids_reverse_match = "false"
  source_route_table_ids                 = [alicloud_vpc.example_hz.route_table_id]
  destination_route_table_ids            = [alicloud_vpc.example_sh.route_table_id]
  source_child_instance_types            = ["VPC"]
  destination_child_instance_types       = ["VPC"]
  destination_cidr_blocks                = [alicloud_vpc.example_sh.cidr_block]
  cidr_match_mode                        = "Include"
  route_types                            = ["System"]
  match_asns                             = ["65501"]
  as_path_match_mode                     = "Include"
  match_community_set                    = ["65501:1"]
  community_match_mode                   = "Include"
  community_operate_mode                 = "Additive"
  operate_community_set                  = ["65501:1"]
  preference                             = "20"
  prepend_as_path                        = ["65501"]
}
```
## Argument Reference

The following arguments are supported:

* `cen_id` - (Required, ForceNew) The ID of the CEN instance.
* `cen_region_id` - (Required) The ID of the region to which the CEN instance belongs.
* `priority` - (Required) The priority of the route map. Value range: 1 to 100. A lower value indicates a higher priority.
* `transmit_direction` - (Required, ForceNew) The direction in which the route map is applied. Valid values: ["RegionIn", "RegionOut"].
* `map_result` - (Required) The action that is performed to a route if the route matches all the match conditions. Valid values: ["Permit", "Deny"].
* `next_priority` - (Optional) The priority of the next route map that is associated with the current route map. Value range: 1 to 100.
* `description` - (Optional) The description of the route map.
* `source_region_ids` - (Optional) A match statement that indicates the list of IDs of the source regions. You can enter a maximum of 32 region IDs.
* `source_instance_ids` - (Optional) A match statement that indicates the list of IDs of the source instances. 
* `source_instance_ids_reverse_match` - (Optional) Indicates whether to enable the reverse match method for the SourceInstanceIds match condition. Valid values: ["false", "true"]. Default to "false".
* `destination_instance_ids` - (Optional) A match statement that indicates the list of IDs of the destination instances.
* `destination_instance_ids_reverse_match` - (Optional) Indicates whether to enable the reverse match method for the DestinationInstanceIds match condition. Valid values: ["false", "true"]. Default to "false".
* `source_route_table_ids` - (Optional) A match statement that indicates the list of IDs of the source route tables. You can enter a maximum of 32 route table IDs. 
* `destination_route_table_ids` - (Optional) A match statement that indicates the list of IDs of the destination route tables. You can enter a maximum of 32 route table IDs. 
* `source_child_instance_types` - (Optional) A match statement that indicates the list of source instance types. Valid values: ["VPC", "VBR", "CCN"].
* `destination_child_instance_types` - (Optional) A match statement that indicates the list of destination instance types. Valid values: ["VPC", "VBR", "CCN", "VPN"].
* `destination_cidr_blocks` - (Optional) A match statement that indicates the prefix list. The prefix is in the CIDR format. You can enter a maximum of 32 CIDR blocks. 
* `cidr_match_mode` - (Optional) A match statement. It indicates the mode in which the prefix attribute is matched. Valid values: ["Include", "Complete"].
* `route_types` - (Optional) A match statement that indicates the list of route types. Valid values: ["System", "Custom", "BGP"].
* `match_asns` - (Optional) A match statement that indicates the AS path list. The AS path is a well-known mandatory attribute, which describes the numbers of the ASs that a BGP route passes through during transmission. 
* `as_path_match_mode` - (Optional) A match statement. It indicates the mode in which the AS path attribute is matched. Valid values: ["Include", "Complete"].
* `match_community_set` - (Optional) A match statement that indicates the community set. The format of each community is nn:nn, which ranges from 1 to 65535. You can enter a maximum of 32 communities. Communities must comply with RFC 1997. Large communities (RFC 8092) are not supported. 
* `community_match_mode` - (Optional) A match statement. It indicates the mode in which the community attribute is matched. Valid values: ["Include", "Complete"].
* `community_operate_mode` - (Optional) An action statement. It indicates the mode in which the community attribute is operated. Valid values: ["Additive", "Replace"].
* `operate_community_set` - (Optional) An action statement that operates the community attribute. The format of each community is nn:nn, which ranges from 1 to 65535. You can enter a maximum of 32 communities. Communities must comply with RFC 1997. Large communities (RFC 8092) are not supported. 
* `preference` - (Optional) An action statement that modifies the priority of the route. Value range: 1 to 100. The default priority of a route is 50. A lower value indicates a higher preference. 
* `prepend_as_path` - (Optional) An action statement that indicates an AS path is prepended when the regional gateway receives or advertises a route.
* `transit_router_route_table_id` - (Optional, ForceNew, Computed, Available in v1.167.0+) The routing table ID of the forwarding router. If you do not enter the routing table ID, the routing policy is automatically associated with the default routing table of the forwarding router.

## Attributes Reference

The RouteMapId attributes are exported:

* `id` - ID of the RouteMap. It formats as `<cen_id>:<route_map_id>`
* `route_map_id` - ID of the RouteMap. It is available in 1.161.0+.
* `status` - (Computed) The status of route map. Valid values: ["Creating", "Active", "Deleting"].

## Import

CEN RouteMap can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_route_map.default <cen_id>:<route_map_id>.
```

