---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_custom_hostname"
description: |-
  Provides a Alicloud ESA Custom Hostname resource.
---

# alicloud_esa_custom_hostname

Provides a ESA Custom Hostname resource.

For information about ESA Custom Hostname and how to use it, see [What is CustomHostname](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateCustomHostname).

-> **NOTE:** Available since v1.246.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-esa-custom-hostname.com"
}

resource "alicloud_esa_rate_plan_instance" "default" {
  type         = "NS"
  auto_renew   = "false"
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = "true"
  plan_name    = "high"
}

resource "alicloud_esa_site" "default" {
  site_name   = var.name
  instance_id = alicloud_esa_rate_plan_instance.default.id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_record" "default" {
  site_id     = alicloud_esa_site.default.id
  record_name = "www.${var.name}"
  record_type = "CNAME"
  source_type = "S3"
  data {
    value = "www.example.com"
  }
  biz_name    = "api"
  host_policy = "follow_hostname"
  ttl         = "100"
}

resource "alicloud_esa_custom_hostname" "default" {
  site_id   = alicloud_esa_site.default.id
  hostname  = "custom.${var.name}"
  record_id = alicloud_esa_record.default.id
  ssl_flag  = "on"
  cert_type = "free"
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required, ForceNew) The website ID. The value can be obtained from the `id` attribute of `alicloud_esa_site`.
* `hostname` - (Required, ForceNew) The user-defined custom host name.
* `record_id` - (Required) The ID of the bound origin record.
* `ssl_flag` - (Required) The state of the SSL switch. Valid values: `on`, `off`.
* `cert_type` - (Required) The certificate type. Valid values: `free`, `upload`, `cas`.
* `certificate` - (Optional, Sensitive) The public key of the uploaded certificate. Required when `cert_type` is `upload`.
* `private_key` - (Optional, Sensitive) The private key of the uploaded certificate. Required when `cert_type` is `upload`.
* `cas_id` - (Optional) The ID of the Alibaba Cloud Security certificate. Required when `cert_type` is `cas`.
* `cas_region` - (Optional, Computed) The region where the Alibaba Cloud Security certificate is located. Required when `cert_type` is `cas`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of the custom hostname. The value is `HostnameId`.
* `hostname_id` - The ID of the custom hostname.
* `status` - The status of the custom hostname. Valid values: `pending`, `active`, `conflicted`, `offline`.
* `site_name` - The website name.
* `record_name` - The name of the bound origin record.
* `cert_id` - The certificate ID.
* `cert_status` - The certificate status.
* `cert_txt_key` - The certificate verification TXT name.
* `cert_txt_value` - The certificate verification TXT content.
* `cert_http_key` - The certificate verification HTTP name.
* `cert_http_value` - The certificate verification HTTP content.
* `cert_not_after` - The certificate expiration time.
* `cert_apply_message` - The free certificate request error description.
* `cert_apply_code` - The free certificate application error code.
* `verify_host` - The attribution verification TXT name.
* `verify_code` - The attribution verification TXT content.
* `offline_reason` - The reason for the offline.
* `conflict_with` - The causes of conflict.
* `create_time` - The creation time.
* `update_time` - The update time.
* `region_id` - The region ID.
* `name_match_type` - The search match pattern for custom hostnames.

## Import

ESA Custom Hostname can be imported using the `id`, e.g.

```shell
terraform import alicloud_esa_custom_hostname.example <hostname_id>
```
