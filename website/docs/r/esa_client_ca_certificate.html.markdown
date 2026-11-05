---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_client_ca_certificate"
description: |-
  Provides a Alicloud ESA Client Ca Certificate resource.
---

# alicloud_esa_client_ca_certificate

Provides a ESA Client Ca Certificate resource.



For information about ESA Client Ca Certificate and how to use it, see [What is Client Ca Certificate](https://next.api.alibabacloud.com/document/ESA/2024-09-10/UploadClientCaCertificate).

-> **NOTE:** Available since v1.244.0.

## Example Usage

Basic Usage

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_client_ca_certificate&exampleId=a04e4aa2-56d7-da96-2234-bc7200118fa91b458e71&activeTab=example&spm=docs.r.esa_client_ca_certificate.0.a04e4aa256&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a><a href="https://api.alibabacloud.com/terraform?resource=alicloud_esa_client_ca_certificate&exampleId=a04e4aa2-56d7-da96-2234-bc7200118fa91b458e71&activeTab=example&spm=docs.r.esa_client_ca_certificate.0.a04e4aa256&intl_lang=EN_US" target="_blank" style="display: inline-block; vertical-align: middle; margin-left: 10px; line-height: 44px; font-size: 14px; color: #1890ff; text-decoration: none;">Open in AlibabaCloud</a>
</div></div>

```terraform
data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
  site_name           = "gositecdn.cn"
}

resource "alicloud_esa_client_ca_certificate" "default" {
  certificate         = "-----BEGIN CERTIFICATE-----\n****-----END CERTIFICATE-----"
  client_ca_cert_name = "example"
  site_id             = data.alicloud_esa_sites.default.sites.0.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_esa_client_ca_certificate&spm=docs.r.esa_client_ca_certificate.example&intl_lang=EN_US) | [VIEW MORE EXAMPLES (International)](https://api.alibabacloud.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_esa_client_ca_certificate&spm=docs.r.esa_client_ca_certificate.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `certificate` - (Required, ForceNew) Certificate content.
* `client_ca_cert_name` - (Optional, ForceNew) The certificate name.
* `client_ca_certificate_hostnames` - (Optional, List, Available since v1.281.0) Client Ca certificate bound domain list.
* `site_id` - (Required, ForceNew) Site Id

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Client Ca Certificate. It formats as `<site_id>:<client_ca_cert_id>`.
* `client_ca_cert_id` - ClientCaCertificate Id.
* `create_time` - Creation time.
* `status` - Certificate status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Client Ca Certificate.
* `update` - (Defaults to 5 mins) Used when update the Client Ca Certificate.
* `delete` - (Defaults to 5 mins) Used when delete the Client Ca Certificate.

## Import

ESA Client Ca Certificate can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_client_ca_certificate.example <site_id>:<client_ca_cert_id>
```
