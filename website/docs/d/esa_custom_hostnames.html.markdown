---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_custom_hostnames"
description: |-
  Provides a list of ESA Custom Hostnames to the user.
---

# alicloud_esa_custom_hostnames

This data source provides the ESA Custom Hostnames of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.246.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_esa_custom_hostnames" "default" {
  site_id = alicloud_esa_site.default.id
  ids     = [alicloud_esa_custom_hostname.default.id]
}

output "first_hostname" {
  value = data.alicloud_esa_custom_hostnames.default.hostnames.0.hostname
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required, ForceNew) The website ID.
* `record_id` - (Optional) The ID of the bound origin record.
* `hostname` - (Optional) The custom host name.
* `hostname_regex` - (Optional) A regex string to filter the result by custom host name.
* `name_match_type` - (Optional) The search match pattern for custom hostnames. Valid values: `exact`, `fuzzy`.
* `status` - (Optional) The status of the custom hostname. Valid values: `pending`, `active`, `conflicted`, `offline`.
* `ids` - (Optional) A list of custom hostname IDs.
* `output_file` - (Optional) File name where to save data source results after running `terraform plan`.

## Attributes Reference

The following attributes are exported in addition to the `arguments` section above:

* `ids` - A list of custom hostname IDs.
* `names` - A list of custom host names.
* `hostnames` - A list of ESA Custom Hostnames. Each element contains the following attributes:
  * `hostname_id` - The ID of the custom hostname.
  * `hostname` - The custom host name.
  * `site_id` - The website ID.
  * `site_name` - The website name.
  * `record_id` - The ID of the bound origin record.
  * `record_name` - The name of the bound origin record.
  * `status` - The status of the custom hostname.
  * `ssl_flag` - The state of the SSL switch.
  * `cert_type` - The certificate type.
  * `cas_id` - The ID of the Alibaba Cloud Security certificate.
  * `cas_region` - The region where the Alibaba Cloud Security certificate is located.
  * `cert_id` - The certificate ID.
  * `cert_status` - The certificate status.
  * `cert_txt_key` - The certificate verification TXT name.
  * `cert_txt_value` - The certificate verification TXT content.
  * `cert_http_key` - The certificate verification HTTP name.
  * `cert_http_value` - The certificate verification HTTP content.
  * `cert_not_after` - The certificate expiration time.
  * `cert_apply_message` - The free certificate request error description.
  * `cert_apply_code` - The free certificate application error code.
  * `verify_host` - The attribution verification TXT name.
  * `verify_code` - The attribution verification TXT content.
  * `offline_reason` - The reason for the offline.
  * `conflict_with` - The causes of conflict.
  * `create_time` - The creation time.
  * `update_time` - The update time.
