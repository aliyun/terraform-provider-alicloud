---
subcategory: "Certificate Management Service (Original SSL Certificate)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_certificates_service_companies"
sidebar_current: "docs-alicloud-datasource-ssl-certificates-service-companies"
description: |-
  Provides a list of SSL Certificate Service Companies owned by an Alibaba Cloud account.
---

# alicloud_ssl_certificates_service_companies

This data source provides SSL Certificate Service Companies available to the user.[What is Company](https://next.api.alibabacloud.com/document/cas/2020-04-07/CreateCompany)

-> **NOTE:** Available since v1.285.0.

## Example Usage

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

data "alicloud_ssl_certificates_service_companies" "default" {
  ids        = ["${alicloud_ssl_certificates_service_company.default.id}"]
  name_regex = "${alicloud_ssl_certificates_service_company.default.company_name}"
}

output "alicloud_ssl_certificates_service_company_example_id" {
  value = "${data.alicloud_ssl_certificates_service_companies.default.companies.0.id}"
}
```

## Argument Reference

The following arguments are supported:
* `company_id` - (Optional, Int) The ID of the company used to filter the results.
* `keyword` - (Optional) The keyword used to filter the companies.
* `ids` - (Optional, Computed) A list of Company IDs.
* `name_regex` - (Optional) A regex string to filter results by Company name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Company IDs.
* `names` - A list of name of Companies.
* `companies` - A list of Company Entries. Each element contains the following attributes:
    * `city` - The city where the company is located.
    * `company_address` - The address of the company.
    * `company_code` - The code of the company.
    * `company_email` - The email address of the company.
    * `company_id` - The ID of the company.
    * `company_name` - The name of the company.
    * `company_phone` - The contact phone number of the company.
    * `company_type` - The type of the company.
    * `country_code` - The country code of the company.
    * `department` - The department of the company.
    * `lang` - The natural language of the content within the request and response.
    * `post_code` - The postal code of the company.
    * `province` - The province where the company is located.
    * `id` - The ID of the Company.
