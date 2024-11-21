---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_peer_attachment"
description: |-
  Provides a Alicloud CEN Transit Router Peer Attachment resource.
---

# alicloud_cen_transit_router_peer_attachment

Provides a CEN transit router peer attachment resource that associate the transit router with the CEN instance. [What is CEN transit router peer attachment](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-createtransitrouterpeerattachment)

-> **NOTE:** Available since v1.128.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_transit_router_peer_attachment&exampleId=c0d36f70-8856-a451-6a1e-ab5f03e359c0dd0321da&activeTab=example&spm=docs.r.cen_transit_router_peer_attachment.0.c0d36f7088&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}
variable "region" {
  default = "cn-hangzhou"
}
variable "peer_region" {
  default = "cn-beijing"
}
provider "alicloud" {
  alias  = "hz"
  region = var.region
}
provider "alicloud" {
  alias  = "bj"
  region = var.peer_region
}

resource "alicloud_cen_instance" "example" {
  provider          = alicloud.bj
  cen_instance_name = var.name
  protection_level  = "REDUCED"
}

resource "alicloud_cen_bandwidth_package" "example" {
  provider                   = alicloud.bj
  bandwidth                  = 5
  cen_bandwidth_package_name = "tf_example"
  geographic_region_a_id     = "China"
  geographic_region_b_id     = "China"
}

resource "alicloud_cen_bandwidth_package_attachment" "example" {
  provider             = alicloud.bj
  instance_id          = alicloud_cen_instance.example.id
  bandwidth_package_id = alicloud_cen_bandwidth_package.example.id
}

resource "alicloud_cen_transit_router" "example" {
  provider = alicloud.hz
  cen_id   = alicloud_cen_bandwidth_package_attachment.example.instance_id
}

resource "alicloud_cen_transit_router" "peer" {
  provider = alicloud.bj
  cen_id   = alicloud_cen_transit_router.example.cen_id
}

resource "alicloud_cen_transit_router_peer_attachment" "example" {
  provider                              = alicloud.hz
  cen_id                                = alicloud_cen_instance.example.id
  transit_router_id                     = alicloud_cen_transit_router.example.transit_router_id
  peer_transit_router_region_id         = var.peer_region
  peer_transit_router_id                = alicloud_cen_transit_router.peer.transit_router_id
  cen_bandwidth_package_id              = alicloud_cen_bandwidth_package_attachment.example.bandwidth_package_id
  bandwidth                             = 5
  transit_router_attachment_description = var.name
  transit_router_attachment_name        = var.name
}
```

## Argument Reference

The following arguments are supported:

* `dry_run` - (Optional) Whether to perform pre-check for this request, including permission, instance status verification, etc.
* `cen_id` - (Required, ForceNew) The ID of the CEN.
* `transit_router_id` - (Optional, ForceNew) The ID of the transit router to attach.
* `peer_transit_router_id` - (Required, ForceNew) The ID of the peer transit router.
* `peer_transit_router_region_id` - (Required, ForceNew) The region ID of peer transit router.
* `resource_type` - (Optional, ForceNew) The resource type to attachment. Only support `VR` and default value is `VR`.
* `cen_bandwidth_package_id` - (Optional) The ID of the bandwidth package. If you do not enter the ID of the package, it means you are using the test. The system default test is 1bps, demonstrating that you test network connectivity
* `bandwidth` - (Optional) The bandwidth of the bandwidth package.
* `auto_publish_route_enabled` - (Optional) Auto publish route enabled. The system default value is `false`.
* `transit_router_attachment_description` - (Optional) The description of transit router attachment. The description is 2~256 characters long and must start with a letter or Chinese, but cannot start with `http://` or `https://`.
* `transit_router_attachment_name` - (Optional) The name of transit router attachment. The name is 2~128 characters in length, starts with uppercase and lowercase letters or Chinese, and can contain numbers, underscores (_) and dashes (-)
* `bandwidth_type` - (Optional, Available since v1.157.0) The method that is used to allocate bandwidth to the cross-region connection. Valid values: `BandwidthPackage` and `DataTransfer`.
  * `DataTransfer` - uses pay-by-data-transfer bandwidth.
  * `BandwidthPackage` - allocates bandwidth from a bandwidth plan.
* `default_link_type` - (Optional, Available since v1.223.1) DefaultLinkType. Valid values: `Platinum` and `Gold`.
* `route_table_association_enabled` - (Deprecated since v1.230.0) Field `route_table_association_enabled` has been deprecated from provider version 1.230.0.
* `route_table_propagation_enabled` - (Deprecated since v1.230.0) Field `route_table_propagation_enabled` has been deprecated from provider version 1.230.0.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Transit Router Peer Attachment. It formats as `<cen_id>:<transit_router_attachment_id>`. 
* `transit_router_attachment_id` - The ID of transit router attachment.
* `status` - The status of the resource.
* `create_time` - The creation time of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Transit Router Peer Attachment.
* `update` - (Defaults to 5 mins) Used when update the Transit Router Peer Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Transit Router Peer Attachment.

## Import

CEN Transit Router Peer Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_peer_attachment.example <cen_id>:<transit_router_attachment_id>
```
