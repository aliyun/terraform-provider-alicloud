---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_cidr"
sidebar_current: "docs-alicloud-resource-cen-transit-router-cidr"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Transit Router Cidr resource.
---

# alicloud_cen_transit_router_cidr

Provides a Cloud Enterprise Network (CEN) Transit Router Cidr resource.

For information about Cloud Enterprise Network (CEN) Transit Router Cidr and how to use it, see [What is Transit Router Cidr](https://www.alibabacloud.com/help/en/cloud-enterprise-network/latest/createtransitroutercidr).

-> **NOTE:** Available since v1.193.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_transit_router_cidr&exampleId=66fdd0bb-a861-8044-3efa-cdd9067d7d023b305102&activeTab=example&spm=docs.r.cen_transit_router_cidr.0.66fdd0bba8&intl_lang=EN_US" target="_blank">
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

resource "alicloud_cen_transit_router_cidr" "example" {
  transit_router_id        = alicloud_cen_transit_router.example.transit_router_id
  cidr                     = "192.168.0.0/16"
  transit_router_cidr_name = "tf_example"
  description              = "tf_example"
  publish_cidr_route       = true
}
```

## Argument Reference

The following arguments are supported:

* `transit_router_id` - (Required, ForceNew) The ID of the transit router.
* `cidr` - (Required) The cidr of the transit router.
* `transit_router_cidr_name` - (Optional) The name of the transit router. The name must be `2` to `128` characters in length, and can contain letters, digits, underscores (_), and hyphens (-). The name must start with a letter but cannot start with `http://` or `https://`.
* `description` - (Optional) The description of the transit router. The description must be `2` to `256` characters in length, and it must start with English letters, but cannot start with `http://` or `https://`.
* `publish_cidr_route` - (Optional, Computed) Whether to allow automatically adding Transit Router Cidr in Transit Router Route Table. Valid values: `true` and `false`. Default value: `true`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Transit Router Cidr. It formats as `<transit_router_id>:<transit_router_cidr_id>`
* `transit_router_cidr_id` - The ID of the transit router cidr.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Transit Router Cidr.
* `update` - (Defaults to 5 mins) Used when update the Transit Router Cidr.
* `delete` - (Defaults to 5 mins) Used when delete the Transit Router Cidr.

## Import

Cloud Enterprise Network (CEN) Transit Router Cidr can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_cidr.default <transit_router_id>:<transit_router_cidr_id>.
```
