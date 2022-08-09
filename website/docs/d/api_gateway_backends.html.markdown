---
subcategory: "API Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_backends"
sidebar_current: "docs-alicloud-datasource-api-gateway-backends"
description: |-
  Provides a list of Api Gateway Backends to the user.
---

# alicloud\_api\_gateway\_backends

This data source provides the Api Gateway Backends of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.181.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_api_gateway_backends" "ids" {}
output "api_gateway_backend_id_1" {
  value = data.alicloud_api_gateway_backends.ids.backends.0.id
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional, ForceNew) A regex string to filter Api Gateway Backends by name.
* `ids` - (Optional, ForceNew) A list of Backends IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Backend names.
* `backends` - A list of Api Gateway Backends. Each element contains the following attributes: 
  * `id` - The resource ID in terraform of Backend.
  * `backend_id` -  The id of the Backend.
  * `backend_type` - The type of the Backend.
  * `backend_name` - The name of the Backend.
  * `create_time` - The created time of the Backend.
  * `description` - The description of the Backend.
  * `modified_time` - The modified time of the Backend.