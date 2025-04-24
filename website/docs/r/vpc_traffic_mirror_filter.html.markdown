---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_traffic_mirror_filter"
sidebar_current: "docs-alicloud-resource-vpc-traffic-mirror-filter"
description: |-
  Provides a Alicloud VPC Traffic Mirror Filter resource.
---

# alicloud_vpc_traffic_mirror_filter

Provides a VPC Traffic Mirror Filter resource. Traffic mirror filter criteria.

For information about VPC Traffic Mirror Filter and how to use it, see [What is Traffic Mirror Filter](https://www.alibabacloud.com/help/doc-detail/207513.htm).

-> **NOTE:** Available since v1.140.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_traffic_mirror_filter&exampleId=d7d035bb-3c44-24ae-8790-5889a1a539a5b44c08a6&activeTab=example&spm=docs.r.vpc_traffic_mirror_filter.0.d7d035bb3c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_resource_manager_resource_group" "default3iXhoa" {
  display_name        = "testname03"
  resource_group_name = var.name
}

resource "alicloud_resource_manager_resource_group" "defaultdNz2qk" {
  display_name        = "testname04"
  resource_group_name = "${var.name}1"
}


resource "alicloud_vpc_traffic_mirror_filter" "default" {
  traffic_mirror_filter_description = "test"
  traffic_mirror_filter_name        = var.name
  resource_group_id                 = alicloud_resource_manager_resource_group.default3iXhoa.id
  egress_rules {
    priority               = 1
    protocol               = "TCP"
    action                 = "accept"
    destination_cidr_block = "32.0.0.0/4"
    destination_port_range = "80/80"
    source_cidr_block      = "16.0.0.0/4"
    source_port_range      = "80/80"
  }
  ingress_rules {
    priority               = 1
    protocol               = "TCP"
    action                 = "accept"
    destination_cidr_block = "10.64.0.0/10"
    destination_port_range = "80/80"
    source_cidr_block      = "10.0.0.0/8"
    source_port_range      = "80/80"
  }
}
```


## Argument Reference

The following arguments are supported:
* `dry_run` - (Optional) Whether to PreCheck only this request. Value:
  - **true**: The check request is sent without creating traffic Image filter conditions. Check items include whether required parameters, request format, and business restrictions are filled in. If the check does not pass, the corresponding error is returned. If the check passes, the error code 'DryRunOperation' is returned '.
  - **false** (default): Sends a normal request, returns a 2xx HTTP status code after passing the check, and directly creates a filter condition.
* `egress_rules` - (Optional, ForceNew, Computed, Available since v1.206.0+) Information about the outbound rule. See the following `Block EgressRules`.
* `ingress_rules` - (Optional, ForceNew, Computed, Available since v1.206.0+) Inward direction rule information. See the following `Block IngressRules`.
* `resource_group_id` - (Optional, Computed, Available since v1.206.0+) The ID of the resource group to which the VPC belongs.
* `tags` - (Optional, Map, Available since v1.206.0+) The tags of this resource.
* `traffic_mirror_filter_description` - (Optional) The description of the TrafficMirrorFilter.
* `traffic_mirror_filter_name` - (Optional) The name of the TrafficMirrorFilter.


#### Block EgressRules

The EgressRules supports the following:
* `action` - (Required, ForceNew) Collection strategy for outbound rules. Value:
  - accept: collects network traffic.
  - drop: No network traffic is collected.
* `destination_cidr_block` - (Optional, ForceNew) DestinationCidrBlock.
* `destination_port_range` - (Optional, ForceNew) The destination port range of the outbound rule network traffic. The port range is 1 to 65535. Use a forward slash (/) to separate the start port and the end Port. The format is 1/200 and 80/80. Among them, - 1/-1 cannot be set separately, which means that the port is not limited.
-> **NOTE:**  When egresrules. N.Protocol is set to ALL or ICMP, this parameter does not need to be configured, indicating that the port is not restricted.
* `priority` - (Optional, ForceNew) Priority.
* `protocol` - (Required, ForceNew) The type of protocol used by the outbound network traffic to be mirrored. Value:
  - ALL: ALL agreements.
  - ICMP: Network Control Message Protocol.
  - TCP: Transmission Control Protocol.
  - UDP: User Datagram Protocol.
* `source_cidr_block` - (Optional, ForceNew) The source address of the outbound rule network traffic.
* `source_port_range` - (Optional, ForceNew) The source port range of the outbound rule network traffic. The port range is 1 to 65535. Use a forward slash (/) to separate the start port and the end Port. The format is 1/200 and 80/80. Among them, - 1/-1 cannot be set separately, which means that the port is not limited.
-> **NOTE:**  When egresrules. N.Protocol is set to ALL or ICMP, this parameter does not need to be configured, indicating that the port is not restricted.

#### Block IngressRules

The IngressRules supports the following:
* `action` - (Required, ForceNew) Collection strategy for outbound rules. Value:
  - accept: collects network traffic.
  - drop: No network traffic is collected.
* `destination_cidr_block` - (Optional, ForceNew) The destination address of the outbound rule network traffic.
* `destination_port_range` - (Optional, ForceNew) The destination port range of the outbound rule network traffic. The port range is 1 to 65535. Use a forward slash (/) to separate the start port and the end Port. The format is 1/200 and 80/80. Among them, - 1/-1 cannot be set separately, which means that the port is not limited.
-> **NOTE:**  When egresrules. N.Protocol is set to ALL or ICMP, this parameter does not need to be configured, indicating that the port is not restricted.
* `priority` - (Optional, ForceNew) The priority of the outbound rule. The smaller the number, the higher the priority. The maximum value of N is 10, that is, a maximum of 10 Outbound rules can be configured for a filter condition.
* `protocol` - (Required, ForceNew) The type of protocol used by the outbound network traffic to be mirrored. Value:
  - ALL: ALL agreements.
  - ICMP: Network Control Message Protocol.
  - TCP: Transmission Control Protocol.
  - UDP: User Datagram Protocol.
* `source_cidr_block` - (Optional, ForceNew) The source address of the outbound rule network traffic.
* `source_port_range` - (Optional, ForceNew) The source port range of the outbound rule network traffic. The port range is 1 to 65535. Use a forward slash (/) to separate the start port and the end Port. The format is 1/200 and 80/80. Among them, - 1/-1 cannot be set separately, which means that the port is not limited.
-> **NOTE:**  When egresrules. N.Protocol is set to ALL or ICMP, this parameter does not need to be configured, indicating that the port is not restricted.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Traffic Mirror Filter.
* `delete` - (Defaults to 5 mins) Used when delete the Traffic Mirror Filter.
* `update` - (Defaults to 5 mins) Used when update the Traffic Mirror Filter.

## Import

VPC Traffic Mirror Filter can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_traffic_mirror_filter.example <id>
```