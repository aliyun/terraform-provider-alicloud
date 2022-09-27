---
subcategory: "Cloud Architect Design Tools"
layout: "alicloud"
page_title: "Alicloud: alicloud_bp_studio_applications"
sidebar_current: "docs-alicloud-datasource-bp-studio-applications"
description: |-
  Provides a list of Cloud Architect Design Tools Applications to the user.
---

# alicloud\_bp\_studio\_applications

This data source provides the Cloud Architect Design Tools Applications of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.192.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_bp_studio_applications" "ids" {
  ids = ["example_id"]
}

output "bp_studio_application_id_1" {
  value = data.alicloud_bp_studio_applications.ids.applications.0.id
}

data "alicloud_bp_studio_applications" "nameRegex" {
  name_regex = "^my-Application"
}

output "bp_studio_application_id_2" {
  value = data.alicloud_bp_studio_applications.nameRegex.applications.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Application IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Application name.
* `keyword` - (Optional, ForceNew) The keyword of the Application.
* `order_type` - (Optional, ForceNew) The order type of the Application. Valid values:
  - `1`: The update time of the Application.
  - `2`: The create time of the Application.
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `status` - (Optional, ForceNew) The status of the Application. Valid values: `success`, `release`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Application names.
* `applications` - A list of Cloud Architect Design Tools Applications. Each element contains the following attributes:
  * `id` - The ID of the Application.
  * `application_id` - The ID of the Application.
  * `application_name` - The name of the Application.
  * `resource_group_id` - The ID of the resource group.
  * `topo_url` - The topo url of the Application.
  * `image_url` - The image url of the Application.
  * `create_time` - The creation time of the Application.
  * `status` - The status of the Application.
