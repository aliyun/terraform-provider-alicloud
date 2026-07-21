---
subcategory: "AliKafka"
layout: "alicloud"
page_title: "Alicloud: alicloud_alikafka_sasl_users"
description: |-
  Provides a list of Alikafka Sasl Users to the user.
---

# alicloud_alikafka_sasl_users

This data source provides the Alikafka Sasl Users of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.66.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "10.4.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_alikafka_instance" "default" {
  name            = var.name
  partition_num   = 50
  disk_type       = "1"
  disk_size       = "500"
  deploy_type     = "5"
  io_max          = "20"
  spec_type       = "professional"
  service_version = "2.2.0"
  vswitch_id      = alicloud_vswitch.default.id
  security_group  = alicloud_security_group.default.id
  config          = <<EOF
  {
    "enable.acl": "true"
  }
  EOF
}

resource "alicloud_alikafka_sasl_user" "default" {
  instance_id = alicloud_alikafka_instance.default.id
  username    = var.name
  password    = "YourPassword1234!"
}

data "alicloud_alikafka_sasl_users" "ids" {
  ids         = [alicloud_alikafka_sasl_user.default.id]
  instance_id = alicloud_alikafka_sasl_user.default.instance_id
}

output "alikafka_sasl_users_id_0" {
  value = data.alicloud_alikafka_sasl_users.ids.users.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of Sasl User IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Sasl User name.
* `instance_id` - (Required, ForceNew) The ID of the instance.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Sasl User names.
* `users` - A list of Sasl Users. Each element contains the following attributes:
  * `id` - (Available since v1.260.0) The resource ID in terraform of Sasl User. It formats as `<instance_id>:<username>`.
  * `username` - The username of the user.
  * `password` - The password of the user.
  * `type` - (Available since v1.260.0) The type of the user.
