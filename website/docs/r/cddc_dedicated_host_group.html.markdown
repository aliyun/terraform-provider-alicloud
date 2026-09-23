---
subcategory: "ApsaraDB for MyBase (CDDC)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cddc_dedicated_host_group"
sidebar_current: "docs-alicloud-resource-cddc-dedicated-host-group"
description: |-
  Provides a Alicloud ApsaraDB for MyBase Dedicated Host Group resource.
---

# alicloud_cddc_dedicated_host_group

Provides a ApsaraDB for MyBase Dedicated Host Group resource.

For information about ApsaraDB for MyBase Dedicated Host Group and how to use it, see [What is Dedicated Host Group](https://www.alibabacloud.com/help/en/apsaradb-for-mybase/latest/creatededicatedhostgroup).

-> **NOTE:** Available since v1.132.0.

-> **DEPRECATED:**  This resource has been [deprecated](https://www.alibabacloud.com/help/en/apsaradb-for-mybase/latest/notice-stop-selling-mybase-hosted-instances-from-august-31-2023) from version `1.225.1`. 

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_cddc_dedicated_host_group" "default" {
  engine                    = "MySQL"
  vpc_id                    = alicloud_vpc.default.id
  cpu_allocation_ratio      = 101
  mem_allocation_ratio      = 50
  disk_allocation_ratio     = 200
  allocation_policy         = "Evenly"
  host_replace_policy       = "Manual"
  dedicated_host_group_desc = var.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cddc_dedicated_host_group&spm=docs.r.cddc_dedicated_host_group.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `allocation_policy` - (Optional) AThe policy that is used to allocate resources in the dedicated cluster. Valid values:`Evenly`,`Intensively`
* `cpu_allocation_ratio` - (Optional) The CPU overcommitment ratio of the dedicated cluster.Valid values: 100 to 300. Default value: 200.
* `dedicated_host_group_desc` - (Optional) The name of the dedicated cluster. The name must be 1 to 64 characters in length and can contain letters, digits, underscores (_), and hyphens (-). It must start with a letter.
* `disk_allocation_ratio` - (Optional) The Disk Allocation Ratio of the Dedicated Host Group. **NOTE:** When `engine = SQLServer`, this attribute does not support to set.
* `engine` - (Required, ForceNew) Database Engine Type.The database engine of the dedicated cluster. Valid values:`Redis`, `SQLServer`, `MySQL`, `PostgreSQL`, `MongoDB`, `alisql`, `tair`, `mssql`. **NOTE:** Since v1.210.0., the `engine = SQLServer` was deprecated.
* `host_replace_policy` - (Optional) The policy based on which the system handles host failures. Valid values:`Auto`,`Manual`
* `mem_allocation_ratio` - (Optional) The Memory Allocation Ratio of the Dedicated Host Group.
* `vpc_id` - (Required, ForceNew) The virtual private cloud (VPC) ID of the dedicated cluster.
* `open_permission` - (Optional, ForceNew, Available since v1.148.0) Whether to enable the feature that allows you to have OS permissions on the hosts in the dedicated cluster. Valid values: `true` and `false`.
**NOTE:** The `open_permission` should be `true` when `engine = "SQLServer"`

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Dedicated Host Group.

## Import

ApsaraDB for MyBase Dedicated Host Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_cddc_dedicated_host_group.example <id>
```
