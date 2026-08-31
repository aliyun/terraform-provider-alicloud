---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_audit_policies"
sidebar_current: "docs-alicloud-datasource-mongodb-audit-policies"
description: |-
  Provides a list of Mongodb Audit Policies to the user.
---

# alicloud_mongodb_audit_policies

This data source provides the Mongodb Audit Policies of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.148.0.

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

* `db_instance_id` - (Required, ForceNew) The id of the db instance.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `policies` - A list of Mongodb Audit Policies. Each element contains the following attributes:
	* `db_instance_id` - The ID of the instance.
	* `id` - The ID of the Audit Policy.
	* `audit_status` - The status of the log audit feature.
	* `service_type` - (Available since v1.284.0) The edition of the audit log, e.g. `Standard` or `V2_Standard`.
	* `storage_period` - (Available since v1.284.0) The audit log retention duration, in days. For `V2_Standard` this is the cold storage duration.
	* `hot_storage_period` - (Available since v1.284.0) The hot storage duration (days) of the V2 audit log.