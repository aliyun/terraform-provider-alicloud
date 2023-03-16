---
subcategory: "Service Catalog"
layout: "alicloud"
page_title: "Alicloud: alicloud_service_catalog_portfolios"
sidebar_current: "docs-alicloud-datasource-service_catalog-portfolios"
description: |-
  Provides a list of Service Catalog Portfolio owned by an Alibaba Cloud account.
---

# alicloud_service_catalog_portfolios

This data source provides Service Catalog Portfolio available to the user.[What is Portfolio](https://www.alibabacloud.com/help/en/servicecatalog/latest/api-doc-servicecatalog-2021-09-01-api-doc-createportfolio)

-> **NOTE:** Available in 1.204.0+

## Example Usage

```terraform
data "alicloud_service_catalog_portfolios" "default" {
  ids        = ["${alicloud_service_catalog_portfolio.default.id}"]
  name_regex = alicloud_service_catalog_portfolio.default.name
}

output "alicloud_service_catalog_portfolio_example_id" {
  value = data.alicloud_service_catalog_portfolios.default.portfolios.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew, Computed) A list of Portfolio IDs.
* `portfolio_names` - (Optional, ForceNew) The name of the Portfolio. You can specify at most 10 names.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `scope` - (Optional, ForceNew) The query scope. Valid values: `Local`(default), `Import`, `All`.
* `sort_by` - (Optional, ForceNew) The field that is used to sort the queried data. The value is fixed as CreateTime, which specifies the creation time of product portfolios.
* `sort_order` - (Optional, ForceNew) The order in which you want to sort the queried data. Valid values: `Asc`, `Desc`.
* `product_id` - (Optional, ForceNew) The ID of the product.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Portfolio IDs.
* `names` - A list of name of Portfolios.
* `portfolios` - A list of Portfolio Entries. Each element contains the following attributes:
  * `create_time` - The creation time of the portfolio
  * `description` - The description of the portfolio
  * `portfolio_arn` - The ARN of the portfolio
  * `id` - The ID of the portfolio
  * `portfolio_id` - The ID of the portfolio
  * `portfolio_name` - The name of the portfolio
  * `provider_name` - The provider name of the portfolio
