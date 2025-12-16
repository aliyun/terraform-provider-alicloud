---
subcategory: "Quotas"
layout: "alicloud"
page_title: "Alicloud: alicloud_quotas_template_service"
description: |-
  Provides a Alicloud Quotas Template Service resource.
---

# alicloud_quotas_template_service

Provides a Quotas Template Service resource.

Quota Template Service.

For information about Quotas Template Service and how to use it, see [What is Template Service](https://www.alibabacloud.com/help/en/quota-center/developer-reference/api-quotas-2020-05-10-modifyquotatemplateservicestatus).

-> **NOTE:** Available since v1.230.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_quotas_template_service&exampleId=2e0ba0ec-4da6-c0bb-2ca3-88940dee1376b5d1ca2c&activeTab=example&spm=docs.r.quotas_template_service.0.2e0ba0ec4d&intl_lang=EN_US" target="_blank">
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


resource "alicloud_quotas_template_service" "default" {
  service_status = "1"
}
```

### Deleting `alicloud_quotas_template_service` or removing it from your configuration

Terraform cannot destroy resource `alicloud_quotas_template_service`. Terraform will remove this resource from the state file, however resources may remain.

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_quotas_template_service&spm=docs.r.quotas_template_service.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `service_status` - (Required, Int) Status of the quota template. Valid values:
  - `-1`: disabled.
  - `1`: enabled.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as ``.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Template Service.
* `update` - (Defaults to 5 mins) Used when update the Template Service.

## Import

Quotas Template Service can be imported using the id, e.g.

```shell
$ terraform import alicloud_quotas_template_service.example 
```