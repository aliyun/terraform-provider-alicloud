---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_shared_resources"
sidebar_current: "docs-alicloud-datasource-resource-manager-shared-resources"
description: |-
  Provides a list of Resource Manager Shared Resources to the user.
---

# alicloud_resource_manager_shared_resources

This data source provides the Resource Manager Shared Resources of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.111.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.ids.0
}

resource "alicloud_resource_manager_resource_share" "default" {
  resource_share_name = var.name
}

resource "alicloud_resource_manager_shared_resource" "default" {
  resource_share_id = alicloud_resource_manager_resource_share.default.id
  resource_id       = data.alicloud_vswitches.default.ids.0
  resource_type     = "VSwitch"
}

data "alicloud_resource_manager_shared_resources" "ids" {
  ids = [format("%s:%s", alicloud_resource_manager_shared_resource.default.resource_id, alicloud_resource_manager_shared_resource.default.resource_type)]
}

output "first_resource_manager_shared_resource_id" {
  value = data.alicloud_resource_manager_shared_resources.ids.resources.0.id
}

data "alicloud_resource_manager_shared_resources" "resourceShareId" {
  resource_share_id = alicloud_resource_manager_shared_resource.default.resource_share_id
}

output "second_resource_manager_shared_resource_id" {
  value = data.alicloud_resource_manager_shared_resources.resourceShareId.resources.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of shared resource IDs.
* `resource_share_id` - (Optional, ForceNew) The resource share ID of resource manager.
* `status` - (Optional, ForceNew) The status of share resource. Valid values: `Associated`, `Associating`, `Disassociated`, `Disassociating` and `Failed`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `resources` - A list of Resource Manager Shared Resources. Each element contains the following attributes:
  * `id` - The ID of the Shared Resource. It formats as `<resource_id>:<resource_type>`.
  * `resource_id` - The ID of the shared resource.
  * `resource_type` - The type of shared resource.
  * `resource_share_id` - The resource share ID of resource manager.
  * `status` - The status of shared resource.
