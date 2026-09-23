---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_web_lock_configs"
sidebar_current: "docs-alicloud-datasource-threat_detection-web-lock-configs"
description: |-
  Provides a list of Threat Detection Web Lock Config owned by an Alibaba Cloud account.
---

# alicloud_threat_detection_web_lock_configs

This data source provides Threat Detection Web Lock Config available to the user.[What is Web Lock Config](https://www.alibabacloud.com/help/en/security-center/latest/api-sas-2018-12-03-describeweblockconfiglist)

-> **NOTE:** Available since v1.195.0.

## Example Usage

```terraform
data "alicloud_threat_detection_web_lock_configs" "default" {
  ids = ["${alicloud_threat_detection_web_lock_config.default.id}"]
}

output "alicloud_threat_detection_web_lock_config_example_id" {
  value = data.alicloud_threat_detection_web_lock_configs.default.configs.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew, Computed) A list of Web Lock Config IDs.
* `lang` - (Optional, ForceNew) The language of the content within the request and the response. Valid values: `zh`, `en`.
* `remark` - (Optional, ForceNew) The string that allows you to search for servers in fuzzy match mode. You can enter a server name or IP address.
* `source_ip` - (Optional, ForceNew) The source IP address of the request.
* `status` - (Optional, ForceNew) The protection status of the server that you want to query. Valid values: `on`, `off`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Web Lock Config IDs.
* `configs` - A list of Web Lock Config Entries. Each element contains the following attributes:
  * `uuid` - The UUID of the server that has web tamper proofing enabled.
  * `id` - The ID of the resource.
  * `defence_mode` - The prevention mode.
  * `dir` - The directory that has web tamper proofing enabled.
  * `exclusive_dir` - The directory that has web tamper proofing disabled.
  * `exclusive_file` - The file that has web tamper proofing disabled. **Note:** If the value of `mode` is `blacklist`, this parameter is returned.
  * `exclusive_file_type` - The type of the file that has web tamper proofing disabled. **Note:** If the value of `mode` is `blacklist`, this parameter is returned.
  * `inclusive_file_type` - The type of the file that has web tamper proofing enabled. **Note:** If the value of `mode` is `whitelist`, this parameter is returned.
  * `local_backup_dir` - The local path to the backup files of the protected directory.
  * `mode` - The protection mode of web tamper proofing. 
