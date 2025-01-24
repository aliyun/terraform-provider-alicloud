---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_page"
description: |-
  Provides a Alicloud ESA Page resource.
---

# alicloud_esa_page

Provides a ESA Page resource.



For information about ESA Page and how to use it, see [What is Page](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.242.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_page&exampleId=7df0db9d-9e50-69bf-52c0-dc0a26b397a52a8cfb71&activeTab=example&spm=docs.r.esa_page.0.7df0db9d9e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
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

## Argument Reference

The following arguments are supported:
* `content` - (Optional) The Content-Type field in the HTTP header. Valid values:

  - text/html
  - application/json
* `content_type` - (Required) The description of the custom error page.
* `description` - (Optional) The name of the custom error page.
* `page_name` - (Required) The ID of the custom error page, which can be obtained by calling the [ListPages](https://www.alibabacloud.com/help/en/doc-detail/2850223.html) operation.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Page.
* `delete` - (Defaults to 5 mins) Used when delete the Page.
* `update` - (Defaults to 5 mins) Used when update the Page.

## Import

ESA Page can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_page.example <id>
```