---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_grey_tag_routes"
sidebar_current: "docs-alicloud-datasource-sae-grey-tag-routes"
description: |-
  Provides a list of Sae GreyTagRoutes to the user.
---

# alicloud\_sae\_grey\_tag\_routes

This data source provides the Sae GreyTagRoutes of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.160.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_sae_grey_tag_routes" "nameRegex" {
  app_id     = "example_id"
  name_regex = "^my-GreyTagRoute"
}
output "sae_grey_tag_routes_id" {
  value = data.alicloud_sae_grey_tag_routes.nameRegex.routes.0.id
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, ForceNew) The ID  of the SAE Application.
* `ids` - (Optional, ForceNew, Computed)  A list of GreyTagRoute IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by GreyTagRoute name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of GreyTagRoute names.
* `routes` - A list of Sae GreyTagRoutes. Each element contains the following attributes:

### `routes`

* `id` - The ID of the GreyTagRoute.
* `dubbo_rules` - The grayscale rule created for Dubbo Application.
* `description` - The description of GreyTagRoute.
* `grey_tag_route_name` - The name of GreyTagRoute.
* `sc_rules` - The grayscale rule created for SpringCloud Application.

### `routes.dubbo_rules`

* `method_name` - The method name
* `service_name` - The service name.
* `version` - The service version.
* `condition` - The conditional Patterns for Grayscale Rules.
* `group` - The service group.
* `items` - A list of conditions items.

### `routes.dubbo_rules.items`

* `index` - The parameter number.
* `expr` - The parameter value gets the expression.
* `cond` - The comparison operator.
* `value` - The value of the parameter.
* `operator` - The operator.

### `routes.sc_rules`

* `path` - The path corresponding to the grayscale rule.
* `condition` - The Conditional Patterns for Grayscale Rules.
* `items` - A list of conditions items.

### `routes.sc_rules.items`

* `name` - The name of the parameter.
* `cond` - The comparison operator.
* `type` - The Compare types.
* `value` - The value of the parameter.
* `operator` - The operator.
