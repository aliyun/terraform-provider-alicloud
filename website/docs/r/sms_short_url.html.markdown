---
subcategory: "Short Message Service (SMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sms_short_url"
sidebar_current: "docs-alicloud-resource-sms-short-url"
description: |-
  Provides a Alicloud SMS Short Url resource.
---

# alicloud_sms_short_url

Provides a SMS Short Url resource.

For information about SMS Short Url and how to use it, see [What is Short Url](https://next.api.alibabacloud.com/api/Dysmsapi/2017-05-25/AddShortUrl).

-> **NOTE:** Available since v1.178.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_sms_short_url&exampleId=fc3f9e23-7e93-5f5c-0732-78dbe89fa072bccdfb79&activeTab=example&spm=docs.r.sms_short_url.0.fc3f9e237e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_sms_short_url" "example" {
  effective_days = 30
  short_url_name = "example_value"
  source_url     = "example_value"
}
```

## Argument Reference

The following arguments are supported:

* `effective_days` - (Required, ForceNew) Short chain service use validity period. Valid values: `30`, `60`, `90`. The unit is days, and the maximum validity period is 90 days.
* `short_url_name` - (Required, ForceNew) The name of the resource.
* `source_url` - (Required, ForceNew) The original link address.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Short Url.
* `status` - Short chain status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Short Url.
* `delete` - (Defaults to 1 mins) Used when delete the Short Url.

## Import

SMS Short Url can be imported using the id, e.g.

```shell
$ terraform import alicloud_sms_short_url.example <id>
```