---
subcategory: "GWLB"
layout: "alicloud"
page_title: "Alicloud: alicloud_gwlb_zones"
sidebar_current: "docs-alicloud-datasource-gwlb-zones"
description: |-
  Provides a list of Gwlb Zone owned by an Alibaba Cloud account.
---

# alicloud_gwlb_zones

This data source provides Gwlb Zone available to the user.[What is Zone](https://www.alibabacloud.com/help/en/)

-> **NOTE:** Available since v1.236.0.

## Example Usage

```terraform
provider "alicloud" {
  region = "cn-wulanchabu"
}

data "alicloud_gwlb_zones" "default" {
}

output "alicloud_gwlb_zone_example_id" {
  value = data.alicloud_gwlb_zones.default.zones.0.id
}
```

## Argument Reference

The following arguments are supported:
* `accept_language` - (ForceNew, Optional) The supported language. Valid values:
  - **zh-CN**: Chinese
  - **en-US** (default): English
  - **ja**: Japanese
* `ids` - (Optional, ForceNew, Computed) A list of Zone IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Zone IDs.
* `zones` - A list of Zone Entries. Each element contains the following attributes:
  * `local_name` - The zone name.
  * `zone_id` - The zone ID.
  * `id` - The zone ID.
