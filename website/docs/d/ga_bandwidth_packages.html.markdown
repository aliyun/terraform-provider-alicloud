---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_bandwidth_packages"
sidebar_current: "docs-alicloud-datasource-ga-bandwidth-packages"
description: |-
  Provides a list of Global Accelerator (GA) Bandwidth Packages to the user.
---

# alicloud\_ga\_bandwidth\_packages

This data source provides the Global Accelerator (GA) Bandwidth Packages of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.112.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ga_bandwidth_packages" "example" {
  ids        = ["example_value"]
  name_regex = "the_resource_name"
}

output "first_ga_bandwidth_package_id" {
  value = data.alicloud_ga_bandwidth_packages.example.packages.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Bandwidth Package IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Bandwidth Package name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the bandwidth plan. Valid values: `active`, `binded`, `binding`, `finacialLocked`, `init`, `unbinding`, `updating`.
* `type` - (Optional, ForceNew) The type of the bandwidth plan. Valid values: `Basic`, `CrossDomain`.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Bandwidth Package names.
* `packages` - A list of Ga Bandwidth Packages. Each element contains the following attributes:
	* `bandwidth` - The bandwidth value of bandwidth packet.
	* `bandwidth_package_id` - The Resource ID of the bandwidth.
	* `bandwidth_package_name` - The name of the bandwidth packet.
	* `bandwidth_type` - The bandwidth type of the bandwidth.
	* `cbn_geographic_region_ida` - Interworking area A of cross domain acceleration package. Only international stations support returning this parameter.
	* `cbn_geographic_region_idb` - Interworking area B of cross domain acceleration package. Only international stations support returning this parameter.
	* `description` - The description of bandwidth package.
	* `expired_time` - Bandwidth package expiration time.
	* `id` - The ID of the Bandwidth Package.
	* `payment_type` - The payment type of the bandwidth.
	* `status` - The status of the bandwidth plan.
	* `type` - The type of the bandwidth packet. China station only supports return to basic.
