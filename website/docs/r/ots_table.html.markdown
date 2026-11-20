---
subcategory: "Table Store (OTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ots_table"
sidebar_current: "docs-alicloud-resource-ots-table"
description: |-
  Provides an OTS (Open Table Service) table resource.
---

# alicloud_ots_table

Provides an OTS table resource.

-> **NOTE:** From Provider version 1.10.0, the provider field 'ots_instance_name' has been deprecated and
you should use resource alicloud_ots_table's new field 'instance_name' and 'table_name' to re-import this resource.

-> **NOTE:** Available since v1.9.2.

-> **NOTE:** A table can have at most one primary key with `auto_increment` enabled, the primary key with `auto_increment` enabled cannot be the first primary key (the partition key) of the table, and `enable_local_txn` cannot be enabled together with a primary key that enables `auto_increment`.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ots_table&exampleId=f6803d4c-3bb0-537e-ce18-c775f579e5d3efe0dc55&activeTab=example&spm=docs.r.ots_table.0.f6803d4c3b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "tf-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_ots_instance" "default" {
  name        = "${var.name}-${random_integer.default.result}"
  description = var.name
  accessed_by = "Any"
  tags = {
    Created = "TF",
    For     = "example",
  }
}

resource "alicloud_ots_table" "default" {
  instance_name = alicloud_ots_instance.default.name
  table_name    = "tf_example"
  time_to_live  = -1
  max_version   = 1
  enable_sse    = true
  sse_key_type  = "SSE_KMS_SERVICE"
  primary_key {
    name = "pk1"
    type = "Integer"
  }
  primary_key {
    name = "pk2"
    type = "String"
  }
  primary_key {
    name = "pk3"
    type = "Binary"
  }

  defined_column {
    name = "col1"
    type = "Integer"
  }
  defined_column {
    name = "col2"
    type = "String"
  }
  defined_column {
    name = "col3"
    type = "Binary"
  }
}
```

### Auto-Increment Primary Key

The following example creates a table whose second primary key column is an auto-increment column. The first primary key column is the partition key and cannot be auto-incremented.

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "tf-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_ots_instance" "default" {
  name        = "${var.name}-${random_integer.default.result}"
  description = var.name
  accessed_by = "Any"
}

resource "alicloud_ots_table" "default" {
  instance_name = alicloud_ots_instance.default.name
  table_name    = "terraform_example_auto_increment"
  time_to_live  = -1
  max_version   = 1

  primary_key {
    name = "pk1"
    type = "Integer"
  }
  primary_key {
    name           = "pk2"
    type           = "Integer"
    auto_increment = true
  }
}
```

### Local Transaction

The following example creates a table with local transactions enabled.

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "tf-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_ots_instance" "default" {
  name        = "${var.name}-${random_integer.default.result}"
  description = var.name
  accessed_by = "Any"
}

resource "alicloud_ots_table" "default" {
  instance_name    = alicloud_ots_instance.default.name
  table_name       = "terraform_example_local_txn"
  time_to_live     = -1
  max_version      = 1
  enable_local_txn = true

  primary_key {
    name = "pk1"
    type = "Integer"
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ots_table&spm=docs.r.ots_table.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required, ForceNew) The name of the OTS instance in which table will located.
* `table_name` - (Required, ForceNew) The table name of the OTS instance. If changed, a new table would be created.
* `primary_key` - (Required, ForceNew) The property of `TableMeta` which indicates the structure information of a table. It describes the attribute value of primary key. The number of `primary_key` should not be less than one and not be more than four. See [`primary_key`](#primary_key) below.
* `enable_local_txn` - (Optional, ForceNew, Bool, Available since v2.3.0) Specifies whether to enable local transactions for the table. Defaults to `false`. This field can only be set when the table is created; it cannot be enabled for an existing table through Terraform, and changing it forces a new resource. The value is not returned by the table describe API, so after `terraform import` it is `false` in state and declaring `true` for an imported table forces a new resource.
* `defined_column` - (Optional, Available since v1.187.0) The property of `TableMeta` which indicates the structure information of a table. It describes the attribute value of defined column. The number of `defined_column` should not be more than 32. See [`defined_column`](#defined_column) below.
* `time_to_live` - (Required) The retention time of data stored in this table (unit: second). The value maximum is 2147483647 and -1 means never expired.
* `max_version` - (Required) The maximum number of versions stored in this table. The valid value is 1-2147483647.
* `allow_update` - (Optional, Available since v1.224.0) Whether allow data update operations. Default value is true. Skipping the resource state refresh step may result in unnecessary execution plan when upgrading from an earlier version.
* `deviation_cell_version_in_sec` - (Optional, Available in 1.42.0+) The max version offset of the table. The valid value is 1-9223372036854775807. Defaults to 86400.
* `enable_sse` - (Optional, Available since v1.172.0) Whether enable OTS server side encryption. Default value is false. **Note: The parameter is immutable after resource creation.
* `sse_key_type` - (Optional, Available since v1.172.0) The key type of OTS server side encryption. `SSE_KMS_SERVICE`, `SSE_BYOK` is allowed. **Note: The parameter is immutable after resource creation.
* `sse_key_id` - (Optional, Available since v1.224.0) . The key ID of secret. `sse_key_id` is valid only when `sse_key_type` is set to `SSE_BYOK`. **Note: The parameter is immutable after resource creation.
* `sse_role_arn` - (Optional, Available since v1.224.0) The arn of role that can access kms service. `sse_role_arn` is valid only when `sse_key_type` is set to `SSE_BYOK`. **Note: The parameter is immutable after resource creation.

### `defined_column`

The defined_column supports the following:

* `name` - (Required, Available since v1.187.0) Name for defined column.
* `type` - (Required, Available since v1.187.0) Type for defined column. `Integer`, `String`, `Binary`, `Double`, `Boolean` is allowed.

### `primary_key`

The primary_key supports the following:

* `name` - (Required, ForceNew) Name for primary key.
* `type` - (Required, ForceNew) Type for primary key. Only `Integer`, `String` or `Binary` is allowed.
* `auto_increment` - (Optional, ForceNew, Bool, Available since v2.3.0) Specifies whether the primary key column is auto-incremented. Defaults to `false`. The primary key column must be of type `Integer`, cannot be the first primary key (the partition key) of the table, and only one primary key column in a table can enable `auto_increment`. Changing this field forces a new resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID. The value is `<instance_name>:<table_name>`.

## Import

OTS table can be imported using id, e.g.

```shell
$ terraform import alicloud_ots_table.table my-ots:ots_table
```
