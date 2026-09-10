---
subcategory: "Elastic Desktop Service (ECD)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_ram_directory"
sidebar_current: "docs-alicloud-resource-ecd-ram-directory"
description: |-
  Provides a Alicloud ECD Ram Directory resource.
---

# alicloud_ecd_ram_directory

Provides a ECD Ram Directory resource.

For information about ECD Ram Directory and how to use it, see [What is Ram Directory](https://www.alibabacloud.com/help/en/wuying-workspace/developer-reference/api-ecd-2020-09-30-createramdirectory).

-> **NOTE:** Available since v1.174.0.

-> **DEPRECATED:** This resource has been deprecated from version `1.239.0`.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "terraform-example"
}
data "alicloud_ecd_zones" "default" {}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_ecd_zones.default.ids.0
  vswitch_name = var.name
}

resource "alicloud_ecd_ram_directory" "default" {
  desktop_access_type = "INTERNET"
  enable_admin_access = true
  ram_directory_name  = var.name
  vswitch_ids         = [alicloud_vswitch.default.id]
}

```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ecd_ram_directory&spm=docs.r.ecd_ram_directory.example&intl_lang=EN_US)

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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Ram Directory.
* `delete` - (Defaults to 1 mins) Used when delete the Ram Directory.


## Import

ECD Ram Directory can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecd_ram_directory.example <id>
```