---
subcategory: "Service Mesh"
layout: "alicloud"
page_title: "Alicloud: alicloud_service_mesh_extension_providers"
sidebar_current: "docs-alicloud-datasource-service-mesh-extension-providers"
description: |-
  Provides a list of Service Mesh Extension Providers to the user.
---

# alicloud\_service\_mesh\_extension\_providers

This data source provides the Service Mesh Extension Providers of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.191.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_service_mesh_extension_providers" "ids" {
  ids             = ["example_id"]
  service_mesh_id = "example_service_mesh_id"
  type            = "httpextauth"
}

output "service_mesh_extension_providers_id_1" {
  value = data.alicloud_service_mesh_extension_providers.ids.providers.0.id
}

data "alicloud_service_mesh_extension_providers" "nameRegex" {
  name_regex      = "^my-ServiceMeshExtensionProvider"
  service_mesh_id = "example_service_mesh_id"
  type            = "httpextauth"
}

output "service_mesh_extension_providers_id_2" {
  value = data.alicloud_service_mesh_extension_providers.nameRegex.providers.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Service Mesh Extension Provider IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Service Mesh Extension Provider name.
* `service_mesh_id` - (Required, ForceNew) The ID of the Service Mesh.
* `type` - (Required, ForceNew) The type of the Service Mesh Extension Provider. Valid values: `httpextauth`, `grpcextauth`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Extension Provider names.
* `providers` - A list of Service Mesh Extension Providers. Each element contains the following attributes:
    * `id` - The ID of the Service Mesh Extension Provider. It formats as `<service_mesh_id>:<type>:<extension_provider_name>`.
    * `service_mesh_id` - The ID of the Service Mesh.
    * `extension_provider_name` - The name of the Service Mesh Extension Provider.
    * `type` - The type of the Service Mesh Extension Provider.
    * `config` - The config of the Service Mesh Extension Provider.
