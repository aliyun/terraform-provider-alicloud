---
subcategory: "Cloud Phone"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_phone_policy"
description: |-
  Provides a Alicloud Cloud Phone Policy resource.
---

# alicloud_cloud_phone_policy

Provides a Cloud Phone Policy resource.

Cloud phone policy.

For information about Cloud Phone Policy and how to use it, see [What is Policy](https://next.api.alibabacloud.com/document/eds-aic/2023-09-30/CreatePolicyGroup).

-> **NOTE:** Available since v1.243.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_phone_policy&exampleId=29c5890a-f2ac-2caf-e6ae-3589f5c071c3777f83e6&activeTab=example&spm=docs.r.cloud_phone_policy.0.29c5890af2&intl_lang=EN_US" target="_blank">
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


resource "alicloud_cloud_phone_policy" "default" {
  policy_group_name = "NewPolicyName"
  resolution_width  = "720"
  lock_resolution   = "on"
  camera_redirect   = "on"
  resolution_height = "1280"
  clipboard         = "read"
  net_redirect_policy {
    net_redirect    = "on"
    custom_proxy    = "on"
    proxy_type      = "socks5"
    host_addr       = "192.168.12.13"
    port            = "8888"
    proxy_user_name = "user1"
    proxy_password  = "123456"
  }
}
```

## Argument Reference

The following arguments are supported:
* `camera_redirect` - (Optional) Whether to turn on local camera redirection.
* `clipboard` - (Optional) Clipboard permissions.
* `lock_resolution` - (Optional) Whether to lock the resolution.
* `net_redirect_policy` - (Optional, List) Network redirection. See [`net_redirect_policy`](#net_redirect_policy) below.
* `policy_group_name` - (Optional, Computed) The policy name.
* `resolution_height` - (Optional, Int) The height of the resolution. Unit: Pixels.
* `resolution_width` - (Optional, Int) The width of the resolution. Unit: Pixels.

### `net_redirect_policy`

The net_redirect_policy supports the following:
* `custom_proxy` - (Optional) Whether to manually configure the transparent proxy.
* `host_addr` - (Optional) The transparent proxy IP address. The format is IPv4 address.
* `net_redirect` - (Optional) Whether to enable network redirection.
* `port` - (Optional) Transparent proxy port. The Port value range is 1\~ 65535.
* `proxy_password` - (Optional) The proxy password. The length range is 1\~ 256. Chinese characters and white space characters are not allowed.
* `proxy_type` - (Optional) Agent protocol type.
* `proxy_user_name` - (Optional) The proxy user name. The length range is 1\~ 256. Chinese characters and white space characters are not allowed.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Policy.
* `update` - (Defaults to 5 mins) Used when update the Policy.

## Import

Cloud Phone Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_phone_policy.example <id>
```