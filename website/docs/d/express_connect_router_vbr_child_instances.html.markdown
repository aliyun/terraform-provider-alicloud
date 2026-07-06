---
subcategory: "Express Connect Router"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_router_vbr_child_instances"
description: |-
  Provides a list of Express Connect Router Vbr Child Instances to the user.
---

# alicloud_express_connect_router_vbr_child_instances

This data source provides the Express Connect Router Vbr Child Instances of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.285.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_account" "default" {
}

data "alicloud_regions" "default" {
  current = true
}

data "alicloud_express_connect_physical_connections" "default" {
  name_regex = "^preserved-NODELETING"
}

resource "alicloud_express_connect_router_express_connect_router" "default" {
  alibaba_side_asn = "65532"
}

resource "alicloud_express_connect_virtual_border_router" "default" {
  physical_connection_id = data.alicloud_express_connect_physical_connections.default.connections.0.id
  vlan_id                = "1000"
  peer_gateway_ip        = "192.168.254.2"
  peering_subnet_mask    = "255.255.255.0"
  local_gateway_ip       = "192.168.254.1"
}

resource "alicloud_express_connect_router_vbr_child_instance" "default" {
  ecr_id                   = alicloud_express_connect_router_express_connect_router.default.id
  child_instance_id        = alicloud_express_connect_virtual_border_router.default.id
  child_instance_type      = "VBR"
  child_instance_owner_id  = data.alicloud_account.default.id
  child_instance_region_id = data.alicloud_regions.default.regions.0.id
  description              = var.name
}

data "alicloud_express_connect_router_vbr_child_instances" "ids" {
  ids    = [alicloud_express_connect_router_vbr_child_instance.default.id]
  ecr_id = alicloud_express_connect_router_vbr_child_instance.default.ecr_id
}

output "express_connect_router_vbr_child_instances_id_0" {
  value = data.alicloud_express_connect_router_vbr_child_instances.ids.instances.0.id
}
```

## Argument Reference

The following attributes are exported:

* `ids` - (Optional, List) A list of Vbr Child Instance IDs.
* `ecr_id` - (Required) The ID of the Express Connect Router instance.
* `child_instance_id` - (Optional) The ID of the network instance to detach.
* `child_instance_type` - (Optional) The type of the network instance. Valid values: `VBR`.
* `child_instance_region_id` - (Optional) The region where the network instance is deployed.
* `status` - (Optional) The deployment status of the associated instance. Valid values: `CREATING`, `ACTIVE`, `ASSOCIATING`, `DISSOCIATING`, `UPDATING`, `DELETING`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `instances` - A list of Vbr Child Instances. Each element contains the following attributes:
  * `id` - The ID of the Vbr Child Instance.
  * `ecr_id` - The ID of the Express Connect Router instance.
  * `child_instance_id` - The ID of the virtual border router instance.
  * `child_instance_type` - The type of the child instance.
  * `child_instance_owner_id` - The Alibaba Cloud account ID of the child instance owner.
  * `child_instance_region_id` - The region of the child instance.
  * `description` - The description of the child instance.
  * `status` - The deployment status of the associated instance.
  * `create_time` - The time when the association was created.
  * `modify_time` - The time when the association was modified.
