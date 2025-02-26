---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_cname"
description: |-
  Provides a Alicloud OSS Bucket Cname resource.
---

# alicloud_oss_bucket_cname

Provides a OSS Bucket Cname resource.

Customizing Bucket domains.

For information about OSS Bucket Cname and how to use it, see [What is Bucket Cname](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.233.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket_cname&exampleId=92762852-0dd7-ec06-735a-62a9c3fc60b9ba3ae5bd&activeTab=example&spm=docs.r.oss_bucket_cname.0.927628520d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_oss_bucket" "CreateBucket" {
  bucket        = var.name
  storage_class = "Standard"
}

resource "alicloud_oss_bucket_cname_token" "defaultZaWJfG" {
  bucket = alicloud_oss_bucket.CreateBucket.bucket
  domain = "tftestacc.com"
}

resource "alicloud_alidns_record" "defaultnHqm5p" {
  status      = "ENABLE"
  line        = "default"
  rr          = "_dnsauth"
  type        = "TXT"
  domain_name = "tftestacc.com"
  priority    = "1"
  value       = alicloud_oss_bucket_cname_token.defaultZaWJfG.token
  ttl         = "600"
  lifecycle {
    ignore_changes = [
      value,
    ]
  }
}

resource "alicloud_oss_bucket_cname" "default" {
  bucket = alicloud_oss_bucket.CreateBucket.bucket
  domain = alicloud_alidns_record.defaultnHqm5p.domain_name
}
```

## Argument Reference

The following arguments are supported:
* `bucket` - (Required, ForceNew) The bucket to which the custom domain name belongs
* `certificate` - (Optional, List) The container for the certificate configuration. See [`certificate`](#certificate) below.
* `delete_certificate` - (Optional) Whether to delete the certificate.
* `domain` - (Required, ForceNew) User-defined domain name
* `force` - (Optional) Whether to force overwrite certificate.
* `previous_cert_id` - (Optional) The current certificate ID. If the Force value is not true, the OSS Server checks whether the value matches the current certificate ID. If the value does not match, an error is reported.

### `certificate`

The certificate supports the following:
* `cert_id` - (Optional, Computed) Certificate Identifier
* `certificate` - (Optional) The certificate public key.
* `private_key` - (Optional) The certificate private key.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<bucket>:<domain>`.
* `certificate` - The container for the certificate configuration.
  * `creation_date` - Certificate creation time
  * `fingerprint` - Certificate Fingerprint
  * `status` - Certificate Status
  * `type` - Certificate Type
  * `valid_end_date` - Certificate validity period end time
  * `valid_start_date` - Certificate validity period start time
* `status` - Cname status

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Cname.
* `delete` - (Defaults to 5 mins) Used when delete the Bucket Cname.
* `update` - (Defaults to 5 mins) Used when update the Bucket Cname.

## Import

OSS Bucket Cname can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_cname.example <bucket>:<domain>
```