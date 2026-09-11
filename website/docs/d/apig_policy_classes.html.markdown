---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_policy_classes"
description: |-
  Provides a list of APIG Policy Class owned by an Alibaba Cloud account.
---

# alicloud_apig_policy_classes

This data source provides the APIG Policy Class of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.288.0.

## Example Usage

```terraform
data "alicloud_apig_policy_classes" "default" {
  type = "Auth"
}

output "apig_policy_class_id_0" {
  value = data.alicloud_apig_policy_classes.default.classes.0.id
}
```

## Argument Reference

The following arguments are supported:
* `type` - (Optional) The type of the policy class used to filter results. Valid values: `Auth`, `FlowControl`, `FlowObservation`, `Security`, `TransportProtocol`.
* `direction` - (Optional) The flow direction used to filter results. Valid values: `Inbound`, `OutBound`, `Both`.
* `attach_resource_type` - (Optional) The attachable resource type used to filter results. Valid values: `HttpApi`, `Operation`, `GatewayRoute`, `Gateway`, `GatewayDomain`.
* `ids` - (Optional, List) A list of Policy Class IDs.
* `name_regex` - (Optional) A regex string to filter results by policy class name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `names` - A list of name of Policy Classes.
* `classes` - A list of Policy Class. Each element contains the following attributes:
  * `id` - The ID of the Policy Class. It is the same as `policy_class_id`.
  * `policy_class_id` - The ID of the policy class.
  * `policy_class_name` - The name of the policy class.
  * `alias` - The alias of the policy class.
  * `version` - The version of the policy class.
  * `description` - The description of the policy class.
  * `type` - The type of the policy class. Valid values: `Auth` (authentication and authorization), `FlowControl` (traffic control), `FlowObservation` (traffic observation), `Security` (security protection), `TransportProtocol` (transport protocol).
  * `direction` - The flow direction of the policy. Valid values: `Inbound`, `OutBound`, `Both`.
  * `attachable_resource_types` - The list of resource types that the policy class can be attached to. Valid values: `HttpApi`, `Operation`, `GatewayRoute`, `Gateway`, `GatewayDomain`.
  * `execute_stage` - The execute stage of the policy. Valid values: `PRE_AUTHN`, `AUTHN`, `POST`, `UNSPECIFIED_PHASE`.
  * `execute_priority` - The execute priority of the policy.
  * `enable_log` - Whether the log is enabled.
  * `config_example` - The configuration example of the policy class.
  * `deprecated` - Whether the policy class is deprecated.
