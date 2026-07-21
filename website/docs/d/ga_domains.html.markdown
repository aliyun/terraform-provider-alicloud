---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_domains"
sidebar_current: "docs-alicloud-datasource-ga-domains"
description: |-
  Provides a list of Ga Domain owned by an Alibaba Cloud account.
---

# alicloud_ga_domains

This data source provides Ga Domain available to the user.[What is Domain](https://www.alibabacloud.com/help/en/global-accelerator/latest/createdomain)

-> **NOTE:** Available since v1.197.0.

## Example Usage

```terraform
data "alicloud_ga_accelerators" "default" {
  status = "active"
}

resource "alicloud_ga_accelerator" "default" {
  count           = length(data.alicloud_ga_accelerators.default.accelerators) > 0 ? 0 : 1
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}

locals {
  accelerator_id = length(data.alicloud_ga_accelerators.default.accelerators) > 0 ? data.alicloud_ga_accelerators.default.accelerators.0.id : alicloud_ga_accelerator.default.0.id
}
data "alicloud_ga_domains" "default" {
  accelerator_id = locals.accelerator_id
  domain         = "your_domain"
}

output "alicloud_ga_domain_example_id" {
  value = data.alicloud_ga_domains.default.domains.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew, Computed) A list of Ga Domain IDs.
* `accelerator_id` - (Optional, ForceNew) The ID of the global acceleration instance.
* `domain` - (Optional, ForceNew) The accelerated domain name to be added. only top-level domain names are supported, such as 'example.com'.
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `illegal`, `inactive`, `active`, `unknown`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `domains` - A list of Domain Entries. Each element contains the following attributes:
  * `id` - The ID of the Ga Domain.
  * `accelerator_id` - The ID of the global acceleration instance.
  * `domain` - The accelerated domain name to be added.
  * `status` - The status of the resource
