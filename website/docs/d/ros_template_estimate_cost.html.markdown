---
subcategory: "ROS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ros_template_estimate_cost"
sidebar_current: "docs-alicloud-datasource-ros-template-estimate-cost"
description: |-
  Queries the estimated prices of the resources created in a ROS or Terraform template.
---

# alicloud\_ros\_template\_estimate\_cost

This data source provides the estimated prices of the resources that are created in a ROS or Terraform template, which helps you evaluate resource costs before deployment.

For information about the resources that support price estimation, see [Estimated resource prices](https://www.alibabacloud.com/help/en/ros/user-guide/estimate-resource-prices).

-> **NOTE:** Available since v1.292.0.

## Example Usage

Estimate the cost of a Terraform template:

```terraform
data "alicloud_ros_template_estimate_cost" "example" {
  template_body = jsonencode({
    "ROSTemplateFormatVersion" = "2015-09-01"
    "Transform"                = "Aliyun::Terraform-v1.2"
    "Workspace" = {
      "main.tf" = <<-EOT
        provider "alicloud" {
          region = "cn-hangzhou"
        }
        resource "alicloud_eip_address" "example" {
          bandwidth            = 5
          internet_charge_type = "PayByTraffic"
          payment_type         = "Subscription"
        }
      EOT
    }
  })
}

output "estimated_resources" {
  value = data.alicloud_ros_template_estimate_cost.example.resources
}
```

Estimate the cost of a template that is stored in ROS:

```terraform
data "alicloud_ros_template_estimate_cost" "example" {
  template_id      = "5ecd1e10-b0e9-4389-a565-e4c15efc****"
  template_version = "v1"
  parameters {
    parameter_key   = "InstanceType"
    parameter_value = "ecs.s6-c1m2.large"
  }
}
```

## Argument Reference

The following arguments are supported:

* `template_body` - (Optional) The structure that contains the template body. The template body must be 1 to 524,288 bytes in length. To estimate the cost of a Terraform template, wrap the `.tf` files in a ROS structure with `Transform` set to `Aliyun::Terraform-v1.2` and the files under `Workspace`. You must specify exactly one of `template_body`, `template_url`, `template_id`, and `template_scratch_id`.
* `template_url` - (Optional) The URL of the file that contains the template body. The URL must point to a template that is located on an HTTP or HTTPS web server or in an Object Storage Service (OSS) bucket, such as `oss://ros/stack-policy/demo`. The template can be up to 524,288 bytes in length. You must specify exactly one of `template_body`, `template_url`, `template_id`, and `template_scratch_id`.
* `template_id` - (Optional) The ID of the template. Shared templates and private templates are supported. You must specify exactly one of `template_body`, `template_url`, `template_id`, and `template_scratch_id`.
* `template_version` - (Optional) The version of the template. This parameter takes effect only when `template_id` is specified.
* `template_scratch_id` - (Optional) The ID of the resource scenario. You must specify exactly one of `template_body`, `template_url`, `template_id`, and `template_scratch_id`. **NOTE:** Resource scenarios of the `ResourceImport` type are not supported for cost estimation.
* `template_scratch_region_id` - (Optional) The region ID of the resource scenario. Default value: the region ID of the current request.
* `stack_id` - (Optional) The ID of the stack. When this parameter is specified, the estimated prices for a stack change (modification) scenario are queried. **NOTE:** If the new template has no priceable changes compared with the stack, an empty `resources` list is returned.
* `client_token` - (Optional) The client token that is used to ensure the idempotence of the request. It can be up to 64 characters in length and can contain letters, digits, hyphens (-), and underscores (_).
* `parameters` - (Optional) The parameters of the template. A maximum of 200 parameters can be specified. If you do not specify the names and values of the parameters that are defined in the template, the default values of the parameters are used. See [`parameters`](#parameters) below.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

### `parameters`

The parameters supports the following:

* `parameter_key` - (Required) The name of the parameter defined in the template.
* `parameter_value` - (Required) The value of the parameter.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `resources` - The estimation details of the resources in the template. Each element contains the following attributes:
  * `resource_name` - The logical name of the resource in the template.
  * `resource_type` - The type of the resource. For Terraform templates, it is the Terraform resource type, such as `alicloud_eip_address`. For ROS templates, it is the ROS resource type, such as `ALIYUN::VPC::EIP`.
  * `alias_type` - The ROS resource type that the Terraform resource type maps to, such as `ALIYUN::VPC::EIP`. It is returned only for Terraform templates.
  * `success` - Indicates whether the price estimation of the resource succeeds.
  * `result` - The estimation result of the resource. Each element contains the following attributes:
    * `inquiry_type` - The type of the price estimation. Valid values: `Buy` (estimation for new purchase), `ModificationBuy` (estimation for specification change).
    * `order` - The order information. Each element contains the following attributes:
      * `currency` - The currency unit. `CNY` applies to the China site (aliyun.com), and `USD` applies to the international site (alibabacloud.com).
      * `original_amount` - The original price.
      * `discount_amount` - The discount amount.
      * `trade_amount` - The final price, which equals the original price minus the discount.
      * `tax_amount` - The tax amount.
      * `total_cost_amount` - The total cost amount.
      * `stand_price` - The list price of the resource.
      * `rule_ids` - The IDs of the promotion rules that are applied.
    * `order_details` - The billing module details of the order. Each element contains the following attributes:
      * `module_code` - The code of the billing module.
      * `module_name` - The name of the billing module.
      * `currency` - The currency unit.
      * `original_amount` - The original price of the module.
      * `discount_amount` - The discount amount of the module.
      * `trade_amount` - The final price of the module.
    * `order_supplement` - The supplementary information of the order. Each element contains the following attributes:
      * `charge_type` - The billing method, such as `PrePaid` (subscription) or `PostPaid` (pay-as-you-go).
      * `period` - The billing duration.
      * `period_unit` - The unit of the billing duration for the subscription billing method. Valid values: `Year`, `Month`.
      * `price_type` - The type of the price.
      * `price_unit` - The unit of the price for the pay-as-you-go billing method, such as `/Hour` or `/GB`.
      * `quantity` - The quantity.
      * `auto_renew` - Indicates whether auto-renewal is enabled.
