---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_audit_policies"
sidebar_current: "docs-alicloud-datasource-mongodb-audit-policies"
description: |-
  Provides a list of Mongodb Audit Policies to the user.
---

# alicloud\_mongodb\_audit\_policies

This data source provides the Mongodb Audit Policies of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.148.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_mongodb_audit_policies" "example" {
  db_instance_id = "example_value"
}
output "mongodb_audit_policy_id_1" {
  value = data.alicloud_mongodb_audit_policies.example.policies.0.id
}

```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Request, ForceNew) The id of the db instance.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `policies` - A list of Mongodb Audit Policies. Each element contains the following attributes:
	* `db_instance_id` - The ID of the instance.
	* `id` - The ID of the Audit Policy.
	* `audit_status` - The status of the log audit feature.