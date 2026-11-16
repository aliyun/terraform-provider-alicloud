---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_site_environments"
description: |-
  Provides a list of Esa Site Environments to the user.
---

# alicloud_esa_site_environments

This data source provides the Esa Site Environments of the current Alibaba Cloud user.

-> **NOTE:** Available since v2.3.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "advancedplan"
}

resource "alicloud_esa_site_environment" "default" {
  site_id               = data.alicloud_esa_sites.default.sites.0.id
  environment_name      = var.name
  rule                  = "(ip.src in {1.1.1.1})"
  next_environment_name = "Production"
  read_only             = true
}

data "alicloud_esa_site_environments" "ids" {
  ids     = [alicloud_esa_site_environment.default.id]
  site_id = alicloud_esa_site_environment.default.site_id
}

output "esa_site_environments_id_0" {
  value = data.alicloud_esa_site_environments.ids.environments.0.id
}
```

## Argument Reference

The following attributes are exported:

* `ids` - (Optional, List) A list of Site Environment IDs.
* `name_regex` - (Optional) A regex string to filter results by Site Environment name.
* `site_id` - (Required) The ID of the Site.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Site Environment names.
* `environments` - A list of Site Environments. Each element contains the following attributes:
  * `id` - The ID of the Site Environment.
  * `environment_name` - The name of the environment.
  * `rule` - The environment rule.
  * `read_only` - Indicates whether the environment is read-only.
  * `site_version` - The version number of the site.
  * `pre_site_version` - The previous version number.
  * `priority` - The priority.
  * `is_default` - Indicates whether the environment is the default environment.
  * `create_time` - The time when the environment was created.
  * `update_time` - The time when the environment was updated.
