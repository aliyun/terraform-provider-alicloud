---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_vbr_ha"
description: |-
  Provides a Alicloud Express Connect Vbr Ha resource.
---

# alicloud_vpc_vbr_ha

Provides a Express Connect Vbr Ha resource.

VBR switching Group.

For information about Express Connect Vbr Ha and how to use it, see [What is Vbr Ha](https://www.alibabacloud.com/help/doc-detail/212629.html).

-> **NOTE:** Available since v1.151.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_vbr_ha&exampleId=51a841d0-a62c-fd12-fbcc-92985eb2b611632d8845&activeTab=example&spm=docs.r.vpc_vbr_ha.0.51a841d0a6&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
}
data "alicloud_regions" "default" {
  current = true
}
data "alicloud_express_connect_physical_connections" "example" {
  name_regex = "^preserved-NODELETING"
}
resource "random_integer" "vlan_id" {
  max = 2999
  min = 1
}
resource "alicloud_express_connect_virtual_border_router" "example" {
  count                      = 2
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.example.connections[count.index].id
  virtual_border_router_name = format("${var.name}-%d", count.index + 1)
  vlan_id                    = random_integer.vlan_id.id + count.index
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}

resource "alicloud_cen_instance" "example" {
  cen_instance_name = var.name
  description       = var.name
  protection_level  = "REDUCED"
}

resource "alicloud_cen_instance_attachment" "example" {
  count                    = 2
  instance_id              = alicloud_cen_instance.example.id
  child_instance_id        = alicloud_express_connect_virtual_border_router.example[count.index].id
  child_instance_type      = "VBR"
  child_instance_region_id = data.alicloud_regions.default.regions.0.id
}

resource "alicloud_vpc_vbr_ha" "example" {
  vbr_id      = alicloud_cen_instance_attachment.example[0].child_instance_id
  peer_vbr_id = alicloud_cen_instance_attachment.example[1].child_instance_id
  vbr_ha_name = var.name
  description = var.name
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional, ForceNew) The description of the VBR switching group.
It must be 2 to 256 characters in length and must start with a letter or Chinese, but cannot start with 'http:// 'or 'https.
* `dry_run` - (Optional) Whether to PreCheck only this request. Value range:
  - *true**: The check request is sent and the instance is not started. Check whether the required parameters, request format, and instance status are filled in. If the check does not pass, the corresponding error is returned. If the check passes, DRYRUN.SUCCESS is returned.
  - *false** (default): Send a normal request and start the instance directly after passing the check.
* `vbr_ha_name` - (Optional, ForceNew) VBR switch group name.
* `peer_vbr_id` - (Required, ForceNew) The instance ID of another VBR in The VBR switching group.
* `vbr_id` - (Required, ForceNew) The VBR instance ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the VBR.
* `region_id` - The ID of the region to which the VBR belongs.
* `status` - Status of VBR switching Group

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 8 mins) Used when create the Vbr Ha.
* `delete` - (Defaults to 8 mins) Used when delete the Vbr Ha.

## Import

Express Connect Vbr Ha can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_vbr_ha.example <id>
```