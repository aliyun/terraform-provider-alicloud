---
layout: "alicloud"
page_title: "Alicloud: alicloud_apigateway_api_group"
sidebar_current: "docs-alicloud-resource-apigateway-api-group"
description: |-
  Provides a Alicloud Apigateway api group resource.
---

# alicloud_apigateway_api_group

Provides an api group resource.

~> **NOTE:** Terraform will auto build api group while it uses `alicloud_apigateway_api_group` to build api group.

## Example Usage

Basic Usage

```
resource "alicloud_apigateway_api_group" "apiGroup" {
  name = "apigatewayGroup"
  description = "description of the api group"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource) The name of the apigateway group. Defaults to null.
* `description` - (Required, Forces new resource) The description of the apigateway group. Defaults to null.

## Attributes Reference

The following attributes are exported:

* `group_id` - The ID of the api group of apigateway.