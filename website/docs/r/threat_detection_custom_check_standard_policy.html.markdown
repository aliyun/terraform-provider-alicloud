---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_custom_check_standard_policy"
sidebar_current: "docs-alicloud-resource-threat-detection-custom-check-standard-policy"
description: |-
  Provides a Alicloud Threat Detection Custom Check Standard Policy resource.
---

# alicloud_threat_detection_custom_check_standard_policy

Provides a Threat Detection Custom Check Standard Policy resource.

Custom check standard policy is a classification policy that groups check items under a standard, requirement, or section in the Security Center (Threat Detection) Cspm module.

For information about Threat Detection Custom Check Standard Policy and how to use it, see [CreateCheckPolicy](https://www.alibabacloud.com/help/en/security-center/latest/api-sas-2018-12-03-createcheckpolicy).

-> **NOTE:** Available since v1.244.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_threat_detection_custom_check_standard_policy" "section" {
  policy_show_name = "tf-example-section"
  policy_type      = "SECTION"
}

resource "alicloud_threat_detection_custom_check_standard_policy" "requirement" {
  policy_show_name    = "tf-example-requirement"
  policy_type         = "REQUIREMENT"
  dependent_policy_id = alicloud_threat_detection_custom_check_standard_policy.section.policy_id
}

resource "alicloud_threat_detection_custom_check_standard_policy" "default" {
  policy_show_name    = "tf-example-standard"
  policy_type         = "STANDARD"
  dependent_policy_id = alicloud_threat_detection_custom_check_standard_policy.requirement.policy_id
  type                = "AISPM"
}
```

## Argument Reference

The following arguments are supported:

* `policy_show_name` - (Required) The name of the custom policy.
* `policy_type` - (Required, ForceNew) The policy category type for custom check rules. Valid values: `STANDARD`, `REQUIREMENT`, `SECTION`.
  * **STANDARD**: Add to a standard.
  * **REQUIREMENT**: Add to a requirement.
  * **SECTION**: Add to a section.
* `dependent_policy_id` - (Optional, ForceNew) The ID of the parent policy. The dependency order from low to high is: Section -> Requirement -> Standard. For a `SECTION` policy, this field is not required.
* `type` - (Optional) The name of the major policy category. Required when `policy_type` is `STANDARD`. Valid values: `AISPM`, `IDENTITY_PERMISSION`, `RISK`, `COMPLIANCE`.
  * **AISPM**: AI Configuration Management (AISPM).
  * **IDENTITY_PERMISSION**: Identity and Permission Management (CIEM).
  * **RISK**: Security Risk.
  * **COMPLIANCE**: Compliance Risk.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in Terraform, formatted as `<policy_id>:<policy_type>`.
* `policy_id` - The ID of the custom policy.
* `check_type` - The check type of the policy.

## Import

Threat Detection Custom Check Standard Policy can be imported using the id with format `<policy_id>:<policy_type>`, e.g.

```shell
$ terraform import alicloud_threat_detection_custom_check_standard_policy.example 123:STANDARD
```
