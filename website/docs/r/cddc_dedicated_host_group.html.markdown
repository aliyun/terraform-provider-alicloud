---
subcategory: "ApsaraDB for MyBase"
layout: "alicloud"
page_title: "Alicloud: alicloud_cddc_dedicated_host_group"
sidebar_current: "docs-alicloud-resource-cddc-dedicated-host-group"
description: |-
  Provides a Alicloud ApsaraDB for MyBase Dedicated Host Group resource.
---

# alicloud\_cddc\_dedicated\_host\_group

Provides a ApsaraDB for MyBase Dedicated Host Group resource.

For information about ApsaraDB for MyBase Dedicated Host Group and how to use it, see [What is Dedicated Host Group](https://www.alibabacloud.com/help/doc-detail/141455.htm).

-> **NOTE:** Available in v1.132.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_vpc" "vpc" {
  vpc_name   = "tf_test_foo"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_cddc_dedicated_host_group" "default" {
  engine                    = "MongoDB"
  vpc_id                    = alicloud_vpc.vpc.id
  cpu_allocation_ratio      = 101
  mem_allocation_ratio      = 50
  disk_allocation_ratio     = 200
  allocation_policy         = "Evenly"
  host_replace_policy       = "Manual"
  dedicated_host_group_desc = "tf-testaccDesc"
}

```

## Argument Reference

The following arguments are supported:

* `allocation_policy` - (Optional) AThe policy that is used to allocate resources in the dedicated cluster. Valid values:`Evenly`,`Intensively`
* `cpu_allocation_ratio` - (Optional) The CPU overcommitment ratio of the dedicated cluster.Valid values: 100 to 300. Default value: 200.
* `dedicated_host_group_desc` - (Optional) The name of the dedicated cluster. The name must be 1 to 64 characters in length and can contain letters, digits, underscores (_), and hyphens (-). It must start with a letter.
* `disk_allocation_ratio` - (Optional) The Disk Allocation Ratio of the Dedicated Host Group. **NOTE:** When `engine = SQLServer`, this attribute does not support to set.
* `engine` - (Required, ForceNew) Database Engine Type.The database engine of the dedicated cluster. Valid values:`Redis`, `SQLServer`, `MySQL`, `PostgreSQL`, `MongoDB`
* `host_replace_policy` - (Optional) The policy based on which the system handles host failures. Valid values:`Auto`,`Manual`
* `mem_allocation_ratio` - (Optional) The Memory Allocation Ratio of the Dedicated Host Group.
* `vpc_id` - (Required, ForceNew) The virtual private cloud (VPC) ID of the dedicated cluster.
* `open_permission` - (Optional, Computed, ForceNew, Available in v1.148.0+) Whether to enable the feature that allows you to have OS permissions on the hosts in the dedicated cluster. Valid values: `true` and `false`.
**NOTE:** The `open_permission` should be `true` when `engine = "SQLServer"`

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Dedicated Host Group.

## Import

ApsaraDB for MyBase Dedicated Host Group can be imported using the id, e.g.

```
$ terraform import alicloud_cddc_dedicated_host_group.example <id>
```
