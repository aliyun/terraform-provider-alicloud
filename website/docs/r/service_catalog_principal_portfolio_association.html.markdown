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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_service_catalog_principal_portfolio_association&exampleId=b8f187a1-9b3c-ad1a-d2b5-23429c69522b34bc14dd&activeTab=example&spm=docs.r.service_catalog_principal_portfolio_association.0.b8f187a19b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_service_catalog_principal_portfolio_association&spm=docs.r.service_catalog_principal_portfolio_association.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `portfolio_id` - (Required, ForceNew) Product Portfolio ID
* `principal_id` - (Required, ForceNew) RAM entity ID
* `principal_type` - (Required, ForceNew) RAM entity type

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<principal_id>:<principal_type>:<portfolio_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Principal Portfolio Association.
* `delete` - (Defaults to 5 mins) Used when delete the Principal Portfolio Association.

## Import

Service Catalog Principal Portfolio Association can be imported using the id, e.g.

```shell
$ terraform import alicloud_service_catalog_principal_portfolio_association.example <principal_id>:<principal_type>:<portfolio_id>
```