---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_inter_region_traffic_qos_policy"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Inter Region Traffic Qos Policy resource.
---

# alicloud_cen_inter_region_traffic_qos_policy

Provides a Cloud Enterprise Network (CEN) Inter Region Traffic Qos Policy resource.



For information about Cloud Enterprise Network (CEN) Inter Region Traffic Qos Policy and how to use it, see [What is Inter Region Traffic Qos Policy](https://next.api.alibabacloud.com/document/Cbn/2017-09-12/CreateCenInterRegionTrafficQosPolicy).

-> **NOTE:** Available since v1.246.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_inter_region_traffic_qos_policy&exampleId=eaad5f98-e1f5-6208-1a44-6b0780caa008475b1cc1&activeTab=example&spm=docs.r.cen_inter_region_traffic_qos_policy.0.eaad5f98e1&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_cen_instance" "defaultpSZB78" {
}

resource "alicloud_cen_transit_router" "defaultUmmxnE" {
  cen_id = alicloud_cen_instance.defaultpSZB78.id
}

resource "alicloud_cen_transit_router" "defaultksqgSa" {
  cen_id = alicloud_cen_instance.defaultpSZB78.id
}

resource "alicloud_cen_transit_router_peer_attachment" "defaultnXZ83y" {
  default_link_type             = "Platinum"
  bandwidth_type                = "DataTransfer"
  cen_id                        = alicloud_cen_instance.defaultpSZB78.id
  peer_transit_router_region_id = alicloud_cen_transit_router.defaultksqgSa.id
  transit_router_id             = alicloud_cen_transit_router.defaultUmmxnE.transit_router_id
  peer_transit_router_id        = alicloud_cen_transit_router.defaultksqgSa.transit_router_id
  bandwidth                     = "10"
}


resource "alicloud_cen_inter_region_traffic_qos_policy" "default" {
  transit_router_attachment_id                = alicloud_cen_transit_router_peer_attachment.defaultnXZ83y.id
  inter_region_traffic_qos_policy_name        = "example1"
  inter_region_traffic_qos_policy_description = "example1"
  bandwidth_guarantee_mode                    = "byBandwidthPercent"
  transit_router_id                           = alicloud_cen_transit_router_peer_attachment.defaultnXZ83y.id
}
```

## Argument Reference

The following arguments are supported:
* `bandwidth_guarantee_mode` - (Optional, ForceNew) Bandwidth guarantee mode. You can select by bandwidth or by bandwidth percentage. The default is by percentage.
* `inter_region_traffic_qos_policy_description` - (Optional) The description information of the traffic scheduling policy.
* `inter_region_traffic_qos_policy_name` - (Optional) The name of the traffic scheduling policy.
* `transit_router_attachment_id` - (Required, ForceNew) Peer Attachment ID.
* `transit_router_id` - (Required, ForceNew) The ID of the forwarding router instance.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - The status of the traffic scheduling policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Inter Region Traffic Qos Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Inter Region Traffic Qos Policy.
* `update` - (Defaults to 5 mins) Used when update the Inter Region Traffic Qos Policy.

## Import

Cloud Enterprise Network (CEN) Inter Region Traffic Qos Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_inter_region_traffic_qos_policy.example <id>
```