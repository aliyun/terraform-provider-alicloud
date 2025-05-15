---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_traffic_marking_policy"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Traffic Marking Policy resource.
---

# alicloud_cen_traffic_marking_policy

Provides a Cloud Enterprise Network (CEN) Traffic Marking Policy resource.



For information about Cloud Enterprise Network (CEN) Traffic Marking Policy and how to use it, see [What is Traffic Marking Policy](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-createtrafficmarkingpolicy).

-> **NOTE:** Available since v1.173.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_traffic_marking_policy&exampleId=be3b43ec-5051-f623-1c57-7ba6b0cb90ba131caf90&activeTab=example&spm=docs.r.cen_traffic_marking_policy.0.be3b43ec50&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_cen_instance" "example" {
  cen_instance_name = "tf_example"
  description       = "an example for cen"
}

resource "alicloud_cen_transit_router" "example" {
  transit_router_name = "tf_example"
  cen_id              = alicloud_cen_instance.example.id
}

resource "alicloud_cen_traffic_marking_policy" "example" {
  marking_dscp                = 1
  priority                    = 1
  traffic_marking_policy_name = "tf_example"
  transit_router_id           = alicloud_cen_transit_router.example.transit_router_id
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) TrafficMarkingPolicyDescription
* `dry_run` - (Optional) Whether to PreCheck only this request. Value:
  - `true`: The check request is sent without creating a traffic marking policy. Check items include whether required parameters, request format, and business restrictions are filled in. If the check does not pass, the corresponding error is returned. If the check passes, the error code 'DryRunOperation' is returned '.
  - `false` (default): A normal request is sent, and a traffic marking policy is directly created after the check is passed.
* `force` - (Optional, Available since v1.248.0) Whether to forcibly delete the traffic marker policy. Valid values:
  - `false` (default value): Before deleting a traffic marking policy, check whether a traffic classification rule exists under the policy. If a flow classification rule exists, deletion is not allowed and the corresponding error is returned.
  - `true`: When you delete a traffic marking policy, all traffic classification rules under the policy are deleted by default.
* `marking_dscp` - (Required, ForceNew, Int) MarkingDscp
* `priority` - (Required, ForceNew, Int) Priority
* `traffic_marking_policy_name` - (Optional) TrafficMarkingPolicyName
* `traffic_match_rules` - (Optional, Set, Available since v1.230.1) List of stream classification rules.
You can add up to 50 stream classification rules at a time. See [`traffic_match_rules`](#traffic_match_rules) below.
* `transit_router_id` - (Required, ForceNew) TransitRouterId

### `traffic_match_rules`

The traffic_match_rules supports the following:
* `address_family` - (Optional, Computed) IP Address Family
* `dst_cidr` - (Optional, Computed) The destination network segment of the traffic message.
The flow classification matches the traffic of the destination IP address in the destination network segment. If the flow classification rule is not set, it means that the flow classification rule matches the traffic of any destination IP address.
* `dst_port_range` - (Optional, Computed, List) The destination port of the traffic message. Valid values: **-1**, `1` to `65535`.

The flow classification rule matches the traffic of the destination port number in the destination port range. If the flow classification rule is not set, it means that the flow classification rule matches the traffic of any destination port number.

The current parameter supports a maximum of 2 port numbers. The input format is described as follows:
  - If you only enter a port number, such as 1, the system defaults to match the traffic with the destination port of 1.
  - If you enter 2 port numbers, such as 1 and 200, the system defaults to match the traffic of the destination port in the range of 1 to 200.
  - If you enter 2 port numbers and one of them is - 1, the other port must also be - 1, indicating that it matches any destination port.
* `match_dscp` - (Optional, Computed, Int) The DSCP value of the traffic message. Valid values: `0` to **63 * *.

The flow classification rule matches the flow with the specified DSCP value. If the flow classification rule is not set, it means that the flow classification rule matches the flow with any DSCP value.

-> **NOTE:**  The current DSCP value refers to the DSCP value that the traffic message has carried before entering the cross-region connection.

* `protocol` - (Optional, Computed) The protocol type of the traffic message.
Stream classification rules can match traffic of multiple protocol types, such as `HTTP`, `HTTPS`, `TCP`, `UDP`, `SSH`, and **Telnet. For more protocol types, please log on to the [Cloud Enterprise Network Management Console](https://cen.console.aliyun.com/cen/list) to view.
* `src_cidr` - (Optional, Computed) The source network segment of the traffic message.
The flow classification rule matches the traffic of the source IP address in the source network segment. If the flow classification rule is not set, it means that the flow classification rule matches the traffic of any source IP address.
* `src_port_range` - (Optional, Computed, List) The source port of the traffic message. Valid values: **-1**, `1` to `65535`.

The flow classification rule matches the traffic of the source port number in the source port range. If it is not set, it means that the flow classification rule matches the traffic of any source port number.

The current parameter supports entering up to two port numbers. The input format is described as follows:
  - If you only enter a port number, such as 1, the system defaults to match the traffic with source port 1.
  - If you enter two port numbers, such as 1 and 200, the system defaults to match the traffic with the source port in the range of 1 to 200.
  - If you enter two port numbers and one of them is - 1, the other port must also be - 1, indicating that it matches any source port.
* `traffic_match_rule_description` - (Optional, Computed) The description information of the stream classification rule.
The description must be 2 to 128 characters in length and can contain numbers, dashes (-), and underscores (_).
* `traffic_match_rule_name` - (Optional, Computed) The name of the stream classification rule.
The name must be 2 to 128 characters in length and can contain numbers, dashes (-), and underscores (_).

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<transit_router_id>:<traffic_marking_policy_id>`.
* `status` - The status of the resource
* `traffic_marking_policy_id` - The first ID of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 8 mins) Used when create the Traffic Marking Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Traffic Marking Policy.
* `update` - (Defaults to 5 mins) Used when update the Traffic Marking Policy.

## Import

Cloud Enterprise Network (CEN) Traffic Marking Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_traffic_marking_policy.example <transit_router_id>:<traffic_marking_policy_id>
```