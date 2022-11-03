---
subcategory: "API Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_app"
sidebar_current: "docs-alicloud-resource-api-gateway-app"
description: |-
  Provides a Alicloud Api Gateway App Resource.
---

# alicloud_api_gateway_app

Provides an app resource.It must create an app before calling a third-party API because the app is the identity used to call the third-party API.

For information about Api Gateway App and how to use it, see [Create An APP](https://www.alibabacloud.com/help/doc-detail/43663.html)

-> **NOTE:** Terraform will auto build api app while it uses `alicloud_api_gateway_app` to build api app.

## Example Usage

Basic Usage

```terraform
resource "alicloud_api_gateway_app" "apiTest" {
  name        = "ApiGatewayAPp"
  description = "description of the app"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the app. 
* `description` - (Optional) The description of the app. Defaults to null.
* `tags` - (Optional, Available in v1.55.3+) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the app of api gateway.

## Import

Api gateway app can be imported using the id, e.g.

```shell
$ terraform import alicloud_api_gateway_app.example "7379660"
```
