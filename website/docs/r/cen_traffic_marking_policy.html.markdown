---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_traffic_marking_policy"
sidebar_current: "docs-alicloud-resource-cen-traffic-marking-policy"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Traffic Marking Policy resource.
---

# alicloud\_cen\_traffic\_marking\_policy

Provides a Cloud Enterprise Network (CEN) Traffic Marking Policy resource.

For information about Cloud Enterprise Network (CEN) Traffic Marking Policy and how to use it, see [What is Traffic Marking Policy](https://help.aliyun.com/document_detail/419025.html).

-> **NOTE:** Available in v1.173.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cen_instance" "example" {
  cen_instance_name = "example_value"
}

resource "alicloud_cen_transit_router" "example" {
  cen_id              = alicloud_cen_instance.example.id
  transit_router_name = "example_value"
}

resource "alicloud_cen_traffic_marking_policy" "example" {
  marking_dscp                = 1
  priority                    = 1
  traffic_marking_policy_name = "example_value"
  transit_router_id           = alicloud_cen_transit_router.example.transit_router_id
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of the Traffic Marking Policy. The description must be 2 to 128 characters in length, and must start with a letter. It can contain digits, underscores (_), and hyphens (-).
* `dry_run` - (Optional) The dry run.
* `marking_dscp` - (Required, ForceNew) The DSCP(Differentiated Services Code Point) of the Traffic Marking Policy. Value range: 0~63.
* `priority` - (Required, ForceNew) The Priority of the Traffic Marking Policy. Value range: 1~100.
* `traffic_marking_policy_name` - (Optional) The name of the Traffic Marking Policy. The name must be 2 to 128 characters in length, and must start with a letter. It can contain digits, underscores (_), and hyphens (-).
* `transit_router_id` - (Required, ForceNew) The ID of the transit router.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource, It is formatted to `<transit_router_id>:<traffic_marking_policy_id>`.
* `status` - The status of the resource.
* `traffic_marking_policy_id` - The ID of the Traffic Marking Policy.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 6 mins) Used when create the Traffic Marking Policy.
* `delete` - (Defaults to 6 mins) Used when delete the Traffic Marking Policy.
* `update` - (Defaults to 6 mins) Used when update the Traffic Marking Policy.

## Import

Cloud Enterprise Network (CEN) Traffic Marking Policy can be imported using the id, e.g.

```
$ terraform import alicloud_cen_traffic_marking_policy.example <transit_router_id>:<traffic_marking_policy_id>
```