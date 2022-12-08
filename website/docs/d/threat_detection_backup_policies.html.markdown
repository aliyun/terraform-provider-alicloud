---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_backup_policies"
sidebar_current: "docs-alicloud-datasource-threat-detection-backup-policies"
description: |-
  Provides a list of Threat Detection Backup Policies to the user.
---

# alicloud\_threat\_detection\_backup\_policies

This data source provides the Threat Detection Backup Policies of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.195.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_threat_detection_backup_policies" "ids" {
  ids = ["example_id"]
}

output "threat_detection_backup_policies_id_1" {
  value = data.alicloud_threat_detection_backup_policies.ids.policies.0.id
}

data "alicloud_threat_detection_backup_policies" "nameRegex" {
  name_regex = "tf-example"
}

output "threat_detection_backup_policies_id_2" {
  value = data.alicloud_threat_detection_backup_policies.nameRegex.policies.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Threat Detection Backup Policies IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Threat Detection Backup Policies name.
* `name` - (Optional, ForceNew) The name of the anti-ransomware policy that you want to query.
* `machine_remark` - (Optional, ForceNew) The information that you want to use to identify the servers protected by the anti-ransomware policy. You can enter the IP address or ID of a server.
* `status` - (Optional, ForceNew) The status of the anti-ransomware policy. Valid Value: `enabled`, `disabled`, `closed`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Threat Detection Backup Policy names.
* `policies` - A list of Threat Detection Backup policies. Each element contains the following attributes:
	* `id` - The ID of the anti-ransomware policy.
	* `backup_policy_id` - The ID of the anti-ransomware policy.
	* `backup_policy_name` - The name of the anti-ransomware policy.
	* `policy` - The configurations of the anti-ransomware policy.
	* `policy_region_id` - The ID of the region that you specified for data backup when you installed the anti-ransomware agent for the server not deployed on Alibaba Cloud.
	* `policy_version` - The version of the anti-ransomware policy.
	* `uuid_list` - The UUIDs of the servers to which the anti-ransomware policy is applied.
	* `status` - The status of the anti-ransomware policy.
	