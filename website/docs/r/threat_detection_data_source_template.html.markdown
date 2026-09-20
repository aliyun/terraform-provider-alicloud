---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_data_source_template"
description: |-
  Provides a Alicloud Threat Detection Data Source Template resource.
---

# alicloud_threat_detection_data_source_template

Provides a Threat Detection Data Source Template resource.

Data Source Template.

For information about Threat Detection Data Source Template and how to use it, see [What is Data Source Template](https://next.api.alibabacloud.com/document/cloud-siem/2024-12-12/ListDataSourceTemplates).

-> **NOTE:** Available since v1.246.0.

-> **NOTE:** The `alicloud_threat_detection_data_source_template` resource does not support create or delete operations because the corresponding API is not available. You must first create a data source template in the console, then use Terraform to manage (adopt and update) its configuration by specifying `data_source_template_id`.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_threat_detection_data_source_template&activeTab=example&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "data_source_template_id" {
  type    = string
  default = "your-existing-data-source-template-id"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_threat_detection_data_source_template" "default" {
  data_source_template_id       = var.data_source_template_id
  data_source_recognize_enabled = true
  auto_scan_new                 = "enabled"
  lang                          = "zh"
  log_user_ids                  = ["user1"]
  log_region_ids                = ["cn-hangzhou"]
}
```

## Argument Reference

The following arguments are supported:

* `data_source_template_id` - (Required, ForceNew) The ID of the data source template. This field identifies an existing template to adopt and manage.
* `data_source_recognize_enabled` - (Optional) Whether to automatically discover new data sources.
* `auto_scan_new` - (Optional) Whether to automatically discover new users. Valid values: `enabled`, `disabled`.
* `data_source_template_name` - (Optional) The data source template name.
* `log_project_pattern` - (Optional) Log Service project name matching rules.
* `log_store_pattern` - (Optional) Log service LogStore name matching rules.
* `log_user_ids` - (Optional) The data batch access user ID list.
* `log_region_ids` - (Optional) List of log storage region IDs.
* `lang` - (Optional) Language.
* `role_for` - (Optional) The user ID that the administrator switches to another member's perspective.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `create_time` - Creation time.
* `update_time` - Update time.
* `data_source_from` - Data sources. Value: `center`, `custom`.
* `data_source_recognizer` - Data source recognizer.
* `region_id` - Region ID.

### Deleting `alicloud_threat_detection_data_source_template` or removing it from your configuration

Terraform cannot destroy resource `alicloud_threat_detection_data_source_template`. Terraform will remove this resource from the state file, however resources may remain.

## Import

Threat Detection Data Source Template can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_data_source_template.example <data_source_template_id>
```
