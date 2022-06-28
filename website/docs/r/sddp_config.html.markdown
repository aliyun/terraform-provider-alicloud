---
subcategory: "Data Security Center"
layout: "alicloud"
page_title: "Alicloud: alicloud_sddp_config"
sidebar_current: "docs-alicloud-resource-sddp-config"
description: |-
  Provides a Alicloud Data Security Center Config resource.
---

# alicloud\_sddp\_config

Provides a Data Security Center Config resource.

For information about Data Security Center Config and how to use it, see [What is Config](https://help.aliyun.com/product/88674.html).

-> **NOTE:** Available in v1.133.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_sddp_config" "default" {
  code  = "access_failed_cnt"
  value = 10
}

```

## Argument Reference

The following arguments are supported:

* `code` - (Optional, ForceNew) Abnormal Alarm General Configuration Module by Using the Encoding. Valid values: `access_failed_cnt`, `access_permission_exprie_max_days`, `log_datasize_avg_days`.
* `description` - (Optional) Abnormal Alarm General Description of the Configuration Item.
* `lang` - (Optional) The language of the request and response. Valid values: `zh`,`en`.
  * `zh`: Chinese.
  * `en`: English.
* `value` - (Optional) The Specified Exception Alarm Generic by Using the Value. Code Different Values for This Parameter the Specific Meaning of Different:
      * `access_failed_cnt`: Value Represents the Non-Authorized Resource Repeatedly Attempts to Access the Threshold. 
      * `access_permission_exprie_max_days`: Value Represents the Permissions during Periods of Inactivity Exceeding a Threshold. 
      * `log_datasize_avg_days`: Value Represents the Date Certain Log Output Is Less than 10 Days before the Average Value of the Threshold.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Config. Its value is same as `code`.

## Import

Data Security Center Config can be imported using the id, e.g.

```
$ terraform import alicloud_sddp_config.example <code>
```
