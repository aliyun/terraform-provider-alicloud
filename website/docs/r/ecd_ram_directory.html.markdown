---
subcategory: "Elastic Desktop Service(EDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_ram_directory"
sidebar_current: "docs-alicloud-resource-ecd-ram-directory"
description: |-
  Provides a Alicloud ECD Ram Directory resource.
---

# alicloud\_ecd\_ram\_directory

Provides a ECD Ram Directory resource.

For information about ECD Ram Directory and how to use it, see [What is Ram Directory](https://help.aliyun.com/document_detail/436216.html).

-> **NOTE:** Available in v1.174.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecd_zones" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_ecd_zones.default.ids.0
}
resource "alicloud_ecd_ram_directory" "default" {
  desktop_access_type    = "INTERNET"
  enable_admin_access    = "true"
  enable_internet_access = "true"
  ram_directory_name     = var.name
  vswitch_ids            = [data.alicloud_vswitches.default.ids.0]
}
```
## Argument Reference

The following arguments are supported:

* `desktop_access_type` - (Optional, ForceNew, Computed) The desktop access type. Valid values: `VPC`, `INTERNET`, `ANY`.
* `enable_admin_access` - (Optional, ForceNew, Computed) Whether to enable public network access.
* `enable_internet_access` - (Optional, ForceNew, Computed) Whether to grant local administrator rights to users who use cloud desktops.
* `ram_directory_name` - (Required, ForceNew) The name of the directory. The name must be 2 to 255 characters in length. It must start with a letter but cannot start with `http://` or `https://`. It can contain letters, digits, colons (:), underscores (_), and hyphens (-).
* `vswitch_ids` - (Required, ForceNew) List of VSwitch IDs in the directory.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Ram Directory.
* `status` - The status of directory.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Ram Directory.
* `delete` - (Defaults to 1 mins) Used when delete the Ram Directory.


## Import

ECD Ram Directory can be imported using the id, e.g.

```
$ terraform import alicloud_ecd_ram_directory.example <id>
```