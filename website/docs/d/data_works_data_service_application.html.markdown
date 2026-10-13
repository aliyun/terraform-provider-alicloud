---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_data_service_application"
sidebar_current: "docs-alicloud-datasource-data-works-data-service-application"
description: |-
  Provides a list of Data Works Data Service Applications to the user.
---

# alicloud\_data\_works\_data\_service\_application

This data source provides the Data Works Data Service Applications of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.230.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_data_works_data_service_application" "default" {
  project_id = "xxxx"
  ids        = ["example_application_id"]
}

output "data_service_application_id" {
  value = data.alicloud_data_works_data_service_application.default.applications.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, Computed) A list of Data Service Application IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `project_id` - (Optional) The ID of the Data Works project. The list of project IDs to query. If `project_id_list` is not set, this field is used as a single-element project ID list.
* `project_id_list` - (Optional) The list of Data Works project IDs to query. One of `project_id` or `project_id_list` must be set.
* `tenant_id` - (Optional) The ID of the tenant.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Data Service Application IDs.
* `applications` - A list of Data Works Data Service Applications. Each element contains the following attributes:
  * `id` - The ID of the Data Service Application.
  * `application_id` - The ID of the Data Service Application.
  * `application_key` - The key of the Data Service Application.
  * `application_name` - The name of the Data Service Application.
  * `application_secret` - The secret of the Data Service Application.
  * `project_id` - The ID of the Data Works project that the application belongs to.
  * `region_id` - The region ID of the Data Service Application.
  * `application_code` - The code of the Data Service Application.
