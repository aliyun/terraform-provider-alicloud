---
subcategory: "SCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_scdn_domain"
sidebar_current: "docs-alicloud-resource-scdn-domain"
description: |-
  Provides a Alicloud SCDN Domain resource.
---

# alicloud\_scdn\_domain

Provides a SCDN Domain resource.

For information about SCDN Domain and how to use it, see [What is Domain](https://help.aliyun.com/document_detail/63672.html).

-> **NOTE:** Available in v1.131.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_scdn_domain" "example" {
  domain_name = "my-Domain"
  sources {
    content  = "xxx.aliyuncs.com"
    enabled  = "online"
    port     = 80
    priority = "20"
    type     = "oss"
  }
}

```

## Argument Reference

The following arguments are supported:

* `biz_name` - (Optional, Deprecated) Attribute perm has been deprecated and suggest removing it from your template.
* `cert_infos` - (Optional) Certificate Information. See the following `Block cert_infos`.
* `check_url` - (Optional) The health check url.
* `domain_name` - (Required, ForceNew) The name of domain.
* `force_set` - (Optional) Whether to set certificate forcibly.
* `resource_group_id` - (Optional, Computed) The resource group id.
* `sources` - (Required) the Origin Server Information. See the following `Block sources`.
* `status` - (Optional, Computed) The status of the resource. Valid values: `offline`, `online`.

#### Block sources

The sources supports the following: 

* `content` - (Required) The Back-to-Source Address.
* `enabled` - (Optional) The source status. Valid values: online, offline.
* `port` - (Required) Port.
* `priority` - (Required) Priority.
* `type` - (Required) The Origin Server Type. Valid Values: 
  * ipaddr: IP Source Station 
  * domain: the Domain Name 
  * oss: OSS Bucket as a Source Station.
  
#### Block function_args

The function_args supports the following: 

* `arg_name` - (Optional) The Configuration Name.
* `arg_value` - (Optional) and Leave the Configuration Values.

#### Block cert_infos

The cert_infos supports the following: 

* `cert_name` - (Optional) If You Enable HTTPS Here Certificate Name.
* `cert_type` - (Optional) Certificate Type. Value Range: 
  * upload: Certificate
  * cas: Certificate Authority Certificate. 
  * free: Free Certificate.
* `ssl_pri` - (Optional) Private Key. Do Not Enable Certificate without Entering a User Name and Configure Certificates Enter Private Key.
* `ssl_protocol` - (Optional) Whether to Enable SSL Certificate. Valid Values: on, off. Valid values: `on`, `off`.
* `ssl_pub` - (Optional) If You Enable HTTPS Here Key.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Domain. Its value is same as `domain_name`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 11 mins) Used when create the Domain.
* `delete` - (Defaults to 1 mins) Used when delete the Domain.
* `update` - (Defaults to 11 mins) Used when update the Domain.

## Import

SCDN Domain can be imported using the id, e.g.

```
$ terraform import alicloud_scdn_domain.example <domain_name>
```
