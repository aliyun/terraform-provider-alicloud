---
subcategory: "Service Mesh"
layout: "alicloud"
page_title: "Alicloud: alicloud_service_mesh_user_permission"
sidebar_current: "docs-alicloud-resource-service-mesh-user-permission"
description: |-
  Provides an Alicloud Service Mesh User Permission resource.
---

# alicloud_service_mesh_user_permission

Provides a Service Mesh UserPermission resource.

For information about Service Mesh User Permission and how to use it, see [What is User Permission](https://www.alibabacloud.com/help/en/alibaba-cloud-service-mesh/latest/api-servicemesh-2020-01-11-grantuserpermissions).

-> **NOTE:** Available since v1.174.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_service_mesh_user_permission&exampleId=c1b96c8e-6c5b-6dbb-ba10-d274c2629366530e8a52&activeTab=example&spm=docs.r.service_mesh_user_permission.0.c1b96c8e6c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tfexample"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

data "alicloud_service_mesh_versions" "default" {
  edition = "Default"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_ram_user" "default" {
  name = var.name
}

resource "alicloud_service_mesh_service_mesh" "default1" {
  service_mesh_name = "${var.name}-${random_integer.default.result}"
  edition           = "Default"
  cluster_spec      = "standard"
  version           = data.alicloud_service_mesh_versions.default.versions.0.version
  network {
    vpc_id        = data.alicloud_vpcs.default.ids.0
    vswitche_list = [data.alicloud_vswitches.default.ids.0]
  }
  load_balancer {
    pilot_public_eip      = false
    api_server_public_eip = false
  }
}

resource "alicloud_service_mesh_user_permission" "default" {
  sub_account_user_id = alicloud_ram_user.default.id
  permissions {
    role_name       = "istio-ops"
    service_mesh_id = alicloud_service_mesh_service_mesh.default1.id
    role_type       = "custom"
    is_custom       = true
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_service_mesh_user_permission&spm=docs.r.service_mesh_user_permission.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `sub_account_user_id` - (Required, ForceNew) The configuration of the Load Balancer. See the following `Block load_balancer`.
* `permissions` - (Optional) List of permissions. **Warning:** The list requires the full amount of permission information to be passed. Adding permissions means adding items to the list, and deleting them or inputting nothing means removing items. See [`permissions`](#permissions) below.

### `permissions`

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


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 15 mins) Used when create the Service Mesh User Permission.
* `update` - (Defaults to 15 mins) Used when update the Service Mesh User Permission.

## Import

Service Mesh User Permission can be imported using the id, e.g.

```shell
$ terraform import alicloud_service_mesh_user_permission.example <id>
```
