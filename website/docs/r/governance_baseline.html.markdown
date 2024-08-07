---
subcategory: "Governance"
layout: "alicloud"
page_title: "Alicloud: alicloud_governance_baseline"
description: |-
  Provides a Alicloud Governance Baseline resource.
---

# alicloud_governance_baseline

Provides a Governance Baseline resource.

Account Factory Baseline.

For information about Governance Baseline and how to use it, see [What is Baseline](https://next.api.aliyun.com/document/governance/2021-01-20/CreateAccountFactoryBaseline).

-> **NOTE:** Available since v1.228.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

variable "item_password_policy" {
  default = "ACS-BP_ACCOUNT_FACTORY_RAM_USER_PASSWORD_POLICY"
}

variable "baseline_name_update" {
  default = "tf-auto-example-baseline-update"
}

variable "item_services" {
  default = "ACS-BP_ACCOUNT_FACTORY_SUBSCRIBE_SERVICES"
}

variable "baseline_name" {
  default = "tf-auto-example-baseline"
}

variable "item_ram_security" {
  default = "ACS-BP_ACCOUNT_FACTORY_RAM_SECURITY_PREFERENCE"
}

resource "alicloud_governance_baseline" "default" {
  baseline_items {
    version = "1.0"
    name    = var.item_password_policy
    config  = jsonencode({ "MinimumPasswordLength" : 8, "RequireLowercaseCharacters" : true, "RequireUppercaseCharacters" : true, "RequireNumbers" : true, "RequireSymbols" : true, "MaxPasswordAge" : 0, "HardExpiry" : false, "PasswordReusePrevention" : 0, "MaxLoginAttempts" : 0 })
  }
  description   = var.name
  baseline_name = "${var.name}-${random_integer.default.result}"
}
```

## Argument Reference

The following arguments are supported:
* `baseline_items` - (Optional) List of baseline items.

  You can invoke [ListAccountFactoryBaselineItems](https://next.api.aliyun.com/document/governance/2021-01-20/ListAccountFactoryBaselineItems) to get a list of account factory baseline items supported by the Cloud Governance Center. See [`baseline_items`](#baseline_items) below.
* `baseline_name` - (Optional) Baseline Name.
* `description` - (Optional) Baseline Description.

### `baseline_items`

The baseline_items supports the following:
* `config` - (Optional) Baseline item configuration. The format is a JSON string.
* `name` - (Optional) The baseline item name.
* `version` - (Optional, Computed) The baseline item version.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Baseline.
* `delete` - (Defaults to 5 mins) Used when delete the Baseline.
* `update` - (Defaults to 5 mins) Used when update the Baseline.

## Import

Governance Baseline can be imported using the id, e.g.

```shell
$ terraform import alicloud_governance_baseline.example <id>
```