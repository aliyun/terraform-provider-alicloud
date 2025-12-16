---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_page"
description: |-
  Provides a Alicloud ESA Page resource.
---

# alicloud_esa_page

Provides a ESA Page resource.



For information about ESA Page and how to use it, see [What is Page](https://www.alibabacloud.com/help/en/edge-security-acceleration/esa/user-guide/customize-page).

-> **NOTE:** Available since v1.242.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_page&exampleId=0c00ee46-f661-77ad-8028-58c5142ad4abbe8a1442&activeTab=example&spm=docs.r.esa_page.0.0c00ee46f6&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "terraform-example"
}

resource "alicloud_esa_page" "default" {
  description  = "example resource html page"
  content_type = "text/html"
  content      = "PCFET0NUWVBFIGh0bWw+CjxodG1sIGxhbmc9InpoLUNOIj4KICA8aGVhZD4KICAgIDx0aXRsZT40MDMgRm9yYmlkZGVuPC90aXRsZT4KICA8L2hlYWQ+CiAgPGJvZHk+CiAgICA8aDE+NDAzIEZvcmJpZGRlbjwvaDE+CiAgPC9ib2R5Pgo8L2h0bWw+"
  page_name    = "resource_example_html_page"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_esa_page&spm=docs.r.esa_page.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `content` - (Optional) The Base64-encoded content of the error page. The content type is specified by the Content-Type field.
* `content_type` - (Required) The Content-Type field in the HTTP header.
* `description` - (Optional) The description of the custom error page.
* `page_name` - (Required) The name of the custom response page.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Page.
* `delete` - (Defaults to 5 mins) Used when delete the Page.
* `update` - (Defaults to 5 mins) Used when update the Page.

## Import

ESA Page can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_page.example <id>
```