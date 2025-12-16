---
subcategory: "ApsaraDB for MyBase (CDDC)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cddc_dedicated_host_account"
sidebar_current: "docs-alicloud-resource-cddc-dedicated-host-account"
description: |-
  Provides a Alicloud ApsaraDB for MyBase Dedicated Host Account resource.
---

# alicloud_cddc_dedicated_host_account

Provides a ApsaraDB for MyBase Dedicated Host Account resource.

For information about ApsaraDB for MyBase Dedicated Host Account and how to use it, see [What is Dedicated Host Account](https://www.alibabacloud.com/help/en/apsaradb-for-mybase/latest/creatededicatedhostaccount).

-> **NOTE:** Available since v1.148.0.

-> **NOTE:** Each Dedicated host can have only one account. Before you create an account for a host, make sure that the existing account is deleted.

-> **DEPRECATED:**  This resource has been [deprecated](https://www.alibabacloud.com/help/en/apsaradb-for-mybase/latest/notice-stop-selling-mybase-hosted-instances-from-august-31-2023) from version `1.225.1`. 

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}
data "alicloud_cddc_zones" "default" {}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_cddc_zones.default.ids.0
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
  open_permission           = true
}

data "alicloud_cddc_host_ecs_level_infos" "default" {
  db_type      = "mysql"
  zone_id      = data.alicloud_cddc_zones.default.ids.0
  storage_type = "cloud_essd"
}

resource "alicloud_cddc_dedicated_host" "default" {
  host_name               = var.name
  dedicated_host_group_id = alicloud_cddc_dedicated_host_group.default.id
  host_class              = data.alicloud_cddc_host_ecs_level_infos.default.infos.0.res_class_code
  zone_id                 = data.alicloud_cddc_zones.default.ids.0
  vswitch_id              = alicloud_vswitch.default.id
  payment_type            = "Subscription"
  tags = {
    Created = "TF"
    For     = "CDDC_DEDICATED"
  }
}

resource "alicloud_cddc_dedicated_host_account" "default" {
  account_name      = var.name
  account_password  = "Password1234"
  dedicated_host_id = alicloud_cddc_dedicated_host.default.dedicated_host_id
  account_type      = "Normal"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cddc_dedicated_host_account&spm=docs.r.cddc_dedicated_host_account.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `account_name` - (Required, ForceNew) The name of the Dedicated host account. The account name must be 2 to 16 characters in length, contain lower case letters, digits, and underscore(_). At the same time, the name must start with a letter and end with a letter or number.
* `account_password` - (Required, Sensitive) The password of the Dedicated host account. The account password must be 6 to 32 characters in length, and can contain letters, digits, and special characters `!@#$%^&*()_+-=`.
* `account_type` - (Optional, ForceNew) The type of the Dedicated host account. Valid values: `Admin`, `Normal`.
* `dedicated_host_id` - (Required, ForceNew) The ID of Dedicated the host.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Dedicated Host Account. The value formats as `<dedicated_host_id>:<account_name>`.

## Import

ApsaraDB for MyBase Dedicated Host Account can be imported using the id, e.g.

```shell
$ terraform import alicloud_cddc_dedicated_host_account.example <dedicated_host_id>:<account_name>
```