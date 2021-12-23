---
subcategory: "Cloud Config"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_configuration_recorders"
sidebar_current: "docs-alicloud-datasource-config-configuration-recorders"
description: |-
    Provides a list of Config Configuration Recorders to the user.
---

# alicloud\_config\_configuration\_recorders

This data source provides the Config Configuration Recorders of the current Alibaba Cloud user.

-> **NOTE:**  Available in 1.99.0+.

-> **NOTE:** The Cloud Config region only support `cn-shanghai` and `ap-southeast-1`.

## Example Usage

```terraform
data "alicloud_config_configuration_recorders" "example" {}

output "list_of_resource_types" {
  value = data.alicloud_config_configuration_recorders.this.recorders.0.resource_types
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `recorders` - A list of Config Configuration Recorders. Each element contains the following attributes:
    * `id` - The ID of the Config Configuration Recorder. Value as the `account_id`.
    * `account_id`- The ID of the Alicloud account.
    * `organization_enable_status` - Enterprise version configuration audit enabled status.
    * `organization_master_id` - The ID of the Enterprise management account.
    * `resource_types` - A list of resource types to be monitored.
    * `status` - Status of resource monitoring.
