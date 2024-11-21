---
subcategory: "Service Catalog"
layout: "alicloud"
page_title: "Alicloud: alicloud_service_catalog_product_version"
description: |-
  Provides a Alicloud Service Catalog Product Version resource.
---

# alicloud_service_catalog_product_version

Provides a Service Catalog Product Version resource.

There can be one or more versions of the product.

For information about Service Catalog Product Version and how to use it, see [What is Product Version](https://www.alibabacloud.com/help/en/service-catalog/developer-reference/api-servicecatalog-2021-09-01-createproductversion).

-> **NOTE:** Available since v1.230.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_service_catalog_product_version&exampleId=84ae7248-ebf0-ef8c-57cf-5206f4d3ccf30faf47df&activeTab=example&spm=docs.r.service_catalog_product_version.0.84ae7248eb&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_service_catalog_product" "defaultmaeTcE" {
  provider_name = var.name
  product_name  = var.name
  product_type  = "Ros"
}


resource "alicloud_service_catalog_product_version" "default" {
  guidance             = "Default"
  template_url         = "oss://servicecatalog-cn-hangzhou/1466115886172051/terraform/template/tpl-bp1x4v3r44u7u7/template.json"
  active               = true
  description          = "产品版本测试"
  product_version_name = var.name
  product_id           = alicloud_service_catalog_product.defaultmaeTcE.id
  template_type        = "RosTerraformTemplate"
}
```

## Argument Reference

The following arguments are supported:
* `active` - (Optional) Whether the version is activated
* `description` - (Optional) Version description
* `guidance` - (Optional) Administrator guidance
* `product_id` - (Required, ForceNew) Product ID
* `product_version_name` - (Required) The name of the resource
* `template_type` - (Required, ForceNew) Template Type
* `template_url` - (Required, ForceNew) Template URL

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Product Version.
* `delete` - (Defaults to 5 mins) Used when delete the Product Version.
* `update` - (Defaults to 5 mins) Used when update the Product Version.

## Import

Service Catalog Product Version can be imported using the id, e.g.

```shell
$ terraform import alicloud_service_catalog_product_version.example <id>
```