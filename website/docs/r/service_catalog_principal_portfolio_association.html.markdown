---
subcategory: "Service Catalog"
layout: "alicloud"
page_title: "Alicloud: alicloud_service_catalog_principal_portfolio_association"
description: |-
  Provides a Alicloud Service Catalog Principal Portfolio Association resource.
---

# alicloud_service_catalog_principal_portfolio_association

Provides a Service Catalog Principal Portfolio Association resource.

Principal portfolio association.

For information about Service Catalog Principal Portfolio Association and how to use it, see [What is Principal Portfolio Association](https://www.alibabacloud.com/help/en/service-catalog/developer-reference/api-servicecatalog-2021-09-01-associateprincipalwithportfolio).

-> **NOTE:** Available since v1.230.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_service_catalog_portfolio" "defaultDaXVxI" {
  provider_name  = var.name
  description    = "desc"
  portfolio_name = var.name
}

resource "alicloud_ram_role" "default48JHf4" {
  name        = var.name
  document    = <<EOF
    {
        "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Effect": "Allow",
            "Principal": {
            "Service": [
                "emr.aliyuncs.com",
                "ecs.aliyuncs.com"
            ]
            }
        }
        ],
        "Version": "1"
    }
    EOF
  description = "this is a role test."
  force       = true
}


resource "alicloud_service_catalog_principal_portfolio_association" "default" {
  principal_id   = alicloud_ram_role.default48JHf4.id
  portfolio_id   = alicloud_service_catalog_portfolio.defaultDaXVxI.id
  principal_type = "RamRole"
}
```

## Argument Reference

The following arguments are supported:
* `portfolio_id` - (Required, ForceNew) Product Portfolio ID
* `principal_id` - (Required, ForceNew) RAM entity ID
* `principal_type` - (Required, ForceNew) RAM entity type

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<principal_id>:<principal_type>:<portfolio_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Principal Portfolio Association.
* `delete` - (Defaults to 5 mins) Used when delete the Principal Portfolio Association.

## Import

Service Catalog Principal Portfolio Association can be imported using the id, e.g.

```shell
$ terraform import alicloud_service_catalog_principal_portfolio_association.example <principal_id>:<principal_type>:<portfolio_id>
```