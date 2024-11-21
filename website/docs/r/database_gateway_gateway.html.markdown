---
subcategory: "Database Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_database_gateway_gateway"
sidebar_current: "docs-alicloud-resource-database-gateway-gateway"
description: |-
  Provides a Alicloud Database Gateway Gateway resource.
---

# alicloud\_database\_gateway\_gateway

Provides a Database Gateway Gateway resource.

For information about Database Gateway Gateway and how to use it, see [What is Gateway](https://www.alibabacloud.com/help/doc-detail/123415.htm).

-> **NOTE:** Available in v1.135.0+.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_database_gateway_gateway&exampleId=23e7836e-753c-bd5f-49f4-397fd8d211e7642eb1a8&activeTab=example&spm=docs.r.database_gateway_gateway.0.23e7836e75&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_database_gateway_gateway" "example" {
  gateway_name = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `gateway_desc` - (Optional) The description of Gateway.
* `gateway_name` - (Required) The name of the Gateway.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Gateway.
* `status` - The status of gateway. Valid values: `EXCEPTION`, `NEW`, `RUNNING`, `STOPPED`.

## Import

Database Gateway Gateway can be imported using the id, e.g.

```shell
$ terraform import alicloud_database_gateway_gateway.example <id>
```
