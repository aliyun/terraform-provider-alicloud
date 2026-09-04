---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_portals"
sidebar_current: "docs-alicloud-datasource-apig-portals"
description: |-
  Provides a list of Apig Portal owned by an Alibaba Cloud account.
---

# alicloud_apig_portals

This data source provides Apig Portal available to the user.[What is Portal](https://next.api.alibabacloud.com/document/APIG/2024-03-27/CreatePortal)

-> **NOTE:** Available since v1.291.0.

## Example Usage

```terraform
data "alicloud_apig_portals" "default" {
}

output "first_portal_id" {
  value = data.alicloud_apig_portals.default.portals.0.id
}
```

## Argument Reference

The following arguments are supported:
* `name` - (ForceNew, Optional) The name of the resource
* `ids` - (Optional, Computed) A list of Portal IDs.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Portal IDs.
* `portals` - A list of Portal Entries. Each element contains the following attributes:
  * `description` - Portal description.
  * `name` - The name of the resource.
  * `portal_id` - The first ID of the resource.
  * `region_id` - **NOTE:** This field is only available when `enable_details` is `true`. The region ID of the resource.
  * `id` - The ID of the resource supplied above.
