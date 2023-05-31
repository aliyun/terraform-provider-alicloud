---
subcategory: "CDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_cdn_domain_new"
sidebar_current: "docs-alicloud-resource-cdn-domain-new"
description: |-
  Provides a Alicloud CDN Domain resource.
---

# alicloud_cdn_domain_new

Provides a CDN Domain resource. CDN domain name.

For information about CDN Domain and how to use it, see [What is Domain](https://www.alibabacloud.com/help/en/alibaba-cloud-cdn/latest/api-doc-cdn-2018-05-10-api-doc-addcdndomain).

-> **NOTE:** Available in v1.34.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}


resource "alicloud_cdn_domain_new" "default" {
  scope       = "domestic"
  domain_name = var.name
  cdn_type    = "web"
  sources {
    type     = "ipaddr"
    content  = "1.1.1.1"
    priority = 20
    port     = 80
    weight   = 15
  }
}
```

## Argument Reference

The following arguments are supported:
* `cdn_type` - (Required, ForceNew) Cdn type of the accelerated domain. Valid values are `web`, `download`, `video`.
* `certificate_config` - (Optional, Computed) Certificate configuration. See the following `Block CertificateConfig`.
* `check_url` - (Optional, ForceNew, Available in v1.206.0+) Health test URL.
* `domain_name` - (Required, ForceNew) Name of the accelerated domain. This name without suffix can have a string of 1 to 63 characters, must contain only alphanumeric characters or "-", and must not begin or end with "-", and "-" must not in the 3th and 4th character positions at the same time. Suffix `.sh` and `.tel` are not supported.
* `resource_group_id` - (Optional, Computed, Available in v1.67.0+) The ID of the resource group.
* `scope` - (Optional, ForceNew, Computed) Scope of the accelerated domain. Valid values are `domestic`, `overseas`, `global`. Default value is `domestic`. This parameter's setting is valid Only for the international users and domestic L3 and above users. Value:
  - **domestic**: Mainland China only.
  - **overseas**: Global (excluding Mainland China).
  - **global**: global.
The default value is **domestic**.
* `sources` - (Required) The source address list of the accelerated domain. Defaults to null. See the following `Block Sources`.
* `tags` - (Optional, Map, Available in v1.55.2+) The tag of the resource.


#### Block CertificateConfig

The CertificateConfig supports the following:
* `cert_id` - (Optional, Computed, Available in v1.206.0+) The ID of the certificate. It takes effect only when CertType = cas.
* `cert_name` - (Optional, Computed) Certificate name, only flyer names are supported.
* `cert_region` - (Optional, Computed, Available in v1.206.0+) The certificate region, which takes effect only when CertType = cas, supports cn-hangzhou (domestic) and ap-southeast-1 (International), and is cn-hangzhou by default.
* `cert_type` - (Optional, Computed) Certificate type. Value:
  - **upload**: upload certificate. 
  - **cas**: Cloud Shield certificate. 
  - **free**: free certificate.
  > If the certificate type is **cas**, **PrivateKey** does not need to pass parameters.
* `private_key` - (Optional, Computed) The content of the private key. If the certificate is not enabled, you do not need to enter the content of the private key. To configure the certificate, enter the content of the private key.
* `server_certificate` - (Optional, Computed) The content of the security certificate. If the certificate is not enabled, you do not need to enter the content of the security certificate. Please enter the content of the certificate to configure the certificate.
* `server_certificate_status` - (Optional, Computed) Whether the HTTPS certificate is enabled. Value:
  - **on**(default): enabled. 
  - **off** : not enabled.

#### Block Sources

The Sources support the following:
* `content` - (Optional) The address of source. Valid values can be ip or doaminName. Each item's `content` can not be repeated.
* `port` - (Optional, Computed) The port of source. Valid values are `443` and `80`. Default value is `80`.
* `priority` - (Optional, Computed) Priority of the source. Valid values are `0` and `100`. Default value is `20`.
* `type` - (Optional) The type of the source. Valid values are `ipaddr`, `domain` and `oss`.
* `weight` - (Optional, Computed) Weight of the source. Valid values are from `0` to `100`. Default value is `10`, but if type is `ipaddr`, the value can only be `10`. .


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `cname` - The CNAME domain name corresponding to the accelerated domain name.
* `status` - The status of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Domain.
* `delete` - (Defaults to 5 mins) Used when delete the Domain.
* `update` - (Defaults to 5 mins) Used when update the Domain.

## Import

CDN Domain can be imported using the id, e.g.

```shell
$ terraform import alicloud_cdn_domain.example <id>
```