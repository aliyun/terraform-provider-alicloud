---
subcategory: "CDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_cdn_domain_configs"
description: |-
  Provides a list of Cdn Domain Configs to the user.
---

# alicloud_cdn_domain_configs

This data source provides the Cdn Domain Configs of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.279.0.

## Example Usage

Basic Usage

```terraform
resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_cdn_domain_new" "default" {
  domain_name = "mycdndomain-${random_integer.default.result}.alicloud-provider.cn"
  cdn_type    = "web"
  scope       = "overseas"
  sources {
    content  = "1.1.1.1"
    type     = "ipaddr"
    priority = "20"
    port     = 80
    weight   = "15"
  }
}

resource "alicloud_cdn_domain_config" "default" {
  domain_name   = alicloud_cdn_domain_new.default.domain_name
  function_name = "ip_allow_list_set"
  function_args {
    arg_name  = "ip_list"
    arg_value = "110.110.110.110"
  }
}

data "alicloud_cdn_domain_configs" "ids" {
  domain_name = alicloud_cdn_domain_config.default.domain_name
  ids         = [alicloud_cdn_domain_config.default.id]
}

output "cdn_domain_configs_id_0" {
  value = data.alicloud_cdn_domain_configs.ids.configs.0.id
}
```

## Argument Reference

The following attributes are exported:

* `ids` - (Optional, List) A list of Domain Config IDs.
* `name_regex` - (Optional) A regex string to filter results by Domain Config name.
* `domain_name` - (Required) The accelerated domain name.
* `function_name` - (Optional) The names of the features. Separate multiple feature names with commas (,).
* `config_id` - (Optional) The ID of the feature configuration.
* `status` - (Optional) The status of the configuration. Valid values: `success`, `testing`, `failed`, `configuring`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Domain Config names.
* `configs` -  A list of Domain Configs. Each element contains the following attributes:
  * `id` - The ID of the Domain Config.
  * `function_name` - The name of the feature.
  * `config_id` - The ID of the feature configuration.
  * `parent_id` - The ID of the rule condition.
  * `status` - The status of the configuration.
  * `function_args` - The args of the domain config.
    * `arg_name` - The name of arg.
    * `arg_value` - The value of arg.
