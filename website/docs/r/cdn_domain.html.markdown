---
subcategory: "CDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_cdn_domain"
sidebar_current: "docs-alicloud-resource-cdn-domain"
description: |-
  Provides a CDN Accelerated Domain resource.
---

# alicloud\_cdn\_domain

-> **DEPRECATED:**  This resource is based on CDN's old version OpenAPI and it has been deprecated from version `1.34.0`.
Please use new resource [alicloud_cdn_domain_new](https://www.terraform.io/docs/providers/alicloud/r/cdn_domain_new) and its
config [alicloud_cdn_domain_config](https://www.terraform.io/docs/providers/alicloud/r/cdn_domain_config) instead.

-> **NOTE:** Because the v2014 API will be offlined on February 25, 2021, the resource also will not work on that.
If your system is still the resource, please update it to the new resource [alicloud_cdn_domain_new](https://www.terraform.io/docs/providers/alicloud/r/cdn_domain_new) as soon as possible to avoid your system unavailability after the older API is offline. If you have any questions, please submit a ticket to feedback.
Provides a CDN Accelerated Domain resource.

## Example Usage

```
# Add a CDN Accelerated Domain with configs.
resource "alicloud_cdn_domain" "domain" {
  domain_name = your_cdn_domain_name
  cdn_type    = "web"
  source_type = "domain"
  sources     = [your_cdn_domain_source1, your_cdn_domain_source2]

  // configs
  optimize_enable      = "off"
  page_compress_enable = "off"
  range_enable         = "off"
  video_seek_enable    = "off"
  block_ips            = ["1.2.3.4", "111.222.111.111"]
  parameter_filter_config {
    enable        = "on"
    hash_key_args = ["hello", "youyouyou"]
  }
  page_404_config {
    page_type       = "other"
    custom_page_url = "http://${your_cdn_domain_name}/notfound/"
  }
  refer_config {
    refer_type  = "block"
    refer_list  = ["www.xxxx.com", "www.xxxx.cn"]
    allow_empty = "off"
  }
  auth_config {
    auth_type  = "type_a"
    master_key = "helloworld1"
    slave_key  = "helloworld2"
  }
  http_header_config {
    header_key   = "Content-Type"
    header_value = "text/plain"
  }
  http_header_config {
    header_key   = "Access-Control-Allow-Origin"
    header_value = "*"
  }
  cache_config {
    cache_content = "/hello/world"
    ttl           = 1000
    cache_type    = "path"
  }
  cache_config {
    cache_content = "/hello/world/youyou"
    ttl           = 1000
    cache_type    = "path"
  }
  cache_config {
    cache_content = "txt,jpg,png"
    ttl           = 2000
    cache_type    = "suffix"
  }
}
```
## Argument Reference

The following arguments are supported:

* `domain_name` - (Required) Name of the accelerated domain. This name without suffix can have a string of 1 to 63 characters, must contain only alphanumeric characters or "-", and must not begin or end with "-", and "-" must not in the 3th and 4th character positions at the same time. Suffix `.sh` and `.tel` are not supported.
* `cdn_type` - (Required) Cdn type of the accelerated domain. Valid values are `web`, `download`, `video`, `liveStream`.
* `source_type` - (Optional) Source type of the accelerated domain. Valid values are `ipaddr`, `domain`, `oss`. You must set this parameter when `cdn_type` value is not `liveStream`.
* `source_port` - (Optional) Source port of the accelerated domain. Valid values are `80` and `443`. Default value is `80`. You must use `80` when the `source_type` is `oss`.
* `sources` - (Optional, Type: list) Sources of the accelerated domain. It's a list of domain names or IP address and consists of at most 20 items. You must set this parameter when `cdn_type` value is not `liveStream`.
* `scope` - (Optional) Scope of the accelerated domain. Valid values are `domestic`, `overseas`, `global`. Default value is `domestic`. This parameter's setting is valid Only for the international users and domestic L3 and above users .

#### Domain config

The config supports the following:

* `optimize_enable` - (Optional) Page Optimize config of the accelerated domain. Valid values are `on` and `off`. Default value is `off`. It can effectively remove the page redundant content, reduce the file size and improve the speed of distribution when this parameter value is `on`.
* `page_compress_enable` - (Optional) Page Compress config of the accelerated domain. Valid values are `on` and `off`. Default value is `off`.
* `range_enable` - (Optional) Range Source config of the accelerated domain. Valid values are `on` and `off`. Default value is `off`.
* `video_seek_enable` - (Optional) Video Seek config of the accelerated domain. Valid values are `on` and `off`. Default value is `off`.

### Block parameter_filter_config

`parameter_filter_config` - (Optional, Type: set) Parameter filter config of the accelerated domain. It's a set and consists of at most one item.
* `enable` - (Optional) This parameter indicates whether or not the `parameter_filter_config` is enable. Valid values are `on` and `off`. Default value is `off`.  
* `hash_key_args` - (Optional, Type: list) Reserved parameters of `parameter_filter_config`. It's a list of string and consists of at most 10 items.

### Block page_404_config

`page_404_config` - (Optional, Type: set) Error Page config of the accelerated domain. It's a set and consists of at most one item.
* `page_type` - (Optional) Page type of the error page. Valid values are `default`, `charity`, `other`. Default value is `default`.
* `custom_page_url` - (Optional) Custom page url of the error page. It must be the full path under the accelerated domain name. It's value must be `http://promotion.alicdn.com/help/oss/error.html` when `page_type` value is `charity` and It can not be set when `page_type` value is `default`.

### Block refer_config

`refer_config` - (Optional, Type: set) Refer anti-theft chain config of the accelerated domain. It's a set and consists of at most 1 item.
* `refer_type` - (Optional) Refer type of the refer config. Valid values are `block` and `allow`. Default value is `block`.
* `refer_list` - (Required, Type: list) A list of domain names of the refer config.
* `allow_empty` - (Optional) This parameter indicates whether or not to allow empty refer access. Valid values are `on` and `off`. Default value is `on`.

### Block auth_config

`auth_config` - (Optional, Type: set)  Auth config of the accelerated domain. It's a set and consist of at most 1 item.
* `auth_type` - (Optional) Auth type of the auth config. Valid values are  `no_auth`, `type_a`, `type_b` and `type_c`. Default value is `no_auth`.
* `master_key` - (Optional) Master authentication key of the auth config. This parameter can have a string of 6 to 32 characters and must contain only alphanumeric characters.
* `slave_key` - (Optional) Slave authentication key of the auth config. This parameter can have a string of 6 to 32 characters and must contain only alphanumeric characters.
* `timeout` - (Optional, Type: int)  Authentication cache time of the auth config. Default value is `1800`. It's value is valid only when the `auth_type` is `type_b` or `type_c`.

### Block certificate_config

`certificate_config` - (Optional, Type: set)  Certificate config of the accelerated domain. It's a set and consist of at most 1 item.
* `server_certificate_status` - (Optional) This parameter indicates whether or not enable https. Valid values are `on` and `off`. Default value is `on`.
* `server_certificate` - (Optional) The SSL server certificate string. This is required if `server_certificate_status` is `on`
* `private_key` - (Optional) The SSL private key. This is required if `server_certificate_status` is `on`

### Block http_header_config

`http_header_config` - (Optional, Type: set) Http header config of the accelerated domain. It's a set and consist of at most 8 items. The `header_key` for each item can not be repeated.
* `header_key` - (Required) Header key of the http header. Valid values are `Content-Type`, `Cache-Control`, `Content-Disposition`, `Content-Language`，`Expires`, `Access-Control-Allow-Origin`, `Access-Control-Allow-Methods` and `Access-Control-Max-Age`.
* `header_value` - (Required) Header value of the http header.

### Block cache_config

`cache_config` - (Optional, Type: set)  Cache config of the accelerated domain. It's a set and each item's `cache_content` can not be repeated.
* `cache_type` - (Required) Cache type of the cache config. Valid values are `suffix` and `path`.
* `cache_content` - (Required) Cache content of the cache config. It's value is a path string when the `cache_type` is `path`. When the `cache_type` is `suffix`, it's value is a string which contains multiple file suffixes separated by commas.
* `ttl` - (Required, Type: int) Cache time of the cache config.
* `weight` - (Optional, Type: int) Weight of the cache config. This parameter's value is between 1 and 99. Default value is `1`. The higher the value, the higher the priority.


## Attributes Reference

The following attributes are exported:

* `domain_name` - The accelerated domain name.
* `sources` - The accelerated domain sources.
* `cdn_type` - The cdn type of the accelerated domain.
* `source_type` - The source type ot the accelerated domain.
* `scope` - The accelerated domain scope.

* `optimize_enable` - The page optimize config of the accelerated domain.
* `page_compress_enable` - The page compress config of the accelerated domain.
* `range_enable` - The range source config of the accelerated domain.
* `video_seek_enable` - The video seek config of the accelerated domain.
* `parameter_filter_config` - The parameter filter config of the accelerated domain.
* `page_404_config` - The error page config of the accelerated domain.
* `refer_config` - The refer config of the accelerated domain.
* `auth_config` - The auth config of the accelerated domain.
* `http_header_config` - The http header configs of the accelerated domain.
* `cache_config` - The cache configs of the accelerated domain.
