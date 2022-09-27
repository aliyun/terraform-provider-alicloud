---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_instance_keywords"
sidebar_current: "docs-alicloud-datasource-rds-instance_keywords"
description: |-
  Provides a list of Rds Instance Keywords to the user.
---

# alicloud\_rds\_instance\_keywords

This data source provides the Rds Instance Keywords of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.175.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_rds_instance_keywords" "example" {
  key = "account"
}
```

## Argument Reference

The following arguments are supported:

* `key` - (Required, ForceNew) The type of reserved keyword to query. Valid values: `account` and `database`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `words` - A list of Rds Instance Keywords.
  
  
  