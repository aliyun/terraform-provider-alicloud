---
subcategory: "Service Catalog"
layout: "alicloud"
page_title: "Alicloud: alicloud_service_catalog_portfolio"
description: |-
  Provides a Alicloud Service Catalog Portfolio resource.
---

# alicloud_service_catalog_portfolio

Provides a Service Catalog Portfolio resource.

For information about Service Catalog Portfolio and how to use it, see [What is Portfolio](https://www.alibabacloud.com/help/en/service-catalog/developer-reference/api-servicecatalog-2021-09-01-createportfolio).

-> **NOTE:** Available since v1.204.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf_example"
}
resource "alicloud_service_catalog_portfolio" "default" {
  portfolio_name = var.name
  provider_name  = var.name
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) The description of the portfolio
* `portfolio_name` - (Required) The name of the portfolio
* `provider_name` - (Required) The provider name of the portfolio

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the portfolio
* `portfolio_arn` - The ARN of the portfolio

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Portfolio.
* `delete` - (Defaults to 5 mins) Used when delete the Portfolio.
* `update` - (Defaults to 5 mins) Used when update the Portfolio.

## Import

Service Catalog Portfolio can be imported using the id, e.g.

```shell
$ terraform import alicloud_service_catalog_portfolio.example <id>
```