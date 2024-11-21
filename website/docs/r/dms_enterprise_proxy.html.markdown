---
subcategory: "DMS Enterprise"
layout: "alicloud"
page_title: "Alicloud: alicloud_dms_enterprise_proxy"
sidebar_current: "docs-alicloud-resource-dms-enterprise-proxy"
description: |-
  Provides a Alicloud DMS Enterprise Proxy resource.
---

# alicloud_dms_enterprise_proxy

Provides a DMS Enterprise Proxy resource.

For information about DMS Enterprise Proxy and how to use it, see [What is Proxy](https://next.api.alibabacloud.com/document/dms-enterprise/2018-11-01/CreateProxy).

-> **NOTE:** Available since v1.188.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_dms_enterprise_proxy&exampleId=3b009d78-a79c-30ae-1c67-744f13732fba6ceb1a26&activeTab=example&spm=docs.r.dms_enterprise_proxy.0.3b009d78a7&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
}

data "alicloud_account" "current" {}
data "alicloud_regions" "default" {
  current = true
}
data "alicloud_dms_user_tenants" "default" {
  status = "ACTIVE"
}

data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones.1.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_db_zones.default.zones.1.id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_db_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  db_instance_storage_type = "cloud_essd"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
  vswitch_id               = alicloud_vswitch.default.id
  instance_name            = var.name
  security_ips             = ["100.104.5.0/24", "192.168.0.6"]
  tags = {
    Created = "TF",
    For     = "example",
  }
}

resource "alicloud_db_account" "default" {
  db_instance_id   = alicloud_db_instance.default.id
  account_name     = "tfexamplename"
  account_password = "Example12345"
  account_type     = "Normal"
}

resource "alicloud_dms_enterprise_instance" "default" {
  tid               = data.alicloud_dms_user_tenants.default.ids.0
  instance_type     = "mysql"
  instance_source   = "RDS"
  network_type      = "VPC"
  env_type          = "dev"
  host              = alicloud_db_instance.default.connection_string
  port              = 3306
  database_user     = alicloud_db_account.default.account_name
  database_password = alicloud_db_account.default.account_password
  instance_name     = var.name
  dba_uid           = data.alicloud_account.current.id
  safe_rule         = "自由操作"
  query_timeout     = 60
  export_timeout    = 600
  ecs_region        = data.alicloud_regions.default.regions.0.id
}

resource "alicloud_dms_enterprise_proxy" "default" {
  instance_id = alicloud_dms_enterprise_instance.default.instance_id
  password    = "Example12345"
  username    = "tfexamplename"
  tid         = data.alicloud_dms_user_tenants.default.ids.0
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the database instance.
* `password` - (Required, Sensitive) The password of the database account.
* `tid` - (Optional) The ID of the tenant.
* `username` - (Required, Sensitive) The username of the database account.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Proxy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Proxy.
* `delete` - (Defaults to 1 mins) Used when delete the Proxy.


## Import

DMS Enterprise Proxy can be imported using the id, e.g.

```shell
$ terraform import alicloud_dms_enterprise_proxy.example <id>
```