---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_traffic_mirror_session"
sidebar_current: "docs-alicloud-resource-vpc-traffic-mirror-session"
description: |-
  Provides a Alicloud VPC Traffic Mirror Session resource.
---

# alicloud\_vpc\_traffic\_mirror\_session

Provides a VPC Traffic Mirror Session resource.

For information about VPC Traffic Mirror Session and how to use it, see [What is Traffic Mirror Session](https://www.alibabacloud.com/help/en/doc-detail/261364.htm).

-> **NOTE:** Available in v1.142.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_vpc_traffic_mirror_filter" "example" {
  traffic_mirror_filter_name = "example_value"
}

resource "alicloud_vpc_traffic_mirror_session" "example" {
  priority                    = "1"
  traffic_mirror_filter_id    = alicloud_vpc_traffic_mirror_filter.example.traffic_mirror_filter_id
  traffic_mirror_session_name = "example_value"
  traffic_mirror_source_ids   = ["example_value"]
  traffic_mirror_target_id    = "example_value"
  traffic_mirror_target_type  = "NetworkInterface"
}

```

## Argument Reference

The following arguments are supported:

* `dry_run` - (Optional) The dry run.
* `enabled` - (Optional) Specifies whether to enable traffic mirror sessions. default to `false`.
* `priority` - (Required) The priority of the traffic mirror session. Valid values: `1` to `32766`. A smaller value indicates a higher priority. You cannot specify the same priority for traffic mirror sessions that are created in the same region with the same Alibaba Cloud account.
* `traffic_mirror_filter_id` - (Required) The ID of the filter.
* `traffic_mirror_session_description` - (Optional) The description of the traffic mirror session. The description must be `2` to `256` characters in length and cannot start with `http://` or `https://`.
* `traffic_mirror_session_name` - (Optional) The name of the traffic mirror session. The name must be `2` to `128` characters in length and can contain digits, underscores (_), and hyphens (-). It must start with a letter.
* `traffic_mirror_source_ids` - (Required) The ID of the mirror source. You can specify only an elastic network interface (ENI) as the mirror source. **NOTE:** Only one mirror source can be added to a traffic mirror session.
* `traffic_mirror_target_id` - (Required) The ID of the mirror destination. You can specify only an ENI or a Server Load Balancer (SLB) instance as a mirror destination.
* `traffic_mirror_target_type` - (Required) The type of the mirror destination. Valid values: `NetworkInterface` or `SLB`. `NetworkInterface`: an ENI. `SLB`: an internal-facing SLB instance
* `virtual_network_id` - (Optional) The VXLAN network identifier (VNI) that is used to distinguish different mirrored traffic. Valid values: `0` to `16777215`. You can specify VNIs for the traffic mirror destination to identify mirrored traffic from different sessions. If you do not specify a VNI, the system randomly allocates a VNI. If you want the system to randomly allocate a VNI, ignore this parameter.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Traffic Mirror Session.
* `status` - The state of the traffic mirror session. Valid values: `Creating`, `Created`, `Modifying` and `Deleting`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Traffic Mirror Session.
* `delete` - (Defaults to 1 mins) Used when delete the Traffic Mirror Session.
* `update` - (Defaults to 1 mins) Used when update the Traffic Mirror Session.

## Import

VPC Traffic Mirror Session can be imported using the id, e.g.

```
$ terraform import alicloud_vpc_traffic_mirror_session.example <id>
```