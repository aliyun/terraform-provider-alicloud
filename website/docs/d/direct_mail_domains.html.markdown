---
subcategory: "Direct Mail"
layout: "alicloud"
page_title: "Alicloud: alicloud_direct_mail_domains"
sidebar_current: "docs-alicloud-datasource-direct-mail-domains"
description: |-
  Provides a list of Direct Mail Domains to the user.
---

# alicloud\_direct\_mail\_domains

This data source provides the Direct Mail Domains of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.134.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_direct_mail_domains" "ids" {
  ids = ["example_id"]
}
output "direct_mail_domain_id_1" {
  value = data.alicloud_direct_mail_domains.ids.domains.0.id
}

data "alicloud_direct_mail_domains" "nameRegex" {
  name_regex = "^my-Domain"
}
output "direct_mail_domain_id_2" {
  value = data.alicloud_direct_mail_domains.nameRegex.domains.0.id
}

data "alicloud_direct_mail_domains" "example" {
  status   = "1"
  key_word = "^my-Domain"
  ids      = ["example_id"]
}
output "direct_mail_domain_id_3" {
  value = data.alicloud_direct_mail_domains.example.domains.0.id
}

```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Domain IDs.
* `key_word` - (Optional, ForceNew) domain, length `1` to `50`, including numbers or capitals or lowercase letters or `.` or `-`
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Domain name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the domain name. Valid values:`0` to `4`. `0`:Available, Passed. `1`: Unavailable, No passed. `2`: Available, cname no passed, icp no passed. `3`: Available, icp no passed. `4`: Available, cname no passed.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Domain names.
* `domains` - A list of Direct Mail Domains. Each element contains the following attributes:
	* `cname_auth_status` - Track verification.
	* `cname_confirm_status` - Indicates whether the CNAME record is successfully verified. Valid values: `0` and `1`. `0`: indicates the verification is successful. `1`: indicates that the verification fails.
	* `cname_record` - The value of the CNAME record.
	* `create_time` - The time when the DNS record was created.
	* `default_domain` - The default domain name.
	* `dns_mx` - The value of the MX record.
	* `dns_spf` - The value of the SPF record.
	* `dns_txt` - The value of the TXT ownership record.
	* `domain_id` - The ID of the domain name.
	* `domain_name` - The domain name.
	* `domain_type` - The type of the domain.
	* `icp_status` - The status of ICP filing. Valid values: `0` and `1`. `0`: indicates that the domain name is not filed. `1`: indicates that the domain name is filed. 
	* `id` - The ID of the Domain.
	* `mx_auth_status` - Indicates whether the MX record is successfully verified. Valid values: `0` and `1`. `0`: indicates the verification is successful. `1`: indicates that the verification fails.
	* `mx_record` - The MX verification record provided by Alibaba Cloud DNS.
	* `spf_auth_status` - Indicates whether the SPF record is successfully verified. Valid values: `0` and `1`. `0`: indicates the verification is successful. `1`: indicates that the verification fails.
	* `spf_record` - The SPF verification record provided by Alibaba Cloud DNS.
	* `status` - The status of the domain name. Valid values:`0` to `4`. `0`:Available, Passed. `1`: Unavailable, No passed. `2`: Available, cname no passed, icp no passed. `3`: Available, icp no passed. `4`: Available, cname no passed.
	* `tl_domain_name` - The primary domain name.
	* `tracef_record` - The CNAME verification record provided by Alibaba Cloud DNS.
