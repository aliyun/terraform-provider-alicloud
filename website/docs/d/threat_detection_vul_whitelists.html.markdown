---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_vul_whitelists"
sidebar_current: "docs-alicloud-datasource-threat-detection-vul-whitelists"
description: |-
  Provides a list of Threat Detection Vul Whitelists to the user.
---

# alicloud\_threat\_detection\_vul\_whitelists

This data source provides Threat Detection Vul Whitelists of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.195.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_threat_detection_vul_whitelists" "default" {
  ids = ["example_id"]
}

output "alicloud_threat_detection_vul_whitelist_example_id" {
  value = data.alicloud_threat_detection_vul_whitelists.default.whitelists.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew, Computed) A list of Threat Detection Vul Whitelist IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `whitelists` - A list of Vul Whitelist Entries. Each element contains the following attributes:
  * `id` - The ID of the Vul Whitelist.
  * `vul_whitelist_id` - The ID of the Vul Whitelist.
  * `whitelist` - Information about the vulnerability to be added to the whitelist.
  * `target_info` - Set the effective range of the whitelist.
  * `reason` - Reason for adding whitelist.
  