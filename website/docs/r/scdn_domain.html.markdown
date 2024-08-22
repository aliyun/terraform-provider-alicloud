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

-> **NOTE:** Alibaba Cloud SCDN has stopped new customer purchases from January 26, 2023, and you can choose to buy Alibaba Cloud DCDN products with more comprehensive acceleration and protection capabilities. If you are already a SCDN customer, you can submit a work order at any time to apply for a smooth migration to Alibaba Cloud DCDN products. In the future, we will provide better acceleration and security protection services in Alibaba Cloud DCDN, thank you for your understanding and cooperation.

-> **DEPRECATED:**  This resource has been [deprecated](https://www.aliyun.com/product/scdn) from version `1.219.0`.

## Example Usage
<div class="oics-button" style="float: right;margin: 0 0 -40px 0;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_scdn_domain&exampleId=b2c708e1-251d-310d-a410-f6411f61d04d093c1215&activeTab=example&spm=docs.r.scdn_domain.0.b2c708e125" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>

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

```shell
$ terraform import alicloud_scdn_domain.example <domain_name>
```
