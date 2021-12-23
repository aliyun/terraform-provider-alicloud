---
subcategory: "SSL Certificates"
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_certificates_service_certificates"
sidebar_current: "docs-alicloud-datasource-ssl-certificates-service-certificates"
description: |-
  Provides a list of Ssl Certificates Service Certificates to the user.
---

# alicloud\_ssl\_certificates\_service\_certificates

This data source provides the Ssl Certificates Service Certificates of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.129.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cas_certificates" "certs" {
  name_regex = "^cas"
  ids        = ["Certificate Id"]
}

output "cert" {
  value = data.alicloud_cas_certificates.certs.certificates.0.id
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Certificate IDs.
* `lang` - (Optional, ForceNew) The lang.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Certificate name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Certificate names.
* `ids`   - A list of Certificate IDs.
* `certificates` - A list of Ssl Certificates Service Certificates. Each element contains the following attributes:
	* `buy_in_aliyun` - The cert is buy from aliyun or not.
	* `cert` - The cert's Cert.
    * `key` - The cert's Keye.
	* `cert_id` - The cert's id.
	* `certificate_name` - The cert's name.
	* `city` - The cert's city.
	* `common` - The cert's common name.
	* `country` - The cert's country.
	* `end_date` - The cert's not valid after time.
	* `expired` - The cert is expired or not.
	* `fingerprint` - The cert's finger.
	* `id` - The cert's id.
	* `issuer` - The cert's Issuer.
	* `org_name` - The cert's organization.
	* `province` - The cert's province.
	* `sans` - The cert's subject alternative name.
	* `start_date` - The cert's not valid before time.
