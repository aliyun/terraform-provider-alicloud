---
subcategory: "DataWorks"
layout: "alicloud"
page_title: "Alicloud: alicloud_dataworks_data_service_authorized_apis"
sidebar_current: "docs-alicloud-datasource-dataworks-data-service-authorized-apis"
description: |-
  Provides a list of DataWorks Data Service authorized APIs to the specified filters.
---

# alicloud_dataworks_data_service_authorized_apis

This data source provides a list of DataWorks Data Service authorized APIs in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available since v1.294.0.

## Example Usage

```terraform
data "alicloud_dataworks_data_service_authorized_apis" "default" {
  project_id       = "123456"
  api_name_keyword = "tf-testacc"
  page_size        = 20
}

output "api_id" {
  value = data.alicloud_dataworks_data_service_authorized_apis.default.apis.0.api_id
}
```

## Argument Reference

The following arguments are supported:

* `api_name_keyword` - (Optional) The keyword used to filter Data Service API names. The API name that contains this keyword is returned.
* `project_id` - (Required) The ID of the DataWorks workspace.
* `tenant_id` - (Optional) The ID of the tenant.
* `page_number` - (Optional) The number of the page to return. Default to `1`.
* `page_size` - (Optional) The number of items to return on each page. Default to `50`. Valid values: `1` to `100`.
* `name_regex` - (Optional) A regex string to filter results by API name.
* `ids` - (Optional) A list of API IDs used to filter the results.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of API IDs.
* `names` - A list of API names.
* `total_count` - The total number of authorized API records.
* `apis` - A list of Data Service authorized APIs. Each element contains the following attributes:
  * `api_id` - The ID of the API.
  * `api_name` - The name of the API.
  * `api_path` - The path of the API.
  * `project_id` - The ID of the DataWorks workspace that the API belongs to.
  * `tenant_id` - The ID of the tenant.
  * `group_id` - The ID of the API group.
  * `region_id` - The region ID of the API.
  * `status` - The status of the authorization.
  * `grant_end_time` - The end time of the authorization.
  * `grant_created_time` - The time when the authorization was created.
  * `grant_operator_id` - The ID of the user who granted the authorization.
  * `create_time` - The time when the API was created.
  * `modified_time` - The time when the API was last modified.
  * `id` - The ID of the API. It is the same as `api_id`.
