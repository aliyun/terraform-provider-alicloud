---
subcategory: "Realtime Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_realtime_compute_variables"
description: |-
  Provides a list of Realtime Compute Variables to the user.
---

# alicloud_realtime_compute_variables

This data source provides the Realtime Compute Variables of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.294.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_oss_buckets" "default" {
}

resource "alicloud_vpc" "default" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default" {
  is_default   = false
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-hangzhou-i"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
}

resource "alicloud_ram_user" "default" {
  name         = var.name
  display_name = "displayname"
  mobile       = "86-18888888888"
  email        = "hello.uuu@aaa.com"
  comments     = "yoyoyo"
}

resource "alicloud_realtime_compute_vvp_instance" "default" {
  vvp_instance_name = var.name
  storage {
    oss {
      bucket = data.alicloud_oss_buckets.default.buckets.0.name
    }
  }
  vpc_id      = alicloud_vpc.default.id
  vswitch_ids = [alicloud_vswitch.default.id]
  resource_spec {
    cpu       = "8"
    memory_gb = "32"
  }
  payment_type = "PayAsYouGo"
  zone_id      = alicloud_vswitch.default.zone_id
}

resource "alicloud_realtime_compute_variable" "default" {
  name        = var.name
  namespace   = "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default"
  workspace   = alicloud_realtime_compute_vvp_instance.default.resource_id
  kind        = "Clear"
  value       = "YourPassword123!"
  description = var.name
}

data "alicloud_realtime_compute_variables" "ids" {
  ids       = [alicloud_realtime_compute_variable.default.id]
  workspace = alicloud_realtime_compute_variable.default.workspace
  namespace = alicloud_realtime_compute_variable.default.namespace
}

output "realtime_compute_variables_id_0" {
  value = data.alicloud_realtime_compute_variables.ids.variables.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, List) A list of Variable IDs. It formats as `<workspace>:<namespace>:<name>`.
* `name_regex` - (Optional) A regex string to filter results by Variable name.
* `namespace` - (Required, ForceNew) The name of the namespace.
* `workspace` - (Required, ForceNew) The ID of the workspace.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Variable names.
* `variables` - A list of Variables. Each element contains the following attributes:
  * `id` - The ID of the Variable.
  * `workspace` - The workspace ID.
  * `namespace` - The name of the namespace.
  * `name` - The name of the variable.
  * `kind` - The kind of the variable, currently supports Plain.
  * `value` - The value of the variable.
  * `description` - The description of the variable.
