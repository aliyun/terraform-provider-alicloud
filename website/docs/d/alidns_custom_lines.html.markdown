---
subcategory: "Alidns"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_custom_lines"
sidebar_current: "docs-alicloud-datasource-alidns-custom-lines"
description: |-
  Provides a list of Alidns Custom Lines to the user.
---

# alicloud_alidns_custom_lines

This data source provides the Alidns Custom Lines of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.151.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_alidns_custom_lines" "ids" {
  enable_details = true
  domain_name    = "your_domain_name"
}
output "alidns_custom_line_id_1" {
  value = data.alicloud_alidns_custom_lines.ids.lines.0.id
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, Computed)  A list of Custom Line IDs.
* `domain_name` - (Required)  The Domain name.
* `lang` - (Optional) The lang.
* `name_regex` - (Optional) A regex string to filter results by Custom Line name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Custom Line names.
* `lines` - A list of Alidns Custom Lines.

---

### `lines`

The lines supports the following attributes:

* `code` - The Custom line Code.
* `custom_line_id` - The first ID of the resource.
* `custom_line_name` - Line name.
* `domain_name` - The Domain name.
* `id` - The ID of the Custom Line.
* `ip_segment_list` - The IP segment list.

---

### `lines.ip_segment_list`

The ip_segment_list supports the following attributes:

* `start_ip` - The start IP address of the CIDR block.
* `end_ip` - The end IP address of the CIDR block.
