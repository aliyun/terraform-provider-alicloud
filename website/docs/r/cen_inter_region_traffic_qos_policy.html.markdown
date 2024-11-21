---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_inter_region_traffic_qos_policy"
sidebar_current: "docs-alicloud-resource-cen-inter-region-traffic-qos-policy"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Inter Region Traffic Qos Policy resource.
---

# alicloud_cen_inter_region_traffic_qos_policy

Provides a Cloud Enterprise Network (CEN) Inter Region Traffic Qos Policy resource.

For information about Cloud Enterprise Network (CEN) Inter Region Traffic Qos Policy and how to use it, see [What is Inter Region Traffic Qos Policy](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-createceninterregiontrafficqospolicy).

-> **NOTE:** Available since v1.195.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_inter_region_traffic_qos_policy&exampleId=4ce5cc26-217d-b415-ebec-e128a69516037c5186cf&activeTab=example&spm=docs.r.cen_inter_region_traffic_qos_policy.0.4ce5cc2621&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  alias  = "bj"
  region = "cn-beijing"
}

provider "alicloud" {
  alias  = "hz"
  region = "cn-hangzhou"
}

resource "alicloud_cen_instance" "default" {
  provider          = alicloud.hz
  cen_instance_name = "tf-example"
}

resource "alicloud_cen_bandwidth_package" "default" {
  provider               = alicloud.hz
  bandwidth              = 5
  geographic_region_a_id = "China"
  geographic_region_b_id = "China"
}

resource "alicloud_cen_bandwidth_package_attachment" "default" {
  provider             = alicloud.hz
  instance_id          = alicloud_cen_instance.default.id
  bandwidth_package_id = alicloud_cen_bandwidth_package.default.id
}

resource "alicloud_cen_transit_router" "hz" {
  provider = alicloud.hz
  cen_id   = alicloud_cen_bandwidth_package_attachment.default.instance_id
}

resource "alicloud_cen_transit_router" "bj" {
  provider = alicloud.bj
  cen_id   = alicloud_cen_transit_router.hz.cen_id
}

resource "alicloud_cen_transit_router_peer_attachment" "default" {
  provider                      = alicloud.hz
  cen_id                        = alicloud_cen_instance.default.id
  transit_router_id             = alicloud_cen_transit_router.hz.transit_router_id
  peer_transit_router_region_id = "cn-beijing"
  peer_transit_router_id        = alicloud_cen_transit_router.bj.transit_router_id
  cen_bandwidth_package_id      = alicloud_cen_bandwidth_package_attachment.default.bandwidth_package_id
  bandwidth                     = 5
}

resource "alicloud_cen_inter_region_traffic_qos_policy" "default" {
  provider                                    = alicloud.hz
  transit_router_id                           = alicloud_cen_transit_router.hz.transit_router_id
  transit_router_attachment_id                = alicloud_cen_transit_router_peer_attachment.default.transit_router_attachment_id
  inter_region_traffic_qos_policy_name        = "tf-example-name"
  inter_region_traffic_qos_policy_description = "tf-example-description"
}
```

## Argument Reference

The following arguments are supported:

* `transit_router_id` - (Required, ForceNew) The ID of the transit router.
* `transit_router_attachment_id` - (Required, ForceNew) The ID of the inter-region connection.
* `inter_region_traffic_qos_policy_name` - (Optional) The name of the QoS policy. The name must be 2 to 128 characters in length, and can contain letters, digits, underscores (_), and hyphens (-). It must start with a letter.
* `inter_region_traffic_qos_policy_description` - (Optional) The description of the QoS policy. The description must be 2 to 128 characters in length, and can contain letters, digits, underscores (_), and hyphens (-). The description must start with a letter.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Inter Region Traffic Qos Policy.
* `status` - The status of the Inter Region Traffic Qos Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Inter Region Traffic Qos Policy.
* `update` - (Defaults to 3 mins) Used when create the Inter Region Traffic Qos Policy.
* `delete` - (Defaults to 3 mins) Used when delete the Inter Region Traffic Qos Policy.

## Import

Cloud Enterprise Network (CEN) Inter Region Traffic Qos Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_inter_region_traffic_qos_policy.example <id>
```
