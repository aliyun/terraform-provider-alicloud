---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_traffic_mirror_filter_ingress_rule"
description: |-
  Provides a Alicloud VPC Traffic Mirror Filter Ingress Rule resource.
---

# alicloud_vpc_traffic_mirror_filter_ingress_rule

Provides a VPC Traffic Mirror Filter Ingress Rule resource. Traffic mirror entry rule.

For information about VPC Traffic Mirror Filter Ingress Rule and how to use it, see [What is Traffic Mirror Filter Ingress Rule](https://www.alibabacloud.com/help/doc-detail/261357.htm).

-> **NOTE:** Available since v1.141.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_traffic_mirror_filter_ingress_rule&exampleId=f8be62fc-aaf9-7f20-e348-7d0002369eccec005f8b&activeTab=example&spm=docs.r.vpc_traffic_mirror_filter_ingress_rule.0.f8be62fcaa&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_vpc_traffic_mirror_filter" "example" {
  traffic_mirror_filter_name = "example_value"
}

resource "alicloud_vpc_traffic_mirror_filter_ingress_rule" "example" {
  traffic_mirror_filter_id = alicloud_vpc_traffic_mirror_filter.example.id
  priority                 = "1"
  action                   = "accept"
  protocol                 = "UDP"
  destination_cidr_block   = "10.0.0.0/24"
  source_cidr_block        = "10.0.0.0/24"
  destination_port_range   = "1/120"
  source_port_range        = "1/120"
}
```

## Argument Reference

The following arguments are supported:
* `action` - (Optional, Available since v1.211.0) The collection policy of the inbound rule. Valid values: `accept` or `drop`. `accept`: collects network traffic. `drop`: does not collect network traffic.
* `destination_cidr_block` - (Required) The destination CIDR block of the inbound traffic.
* `destination_port_range` - (Optional, Computed) The destination CIDR block of the inbound traffic. Valid values: `1` to `65535`. Separate the first port and last port with a forward slash (/), for example, `1/200` or `80/80`. A value of `-1/-1` indicates that all ports are available. Therefore, do not set the value to `-1/-1`. **NOTE:** When `protocol` is `ICMP`, this parameter is invalid.
* `dry_run` - (Optional) Whether to PreCheck this request only. Value:
  - **true**: sends a check request and does not create inbound or outbound rules. Check items include whether required parameters are filled in, request format, and restrictions. If the check fails, the corresponding error is returned. If the check passes, the error code 'DryRunOperation' is returned '.
  - **false** (default): Sends a normal request and directly creates an inbound or outbound direction rule after checking.
* `priority` - (Required) The priority of the inbound rule. A smaller value indicates a higher priority. The maximum value is `10`, which indicates that you can configure at most 10 inbound rules for a filter.
* `protocol` - (Required) The transport protocol used by inbound traffic that needs to be mirrored. Valid values: `ALL`, `ICMP`, `TCP`, `UDP`.
* `source_cidr_block` - (Required) The source CIDR block of the inbound traffic.
* `source_port_range` - (Optional, Computed) The source port range of the inbound traffic. Valid values: `1` to `65535`. Separate the first port and last port with a forward slash (/), for example, `1/200` or `80/80`. A value of `-1/-1` indicates that all ports are available. Therefore, do not set the value to `-1/-1`. **NOTE:** When `protocol` is `ICMP`, this parameter is invalid.
* `traffic_mirror_filter_id` - (Required, ForceNew) The ID of the filter.

The following arguments will be discarded. Please use new fields as soon as possible:
* `rule_action` - (Deprecated since v1.211.0). Field 'rule_action' has been deprecated from provider version 1.211.0. New field 'action' instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<traffic_mirror_filter_id>:<traffic_mirror_filter_ingress_rule_id>`.
* `status` - The state of the inbound rule. `Creating`, `Created`, `Modifying` and `Deleting`.
* `traffic_mirror_filter_ingress_rule_id` - The ID of the outbound rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Traffic Mirror Filter Ingress Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Traffic Mirror Filter Ingress Rule.
* `update` - (Defaults to 5 mins) Used when update the Traffic Mirror Filter Ingress Rule.

## Import

VPC Traffic Mirror Filter Ingress Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_traffic_mirror_filter_ingress_rule.example <traffic_mirror_filter_id>:<traffic_mirror_filter_ingress_rule_id>
```