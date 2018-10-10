---
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_group"
sidebar_current: "docs-alicloud-resource-api-gateway-group"
description: |-
  Provides a Alicloud Api Gateway Group Resource.
---

# alicloud_api_gateway_group

Provides an api group resource.

~> **NOTE:** Terraform will auto build api group while it uses `alicloud_api_gateway_group` to build api group.

## Example Usage

Basic Usage

```
resource "alicloud_api_gateway_group" "apiGroup" {
  name = "ApiGatewayGroup"
  description = "description of the api group"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource) The name of the api gateway group. Defaults to null.
* `description` - (Required, Forces new resource) The description of the api gateway group. Defaults to null.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the api group of api gateway.