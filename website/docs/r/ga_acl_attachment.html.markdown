---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_acl_attachment"
sidebar_current: "docs-alicloud-resource-ga-acl-attachment"
description: |-
  Provides a Alicloud Global Accelerator (GA) Acl Attachment resource.
---

# alicloud_ga_acl_attachment

Provides a Global Accelerator (GA) Acl Attachment resource.

For information about Global Accelerator (GA) Acl Attachment and how to use it, see [What is Acl Attachment](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-associateaclswithlistener).

-> **NOTE:** Available since v1.150.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ga_acl_attachment&exampleId=63ba4e5a-0874-fd03-ee1b-3964e0a191ca3b17a5f5&activeTab=example&spm=docs.r.ga_acl_attachment.0.63ba4e5a08&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_ga_accelerator" "default" {
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}

resource "alicloud_ga_bandwidth_package" "default" {
  bandwidth      = 100
  type           = "Basic"
  bandwidth_type = "Basic"
  payment_type   = "PayAsYouGo"
  billing_type   = "PayBy95"
  ratio          = 30
}

resource "alicloud_ga_bandwidth_package_attachment" "default" {
  accelerator_id       = alicloud_ga_accelerator.default.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}

resource "alicloud_ga_listener" "default" {
  accelerator_id = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  port_ranges {
    from_port = 80
    to_port   = 80
  }
}

resource "alicloud_ga_acl" "default" {
  acl_name           = "terraform-example"
  address_ip_version = "IPv4"
}

resource "alicloud_ga_acl_entry_attachment" "default" {
  acl_id            = alicloud_ga_acl.default.id
  entry             = "192.168.1.1/32"
  entry_description = "terraform-example"
}

resource "alicloud_ga_acl_attachment" "default" {
  listener_id = alicloud_ga_listener.default.id
  acl_id      = alicloud_ga_acl.default.id
  acl_type    = "white"
}
```

## Argument Reference

The following arguments are supported:

* `listener_id` - (Required, ForceNew) The ID of the listener.
* `acl_id` - (Required, ForceNew) The ID of an ACL.
* `acl_type` - (Required, ForceNew) The type of the ACL. Valid values:
  - `white`: Only requests from IP addresses or address segments in the selected access control list are forwarded. The whitelist applies to scenarios where applications only allow specific IP addresses. There are certain business risks in setting up a whitelist. Once the whitelist is set, only the IP addresses in the whitelist can access global acceleration listeners. If whitelist access is enabled, but no IP is added to the access policy group, the global acceleration listener will not forward the request.
  - `black`: All requests from IP addresses or address segments in the selected access control list are not forwarded. Blacklists are applicable to scenarios where applications restrict access to specific IP addresses. If blacklist access is enabled and no IP is added to the access policy group, the global acceleration listener forwards all requests.
* `dry_run` - (Optional, Bool) The dry run.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Acl Attachment. The value formats as `<listener_id>:<acl_id>`.
* `status` - The status of the Acl Attachment. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Acl Attachment.
* `delete` - (Defaults to 10 mins) Used when delete the Acl Attachment.

## Import

Global Accelerator (GA) Acl Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_acl_attachment.example <listener_id>:<acl_id>
```
