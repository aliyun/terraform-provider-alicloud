---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_grey_tag_route"
sidebar_current: "docs-alicloud-resource-sae-grey_tag_route"
description: |-
  Provides a Alicloud Serverless App Engine (SAE) GreyTagRoute resource.
---

# alicloud\_sae\_grey\_tag\_route

Provides a Serverless App Engine (SAE) GreyTagRoute resource.

For information about Serverless App Engine (SAE) GreyTagRoute and how to use it, see [What is GreyTagRoute](https://help.aliyun.com/document_detail/97792.html).

-> **NOTE:** Available in v1.160.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testacc"
}

variable "region" {
  default = "cn-hangzhou"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_sae_namespace" "default" {
  namespace_description = var.name
  namespace_id          = join(":", [var.region, var.name])
  namespace_name        = var.name
}

resource "alicloud_sae_application" "default" {
  app_description = var.name
  app_name        = var.name
  namespace_id    = alicloud_sae_namespace.default.namespace_id
  image_url       = "registry-vpc.cn-hangzhou.aliyuncs.com/lxepoo/apache-php5"
  package_type    = "Image"
  jdk             = "Open JDK 8"
  vswitch_id      = data.alicloud_vswitches.default.ids.0
  vpc_id          = data.alicloud_vpcs.default.ids.0
  timezone        = "Asia/Shanghai"
  replicas        = "5"
  cpu             = "500"
  memory          = "2048"
}

resource "alicloud_sae_grey_tag_route" "default" {
  grey_tag_route_name = var.name
  description         = var.name
  app_id              = alicloud_sae_application.default.id
  sc_rules {
    items {
      type     = "param"
      name     = "tftest"
      operator = "rawvalue"
      value    = "test"
      cond     = "=="
    }
    path      = "/tf/test"
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
    method_name  = "test"
    service_name = "com.test.service"
    version      = "1.0.0"
  }
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, ForceNew) The ID  of the SAE Application.
* `description` - (Optional) The description of GreyTagRoute.
* `grey_tag_route_name` - (Required, ForceNew) The name of GreyTagRoute.
* `dubbo_rules` - (Optional) The grayscale rule created for Dubbo Application. The details see Block `dubbo_rules`.
* `sc_rules` - (Optional) The grayscale rule created for SpringCloud Application. The details see Block `sc_rules`.

### dubbo_rules

The `dubbo_rules` supports the following:
* `method_name` - (Optional) The method name
* `service_name` - (Optional) The service name.
* `version` - (Optional) The service version.
* `condition` - (Optional) The Conditional Patterns for Grayscale Rules. Valid values: `AND`, `OR`.
* `group` - (Optional) The service group.
* `items` - (Optional) A list of conditions items. The details see Block `dubbo_rules_items`.

#### dubbo_rules_items

The `dubbo_rules_items` supports the following:
* `index` - (Optional) The parameter number.
* `expr` - (Optional) The parameter value gets the expression.
* `cond` - (Optional) The comparison operator. Valid values: `>`, `<`, `>=`, `<=`, `==`, `!=`.
* `value` - (Optional) The value of the parameter.
* `operator` - (Optional) The operator. Valid values: `rawvalue`, `list`, `mod`, `deterministic_proportional_steaming_division`

### sc_rules

The `sc_rules` supports the following:
* `path` - (Optional) The path corresponding to the grayscale rule.
* `condition` - (Optional) The conditional Patterns for Grayscale Rules. Valid values: `AND`, `OR`.
* `items` - (Optional) A list of conditions items. The details see Block `sc_rules_items`.

#### sc_rules_items

The `sc_rules_items` supports the following:
* `name` - (Optional) The name of the parameter.
* `type` - (Optional) The compare types. Valid values: `param`, `cookie`, `header`.
* `cond` - (Optional) The comparison operator. Valid values: `>`, `<`, `>=`, `<=`, `==`, `!=`.
* `value` - (Optional) The value of the parameter.
* `operator` - (Optional) The operator. Valid values: `rawvalue`, `list`, `mod`, `deterministic_proportional_steaming_division`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of GreyTagRoute.


#### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the resource.
* `update` - (Defaults to 1 mins) Used when update the resource.
* `delete` - (Defaults to 1 mins) Used when delete the resource.

## Import

Serverless App Engine (SAE) GreyTagRoute can be imported using the id, e.g.

```
$ terraform import alicloud_sae_grey_tag_route.example <id>
```
