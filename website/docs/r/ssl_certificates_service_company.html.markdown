---
subcategory: "Certificate Management Service (Original SSL Certificate)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_certificates_service_company"
description: |-
  Provides an Alicloud Certificate Management Service (Original SSL Certificate) Company resource.
---

# alicloud_ssl_certificates_service_company

Provides a Certificate Management Service (Original SSL Certificate) Company resource.



For information about Certificate Management Service (Original SSL Certificate) Company and how to use it, see [What is Company](https://next.api.alibabacloud.com/document/cas/2020-04-07/CreateCompany).

-> **NOTE:** Available since v1.285.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ssl_certificates_service_company&exampleId=2d35fd73-1ce5-fa29-1c7b-04c7f03f615b86e193d9&activeTab=example&spm=docs.r.ssl_certificates_service_company.0.2d35fd731c&intl_lang=EN_US" target="_blank">
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


resource "alicloud_ssl_certificates_service_company" "default" {
  company_address = "西安市"
  company_name    = "example公司1"
  department      = "example部门1"
  city            = "西安"
  company_type    = "1"
  country_code    = "111122"
  post_code       = "11112233"
  company_code    = "12312311"
  company_phone   = "15101081174"
  province        = "陕西"
  lang            = "zh"
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ssl_certificates_service_company&spm=docs.r.ssl_certificates_service_company.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `city` - (Required) The city where the company is located.
* `company_address` - (Required) The address of the company.
* `company_code` - (Required) The code of the company.
* `company_email` - (Optional) The email address of the company.
* `company_name` - (Required) The name of the company.
* `company_phone` - (Required) The contact phone number of the company.
* `company_type` - (Required, Int) The type of the company.
* `country_code` - (Required) The country code of the company.
* `department` - (Optional) The department of the company.
* `lang` - (Required) The natural language of the content within the request and response.
* `post_code` - (Required) The postal code of the company.
* `province` - (Required) The province where the company is located.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Company.
* `delete` - (Defaults to 5 mins) Used when delete the Company.
* `update` - (Defaults to 5 mins) Used when update the Company.

## Import

Certificate Management Service (Original SSL Certificate) Company can be imported using the id, e.g.

```shell
$ terraform import alicloud_ssl_certificates_service_company.example <company_id>
```