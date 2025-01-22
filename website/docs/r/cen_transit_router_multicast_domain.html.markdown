---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_multicast_domain"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Transit Router Multicast Domain resource.
---

# alicloud_cen_transit_router_multicast_domain

Provides a Cloud Enterprise Network (CEN) Transit Router Multicast Domain resource.



For information about Cloud Enterprise Network (CEN) Transit Router Multicast Domain and how to use it, see [What is Transit Router Multicast Domain](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-createtransitroutermulticastdomain).

-> **NOTE:** Available since v1.195.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_transit_router_multicast_domain&exampleId=8f7dfc4e-1030-97e2-1c03-b042623c82ae551ab1cf&activeTab=example&spm=docs.r.cen_transit_router_multicast_domain.0.8f7dfc4e10&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_cen_instance" "example" {
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "example" {
  transit_router_name = var.name
  cen_id              = alicloud_cen_instance.example.id
  support_multicast   = true
}

resource "alicloud_cen_transit_router_multicast_domain" "default" {
  transit_router_id                           = alicloud_cen_transit_router.example.transit_router_id
  transit_router_multicast_domain_name        = var.name
  transit_router_multicast_domain_description = var.name
  options {
    igmpv2_support = "disable"
  }
}
```

## Argument Reference

The following arguments are supported:
* `options` - (Optional, Set, Available since v1.242.0) The function options of the multicast domain. See [`options`](#options) below.
* `tags` - (Optional, Map) A mapping of tags to assign to the resource.
* `transit_router_id` - (Required, ForceNew) The ID of the forwarding router instance.
* `transit_router_multicast_domain_description` - (Optional) The description of the multicast domain.
* `transit_router_multicast_domain_name` - (Optional) The name of the multicast domain.

### `options`

The options supports the following:
* `igmpv2_support` - (Optional) Whether to enable IGMP function for multicast domain. Default value: `disable`. Valid values: `enable`, `disable`.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Transit Router Multicast Domain.
* `region_id` - (Available since v1.242.0) The region ID of the transit router.
* `status` - The status of the Transit Router Multicast Domain.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Transit Router Multicast Domain.
* `delete` - (Defaults to 5 mins) Used when delete the Transit Router Multicast Domain.
* `update` - (Defaults to 5 mins) Used when update the Transit Router Multicast Domain.

## Import

Cloud Enterprise Network (CEN) Transit Router Multicast Domain can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_multicast_domain.example <id>
```
