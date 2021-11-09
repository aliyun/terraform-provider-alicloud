---
subcategory: "Data Security Center"
layout: "alicloud"
page_title: "Alicloud: alicloud_sddp_configs"
sidebar_current: "docs-alicloud-datasource-sddp-configs"
description: |-
  Provides a list of Sddp Configs to the user.
---

# alicloud\_sddp\_configs

This data source provides the Sddp Configs of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.133.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_sddp_config" "default" {
  code  = "access_failed_cnt"
  value = 10
}
data "alicloud_sddp_configs" "default" {
  ids         = [alicloud_sddp_config.default.id]
  output_file = "./t.json"
}
output "sddp_config_id" {
  value = data.alicloud_sddp_configs.default.ids
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Config IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `configs` - A list of Sddp Configs. Each element contains the following attributes:
    * `code` - Abnormal Alarm General Configuration Module by Using the Encoding.Valid values: `access_failed_cnt`, `access_permission_exprie_max_days`, `log_datasize_avg_days`.
    * `config_id` - Configure the Number.
    * `default_value` - Default Value.
    * `description` - Abnormal Alarm General Description of the Configuration Item.
    * `id` - The ID of the Config.
    * `value` - The Specified Exception Alarm Generic by Using the Value. Code Different Values for This Parameter the Specific Meaning of Different.
      * `access_failed_cnt`: Value Represents the Non-Authorized Resource Repeatedly Attempts to Access the Threshold. 
      * `access_permission_exprie_max_days`: Value Represents the Permissions during Periods of Inactivity Exceeding a Threshold. 
      * `log_datasize_avg_days`: Value Represents the Date Certain Log Output Is Less than 10 Days before the Average Value of the Threshold.
