---
subcategory: "Simple Application Server"
layout: "alicloud"
page_title: "Alicloud: alicloud_simple_application_server_plans"
sidebar_current: "docs-alicloud-datasource-simple-application-server-plans"
description: |-
  Provides a list of Simple Application Server Plans to the user.
---

# alicloud\_simple\_application\_server\_plans

This data source provides the Simple Application Server Plans of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.135.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_simple_application_server_plans" "example" {
  memory    = 1
  bandwidth = 3
  disk_size = 40
  flow      = 6
  core      = 2
}
output "simple_application_server_plan_id_1" {
  value = data.alicloud_simple_application_server_plans.ids.plans.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Instance Plan IDs.
* `bandwidth` - The peak bandwidth. Unit: Mbit/s.
* `core` - The number of CPU cores.
* `disk_size` - The size of the enhanced SSD (ESSD). Unit: GB.
* `flow` - The monthly data transfer quota. Unit: GB.
* `memory` - The memory size. Unit: GB.
* `platform` - (Available in v1.161.0) The platform of Plan supported. Valid values: ["Linux", "Windows"].
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `plans` - A list of Simple Application Server Plans. Each element contains the following attributes:
	* `bandwidth` - The peak bandwidth. Unit: Mbit/s.
	* `core` - The number of CPU cores.
	* `disk_size` - The size of the enhanced SSD (ESSD). Unit: GB.
	* `flow` - The monthly data transfer quota. Unit: GB.
	* `id` - The ID of the Instance Plan.
	* `plan_id` - The ID of the Instance Plan.
	* `memory` - The memory size. Unit: GB.
	* `support_platform` - (Available in v1.161.0) The platform of Plan supported.
