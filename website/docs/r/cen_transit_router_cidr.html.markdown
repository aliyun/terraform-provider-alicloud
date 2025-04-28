---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_cidr"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Transit Router Cidr resource.
---

# alicloud_cen_transit_router_cidr

Provides a Cloud Enterprise Network (CEN) Transit Router Cidr resource.

Used for Vpn Attachment, Connect Attachment, etc. Assign address segments.

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
* `cidr` - (Required) The new CIDR block of the transit router.
* `description` - (Optional) The new description of the transit router CIDR block.
The description must be 1 to 256 characters in length, and cannot start with http:// or https://. You can also leave this parameter empty.
* `publish_cidr_route` - (Optional) Specifies whether to allow the system to automatically add a route that points to the CIDR block to the route table of the transit router.

  - `true` (default)

    If you set the value to true, after you create a VPN attachment on a private VPN gateway and enable route learning for the VPN attachment, the system automatically adds the following route to the route table of the transit router that is in route learning relationship with the VPN attachment:

    A blackhole route whose destination CIDR block is the transit router CIDR block, which refers to the CIDR block from which gateway IP addresses are allocated to the IPsec-VPN connection. The blackhole route is advertised only to the route tables of virtual border routers (VBRs) connected to the transit router.

  - `false`
* `transit_router_cidr_name` - (Optional) The new name of the transit router CIDR block.
The name must be 1 to 128 characters in length, and cannot start with http:// or https://. You can also leave this parameter empty.
* `transit_router_id` - (Required, ForceNew) The ID of the transit router.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<transit_router_id>:<transit_router_cidr_id>`.
* `transit_router_cidr_id` - The ID of the CIDR block.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Transit Router Cidr.
* `delete` - (Defaults to 5 mins) Used when delete the Transit Router Cidr.
* `update` - (Defaults to 5 mins) Used when update the Transit Router Cidr.

## Import

Cloud Enterprise Network (CEN) Transit Router Cidr can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_cidr.example <transit_router_id>:<transit_router_cidr_id>
```