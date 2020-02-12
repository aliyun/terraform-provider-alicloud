---
subcategory: "ADB"
layout: "alicloud"
page_title: "Alicloud: alicloud_adb_account"
sidebar_current: "docs-alicloud-resource-adb-account"
description: |-
  Provides a ADB account resource.
---

# alicloud\_adb\_account

Provides a ADB account resource and used to manage databases.

-> **NOTE:** Available in v1.67.0+. 

## Example Usage

```
variable "creation" {
  default = "ADB"
}

variable "name" {
  default = "adbaccountmysql"
}

data "alicloud_zones" "default" {
  available_resource_creation = "${var.creation}"
}

resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alicloud_adb_cluster" "cluster" {
  db_cluster_version      = "3.0"
  db_cluster_category     = "Cluster"
  db_cluster_network_type = "VPC"
  db_node_class           = "C8"
  db_node_count           = 2
  db_node_storage         = 200
  pay_type                = "PostPaid"
  vswitch_id              = "${alicloud_vswitch.default.id}"
  description             = "${var.name}"
}

resource "alicloud_db_account" "account" {
  db_cluster_id         = "${alicloud_db_instance.instance.id}"
  account_name          = "tftestnormal"
  account_password      = "Test12345"
  account_description   = "${var.name}"
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The Id of cluster in which account belongs.
* `account_name` - (Required, ForceNew) Operation account requiring a uniqueness check. It may consist of lower case letters, numbers, and underlines, and must start with a letter and have no more than 16 characters.
* `account_password` - (Required) Operation password. It may consist of letters, digits, or underlines, with a length of 6 to 32 characters.
* `account_description` - (Optional) Account description. It cannot begin with https://. It must start with a Chinese character or English letter. It can include Chinese and English characters, underlines (_), hyphens (-), and numbers. The length may be 2-256 characters.

## Attributes Reference

The following attributes are exported:

* `id` - The current account resource ID. Composed of instance ID and account name with format `<instance_id>:<name>`.

## Import

ADB account can be imported using the id, e.g.

```
$ terraform import alicloud_adb_account.example "am-12345:tf_account"
```