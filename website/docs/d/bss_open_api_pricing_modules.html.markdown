---
subcategory: "Bss Open Api"
layout: "alicloud"
page_title: "Alicloud: alicloud_bss_open_api_pricing_modules"
sidebar_current: "docs-alicloud-datasource-bss-openapi-pricing-modules"
description: |-
  Provides a list of Bss Open Api Pricing Module owned by an Alibaba Cloud account.
---

# alicloud_bssopenapi_pricing_modules

This data source provides Bss Open Api Pricing Module available to the user.[What is Pricing Module](https://www.alibabacloud.com/help/en/bss-openapi/latest/describepricingmodule#doc-api-BssOpenApi-DescribePricingModule)

-> **NOTE:** Available in 1.195.0+

## Example Usage

```
data "alicloud_bss_open_api_pricing_modules" "default" {
  name_regex        = "国内月均日峰值带宽"
  product_code      = "cdn"
  product_type      = "CDN"
  subscription_type = "PayAsYouGo"
}

output "alicloud_bss_openapi_pricing_module_example_id" {
  value = data.alicloud_bss_open_api_pricing_modules.default.modules.0.code
}
```

## Argument Reference

The following arguments are supported:
* `product_code` - (Required,ForceNew) The product code. 
* `product_type` - (ForceNew,Optional) The product type. 
* `subscription_type` - (Required,ForceNew) Subscription type. Value:
  * Subscription: Prepaid.
  * PayAsYouGo: postpaid.
* `id` - (Optional, ForceNew) A list of Price Module IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Property name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `names` - A list of name of Pricing Modules.
* `modules` - A list of Pricing Module Entries. Each element contains the following attributes:
    * `code` - Property Code.
    * `pricing_module_name` - Attribute name.
    * `unit` - Attribute unit.
    * `values` - Property.
        * `name` - The module Code corresponds to the attribute value.
        * `remark` - Module value description information.
        * `type` - The attribute value type corresponding to the module Code. Value:
          * single_float: single value type.
          * range_float: range value type.
        * `value` - The module Code corresponds to the attribute value.
          > format 1024-1024000 when Type = range_float: 1024 means from 1024 to 1024000, step size 1024.
