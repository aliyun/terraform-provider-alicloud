---
subcategory: "Express Connect Router"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_router_vbr_child_instance"
description: |-
  Provides a Alicloud Express Connect Router Express Connect Router Vbr Child Instance resource.
---

# alicloud_express_connect_router_vbr_child_instance

Provides a Express Connect Router Express Connect Router Vbr Child Instance resource.

For information about Express Connect Router Express Connect Router Vbr Child Instance and how to use it, see [What is Express Connect Router Vbr Child Instance](https://next.api.alibabacloud.com/api/ExpressConnectRouter/2023-09-01/AttachExpressConnectRouterChildInstance).

-> **NOTE:** Available since v1.224.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_express_connect_router_vbr_child_instance&exampleId=0032239e-111f-a619-ea6a-5d6f876f283f888d35a7&activeTab=example&spm=docs.r.express_connect_router_vbr_child_instance.0.0032239e11&intl_lang=EN_US" target="_blank">
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

data "alicloud_express_connect_physical_connections" "nameRegex" {
  name_regex = "^preserved-NODELETING"
}

resource "alicloud_express_connect_virtual_border_router" "defaultydbbk3" {
  physical_connection_id = data.alicloud_express_connect_physical_connections.nameRegex.connections.0.id
  vlan_id                = "1000"
  peer_gateway_ip        = "192.168.254.2"
  peering_subnet_mask    = "255.255.255.0"
  local_gateway_ip       = "192.168.254.1"
}

resource "alicloud_express_connect_router_express_connect_router" "defaultAAlhUy" {
  alibaba_side_asn = "65532"
}

data "alicloud_account" "current" {
}

resource "alicloud_express_connect_router_vbr_child_instance" "default" {
  child_instance_id        = alicloud_express_connect_virtual_border_router.defaultydbbk3.id
  child_instance_region_id = "cn-hangzhou"
  ecr_id                   = alicloud_express_connect_router_express_connect_router.defaultAAlhUy.id
  child_instance_type      = "VBR"
  child_instance_owner_id  = data.alicloud_account.current.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_express_connect_router_vbr_child_instance&spm=docs.r.express_connect_router_vbr_child_instance.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `child_instance_id` - (Required, ForceNew) The ID of the leased line gateway subinstance.
* `child_instance_owner_id` - (Optional, ForceNew, Computed, Int) The ID of the Alibaba Cloud account (primary account) to which the VBR instance belongs.

-> **NOTE:**  This parameter is required if you want to load a cross-account network instance.

* `child_instance_region_id` - (Required, ForceNew) Region of the leased line gateway sub-instance
* `child_instance_type` - (Required, ForceNew) The type of the network instance. Value: `VBR`: VBR instance.
* `description` - (Optional, Available since v1.235.0) Resource attribute fields that represent descriptive information
* `ecr_id` - (Required, ForceNew) ID of the representative leased line gateway instance.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<ecr_id>:<child_instance_id>:<child_instance_type>`.
* `create_time` - The creation time of the resource.
* `status` - Binding relationship status of leased line gateway subinstances.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Express Connect Router Vbr Child Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Express Connect Router Vbr Child Instance.
* `update` - (Defaults to 5 mins) Used when update the Express Connect Router Vbr Child Instance.

## Import

Express Connect Router Express Connect Router Vbr Child Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_express_connect_router_vbr_child_instance.example <ecr_id>:<child_instance_id>:<child_instance_type>
```