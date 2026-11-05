---
subcategory: "ApsaraVideo VoD (VOD)"
layout: "alicloud"
page_title: "Alicloud: alicloud_vod_domain"
description: |-
  Provides a Alicloud VOD Domain resource.
---

# alicloud_vod_domain

Provides a VOD Domain resource.

For information about VOD Domain and how to use it, see [What is Domain](https://www.alibabacloud.com/help/product/29932.html).

-> **NOTE:** Available since v1.136.0.

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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_vod_domain&spm=docs.r.vod_domain.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, ForceNew) The domain name for CDN that you want to add to ApsaraVideo VOD. Wildcard domain names are supported. Start the domain name with a period (.). Example: `.example.com.`.
* `scope` - (Optional, ForceNew) This parameter is applicable to users of level 3 or higher in mainland China and users outside mainland China. Default value: `domestic`. Valid values:
  - `domestic`: mainland China.
  - `overseas`: outside mainland China.
  - `global`: regions in and outside mainland China.
* `ssl_protocol` - (Optional, Available since v2.2.0) Specifies whether to enable the HTTPS certificate. Default value: `off`. Valid values:
  - `on`: Enabled.
  - `off`: Disabled.
* `cert_id` - (Optional, Available since v2.2.0) The ID of the certificate.
* `cert_name` - (Optional, Available since v2.2.0) The name of the certificate.
* `cert_type` - (Optional, Available since v2.2.0) The type of the certificate. Valid values:
  - `upload`: Uploaded certificate.
  - `cas`: Certificate from SSL Certificates Service.
* `cert_region` - (Optional, Available since v2.2.0) The region of the certificate. Valid values: `cn-hangzhou`, `ap-southeast-1`.
* `ssl_pub` - (Optional, Available since v2.2.0) The content of the security certificate.
* `ssl_pri` - (Optional, Sensitive, Available since v2.2.0) The content of the private key.
* `check_url` - (Optional) The URL for health checks.
* `top_level_domain` - (Optional) The top-level domain name.
* `sources` - (Required, Set) The list of origin addresses. See [`sources`](#sources) below.
* `tags` - (Optional) A mapping of tags to assign to the resource.
  - `Key`: It can be up to 64 characters in length. It cannot be a null string.
  - `Value`: It can be up to 128 characters in length. It can be a null string.

### `sources`

The sources supports the following:

* `source_type` - (Required) The origin server type. Valid values:
  * `ipaddr` - a server that you can access by using an IP address.
  * `domain` - a server that you can access by using a domain name.
  * `oss` - an Object Storage Service (OSS) bucket.
* `source_content` - (Required) The address of the origin server. You can specify an IP address or a domain name.
* `source_port` - (Required) The port number. You can specify port 443 or 80. If you specify port 443, Alibaba Cloud CDN communicates with the origin server over HTTPS. You can also customize a port.
* `source_priority` - (Optional) The priority of the origin server if multiple origin servers are specified. Valid values: `20`, `30`. Default value: `20`. A value of 20 indicates that the origin server is the primary origin server. A value of 30 indicates that the origin server is a secondary origin server.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Domain. Its value is same as `domain_name`.
* `cname` - The CNAME that is assigned to the domain name for CDN. You must add a CNAME record in the system of your Domain Name System (DNS) service provider to map the domain name for CDN to the CNAME.
* `weight` - The weight of the origin server.
* `description` - The description of the ApsaraVideo VOD accelerated domain name.
* `status` - The running status of the accelerated domain name.
* `gmt_created` - The creation time. The format is yyyy-MM-ddTHH:mm:ssZ (UTC).
* `gmt_modified` - The last modification time. The format is yyyy-MM-ddTHH:mm:ssZ (UTC).

## Import

VOD Domain can be imported using the id, e.g.

```shell
$ terraform import alicloud_vod_domain.example <domain_name>
```
