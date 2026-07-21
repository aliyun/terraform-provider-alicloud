---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_character_set_names"
sidebar_current: "docs-alicloud-datasource-rds-character-set-names"
description: |-
  Provide character sets supported by available RDS instances for Alibaba Cloud accounts.
---

# alicloud\_rds\_character\_set\_names

This data source is the character set supported by querying RDS instances.

-> **NOTE:** Available in v1.198.0+.

## Example Usage

```
# Declare the data source
data "alicloud_rds_character_set_names" "names" {
    engine = "MySQL"
}

output "first_rds_character_set_names" {
  value = data.alicloud_rds_character_set_names.names.names.0
}
```

## Argument Reference

The following arguments are supported:

* `engine` - (Required, ForceNew) Database type. Options are `MySQL`, `SQLServer`, `PostgreSQL`, `MariaDB`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - The list of supported character sets.
