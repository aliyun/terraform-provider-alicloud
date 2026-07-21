---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_vpc_endpoint_linked_vpcs"
sidebar_current: "docs-alicloud-datasource-cr-vpc-endpoint-linked-vpcs"
description: |-
  Provides a list of CR Vpc Endpoint Linked Vpcs to the user.
---

# alicloud\_cr\_vpc\_endpoint\_linked\_vpcs

This data source provides the CR Vpc Endpoint Linked Vpcs of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.199.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cr_vpc_endpoint_linked_vpcs" "ids" {
  ids         = ["example_id"]
  instance_id = "your_cr_instance_id"
  module_name = "Registry"
}

output "alicloud_cr_vpc_endpoint_linked_vpcs_id_1" {
  value = data.alicloud_cr_vpc_endpoint_linked_vpcs.ids.vpc_endpoint_linked_vpcs.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of CR Vpc Endpoint Linked Vpc IDs.
* `instance_id` - (Required, ForceNew) The ID of the instance.
* `module_name` - (Required, ForceNew) The name of the module that you want to access. Valid Values:
  - `Registry`: the image repository.
  - `Chart`: a Helm chart.
* `status` - (Optional, ForceNew) The status of the Vpc Endpoint Linked Vpc. Valid Values: `CREATING`, `RUNNING`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `vpc_endpoint_linked_vpcs` - A list of CR Vpc Endpoint Linked Vpcs. Each element contains the following attributes:
  * `id` - The ID of the Vpc Endpoint Linked Vpc. It formats as `<instance_id>:<vpc_id>:<vswitch_id>:<module_name>`.
  * `instance_id` - The ID of the instance.
  * `vpc_id` - The ID of the VPC.
  * `vswitch_id` - The ID of the vSwitch.
  * `module_name` - The name of the module that you want to access.
  * `ip` - IP address.
  * `default_access` - Indicates whether the default policy is used to access the instance.
  * `status` - The status of the Vpc Endpoint Linked Vpc.
  