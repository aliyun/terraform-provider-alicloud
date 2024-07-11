---
subcategory: "Direct Mail"
layout: "alicloud"
page_title: "Alicloud: alicloud_direct_mail_domains"
sidebar_current: "docs-alicloud-datasource-direct-mail-domains"
description: |-
  Provides a list of Direct Mail Domains to the user.
---

# alicloud_direct_mail_domains

This data source provides the Direct Mail Domains of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.134.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example.pop.com"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_direct_mail_domain" "default" {
  domain_name = var.name
}

data "alicloud_direct_mail_domains" "ids" {
  ids = [alicloud_direct_mail_domain.default.id]
}

output "direct_mail_domains_id_0" {
  value = data.alicloud_direct_mail_domains.ids.domains.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of Domain IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Domain name.
* `key_word` - (Optional, ForceNew) The domain name. It must be 1 to 50 characters in length and can contain digits, letters, periods (.), and hyphens (-).
* `status` - (Optional, ForceNew) The status of the domain name. Valid values:
  - `0`: Indicates that the domain name is verified and available.
  - `1`: Indicates that the domain name fails to be verified and is unavailable.
  - `2`: Indicates that the domain name is available, but not filed or configured with a CNAME record.
  - `3`: Indicates that the domain name is available but not filed.
  - `4`: Indicates that the domain name is available but not configured with a CNAME record.
* `enable_details` - (Optional, Bool) Whether to query the detailed list of resource attributes. Default value: `false`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Domain names.
* `domains` - A list of Domains. Each element contains the following attributes:
  * `id` - The ID of the Domain.
  * `domain_id` - The ID of the domain name.
  * `domain_name` - The domain name.
  * `domain_record` - (Available since v1.228.0) The value of the Domain record.
  * `domain_type` - The type of the domain. **Note:** `domain_type` takes effect only if `enable_details` is set to `true`.
  * `cname_auth_status` - Indicates whether your ownership of the domain is verified.
  * `cname_confirm_status` - Indicates whether the CNAME record is successfully verified. **Note:** `cname_confirm_status` takes effect only if `enable_details` is set to `true`.
  * `cname_record` - The value of the CNAME record. **Note:** `cname_record` takes effect only if `enable_details` is set to `true`.
  * `icp_status` - The status of ICP filing.
  * `mx_auth_status` - Indicates whether the MX record is successfully verified.
  * `mx_record` - The MX verification record provided by the Direct Mail console. **Note:** `mx_record` takes effect only if `enable_details` is set to `true`.
  * `spf_auth_status` - Indicates whether the SPF record is successfully verified.
  * `spf_record` - The SPF verification record provided by the Direct Mail console. **Note:** `spf_record` takes effect only if `enable_details` is set to `true`.
  * `default_domain` - The default domain name. **Note:** `default_domain` takes effect only if `enable_details` is set to `true`.
  * `host_record` - (Available since v1.228.0) The value of the host record. **Note:** `host_record` takes effect only if `enable_details` is set to `true`.
  * `dns_mx` - The MX record value resolved through public DNS. **Note:** `dns_mx` takes effect only if `enable_details` is set to `true`.
  * `dns_txt` - The TXT record value resolved through public DNS. **Note:** `dns_txt` takes effect only if `enable_details` is set to `true`.
  * `dns_spf` - The SPF record value resolved through public DNS. **Note:** `dns_spf` takes effect only if `enable_details` is set to `true`.
  * `dns_dmarc` - (Available since v1.228.0) The DMARC record value resolved through public DNS. **Note:** `dns_dmarc` takes effect only if `enable_details` is set to `true`.
  * `dkim_auth_status` - (Available since v1.228.0) The DKIM validation flag. **Note:** `dkim_auth_status` takes effect only if `enable_details` is set to `true`.
  * `dkim_rr` - (Available since v1.228.0) The DKIM Host Record. **Note:** `dkim_rr` takes effect only if `enable_details` is set to `true`.
  * `dkim_public_key` - (Available since v1.228.0) The DKIM public key. **Note:** `dkim_public_key` takes effect only if `enable_details` is set to `true`.
  * `dmarc_auth_status` - (Available since v1.228.0) The DMARC validation flag. **Note:** `dmarc_auth_status` takes effect only if `enable_details` is set to `true`.
  * `dmarc_record` - (Available since v1.228.0) The DMARC record. **Note:** `dmarc_record` takes effect only if `enable_details` is set to `true`.
  * `dmarc_host_record` - (Available since v1.228.0) The DMARC Host Record. **Note:** `dmarc_host_record` takes effect only if `enable_details` is set to `true`.
  * `tl_domain_name` - The primary domain name. **Note:** `tl_domain_name` takes effect only if `enable_details` is set to `true`.
  * `tracef_record` - The CNAME verification record provided by the Direct Mail console. **Note:** `tracef_record` takes effect only if `enable_details` is set to `true`.
  * `status` - The status of the domain name.
  * `create_time` - The time when the DNS record was created.
