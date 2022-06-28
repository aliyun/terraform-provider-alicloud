---
subcategory: "ApsaraDB for MyBase"
layout: "alicloud"
page_title: "Alicloud: alicloud_cddc_dedicated_host_accounts"
sidebar_current: "docs-alicloud-datasource-cddc-dedicated-host-accounts"
description: |-
  Provides a list of Cddc Dedicated Host Accounts to the user.
---

# alicloud\_cddc\_dedicated\_host\_accounts

This data source provides the Cddc Dedicated Host Accounts of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.148.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cddc_dedicated_host_accounts" "ids" {}
output "cddc_dedicated_host_account_id_1" {
  value = data.alicloud_cddc_dedicated_host_accounts.ids.accounts.0.id
}
```

## Argument Reference

The following arguments are supported:

* `dedicated_host_id` - (Optional, ForceNew) The ID of the host.
* `ids` - (Optional, ForceNew, Computed)  A list of Dedicated Host Account IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Account name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Account names.
* `accounts` - A list of Cddc Dedicated Host Accounts. Each element contains the following attributes:
	* `account_name` - The name of the Dedicated host account.
	* `dedicated_host_id` - The ID of the Dedicated host.
	* `id` - The ID of the Dedicated Host Account. The value formats as `<dedicated_host_id>:<account_name>`.