---
subcategory: "DCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_dcdn_domain"
description: |-
  Provides a Alicloud DCDN Domain resource.
---

# alicloud_dcdn_domain

Provides a DCDN Domain resource.

Full station accelerated domain name.

For information about DCDN Domain and how to use it, see [What is Domain](https://www.alibabacloud.com/help/en/doc-detail/130628.htm).

-> **NOTE:** Available since v1.94.0.

-> **NOTE:** Field `force_set`, `security_token` has been removed from provider version 1.227.1.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_dcdn_domain&exampleId=1f54f4d6-43ac-4cf3-0c9c-296f5bccc569023143f4&activeTab=example&spm=docs.r.dcdn_domain.0.1f54f4d643&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "domain_name" {
  default = "tf-example.com"
}
resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_dcdn_domain" "example" {
  domain_name = "${var.domain_name}-${random_integer.default.result}"
  scope       = "overseas"
  sources {
    content  = "1.1.1.1"
    port     = "80"
    priority = "20"
    type     = "ipaddr"
    weight   = "10"
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_dcdn_domain&spm=docs.r.dcdn_domain.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `cert_id` - (Optional) The certificate ID. This parameter is required and valid only when `CertType` is set to `cas`. If you specify this parameter, an existing certificate is used. 
* `cert_name` - (Optional) The name of the new certificate. You can specify only one certificate name. This parameter is optional and valid only when `CertType` is set to `upload`. 
* `cert_region` - (Optional, Available since v1.227.1) The region of the SSL certificate. This parameter takes effect only when `CertType` is set to `cas`. Default value: **cn-hangzhou**. Valid values: **cn-hangzhou** and **ap-southeast-1**. 
* `cert_type` - (Optional) The certificate type.
  * `upload`: a user-uploaded SSL certificate.
  * `cas`: a certificate that is acquired through Certificate Management Service.

-> **NOTE:**  If the value of the `cert_type` parameter is `cas`, the `ssl_pri` parameter is not required.

* `check_url` - (Optional) The URL that is used for health checks. 
* `domain_name` - (Required, ForceNew) The accelerated domain name. You can specify multiple domain names and separate them with commas (,). You can specify up to 500 domain names in each request. The query results of multiple domain names are aggregated. If you do not specify this parameter, data of all accelerated domain names under your account is queried. 
* `env` - (Optional, Available since v1.227.1) Specifies whether the certificate is issued in canary releases. If you set this parameter to `staging`, the certificate is issued in canary releases. If you do not specify this parameter or set this parameter to other values, the certificate is officially issued. 
* `function_type` - (Optional, ForceNew, Available since v1.227.1) Computing service type. Valid values:
  - `routine`
  - `image`
  - `cloudFunction`

* `resource_group_id` - (Optional, Computed) The ID of the resource group. If you do not specify a value for this parameter, the system automatically assigns the ID of the default resource group. 
* `scene` - (Optional, ForceNew, Available since v1.227.1) The Acceleration scen. Supported:
  - `apiscene`: API acceleration.
  - `webservicescene`: accelerate website business.
  - `staticscene`: video and graphic acceleration.
  - (Empty): no scene.
* `scope` - (Optional) The region where the acceleration service is deployed. Valid values:
  - `domestic`: Chinese mainland
  - `overseas`: global (excluding mainland China)
  - `global`: global

* `sources` - (Optional) Source  See [`sources`](#sources) below.
* `ssl_pri` - (Optional) The private key. Specify the private key only if you want to enable the SSL certificate. 
* `ssl_protocol` - (Optional) Specifies whether to enable the SSL certificate. Valid values:
  - `on`
  - `off`

* `ssl_pub` - (Optional) The content of the SSL certificate. Specify the content of the SSL certificate only if you want to enable the SSL certificate. 
* `status` - (Optional, Computed) The status of the domain name. Valid values:
  - `online`: enabled
  - `offline`: disabled
  - `configuring`: configuring
  - `configure_failed`: configuration failed
  - `checking`: reviewing
  - `check_failed`: review failed

* `tags` - (Optional, ForceNew, Map, Available since v1.204.1) The tag of the resource
* `top_level_domain` - (Optional) The top-level domain. 

### `sources`

The sources supports the following:
* `content` - (Optional) The address of the source station.
* `port` - (Optional, Computed) The port number. Valid values: `443` and `80`. Default to `80`.
* `priority` - (Optional, Computed) The priority of the origin if multiple origins are specified. Default to `20`.
* `type` - (Optional) The type of the origin. Valid values:
  - `ipaddr`: The origin is configured using an IP address.
  - `domain`: The origin is configured using a domain name.
  - `oss`: The origin is configured using the Internet domain name of an Alibaba Cloud Object Storage Service (OSS) bucket.
* `weight` - (Optional) The weight of the origin if multiple origins are specified. Default to `10`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. It is the same as the `domain_name`.
* `cname` - The CNAME domain name corresponding to the accelerated domain name.
* `create_time` - The time when the accelerated domain name was created.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 30 mins) Used when create the Domain.
* `delete` - (Defaults to 30 mins) Used when delete the Domain.
* `update` - (Defaults to 30 mins) Used when update the Domain.

## Import

DCDN Domain can be imported using the id, e.g.

```shell
$ terraform import alicloud_dcdn_domain.example <id>
```