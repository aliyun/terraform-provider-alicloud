---
subcategory: "Governance"
layout: "alicloud"
page_title: "Alicloud: alicloud_governance_baselines"
sidebar_current: "docs-alicloud-datasource-governance-baselines"
description: |-
  Provides a list of Governance Baseline owned by an Alibaba Cloud account.
---

# alicloud_governance_baselines

This data source provides Governance Baseline available to the user.[What is Baseline](https://next.api.aliyun.com/document/governance/2021-01-20/CreateAccountFactoryBaseline)

-> **NOTE:** Available since v1.228.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform_example"
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

variable "item_services" {
  default = "ACS-BP_ACCOUNT_FACTORY_SUBSCRIBE_SERVICES"
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

data "alicloud_governance_baselines" "default" {
  ids        = ["${alicloud_governance_baseline.default.id}"]
  name_regex = alicloud_governance_baseline.default.baseline_name
}

output "alicloud_governance_baseline_example_id" {
  value = data.alicloud_governance_baselines.default.baselines.0.baseline_id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew, Computed) A list of Baseline IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Baseline IDs.
* `names` - A list of name of Baselines.
* `baselines` - A list of Baseline Entries. Each element contains the following attributes:
  * `baseline_id` - Baseline ID
  * `baseline_name` - Baseline Name.
  * `description` - Baseline Description.
