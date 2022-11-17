---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_app_attachment"
sidebar_current: "docs-alicloud-resource-api-gateway-app-attachment"
description: |-
  Provides a Alicloud Api Gateway App Attachment Resource.
---

# alicloud_api_gateway_app_attachment

Provides an app attachment resource.It is used for authorizing a specific api to an app accessing. 

For information about Api Gateway App attachment and how to use it, see [Add specified API access authorities](https://www.alibabacloud.com/help/doc-detail/43673.htm)

-> **NOTE:** Terraform will auto build app attachment while it uses `alicloud_api_gateway_app_attachment` to build.

## Example Usage

Basic Usage

```terraform
resource "alicloud_api_gateway_app_attachment" "foo" {
  api_id     = "d29d25b9cfdf4742b1a3f6537299a749"
  group_id   = "aaef8cdbb404420f9398a74ed1db7fff"
  app_id     = "20898181"
  stage_name = "PRE"
}
```

## Argument Reference

The following arguments are supported:

* `api_id` - (Required，ForceNew) The api_id that app apply to access.
* `group_id` - (Required，ForceNew) The group that the api belongs to.
* `app_id` - (Required，ForceNew) The app that apply to the authorization.
* `stage_name` - (Required，ForceNew) Stage that the app apply to access.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the app attachment of api gateway., formatted as `<group_id>:<api_id>:<app_id>:<stage_name>`.
