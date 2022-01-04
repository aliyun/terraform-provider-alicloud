---
subcategory: "DNS"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_custom_line"
sidebar_current: "docs-alicloud-resource-alidns-custom-line"
description: |-
  Provides a Alicloud Alidns Custom Line resource.
---

# alicloud\_alidns\_custom\_line

Provides a Alidns Custom Line resource.

For information about Alidns Custom Line and how to use it, see [What is Custom Line](https://www.alibabacloud.com/help/en/doc-detail/145059.html).

-> **NOTE:** Available in v1.151.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_alidns_custom_line" "default" {
  custom_line_name = "tf-testacc"
  domain_name      = "your_domain_name"
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
* `ip_segment_list` - (Required) The IP segment list. See the following `Block ip_segment_list`.
* `lang` - (Optional) The lang.

### Block ip_segment_list

The ip_segment_list supports the following:

* `start_ip` - (Required) The start IP address of the CIDR block.
* `end_ip` - (Required) The end IP address of the CIDR block.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Custom Line.

## Import

Alidns Custom Line can be imported using the id, e.g.

```
$ terraform import alicloud_alidns_custom_line.example <id>
```