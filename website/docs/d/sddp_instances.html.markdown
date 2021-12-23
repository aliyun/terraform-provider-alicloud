---
subcategory: "Data Security Center"
layout: "alicloud"
page_title: "Alicloud: alicloud_sddp_instances"
sidebar_current: "docs-alicloud-datasource-sddp-instances"
description: |-
  Provides a list of Sddp Instances to the user.
---

# alicloud\_sddp\_instances

This data source provides the Sddp Instances of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.136.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_sddp_instances" "default" {}
output "sddp_instance_id" {
  value = data.alicloud_sddp_instances.default.instances.0
}

```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `instances` - A list of Sddp Instances. Each element contains the following attributes:
	* `authed` - Whether the required RAM authorization is configured.
	* `instance_id` - The ID of the instance.
	* `instance_num` - The number of instances.
	* `odps_set` - Whether the authorized MaxCompute (ODPS) assets.
	* `oss_bucket_set` - Whether the authorized oss assets.
	* `oss_size` - The OSS size of the instance.
	* `payment_type` - The payment type of the resource. Valid values: `Subscription`.
	* `rds_set` - Whether the authorized rds assets.
	* `status` - The status of the resource.
