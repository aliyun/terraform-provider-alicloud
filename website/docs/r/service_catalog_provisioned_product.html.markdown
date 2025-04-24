---
subcategory: "Service Catalog"
layout: "alicloud"
page_title: "Alicloud: alicloud_service_catalog_provisioned_product"
sidebar_current: "docs-alicloud-resource-service-catalog-provisioned-product"
description: |-
  Provides an Alicloud Service Catalog Provisioned Product resource.
---

# alicloud_service_catalog_provisioned_product

Provides a Service Catalog Provisioned Product resource.

For information about Service Catalog Provisioned Product and how to use it, see [What is Provisioned Product](https://www.alibabacloud.com/help/en/service-catalog/developer-reference/api-servicecatalog-2021-09-01-launchproduct).

-> **NOTE:** Available in v1.196.0+.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_service_catalog_provisioned_product&exampleId=0409ca43-bb22-2c6d-b67d-d8f9501ec5b887d143d2&activeTab=example&spm=docs.r.service_catalog_provisioned_product.0.0409ca43bb&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-testAccServiceCatalogProvisionedProduct"
}

resource "alicloud_service_catalog_provisioned_product" "default" {
  provisioned_product_name = var.name
  stack_region_id          = "cn-hangzhou"
  product_version_id       = "pv-bp1d7dxy2pcc1g"
  product_id               = "prod-bp1u3dkc282cwd"
  portfolio_id             = "port-bp119dvn27jccw"
  tags = {
    "v1" = "tf-test"
  }
  parameters {
    parameter_key   = "role_name"
    parameter_value = var.name
  }
}
```

## Argument Reference

The following arguments are supported:
* `parameters` - (Optional) Template parameters entered by the user.The maximum value of N is 200.See the following `Block Parameters`.
* `portfolio_id` - (Optional) Product mix ID.> When there is a default Startup option, there is no need to fill in the portfolio. When there is no default Startup option, you must fill in the portfolio. 
* `product_id` - (Required) Product ID.
* `product_version_id` - (Required) Product version ID.
* `provisioned_product_name` - (Required,ForceNew) The name of the instance.The length is 1~128 characters.
* `stack_region_id` - (Required,ForceNew) The ID of the region to which the resource stack of the Alibaba Cloud resource orchestration service (ROS) belongs.
* `tags` - (Optional) A mapping of tags to assign to the resource.

#### Block Parameters

The Parameters support the following:
* `parameter_key` - (Optional) The name of the parameter defined in the template.
* `parameter_value` - (Optional) The Template parameter value entered by the user.


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `create_time` - The creation time of the product instance
* `last_provisioning_task_id` - The ID of the last instance operation task
* `last_successful_provisioning_task_id` - The ID of the last successful instance operation task
* `last_task_id` - The ID of the last task
* `outputs` - The output value of the template.
    * `description` - Description of the output value defined in the template.
    * `output_key` - The name of the output value defined in the template.
    * `output_value` - The content of the output value defined in the template.
* `owner_principal_id` - The RAM entity ID of the owner
* `owner_principal_type` - The RAM entity type of the owner
* `product_name` - The name of the product
* `product_version_name` - The name of the product version
* `provisioned_product_arn` - The ARN of the product instance
* `provisioned_product_id` - The ID of the instance.
* `provisioned_product_type` - Instance type.The value is RosStack, which indicates the stack of Alibaba Cloud resource orchestration service (ROS).
* `stack_id` - The ID of the ROS stack
* `status` - Instance status
* `status_message` - The status message of the product instance

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 3 mins) Used when create the Provisioned Product.
* `delete` - (Defaults to 24 mins) Used when delete the Provisioned Product.
* `update` - (Defaults to 24 mins) Used when update the Provisioned Product.

## Import

Service Catalog Provisioned Product can be imported using the id, e.g.

```shell
$terraform import alicloud_service_catalog_provisioned_product.example <id>
```