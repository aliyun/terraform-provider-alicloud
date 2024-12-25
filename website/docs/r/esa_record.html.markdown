---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_record"
description: |-
  Provides a Alicloud ESA Record resource.
---

# alicloud_esa_record

Provides a ESA Record resource.



For information about ESA Record and how to use it, see [What is Record](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.240.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_esa_rate_plan_instance" "default" {
  type         = "NS"
  auto_renew   = "false"
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = "true"
  plan_name    = "high"
}

resource "alicloud_esa_site" "default" {
  site_name   = "idlexamplerecord.com"
  instance_id = alicloud_esa_rate_plan_instance.default.id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_record" "default" {
  data {
    value    = "www.eerrraaa.com"
    weight   = "1"
    priority = "1"
    port     = "80"
  }

  ttl         = "100"
  record_name = "_udp._sip.idlexamplerecord.com"
  comment     = "This is a remark"
  site_id     = alicloud_esa_site.default.id
  record_type = "SRV"
}
```

## Argument Reference

The following arguments are supported:
* `auth_conf` - (Optional, List) The origin authentication information of the CNAME record. See [`auth_conf`](#auth_conf) below.
* `biz_name` - (Optional) The business scenario of the record for acceleration. Valid values:

  - `image_video`: video and image.
  - `api`: API.
  - `web`: web page.
* `comment` - (Optional) The comments of the record.
* `data` - (Required, List) The DNS record information. The format of this field varies based on the record type. For more information, see [Add DNS records](https://www.alibabacloud.com/help/doc-detail/2708761.html). See [`data`](#data) below.
* `host_policy` - (Optional) The origin host policy. This policy takes effect when the record type is CNAME. You can set the policy in two modes:

  - `follow_hostname`: match the requested domain name.
  - `follow_origin_domain`: match the origin's domain name.
* `proxied` - (Optional) Filters by whether the record is proxied. Valid values:

  - `true`
  - `false`
* `record_name` - (Required, ForceNew) The record name. This parameter specifies a filter condition for the query.
* `record_type` - (Required, ForceNew) The DNS record type.
* `site_id` - (Required, ForceNew, Int) The website ID, which can be obtained by calling the [ListSites](https://www.alibabacloud.com/help/en/doc-detail/2850189.html) operation.
* `source_type` - (Optional) The origin type of the record. Only CNAME records can be filtered by using this field. Valid values:

  - `OSS`: OSS bucket.
  - `S3`: S3 bucket.
  - `LB`: load balancer.
  - `OP`: origin pool.
  - `Domain`: domain name.
* `ttl` - (Optional, Int) The TTL of the record. Unit: seconds. The range is 30 to 86,400, or 1. If the value is 1, the TTL of the record is determined by the system.

### `auth_conf`

The auth_conf supports the following:
* `access_key` - (Optional) The access key of the account to which the origin server belongs. This parameter is required when the SourceType is OSS, and AuthType is private_same_account, or when the SourceType is S3 and AuthType is private.
* `auth_type` - (Optional) The authentication type of the origin server. Different origins support different authentication types. The type of origin refers to the SourceType parameter in this operation. If the type of origin is OSS or S3, you must specify the authentication type of the origin. Valid values:

  - `public`: public read. Select this value when the origin type is OSS or S3 and the origin access is public read.
  - `private`: private read. Select this value when the origin type is S3 and the origin access is private read.
  - `private_same_account`: private read under the same account. Select this value when the origin type is OSS, the origins belong to the same Alibaba Cloud account, and the origins have private read access.
* `region` - (Optional) The version of the signature algorithm. This parameter is required when the origin type is S3 and AuthType is private. The following two types are supported:

  - `v2`
  - `v4`

  If you leave this parameter empty, the default value v4 is used.
* `secret_key` - (Optional) The secret access key of the account to which the origin server belongs. This parameter is required when the SourceType is OSS, and AuthType is private_same_account, or when the SourceType is S3 and AuthType is private.
* `version` - (Optional) The region of the origin. If the origin type is S3, you must specify this value. You can get the region information from the official website of S3.

### `data`

The data supports the following:
* `algorithm` - (Optional, Int) The encryption algorithm used for the record, specified within the range from 0 to 255. This parameter is required when you add CERT or SSHFP records.
* `certificate` - (Optional) The public key of the certificate. This parameter is required when you add CERT, SMIMEA, or TLSA records.
* `fingerprint` - (Optional) The public key fingerprint of the record. This parameter is required when you add a SSHFP record.
* `flag` - (Optional, Int) The flag bit of the record. The Flag for a CAA record indicates its priority and how it is processed, specified within the range of 0 to 255. This parameter is required when you add a CAA record.
* `key_tag` - (Optional, Int) The public key identification for the record, specified within the range of 0 to 65,535. This parameter is required when you add a CAA record.
* `matching_type` - (Optional, Int) The algorithm policy used to match or validate the certificate, specified within the range 0 to 255. This parameter is required when you add SMIMEA or TLSA records.
* `port` - (Optional, Int) The port of the record, specified within the range of 0 to 65,535. This parameter is required when you add an SRV record.
* `priority` - (Optional, Int) The priority of the record, specified within the range of 0 to 65,535. A smaller value indicates a higher priority. This parameter is required when you add MX, SRV, and URI records.
* `selector` - (Optional, Int) The type of certificate or public key, specified within the range of 0 to 255. This parameter is required when you add SMIMEA or TLSA records.
* `tag` - (Optional) The label of the record. The Tag of a CAA record indicate its specific type and usage. This parameter is required when you add a CAA record.
* `type` - (Optional, Int) The certificate type of the record (in CERT records), or the public key type (in SSHFP records). This parameter is required when you add CERT or SSHFP records.
* `usage` - (Optional, Int) The usage identifier of the record, specified within the range of 0 to 255. This parameter is required when you add SMIMEA or TLSA records.
* `value` - (Optional) The record value or part of the record content. This parameter is required when you add A/AAAA, CNAME, NS, MX, TXT, CAA, SRV, and URI records. It has different meanings based on different types of records:

  - **A/AAAA**: the IP address(es). Separate multiple IPs with commas (,). You must have at least one IPv4 address.
  - `CNAME`: the target domain name.
  - `NS`: the name servers for the domain name.
  - `MX`: a valid domain name of the target mail server.
  - `TXT`: a valid text string.
  - `CAA`: a valid domain name of the certificate authority.
  - `SRV`: a valid domain name of the target host.
  - `URI`: a valid URI string.
* `weight` - (Optional, Int) The weight of the record, specified within the range of 0 to 65,535. This parameter is required when you add SRV or URI records.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the record. The date format follows ISO8601 notation and uses UTC time. The format is yyyy-MM-ddTHH:mm:ssZ.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Record.
* `delete` - (Defaults to 5 mins) Used when delete the Record.
* `update` - (Defaults to 5 mins) Used when update the Record.

## Import

ESA Record can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_record.example <id>
```