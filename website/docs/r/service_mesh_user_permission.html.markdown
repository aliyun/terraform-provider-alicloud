---
subcategory: "Service Mesh"
layout: "alicloud"
page_title: "Alicloud: alicloud_service_mesh_user_permission"
sidebar_current: "docs-alicloud-resource-service-mesh-user-permission"
description: |-
  Provides an Alicloud Service Mesh User Permission resource.
---

# alicloud\_service\_mesh\_user\_permission

Provides a Service Mesh UserPermission resource.

For information about Service Mesh User Permission and how to use it, see [What is User Permission](https://help.aliyun.com/document_detail/171622.html).

-> **NOTE:** Available in v1.174.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "servicemesh"
}
data "alicloud_service_mesh_versions" "default" {
  edition = "Default"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_ram_user" "default" {
  name = var.name
}

resource "alicloud_service_mesh_service_mesh" "default1" {
  service_mesh_name = var.name
  edition           = "Default"
  version           = data.alicloud_service_mesh_versions.default.versions.0.version
  cluster_spec      = "standard"
  network {
    vpc_id        = data.alicloud_vpcs.default.ids.0
    vswitche_list = [data.alicloud_vswitches.default.ids.0]
  }
  load_balancer {
    pilot_public_eip      = false
    api_server_public_eip = false
  }
}

resource "alicloud_service_mesh_user_permission" "example" {
  sub_account_user_id = alicloud_ram_user.default.id
  permissions {
    role_name       = "istio-admin"
    service_mesh_id = alicloud_service_mesh_service_mesh.default1.id
    role_type       = "custom"
    is_custom       = true
    is_ram_role     = false
  }
}

```

## Argument Reference

The following arguments are supported:

* `sub_account_user_id` - (Required, ForceNew) The configuration of the Load Balancer. See the following `Block load_balancer`.
* `permissions` - (Optional) List of permissions. **Warning:** The list requires the full amount of permission information to be passed. Adding permissions means adding items to the list, and deleting them or inputting nothing means removing items. See the following `Block permissions`.


#### Block permissions

The permissions supports the following:

* `role_name` - (Optional) The permission name. Valid values: `istio-admin`, `istio-ops`, `istio-readonly`.
  - `istio-admin`:  The administrator.
  - `istio-ops`: The administrator of the service mesh resource.
  - `istio-readonly`: The read only permission.
* `service_mesh_id` - (Optional) The service mesh id.
* `role_type` - (Optional) The role type. Valid Value: `custom`.
* `is_custom` - (Optional) Whether the grant object is a RAM role.
* `is_ram_role` - (Optional) Whether the grant object is an entity.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of User Permission. The value is same as `sub_account_user_id`.


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 15 mins) Used when create the Service Mesh User Permission.
* `update` - (Defaults to 15 mins) Used when update the Service Mesh User Permission.

## Import

Service Mesh User Permission can be imported using the id, e.g.

```
$ terraform import alicloud_service_mesh_user_permission.example <id>
```
