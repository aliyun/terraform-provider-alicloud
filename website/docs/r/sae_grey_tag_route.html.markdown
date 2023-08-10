---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_grey_tag_route"
sidebar_current: "docs-alicloud-resource-sae-grey_tag_route"
description: |-
  Provides a Alicloud Serverless App Engine (SAE) GreyTagRoute resource.
---

# alicloud_sae_grey_tag_route

Provides a Serverless App Engine (SAE) GreyTagRoute resource.

For information about Serverless App Engine (SAE) GreyTagRoute and how to use it, see [What is GreyTagRoute](https://www.alibabacloud.com/help/en/sae/latest/create-grey-tag-route).

-> **NOTE:** Available since v1.160.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_regions" "default" {
  current = true
}
resource "random_integer" "default" {
  max = 99999
  min = 10000
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}
resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_sae_namespace" "default" {
  namespace_id              = "${data.alicloud_regions.default.regions.0.id}:example${random_integer.default.result}"
  namespace_name            = var.name
  namespace_description     = var.name
  enable_micro_registration = false
}

resource "alicloud_sae_application" "default" {
  app_description   = var.name
  app_name          = var.name
  namespace_id      = alicloud_sae_namespace.default.id
  image_url         = "registry-vpc.${data.alicloud_regions.default.regions.0.id}.aliyuncs.com/sae-demo-image/consumer:1.0"
  package_type      = "Image"
  security_group_id = alicloud_security_group.default.id
  vpc_id            = alicloud_vpc.default.id
  vswitch_id        = alicloud_vswitch.default.id
  timezone          = "Asia/Beijing"
  replicas          = "5"
  cpu               = "500"
  memory            = "2048"
}

resource "alicloud_sae_grey_tag_route" "default" {
  grey_tag_route_name = var.name
  description         = var.name
  app_id              = alicloud_sae_application.default.id
  sc_rules {
    items {
      type     = "param"
      name     = "tfexample"
      operator = "rawvalue"
      value    = "example"
      cond     = "=="
    }
    path      = "/tf/example"
    condition = "AND"
  }

  dubbo_rules {
    items {
      cond     = "=="
      expr     = ".key1"
      index    = "1"
      operator = "rawvalue"
      value    = "value1"
    }
    condition    = "OR"
    group        = "DUBBO"
    method_name  = "example"
    service_name = "com.example.service"
    version      = "1.0.0"
  }
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, ForceNew) The ID  of the SAE Application.
* `description` - (Optional) The description of GreyTagRoute.
* `grey_tag_route_name` - (Required, ForceNew) The name of GreyTagRoute.
* `dubbo_rules` - (Optional) The grayscale rule created for Dubbo Application. See [`dubbo_rules`](#dubbo_rules) below.
* `sc_rules` - (Optional) The grayscale rule created for SpringCloud Application. See [`sc_rules`](#sc_rules) below.

### `dubbo_rules`

The `dubbo_rules` supports the following:
* `method_name` - (Optional) The method name
* `service_name` - (Optional) The service name.
* `version` - (Optional) The service version.
* `condition` - (Optional) The Conditional Patterns for Grayscale Rules. Valid values: `AND`, `OR`.
* `group` - (Optional) The service group.
* `items` - (Optional) A list of conditions items. See [`items`](#dubbo_rules-items) below.

### `dubbo_rules-items`

The `items` supports the following:
* `index` - (Optional) The parameter number.
* `expr` - (Optional) The parameter value gets the expression.
* `cond` - (Optional) The comparison operator. Valid values: `>`, `<`, `>=`, `<=`, `==`, `!=`.
* `value` - (Optional) The value of the parameter.
* `operator` - (Optional) The operator. Valid values: `rawvalue`, `list`, `mod`, `deterministic_proportional_steaming_division`

### `sc_rules`

The `sc_rules` supports the following:
* `path` - (Optional) The path corresponding to the grayscale rule.
* `condition` - (Optional) The conditional Patterns for Grayscale Rules. Valid values: `AND`, `OR`.
* `items` - (Optional) A list of conditions items.See [`items`](#sc_rules-items) below.

### `sc_rules-items`

The `items` supports the following:
* `name` - (Optional) The name of the parameter.
* `type` - (Optional) The compare types. Valid values: `param`, `cookie`, `header`.
* `cond` - (Optional) The comparison operator. Valid values: `>`, `<`, `>=`, `<=`, `==`, `!=`.
* `value` - (Optional) The value of the parameter.
* `operator` - (Optional) The operator. Valid values: `rawvalue`, `list`, `mod`, `deterministic_proportional_steaming_division`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of GreyTagRoute.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the resource.
* `update` - (Defaults to 1 mins) Used when update the resource.
* `delete` - (Defaults to 1 mins) Used when delete the resource.

## Import

Serverless App Engine (SAE) GreyTagRoute can be imported using the id, e.g.

```shell
$ terraform import alicloud_sae_grey_tag_route.example <id>
```
