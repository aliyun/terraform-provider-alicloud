---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_vpc_endpoint_service_user"
description: |-
  Provides a Alicloud Private Link Vpc Endpoint Service User resource.
---

# alicloud_privatelink_vpc_endpoint_service_user

Provides a Private Link Vpc Endpoint Service User resource.

Endpoint service user whitelist.

For information about Private Link Vpc Endpoint Service User and how to use it, see [What is Vpc Endpoint Service User](https://www.alibabacloud.com/help/en/privatelink/latest/api-privatelink-2020-04-15-addusertovpcendpointservice).

-> **NOTE:** Available since v1.110.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_privatelink_vpc_endpoint_service_user&exampleId=4dd99c3d-edbb-793b-9ab6-dc2f06139de5919ea744&activeTab=example&spm=docs.r.privatelink_vpc_endpoint_service_user.0.4dd99c3ded&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tfexampleuser"
}

resource "alicloud_privatelink_vpc_endpoint_service" "example" {
  service_description    = var.name
  connect_bandwidth      = 103
  auto_accept_connection = false
}

resource "alicloud_ram_user" "example" {
  name         = var.name
  display_name = "user_display_name"
  mobile       = "86-18688888888"
  email        = "hello.uuu@aaa.com"
  comments     = "yoyoyo"
}

resource "alicloud_privatelink_vpc_endpoint_service_user" "example" {
  service_id = alicloud_privatelink_vpc_endpoint_service.example.id
  user_id    = alicloud_ram_user.example.id
}
```

## Argument Reference

The following arguments are supported:
* `dry_run` - (Optional) Specifies whether to perform only a dry run, without performing the actual request. Valid values:
  - `true`: performs only a dry run. The system checks the request for potential issues, including missing parameter values, incorrect request syntax, and service limits. If the request fails the dry run, an error message is returned. If the request passes the dry run, the DryRunOperation error code is returned.
  - **false (default)**: performs a dry run and performs the actual request. If the request passes the dry run, a 2xx HTTP status code is returned and the operation is performed.
* `service_id` - (Required, ForceNew) The endpoint service ID.
* `user_arn` - (Optional, Available since v1.232.0) The whitelist in the format of ARN.
* `user_id` - (Required, ForceNew) The ID of the Alibaba Cloud account in the whitelist of the endpoint service.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<service_id>:<user_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vpc Endpoint Service User.
* `delete` - (Defaults to 5 mins) Used when delete the Vpc Endpoint Service User.

## Import

Private Link Vpc Endpoint Service User can be imported using the id, e.g.

```shell
$ terraform import alicloud_privatelink_vpc_endpoint_service_user.example <service_id>:<user_id>
```