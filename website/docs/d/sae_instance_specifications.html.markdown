---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_instance_specifications"
sidebar_current: "docs-alicloud-datasource-sae-instance-specifications"
description: |-
  Provides a list of Sae Instance Specifications to the user.
---

# alicloud\_sae\_instance\_specifications

This data source provides the Sae Instance Specifications of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.139.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_sae_instance_specifications" "ids" {}
output "sae_instance_specification_id_1" {
  value = data.alicloud_sae_instance_specifications.ids.specifications.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Instance Specification IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `specifications` - A list of Sae Instance Specifications. Each element contains the following attributes:
	* `cpu` - CPU Size, Specifications for Micronucleus.
	* `enable` - Whether the instance is available. The value description is as follows:
	  * `true` - indicates that it is available.
	  * `false` - means unavailable.
	* `id` - The ID of the Instance Specification.
	* `instance_specification_id` - The first ID of the resource.
	* `memory` - The Memory specifications for the MB.
	* `spec_info` - The specification configuration name.
	* `version` - The specification configuration version.
