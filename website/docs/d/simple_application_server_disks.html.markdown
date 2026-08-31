---
subcategory: "Simple Application Server"
layout: "alicloud"
page_title: "Alicloud: alicloud_simple_application_server_disks"
sidebar_current: "docs-alicloud-datasource-simple-application-server-disks"
description: |-
  Provides a list of Simple Application Server Disks to the user.
---

# alicloud\_simple\_application\_server\_disks

This data source provides the Simple Application Server Disks of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.143.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_simple_application_server_disks" "ids" {
  ids = ["example_id"]
}
output "simple_application_server_disk_id_1" {
  value = data.alicloud_simple_application_server_disks.ids.disks.0.id
}

data "alicloud_simple_application_server_disks" "nameRegex" {
  name_regex = "^my-Disk"
}
output "simple_application_server_disk_id_2" {
  value = data.alicloud_simple_application_server_disks.nameRegex.disks.0.id
}

data "alicloud_simple_application_server_disks" "status" {
  status = "In_use"
}
output "simple_application_server_disk_id_3" {
  value = data.alicloud_simple_application_server_disks.status.disks.0.id
}

data "alicloud_simple_application_server_disks" "instanceId" {
  instance_id = "example_value"
}
output "simple_application_server_disk_id_4" {
  value = data.alicloud_simple_application_server_disks.instanceId.disks.0.id
}

data "alicloud_simple_application_server_disks" "diskType" {
  disk_type = "System"
}
output "simple_application_server_disk_id_5" {
  value = data.alicloud_simple_application_server_disks.diskType.disks.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Disk IDs.
* `instance_id` - (Optional, ForceNew) The ID of the simple application server to which the disk is attached.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Disk name.
* `disk_type` - (Optional, ForceNew) The type of the disk. Possible values: `System`, `Data`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the disk. Valid values: `ReIniting`, `Creating`, `In_Use`, `Available`, `Attaching`, `Detaching`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Disk names.
* `disks` - A list of Simple Application Server Disks. Each element contains the following attributes:
	* `category` - Disk type. Possible values: `ESSD`, `SSD`.
	* `create_time` - The time when the disk was created. The time follows the ISO 8601 standard in the `yyyy-MM-ddTHH:mm:ssZ` format. The time is displayed in UTC.
	* `device` - The device name of the disk on the simple application server.
	* `disk_id` - The first ID of the resource.
	* `disk_name` - The name of the resource.
	* `disk_type` - The type of the disk. Possible values: `System`, `Data`.
	* `id` - The ID of the Disk.
	* `instance_id` - Alibaba Cloud simple application server instance ID.
	* `payment_type` - The payment type of the resource. Valid values: `PayAsYouGo`, `Subscription`.
	* `size` - The size of the disk. Unit: `GB`.
	* `status` - The status of the disk. Valid values: `ReIniting`, `Creating`, `In_Use`, `Available`, `Attaching`, `Detaching`.