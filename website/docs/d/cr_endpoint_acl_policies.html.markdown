---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_endpoint_acl_policies"
sidebar_current: "docs-alicloud-datasource-cr-endpoint-acl-policies"
description: |-
  Provides a list of Cr Endpoint Acl Policies to the user.
---

# alicloud\_cr\_endpoint\_acl\_policies

This data source provides the Cr Endpoint Acl Policies of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.139.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cr_endpoint_acl_policies" "ids" {
  instance_id   = "example_value"
  endpoint_type = "example_value"
  ids           = ["example_value-1", "example_value-2"]
}
output "cr_endpoint_acl_policy_id_1" {
  value = data.alicloud_cr_endpoint_acl_policies.ids.policies.0.id
}
```

## Argument Reference

The following arguments are supported:

* `endpoint_type` - (Required, ForceNew)  The type of endpoint. Valid values: `internet`.
* `ids` - (Optional, ForceNew, Computed)  A list of Endpoint Acl Policy IDs.
* `instance_id` - (Required, ForceNew)  The ID of the CR Instance.
* `module_name` - (Optional, ForceNew) The module that needs to set the access policy. Valid values: `Registry`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `policies` - A list of Cr Endpoint Acl Policies. Each element contains the following attributes:
	* `description` - The description of the entry.
	* `endpoint_type` - The type of endpoint.
	* `entry` - The IP segment that allowed to access.
	* `id` - The ID of the Endpoint Acl Policy.
	* `instance_id` - The ID of the CR Instance.
