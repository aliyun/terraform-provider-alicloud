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