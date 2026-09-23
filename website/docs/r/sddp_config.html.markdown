---
subcategory: "Data Security Center (SDDP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sddp_config"
sidebar_current: "docs-alicloud-resource-sddp-config"
description: |-
  Provides a Alicloud Data Security Center Config resource.
---

# alicloud_sddp_config

Provides a Data Security Center Config resource.

For information about Data Security Center Config and how to use it, see [What is Config](https://www.alibabacloud.com/help/en/data-security-center/latest/api-sddp-2019-01-03-createconfig).

-> **NOTE:** Available since v1.133.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_sddp_config&exampleId=614c07dc-2a51-7a56-e455-aaf1a3b7bcffcaa6f644&activeTab=example&spm=docs.r.sddp_config.0.614c07dc2a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_sddp_config" "default" {
  code  = "access_failed_cnt"
  value = 10
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_sddp_config&spm=docs.r.sddp_config.example&intl_lang=EN_US)

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

```shell
$ terraform import alicloud_sddp_config.example <code>
```
