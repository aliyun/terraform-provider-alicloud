---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_ips_config"
description: |-
  Provides a Alicloud Cloud Firewall IPS Config resource.
---

# alicloud_cloud_firewall_ips_config

Provides a Cloud Firewall IPS Config resource.

Support interception mode modification.

For information about Cloud Firewall IPS Config and how to use it, see [What is IPS Config](https://next.api.alibabacloud.com/document/Cloudfw/2017-12-07/DescribeDefaultIPSConfig).

-> **NOTE:** Available since v1.249.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_cloud_firewall_instance" "default" {
  payment_type = "PayAsYouGo"
}

resource "alicloud_cloud_firewall_ips_config" "default" {
  lang = "zh"
  depends_on = [
    "alicloud_cloud_firewall_instance.default"
  ]
  max_sdl     = "1000"
  basic_rules = "1"
  run_mode    = "1"
  cti_rules   = "0"
  patch_rules = "0"
  rule_class  = "1"
}
```

### Deleting `alicloud_cloud_firewall_ips_config` or removing it from your configuration

Terraform cannot destroy resource `alicloud_cloud_firewall_ips_config`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `basic_rules` - (Optional, Int) Basic rule switch. Value:
  - 1: Open.
  - 0: Closed (Default).
* `cti_rules` - (Optional, Int) Threat intelligence. Value:
  - 1: Open.
  - 0: Closed (Default).
* `lang` - (Optional) Language
* `max_sdl` - (Optional, Int) Sensitive data detection Daily detection traffic limit. Defaults to 0.
* `patch_rules` - (Optional, Int) Virtual patch switch. Value:
  - 1: Open.
  - 0: Closed (Default).
* `rule_class` - (Optional, Int) The IPS rule Group. Value:
  - 1: loose rule Group.
  - 2: Medium rule Group.
  - 3: Strict rule groups.
* `run_mode` - (Optional, Int) IPS defense mode. Value:
  - 1: Intercept mode.
  - 0: Observation mode (Default).

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as ``.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `update` - (Defaults to 5 mins) Used when update the IPS Config.

## Import

Cloud Firewall IPS Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_firewall_ips_config.example 
```