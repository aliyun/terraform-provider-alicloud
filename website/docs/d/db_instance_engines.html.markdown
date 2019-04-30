---
layout: "alicloud"
page_title: "Alicloud: alicloud_db_instance_engines"
sidebar_current: "docs-alicloud-datasource-db-instance-engines"
description: |-
    Provides a list of RDS instacne engines resource info.
---

# alicloud\_db\_instances\_engines

This data source provides the RDS instance engines resource available info of Alibaba Cloud.

-> **NOTE:** Available in v1.46.0+

## Example Usage

```tf
data "alicloud_db_instance_engines" "resources" {
  instance_charge_type = "PostPaid"
  engine               = "MySQL"
  engine_version       = "5.6"
  output_file          = "./engines.txt"
}

output "first_db_category" {
  value = "${data.alicloud_db_instance_engines.resources.instance_engines.0.category}"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Optional) The Zone to launch the DB instance.
* `instance_charge_type` - (Optional) Filter the results by charge type. Valid values: `PrePaid` and `PostPaid`. Default to `PostPaid`.
* `engine` - (Optional) Database type. Options are `MySQL`, `SQLServer`, `PostgreSQL` and `PPAS`. If no value is specified, all types are returned.
* `engine_version` - (Optional) Database version required by the user. Value options can refer to the latest docs [detail info](https://www.alibabacloud.com/help/doc-detail/26228.htm) `EngineVersion`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform apply`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `instance_engines` - A list of Rds available resource. Each element contains the following attributes:
  * `zone_id` - The Zone to launch the DB instance.
  * `engine` - Database type.
  * `engine_version` - DB Instance version.
  * `category` - DB Instance category.
  