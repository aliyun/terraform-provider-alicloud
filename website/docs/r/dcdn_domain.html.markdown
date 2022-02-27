---
subcategory: "DCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_dcdn_domain"
sidebar_current: "docs-alicloud-dcdn-domain"
description: |-
  Provides an Alicloud DCDN Domain resource.
---

# alicloud\_dcdn\_domain

You can use DCDN to improve the overall performance of your website and accelerate content delivery to improve user experience. For information about Alicloud DCDN Domain and how to use it, see [What is Resource Alicloud DCDN Domain](https://www.alibabacloud.com/help/en/doc-detail/130628.htm).

-> **NOTE:** Available in v1.94.0+.

-> **NOTE:** You must activate the Dynamic Route for CDN (DCDN) service before you create an accelerated domain.

-> **NOTE:** Make sure that you have obtained an Internet content provider (ICP) filling for the accelerated domain.

-> **NOTE:** If the origin content is not saved on Alibaba Cloud, the content must be reviewed by Alibaba Cloud. The review will be completed by the next working day after you submit the application.

## Example Usage

Basic Usage

```
resource "alicloud_dcdn_domain" "example" {
  domain_name = "example.com"
  sources {
    content  = "1.1.1.1"
    port     = "80"
    priority = "20"
    type     = "ipaddr"
  }
  scope = "overseas"
}
```
## Argument Reference

The following arguments are supported:

* `cert_name` - (Optional) Indicates the name of the certificate if the HTTPS protocol is enabled.
* `cert_type` - (Optional) The type of the certificate. Valid values:
    `free`: a free certificate.
    `cas`: a certificate purchased from Alibaba Cloud SSL Certificates Service.
    `upload`: a user uploaded certificate.
* `check_url` - (Optional, ForceNew) The URL that is used to test the accessibility of the origin.
* `domain_name` - (Required, ForceNew) The name of the accelerated domain.
* `force_set` - (Optional) Specifies whether to check the certificate name for duplicates. If you set the value to 1, the system does not perform the check and overwrites the information of the existing certificate with the same name.
* `resource_group_id` - (Optional) The ID of the resource group.
* `ssl_protocol` - (Optional) Indicates whether the SSL certificate is enabled. Valid values: `on` enabled, `off` disabled.
* `ssl_pri` - (Optional) The private key. Specify this parameter only if you enable the SSL certificate.
* `ssl_pub` - (Optional) Indicates the public key of the certificate if the HTTPS protocol is enabled.
* `scope` - (Optional) The acceleration region.
* `sources` - (Required) The origin information.
* `status` - (Optional) The status of DCDN Domain. Valid values: `online`, `offline`. Default to `online`.
* `top_level_domain` - (Optional) The top-level domain name.
* `security_token` - (Optional) The top-level domain name.

### Block sources
* `content` - (Required) The origin address.
* `type` - (Required) The type of the origin. Valid values:
    `ipaddr`: The origin is configured using an IP address.
    `domain`: The origin is configured using a domain name.
    `oss`: The origin is configured using the Internet domain name of an Alibaba Cloud Object Storage Service (OSS) bucket.
* `port` - (Optional) The port number. Valid values: `443` and `80`. Default to `80`.
* `priority` - (Optional) The priority of the origin if multiple origins are specified. Default to `20`.
* `weight` - (Optional) The weight of the origin if multiple origins are specified. Default to `10`.

## Attributes Reference

* `id` -The id of the DCDN Domain. It is the same as its domain name.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when Creating DCDN domain instance. 
* `update` - (Defaults to 5 mins) Used when Creating DCDN domain instance. 
* `delete` - (Defaults to 10 mins) Used when terminating the DCDN domain instance. 

## Import

DCDN Domain can be imported using the id or DCDN Domain name, e.g.

```
$ terraform import alicloud_dcdn_domain.example example.com
```
