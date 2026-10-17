---
subcategory: "PAIPlugin"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_plugin_content_generation_job_config"
description: |-
  Provides a Alicloud PAIPlugin Content Generation Job Config resource.
---

# alicloud_pai_plugin_content_generation_job_config

Provides a PAIPlugin Content Generation Job Config resource.

For information about PAIPlugin Content Generation Job Config and how to use it, see
[What is Content Generation Job Config](https://www.alibabacloud.com/help/en/pai/developer-reference/api-paiplugin-2022-01-12-createcontentgenerationjobconfig).

-> **NOTE:** Available since v1.249.0.

## Example Usage

```terraform
resource "alicloud_pai_plugin_content_generation_job_config" "default" {
  name                 = "tf-demo-cfg"
  type                 = "Template"
  model                = "QWen"
  industry_code        = "100"
  signature_id         = "sig-xxx"
  recommended_keywords = ["a", "b"]
  custom_keywords      = ["c"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the config.
* `type` - (Required, ForceNew) The content generation type. Currently only `Template` is supported.
* `model` - (Required, ForceNew) The model used by the config. Currently only `QWen` is supported.
* `industry_code` - (Required, ForceNew) The industry code.
* `signature_id` - (Required, ForceNew) The signature id.
* `recommended_keywords` - (Required, ForceNew) The recommended keywords list.
* `purpose` - (Optional, ForceNew) The purpose of the config. Valid values: `UserTouch`, `SmsRecall`.
* `custom_keywords` - (Optional, ForceNew) The custom keywords list.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. Same as `content_generation_job_config_id`.
* `content_generation_job_config_id` - The config id.
* `industry_name` - The industry name.
* `signature_name` - The signature name.
* `job_count` - The job count.
* `gmt_create_time` - The create time.
* `gmt_modified_time` - The modified time.

## Import

PAIPlugin Content Generation Job Config can be imported using the id, e.g.

```shell
terraform import alicloud_pai_plugin_content_generation_job_config.default <content_generation_job_config_id>
```
