---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_traffic_marking_policy"
sidebar_current: "docs-alicloud-resource-cen-traffic-marking-policy"
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
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_cen_traffic_marking_policy&exampleId=be3b43ec-5051-f623-1c57-7ba6b0cb90ba131caf90&activeTab=example&spm=docs.r.cen_traffic_marking_policy.0.be3b43ec50&intl_lang=EN_US" target="_blank">
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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 6 mins) Used when create the Traffic Marking Policy.
* `delete` - (Defaults to 6 mins) Used when delete the Traffic Marking Policy.
* `update` - (Defaults to 6 mins) Used when update the Traffic Marking Policy.

## Import

Cloud Enterprise Network (CEN) Traffic Marking Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_traffic_marking_policy.example <transit_router_id>:<traffic_marking_policy_id>
```