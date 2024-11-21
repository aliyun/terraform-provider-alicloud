---
subcategory: "Alidns"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_custom_line"
sidebar_current: "docs-alicloud-resource-alidns-custom-line"
description: |-
  Provides a Alicloud Alidns Custom Line resource.
---

# alicloud_alidns_custom_line

Provides a Alidns Custom Line resource.

For information about Alidns Custom Line and how to use it, see [What is Custom Line](https://www.alibabacloud.com/help/en/doc-detail/145059.html).

-> **NOTE:** Available since v1.151.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alidns_custom_line&exampleId=63236334-d581-95cb-b9a5-d31c6d94d2a9b98995ae&activeTab=example&spm=docs.r.alidns_custom_line.0.63236334d5&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_alidns_custom_line" "default" {
  custom_line_name = "tf-example"
  domain_name      = "alicloud-provider.com"
  ip_segment_list {
    start_ip = "192.0.2.123"
    end_ip   = "192.0.2.125"
  }
}
```

## Argument Reference

The following arguments are supported:
* `custom_line_name` - (Required) The name of the Custom Line.
* `domain_name` - (Required, ForceNew) The Domain name.
* `ip_segment_list` - (Required) The IP segment list. See [`ip_segment_list`](#ip_segment_list) below for details.
* `lang` - (Optional) The lang.

### `ip_segment_list`

The ip_segment_list supports the following:

* `start_ip` - (Required) The start IP address of the CIDR block.
* `end_ip` - (Required) The end IP address of the CIDR block.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Custom Line.

## Import

Alidns Custom Line can be imported using the id, e.g.

```shell
$ terraform import alicloud_alidns_custom_line.example <id>
```