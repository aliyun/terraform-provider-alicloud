---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_instance_keywords"
sidebar_current: "docs-alicloud-datasource-instance-keywords"
description: |-
  Operation to query the reserved keywords of an ApsaraDB RDS instance. The reserved keywords cannot be used for the usernames of accounts or the names of databases.
---

# alicloud\_instance\_keywords

Operation to query the reserved keywords of an ApsaraDB RDS instance. The reserved keywords cannot be used for the usernames of accounts or the names of databases.

-> **NOTE:** Available in v1.196.0+

## Example Usage

```
data "alicloud_instance_keywords" "resources" {
  key         = "account"
  output_file = "./classes.txt"
}

output "account_keywords" {
  value = "${data.alicloud_instance_keywords.resources.keywords.0}"
}
```

## Argument Reference

The following arguments are supported:

* `key` - (Required, ForceNew) The type of reserved keyword to query. Valid values: `account`, `database`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform apply`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of keywords.
* `keywords` - An array that consists of reserved keywords.