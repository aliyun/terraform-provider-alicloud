---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_common_bandwidth_packages"
sidebar_current: "docs-alicloud-datasource-ens-common-bandwidth-packages"
description: |-
  Provides a list of ENS Common Bandwidth Packages to the user.
---

# alicloud_ens_common_bandwidth_packages

This data source provides the ENS Common Bandwidth Packages of the current Alibaba Cloud user.

For information about ENS Common Bandwidth Package, see [What is Common Bandwidth Package](https://next.api.alibabacloud.com/document/Ens/2017-11-10/DescribeCommonBandwidthPackages).

-> **NOTE:** Available since v1.289.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_ens_common_bandwidth_packages" "default" {
  ens_region_id = "cn-xxx-7fj20x"
  name          = "tf-example"
}

output "first_package_id" {
  value = data.alicloud_ens_common_bandwidth_packages.default.packages.0.bandwidth_package_id
}
```

## Argument Reference

The following arguments are supported:

* `bandwidth_package_id` - (Optional, ForceNew) The ID of the Common Bandwidth Package.
* `ens_region_id` - (Optional, ForceNew) The ID of the ENS node.
* `ids` - (Optional) A list of Common Bandwidth Package IDs.
* `name` - (Optional, ForceNew) The name of the Common Bandwidth Package.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported:

* `id` - The data source ID.
* `ids` - A list of Common Bandwidth Package IDs.
* `packages` - A list of Common Bandwidth Packages. Each element contains the following attributes:
  * `bandwidth` - Speed limit bandwidth value, unit: Mbps.
  * `bandwidth_package_id` - The ID of the Common Bandwidth Package.
  * `creation_time` - The creation time of the resource.
  * `description` - The description of the Common Bandwidth Package.
  * `ens_region_id` - The ID of the ENS node.
  * `id` - The ID of the Common Bandwidth Package.
  * `name` - The name of the Common Bandwidth Package.
  * `status` - The status of the resource.
