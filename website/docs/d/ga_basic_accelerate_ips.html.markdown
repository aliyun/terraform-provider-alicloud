---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_basic_accelerate_ips"
sidebar_current: "docs-alicloud-datasource-ga-basic-accelerate-ips"
description: |-
  Provides a list of Global Accelerator (GA) Basic Accelerate IPs to the user.
---

# alicloud_ga_basic_accelerate_ips

This data source provides the Global Accelerator (GA) Basic Accelerate IPs of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.194.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_ga_basic_accelerate_ips" "ids" {
  ids       = ["example_id"]
  ip_set_id = "example_ip_set_id"
}

output "ga_basic_accelerate_ip_id_1" {
  value = data.alicloud_ga_basic_accelerate_ips.ids.ips.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Global Accelerator Basic Accelerate IP IDs.
* `ip_set_id` - (Required, ForceNew) The ID of the Basic Ip Set.
* `accelerate_ip_id` - (Optional, ForceNew) The id of the Basic Accelerate IP.
* `accelerate_ip_address` - (Optional, ForceNew) The address of the Basic Accelerate IP.
* `status` - (Optional, ForceNew) The status of the Global Accelerator Basic Accelerate IP instance. Valid Value: `active`, `binding`, `bound`, `unbinding`, `deleting`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ips` - A list of Global Accelerator Basic Accelerate IPs. Each element contains the following attributes:
  * `id` - The id of the Basic Accelerate IP.
  * `accelerate_ip_id` - The id of the Basic Accelerate IP.
  * `accelerate_ip_address` - The address of the Basic Accelerate IP.
  * `accelerator_id` - The id of the Global Accelerator Basic Accelerator instance.
  * `ip_set_id` - The ID of the Basic Ip Set.
  * `status` - The status of the Basic Accelerate IP instance.
	