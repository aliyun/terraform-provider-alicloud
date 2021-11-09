---
subcategory: "Elastic Desktop Service(EDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_users"
sidebar_current: "docs-alicloud-datasource-ecd-users"
description: |-
  Provides a list of Elastic Desktop Service(EDS) Users to the user.
---

# alicloud\_ecd\_users

This data source provides the Elastic Desktop Service(EDS) Users of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.142.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecd_user" "default" {
  end_user_id = "example_value"
  email       = "your_email"
  phone       = "your_phone"
  password    = "your_password"
}
data "alicloud_ecd_users" "ids" {}
output "ecd_user_id_1" {
  value = data.alicloud_ecd_users.ids.users.0.id
}

```

## Argument Reference

The following arguments are supported:


* `ids` - (Optional, ForceNew, Computed)  A list of User IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `Unlocked`, `Locked`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `users` - A list of Ecd Users. Each element contains the following attributes:
	* `email` - The email of the user email.
	* `end_user_id` - The Username. The custom setting is composed of lowercase letters, numbers and underscores, and the length is 3~24 characters.
	* `id` - The ID of the user id.
	* `phone` - The phone of the mobile phone number.
	* `status` - The status of the resource.
	
