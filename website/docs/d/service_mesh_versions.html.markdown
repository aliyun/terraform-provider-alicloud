---
subcategory: "Service Mesh"
layout: "alicloud"
page_title: "Alicloud: alicloud_service_mesh_service_meshes"
sidebar_current: "docs-alicloud-datasource-service-mesh-service-meshes"
description: |-
  Provides a list of Service Mesh Versions to the user.
---

# alicloud\_service\_mesh\_versions

This data source provides ASM available versions in the specified region.

-> **NOTE:** Available in v1.161.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_service_mesh_versions" "default" {
  edition = "Default"
}
output "service_mesh_version" {
  value = data.alicloud_service_mesh_versions.versions.0.version
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) The ASM versions available for the ASM instance of the current edition. Its element formats as `<edition>:<version>`.
* `edition` - (Optional, ForceNew) The edition of the ASM instance. Valid values:
  - Default: Standard Edition 
  - Pro: Professional Edition
    
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of ASM versions. Its element formats as `<edition>:<version>`.
* `version` - A list of Service Mesh Service Meshes. Each element contains the following attributes:
    * `id` - The ASM version id. It formats as `<edition>:<version>`.
    * `version` - The AMS version.
    * `edition` - The edition of the ASM instance.
