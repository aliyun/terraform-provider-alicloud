---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_check_structures"
sidebar_current: "docs-alicloud-datasource-threat-detection-check-structures"
description: |-
  Provides a list of Threat Detection Check Structure owned by an Alibaba Cloud account.
---

# alicloud_threat_detection_check_structures

This data source provides Threat Detection Check Structure available to the user.[What is Check Structure](https://next.api.alibabacloud.com/document/Sas/2018-12-03/GetCheckStructure)

-> **NOTE:** Available since v1.267.0.

## Example Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_threat_detection_check_structures" "default" {

}

output "alicloud_threat_detection_check_structure_example_standard_type" {
  value = data.alicloud_threat_detection_check_structures.default.structures.0.standard_type
}
```

## Argument Reference

The following arguments are supported:
* `current_page` - (ForceNew, Optional) The page number.
* `lang` - (ForceNew, Optional) The language of the content within the request and response. Default value: zh. Valid values:- **zh**: Chinese- **en**: English
* `task_sources` - (ForceNew, Optional) List of task sources.
* `ids` - (Optional, ForceNew, Computed) A list of Check Structure IDs.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Check Structure IDs.
* `structures` - A list of Check Structure Entries. Each element contains the following attributes:
  * `standard_type` - The type of the check item.- **RISK**: security risk.- **IDENTITY_PERMISSION**: Cloud Infrastructure Entitlement Management (CIEM).- **COMPLIANCE**: security compliance.
  * `standards` - The structure information about the check items of the business type.
    * `id` - The standard ID of the check item.
    * `requirements` - The standards of the check items.
      * `id` - The ID of the requirement item for the check item.
      * `sections` - The information about the sections of check items.
        * `id` - The ID of the section for the check item.
        * `show_name` - The display name of the section for the check item.
      * `show_name` - The display name of the requirement item for the check item.
      * `total_check_count` - The total number of check items for the requirement.
    * `show_name` - The display name of the standard for the check item.
    * `type` - The standard type of the check item. Valid values:- **RISK**: security risk.- **IDENTITY_PERMISSION**: CIEM.- **COMPLIANCE**: security compliance.
