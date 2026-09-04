---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_global_config"
description: |-
  Provides a Alicloud Threat Detection Global Config resource.
---

# alicloud_threat_detection_global_config

Provides a Threat Detection Global Config resource.

Cloud Security Information Event Management (SIEM) global configuration. The resource supports reading and updating a named global configuration entry; creation and deletion are not supported by the API, so `create` issues an update to set the value and `delete` removes the resource from state only.

For information about Threat Detection Global Config and how to use it, see [What is Global Config](https://www.alibabacloud.com/help/en/security-center/developer-reference/).

-> **NOTE:** This is a singleton-style resource. The configuration entry is identified by `global_config_name`; Terraform cannot destroy it, only remove it from the state file.

-> **NOTE:** Available since v1.241.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_threat_detection_global_config" "default" {
  global_config_name  = "alert-config"
  global_config_value = "enabled"
  lang                = "zh"
  role_for            = 0
}
```

## Argument Reference

The following arguments are supported:

* `global_config_name` - (Required, ForceNew) The name of the global configuration entry. It identifies which named configuration to manage; changing it creates a new resource.
* `global_config_value` - (Optional) The value of the global configuration entry.
* `lang` - (Optional) The language type. Valid values: `zh`, `en`.
* `role_for` - (Optional, Int) The role for user ID. Default to `0`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, which equals to the `global_config_name`.

## Import

Threat Detection Global Config can be imported using the id, e.g. the `global_config_name`, e.g.

```shell
$ terraform import alicloud_threat_detection_global_config.example alert-config
```
