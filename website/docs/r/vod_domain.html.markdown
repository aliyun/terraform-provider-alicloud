---
subcategory: "ApsaraVideo VoD (VOD)"
layout: "alicloud"
page_title: "Alicloud: alicloud_vod_domain"
sidebar_current: "docs-alicloud-resource-vod-domain"
description: |-
  Provides a Alicloud VOD Domain resource.
---

# alicloud_vod_domain

Provides a VOD Domain resource.

For information about VOD Domain and how to use it, see [What is Domain](https://www.alibabacloud.com/help/product/29932.html).

-> **NOTE:** Available since v1.136.0+.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vod_domain&exampleId=181ae1f4-55e8-2611-f439-041abeec270ce6d3e821&activeTab=example&spm=docs.r.vod_domain.0.181ae1f455&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-shanghai"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_vod_domain" "default" {
  domain_name = "example-${random_integer.default.result}.com"
  scope       = "domestic"
  sources {
    source_type    = "domain"
    source_content = "outin-c7405446108111ec9a7100163e0eb78b.oss-cn-beijing.aliyuncs.com"
    source_port    = "443"
  }
  tags = {
    Created = "terraform"
    For     = "example"
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, ForceNew) The domain name for CDN that you want to add to ApsaraVideo VOD. Wildcard domain names are supported. Start the domain name with a period (.). Example: `.example.com.`.
* `sources` - (Required) The information about the address of the origin server. For more information about the Sources parameter, See the following `Block sources`.
* `check_url` - (Optional) The URL that is used for health checks.
* `scope` - (Optional, ForceNew) This parameter is applicable to users of level 3 or higher in mainland China and users outside mainland China. Valid values: 
  * `domestic` - mainland China. This is the default value.
  * `overseas` - outside mainland China.
  * `global` - regions in and outside mainland China.
* `top_level_domain` - (Optional) The top-level domain name.
* `tags` - (Optional) A mapping of tags to assign to the resource.
  * `Key`: It can be up to 64 characters in length. It cannot be a null string.
  * `Value`: It can be up to 128 characters in length. It can be a null string.

#### Block sources

The sources supports the following: 

* `source_content` - The address of the origin server. You can specify an IP address or a domain name.
* `source_port` - The port number. You can specify port 443 or 80. **Default value: 80**. If you specify port 443, Alibaba Cloud CDN communicates with the origin server over HTTPS. You can also customize a port.
* `source_priority` - The priority of the origin server if multiple origin servers are specified. Valid values: `20` and `30`. **Default value: 20**. A value of 20 indicates that the origin server is the primary origin server. A value of 30 indicates that the origin server is a secondary origin server.
* `source_type` - The type of the origin server. Valid values:
  * `ipaddr` - a server that you can access by using an IP address.
  * `domain` - a server that you can access by using a domain name.
  * `oss` - an Object Storage Service (OSS) bucket.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Domain. Its value is same as `domain_name`.
* `domain_name` - The domain name for CDN.
* `description` - The description of the domain name for CDN.
* `cert_name` - The name of the certificate. The value of this parameter is returned if HTTPS is enabled.
* `cname` - The CNAME that is assigned to the domain name for CDN. You must add a CNAME record in the system of your Domain Name System (DNS) service provider to map the domain name for CDN to the CNAME.
* `gmt_created` - The time when the domain name for CDN was added. The time follows the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time is displayed in UTC.
* `gmt_modified` - The last time when the domain name for CDN was modified. The time follows the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time is displayed in UTC.
* `ssl_protocol` - Indicates whether the Secure Sockets Layer (SSL) certificate is enabled. Valid values: `on`,`off`.
* `ssl_pub` - The public key of the certificate. The value of this parameter is returned if HTTPS is enabled.
* `scope` - This parameter is applicable to users of level 3 or higher in mainland China and users outside mainland China. Valid values:
  * `domestic` - mainland China. This is the default value.
  * `overseas` - outside mainland China.
  * `global` - regions in and outside mainland China.
* `sources` - The information about the address of the origin server. For more information about the Sources parameter, See the following `Block sources`.
* `weight` - The weight of the origin server.
* `status` - The status of the domain name for CDN. Valid values:
  * `online` - indicates that the domain name is enabled. 
  * `offline` - indicates that the domain name is disabled.
  * `configuring` - indicates that the domain name is being configured.
  * `configure_failed` - indicates that the domain name failed to be configured.
  * `checking` - indicates that the domain name is under review.
  * `check_failed` - indicates that the domain name failed the review.

## Import

VOD Domain can be imported using the id, e.g.

```shell
$ terraform import alicloud_vod_domain.example <domain_name>
```
