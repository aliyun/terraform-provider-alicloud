---
subcategory: "API Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_group"
sidebar_current: "docs-alicloud-resource-api-gateway-group"
description: |-
  Provides a Alicloud Api Gateway Group Resource.
---

# alicloud_api_gateway_group

Provides an api group resource.To create an API, you must firstly create a group which is a basic attribute of the API.

For information about Api Gateway Group and how to use it, see [Create An Api Group](https://www.alibabacloud.com/help/doc-detail/43611.html)

-> **NOTE:** Terraform will auto build api group while it uses `alicloud_api_gateway_group` to build api group.

## Example Usage

Basic Usage

```
resource "alicloud_api_gateway_group" "apiGroup" {
  name        = "ApiGatewayGroup"
  description = "description of the api group"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the api gateway group. Defaults to null.
* `description` - (Required) The description of the api gateway group. Defaults to null.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the api group of api gateway.

## Import

Api gateway group can be imported using the id, e.g.

```
$ terraform import alicloud_api_gateway_group.example "ab2351f2ce904edaa8d92a0510832b91"
```
