---
subcategory: "Function Compute Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_fc_services"
sidebar_current: "docs-alicloud-datasource-fc-services"
description: |-
    Provides a list of FC services to the user.
---

# alicloud\_fc_services

This data source provides the Function Compute services of the current Alibaba Cloud user.

## Example Usage

```
data "alicloud_fc_services" "fc_services_ds" {
  name_regex = "sample_fc_service"
}

output "first_fc_service_name" {
  value = "${data.alicloud_fc_services.fc_services_ds.services.0.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by FC service name.
* `ids` (Optional, Available in 1.53.0+) - A list of FC services ids.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of FC services ids.
* `names` - A list of FC services names.
* `services` - A list of FC services. Each element contains the following attributes:
  * `id` - FC service ID.
  * `name` - FC service name.
  * `description` - FC service description.
  * `role` - FC service role ARN.
  * `internet_access` - Indicate whether the service can access to internet or not.
  * `creation_time` - FC service creation time.
  * `last_modification_time` - FC service last modification time.
  * `log_config` - A list of one element containing information about the associated log store. It contains the following attributes:
    * `project` - Log Service project name.
    * `logstore` - Log Service store name.
  * `vpc_config` - A list of one element containing information about accessible VPC resources. It contains the following attributes:
    * `vpc_id` - Associated VPC ID.
    * `vswitch_ids` - Associated VSwitch IDs.
    * `security_group_id` - Associated security group ID.
  * `nas_config` - A list of one element about the nas configuration.
    * `user_id` - The user id of the NAS file system.
    * `group_id` - The group id of the NAS file system.
    * `mount_points` - The mount points configuration, including following attributes:
      * `server_addr` - The address of the remote NAS directory.
      * `mount_dir` - The local address where to mount your remote NAS directory.
